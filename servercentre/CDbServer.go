package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CDbServer struct {
	m_conn net.Conn
}

func (p *CDbServer) Init() {
	netcmd.AddCmdData(netcmd.NETDIALOK_2, p.LoginOk)

	netcmd.AddCmdData(int(netproto.CmdDefine_d2t_login), p.NetCmdLogin)

	netcmd.AddCmdData(int(netproto.CmdDefine_d2t_player_login_ok), p.NetCmdPlayerLoginOk)

	Println("login db", g_strDbIp, g_iDbPort)
	netcmd.CmdDial(g_strDbIp, g_iDbPort, netcmd.NETDIALOK_2)
}

func (p *CDbServer) Write(b []byte) {
	if p.m_conn != nil {
		p.m_conn.Write(b)
	}
}

func (p *CDbServer) CmdError(conn net.Conn, data *netcmd.CmdData) {
	//断线重连
	if conn != p.m_conn {
		return
	}
	p.m_conn = nil
	netcmd.CmdDial(g_strDbIp, g_iDbPort, netcmd.NETDIALOK_2)
	PrintfWarning("Login again db server ip=%s", g_strDbIp)
}

func (p *CDbServer) LoginOk(conn net.Conn, data *netcmd.CmdData) error {
	Println("CDbServer LoginOk")
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Center
	login.ServerId = int32(-1)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CDbServer ServerLogin %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_t2d_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())
	return nil
}

func (p *CDbServer) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CDbServer NetCmdLogin %t", err)
	}
	loginOk := &netproto.ServerLoginOk{}
	err = proto.Unmarshal(bytes, loginOk)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CDbServer NetCmdLogin %t", err)
	}

	p.m_conn = conn

	Println("login db ok", g_strDbIp, g_iDbPort)
	return nil
}

func (p *CDbServer) NetCmdPlayerLoginOk(conn net.Conn, data *netcmd.CmdData) error {
	serverid := data.GetInt()
	playerid := data.GetInt()
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLoginOk %t", err)
	}
	loginok := &netproto.LoginOk{}
	err = proto.Unmarshal(bytes, loginok)
	if err != nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLoginOk %t", err)
	}

	gate := g_GateServer.FindServer(serverid)
	if gate == nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLoginOk gate=nil,id=%d", serverid)
	}

	if loginok.Login_Error == 1 {
		player := &CPlayerInfo{}
		player.m_iId = int(loginok.UserId)
		player.m_strName = loginok.PlayerName
		player.m_strRoleName = loginok.UserName
		g_PlayerManage.AddPlayer(player)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_t2e_player_login_ok))
	buffer.AddInt(playerid)
	buffer.AddBytes(bytes)
	gate.Write(buffer.Data())

	Println("player login ok", loginok.UserName)
	return nil
}
