package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CGameServer struct {
	m_iServerId int
	m_strIp     string
	m_iPort     int
	m_conn      net.Conn
}

func (p *CGameServer) Init() {
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_2)
	PrintfWarning("Login game server ip=%s", p.m_strIp)
}

type CGameServerManage struct {
	m_oGameServer   map[net.Conn]*CGameServer
	m_oGameServerId map[int]*CGameServer
}

func (p *CGameServerManage) Init() {
	p.m_oGameServer = make(map[net.Conn]*CGameServer)
	p.m_oGameServerId = make(map[int]*CGameServer)
	netcmd.AddCmdData(netcmd.NETDIALOK_2, p.LoginOk)

	netcmd.AddCmdData(int(netproto.CmdDefine_g2e_login), p.NetCmdLogin)

}
func (p *CGameServerManage) FindServer(id int) *CGameServer {
	server, ok := p.m_oGameServerId[id]
	if ok {
		return server
	}
	return nil
}

func (p *CGameServerManage) AddNewServer(id int, ip string, port int) {
	if p.FindServer(id) != nil {
		Println("have game server login id=", id)
		return
	}
	var server CGameServer
	server.m_iServerId = id
	server.m_strIp = ip
	server.m_iPort = port

	p.m_oGameServerId[id] = &server
	server.Init()
	Println("login game ", ip, port)
}

func (p *CGameServerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {
	server, ok := p.m_oGameServer[conn]
	if ok {
		delete(p.m_oGameServerId, server.m_iServerId)
		delete(p.m_oGameServer, conn)
		Printf("game server %d out line count=%d\n", server.m_iServerId, len(p.m_oGameServer))
	}
}

func (p *CGameServerManage) LoginOk(conn net.Conn, data *netcmd.CmdData) error {
	Println("CGameServerManage LoginOk")
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Gate
	login.ServerId = int32(g_iServerId)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage LoginOk %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_e2g_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())
	return nil
}

func (p *CGameServerManage) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	Println("CGameServerManage NetCmdLogin")
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}

	login := &netproto.ServerLogin{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}

	if login.ServerType != netproto.ServerType_Type_Game {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin not game")
	}
	if login.ServerPWD != "123456" {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin pwd is not")
	}
	gameserver := p.FindServer(int(login.ServerId))
	if gameserver != nil {
		gameserver.m_conn = conn
		p.m_oGameServer[conn] = gameserver
		Printf("game server id=%d login ok\n", gameserver.m_iServerId)
	} else {
		conn.Close()
		return fmt.Errorf("game have server id=%d \n", login.ServerId)
	}

	return nil
}
