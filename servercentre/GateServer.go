package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CGateServer struct {
	m_iServerId int
	m_strIp     string
	m_iPort     int
	m_conn      net.Conn
}

func (p *CGateServer) Write(bytes []byte) {
	if p.m_conn != nil {
		p.m_conn.Write(bytes)
	}
}

type CGateServerManage struct {
	m_oGateServer   map[net.Conn]*CGateServer
	m_oGateServerId map[int]*CGateServer
}

func (p *CGateServerManage) Init() {
	p.m_oGateServer = make(map[net.Conn]*CGateServer)
	p.m_oGateServerId = make(map[int]*CGateServer)
	netcmd.AddCmdData(int(netproto.CmdDefine_e2t_login), p.NetCmdLogin)
	netcmd.AddCmdData(int(netproto.CmdDefine_e2t_player_login), p.NetCmdPlayerLogin)
}
func (p *CGateServerManage) FindServer(id int) *CGateServer {
	server, ok := p.m_oGateServerId[id]
	if ok {
		return server
	}
	return nil
}
func (p *CGateServerManage) IsServerLogin(conn net.Conn) bool {
	_, ok := p.m_oGateServer[conn]
	if ok {
		return true
	}
	conn.Close()
	return false
}
func (p *CGateServerManage) FindServerByConn(conn net.Conn) *CGateServer {
	server, ok := p.m_oGateServer[conn]
	if ok {
		return server
	}
	conn.Close()
	return nil
}

func (p *CGateServerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {
	server, ok := p.m_oGateServer[conn]
	if ok {
		delete(p.m_oGateServerId, server.m_iServerId)
		delete(p.m_oGateServer, conn)
		Printf("gate server %d out line count=%d\n", server.m_iServerId, len(p.m_oGateServer))
	}
}
func (p *CGateServerManage) SendAllGate(bytes []byte) {
	for _, gate := range p.m_oGateServer {
		gate.Write(bytes)
	}
}

func (p *CGateServerManage) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGateServerManage NetCmdLogin %t", err)
	}

	login := &netproto.ServerLogin{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGateServerManage NetCmdLogin %t", err)
	}

	if login.ServerType != netproto.ServerType_Type_Gate {
		conn.Close()
		return fmt.Errorf("CGateServerManage NetCmdLogin is not gate")
	}
	if login.ServerPWD != "123456" {
		conn.Close()
		return fmt.Errorf("CGateServerManage NetCmdLogin is not pwd")
	}
	if p.FindServer(int(login.ServerId)) != nil {
		conn.Close()
		return fmt.Errorf("CGateServerManage NetCmdLogin have server id=%d", login.ServerId)
	}

	var gateserver CGateServer
	gateserver.m_iServerId = int(login.ServerId)
	gateserver.m_strIp = login.ServerIp
	gateserver.m_iPort = int(login.ServerPort)
	gateserver.m_conn = conn
	p.m_oGateServer[conn] = &gateserver
	p.m_oGateServerId[gateserver.m_iServerId] = &gateserver

	gateLoginOk := &netproto.ServerLoginOk{}
	gateLoginOk.ServerType = netproto.ServerType_Type_Center
	g_GameServer.SendGameServerList(gateLoginOk)
	res, err := proto.Marshal(gateLoginOk)
	if err != nil {
		return fmt.Errorf("CGateServerManage NetCmdLogin %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_t2e_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())

	Printf("gate server id=%d login ok\n", gateserver.m_iServerId)
	return nil
}

func (p *CGateServerManage) NetCmdPlayerLogin(conn net.Conn, data *netcmd.CmdData) error {
	server := p.FindServerByConn(conn)
	if server == nil {
		return fmt.Errorf("CGateServerManage NetCmdPlayerLogin no gate login")
	}
	playerid := data.GetInt()
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CGateServerManage NetCmdPlayerLogin %t", err)
	}
	login := &netproto.Login{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		return fmt.Errorf("CPlayerManage NetCmdLogin %t", err)
	}

	//验证登入,现在直接成功
	Println("client login", login.PlayerName)

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_t2d_player_login))
	buffer.AddInt(server.m_iServerId)
	buffer.AddInt(playerid)
	buffer.AddBytes(bytes)
	g_DbServer.Write(buffer.Data())

	return nil
}
