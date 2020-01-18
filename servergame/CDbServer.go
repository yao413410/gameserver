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
	netcmd.AddCmdData(netcmd.NETDIALOK_3, p.LoginOk)

	netcmd.AddCmdData(int(netproto.CmdDefine_d2g_login), p.NetCmdLogin)

	Println("login db", g_strDbIp, g_iDbPort)
	netcmd.CmdDial(g_strDbIp, g_iDbPort, netcmd.NETDIALOK_3)
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
	netcmd.CmdDial(g_strDbIp, g_iDbPort, netcmd.NETDIALOK_3)
	PrintfWarning("Login again db server ip=%s", g_strDbIp)
}

func (p *CDbServer) LoginOk(conn net.Conn, data *netcmd.CmdData) error {
	Println("CDbServer LoginOk")
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Game
	login.ServerId = int32(g_iServerId)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CDbServer ServerLogin %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_g2d_login))
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
