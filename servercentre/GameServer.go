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

type CGameServerManage struct {
	m_oGameServer   map[net.Conn]*CGameServer
	m_oGameServerId map[int]*CGameServer
}

func (p *CGameServerManage) Init() {
	p.m_oGameServer = make(map[net.Conn]*CGameServer)
	p.m_oGameServerId = make(map[int]*CGameServer)
	netcmd.AddCmdData(int(netproto.CmdDefine_g2t_login), p.NetCmdLogin)
	Println("CGameServerManage Init")
}
func (p *CGameServerManage) FindServer(id int) *CGameServer {
	server, ok := p.m_oGameServerId[id]
	if ok {
		return server
	}
	return nil
}
func (p *CGameServerManage) IsServerLogin(conn net.Conn) bool {
	_, ok := p.m_oGameServer[conn]
	if ok {
		return true
	}
	conn.Close()
	return false
}

func (p *CGameServerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {
	server, ok := p.m_oGameServer[conn]
	if ok {
		delete(p.m_oGameServerId, server.m_iServerId)
		delete(p.m_oGameServer, conn)
		Printf("game server %d out line count=%d\n", server.m_iServerId, len(p.m_oGameServer))
	}
}

func (p *CGameServerManage) SendGameServerList(gateLoginOk *netproto.ServerLoginOk) {
	for _, game := range p.m_oGameServer {
		netGame := &netproto.ServerGame{}
		netGame.ServerId = int32(game.m_iServerId)
		netGame.ServerIp = game.m_strIp
		netGame.ServerPort = int32(game.m_iPort)
		netGame.ServerPWD = "123456"
		gateLoginOk.ServerGameList = append(gateLoginOk.ServerGameList, netGame)
	}
}

func (p *CGameServerManage) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
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
		return fmt.Errorf("CGameServerManage NetCmdLogin is not game")
	}
	if login.ServerPWD != "123456" {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin is not pwd", err)
	}
	if p.FindServer(int(login.ServerId)) != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin have server id=%d", login.ServerId)
	}
	var gameserver CGameServer
	gameserver.m_iServerId = int(login.ServerId)
	gameserver.m_strIp = login.ServerIp
	gameserver.m_iPort = int(login.ServerPort)
	gameserver.m_conn = conn
	p.m_oGameServer[conn] = &gameserver
	p.m_oGameServerId[gameserver.m_iServerId] = &gameserver

	netok := &netproto.ServerLoginOk{}
	netok.ServerType = netproto.ServerType_Type_Center
	buf, err := proto.Marshal(netok)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CGameServerManage NetCmdLogin %t", err)
	}
	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_t2g_login))
	buffer.AddBytes(buf)
	conn.Write(buffer.Data())

	netGame := &netproto.ServerGame{}
	netGame.ServerId = int32(gameserver.m_iServerId)
	netGame.ServerIp = gameserver.m_strIp
	netGame.ServerPort = int32(gameserver.m_iPort)
	netGame.ServerPWD = "123456"
	res, err := proto.Marshal(netGame)
	if err == nil {
		var buffer netcmd.CmdData
		buffer.AddCmdID(int(netproto.CmdDefine_t2e_game_login))
		buffer.AddBytes(res)
		g_GateServer.SendAllGate(buffer.Data())
	}

	Printf("game server id=%d login ok\n", gameserver.m_iServerId)
	return nil
}
