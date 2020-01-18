package main

import (
	"fmt"
	"net"
	"netcmd"
	"strconv"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CServer struct {
	m_iServerId int
	m_strIp     string
	m_iPort     int
	m_conn      net.Conn
}

type CServerManage struct {
	m_oServer   map[net.Conn]*CServer
	m_oServerId map[int]*CServer
}

func (p *CServerManage) Init() {
	p.m_oServer = make(map[net.Conn]*CServer)
	p.m_oServerId = make(map[int]*CServer)
	netcmd.AddCmdData(int(netproto.CmdDefine_t2d_login), p.NetCmdCentreLogin)
	netcmd.AddCmdData(int(netproto.CmdDefine_g2d_login), p.NetCmdGameLogin)

	netcmd.AddCmdData(int(netproto.CmdDefine_t2d_player_login), p.NetCmdPlayerLogin)
	Println("CServerManage Init")
}
func (p *CServerManage) FindServer(id int) *CServer {
	server, ok := p.m_oServerId[id]
	if ok {
		return server
	}
	return nil
}
func (p *CServerManage) IsServerLogin(conn net.Conn) bool {
	_, ok := p.m_oServer[conn]
	if ok {
		return true
	}
	conn.Close()
	return false
}

func (p *CServerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {
	server, ok := p.m_oServer[conn]
	if ok {
		delete(p.m_oServerId, server.m_iServerId)
		delete(p.m_oServer, conn)
		Printf("db server %d out line count=%d\n", server.m_iServerId, len(p.m_oServer))
	}
}

func (p *CServerManage) NetCmdCentreLogin(conn net.Conn, data *netcmd.CmdData) error {
	Println("NetCmdCentreLogin")
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin %t", err)
	}

	login := &netproto.ServerLogin{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin %t", err)
	}

	if login.ServerType != netproto.ServerType_Type_Center {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin is not game")
	}
	if login.ServerPWD != "123456" {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin is not pwd", err)
	}
	if p.FindServer(int(login.ServerId)) != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin have server id=%d", login.ServerId)
	}
	var server CServer
	server.m_iServerId = int(login.ServerId)
	server.m_strIp = login.ServerIp
	server.m_iPort = int(login.ServerPort)
	server.m_conn = conn
	p.m_oServer[conn] = &server
	p.m_oServerId[server.m_iServerId] = &server

	netok := &netproto.ServerLoginOk{}
	netok.ServerType = netproto.ServerType_Type_Db
	buf, err := proto.Marshal(netok)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}
	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_d2t_login))
	buffer.AddBytes(buf)
	conn.Write(buffer.Data())

	Printf("game server id=%d login ok\n", server.m_iServerId)
	return nil
}

func (p *CServerManage) NetCmdGameLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin %t", err)
	}

	login := &netproto.ServerLogin{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin %t", err)
	}

	if login.ServerType != netproto.ServerType_Type_Game {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin is not game")
	}
	if login.ServerPWD != "123456" {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin is not pwd", err)
	}
	if p.FindServer(int(login.ServerId)) != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdCentreLogin have server id=%d", login.ServerId)
	}
	var server CServer
	server.m_iServerId = int(login.ServerId)
	server.m_strIp = login.ServerIp
	server.m_iPort = int(login.ServerPort)
	server.m_conn = conn
	p.m_oServer[conn] = &server
	p.m_oServerId[server.m_iServerId] = &server

	netok := &netproto.ServerLoginOk{}
	netok.ServerType = netproto.ServerType_Type_Db
	buf, err := proto.Marshal(netok)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}
	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_d2g_login))
	buffer.AddBytes(buf)
	conn.Write(buffer.Data())

	Printf("game server id=%d login ok\n", server.m_iServerId)
	return nil
}

var idddd int32 = 0

func GetId() int32 {
	idddd++
	return idddd
}
func (p *CServerManage) NetCmdPlayerLogin(conn net.Conn, data *netcmd.CmdData) error {
	serverid := data.GetInt()
	playerid := data.GetInt()
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLogin %t", err)
	}
	login := &netproto.Login{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		return fmt.Errorf("CServerManage NetCmdPlayerLogin %t", err)
	}

	loginok := &netproto.LoginOk{}
	loginok.Login_Error = 1
	loginok.UserId = GetId()
	loginok.UserName = "user" + strconv.Itoa(int(loginok.UserId))
	loginok.PlayerName = login.PlayerName
	buf, err := proto.Marshal(loginok)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}
	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_d2t_player_login_ok))
	buffer.AddInt(serverid)
	buffer.AddInt(playerid)
	buffer.AddBytes(buf)
	conn.Write(buffer.Data())

	Printf("game server player login serverid=%d,name=%s,id=%s\n", serverid, login.PlayerName, loginok.UserName)
	return nil
}
