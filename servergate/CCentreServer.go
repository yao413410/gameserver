package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CCentreServer struct {
	m_strIp string
	m_iPort int
	m_conn  net.Conn
}

func (p *CCentreServer) Init() {
	netcmd.AddCmdData(netcmd.NETDIALOK_1, p.LoginOk)

	netcmd.AddCmdData(int(netproto.CmdDefine_t2e_login), p.NetCmdLogin)
	netcmd.AddCmdData(int(netproto.CmdDefine_t2e_game_login), p.NetCmdGameLogin)

	netcmd.AddCmdData(int(netproto.CmdDefine_t2e_player_login_ok), p.NetCmdPlayerLoginOk)

	Println("login centre", p.m_strIp, p.m_iPort)
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_1)
}

func (p *CCentreServer) Write(b []byte) {
	if p.m_conn != nil {
		p.m_conn.Write(b)
	}
}

func (p *CCentreServer) CmdError(conn net.Conn, data *netcmd.CmdData) {
	//断线重连
	if conn != p.m_conn {
		return
	}
	p.m_conn = nil
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_1)
	PrintfWarning("Login again centre server ip=%s", p.m_strIp)
}

func (p *CCentreServer) LoginOk(conn net.Conn, data *netcmd.CmdData) error {
	Println("CCentreServer LoginOk")
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Gate
	login.ServerId = int32(g_iServerId)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CCentreServer ServerLogin %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_e2t_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())
	return nil
}

func (p *CCentreServer) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CCentreServer NetCmdLogin %t", err)
	}
	loginOk := &netproto.ServerLoginOk{}
	err = proto.Unmarshal(bytes, loginOk)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CCentreServer NetCmdLogin %t", err)
	}

	Println("CCentreServer NetCmdLogin")
	if len(loginOk.ServerGameList) > 0 {
		for _, game := range loginOk.ServerGameList {
			if game.ServerPWD == "123456" {
				g_GameServer.AddNewServer(int(game.ServerId), game.ServerIp, int(game.ServerPort))
			}
		}
	}

	p.m_conn = conn

	Println("login centre ok", p.m_strIp, p.m_iPort)
	return nil
}
func (p *CCentreServer) NetCmdGameLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CCentreServer NetCmdGameLogin %t", err)
	}
	game := &netproto.ServerGame{}
	err = proto.Unmarshal(bytes, game)
	if err != nil {
		return fmt.Errorf("CCentreServer NetCmdGameLogin %t", err)
	}
	if game.ServerPWD != "123456" {
		return fmt.Errorf("CCentreServer NetCmdGameLogin pwd is not")
	}
	g_GameServer.AddNewServer(int(game.ServerId), game.ServerIp, int(game.ServerPort))
	return nil
}

func (p *CCentreServer) NetCmdPlayerLoginOk(conn net.Conn, data *netcmd.CmdData) error {
	playerid := data.GetInt()
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CCentreServer NetCmdGameLogin %t", err)
	}

	loginok := &netproto.LoginOk{}
	err = proto.Unmarshal(bytes, loginok)
	if err != nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLoginOk %t", err)
	}
	player := g_PlayerManage.FindTmpPlayer(playerid)
	if player == nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLoginOk no player", playerid)
	}

	if player != nil {
		if loginok.Login_Error == 1 && g_PlayerManage.ReplacePlayer(playerid, int(loginok.UserId)) {
			var buffer netcmd.CmdData
			buffer.AddCmdID(int(netproto.CmdDefine_r2c_login_ok))
			buffer.AddBytes(bytes)
			player.Write(buffer.Data())
		} else {
			if loginok.Login_Error == 1 {
				loginok.Login_Error = 2
				bytes, _ = proto.Marshal(loginok)
			}
			var buffer netcmd.CmdData
			buffer.AddCmdID(int(netproto.CmdDefine_r2c_login_ok))
			buffer.AddBytes(bytes)
			player.Write(buffer.Data())
		}
	}

	Println("player login ok", loginok.UserName)
	return nil
}
