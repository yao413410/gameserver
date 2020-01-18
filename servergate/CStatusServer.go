package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CStatusServer struct {
	m_strIp string
	m_iPort int
	m_conn  net.Conn
}

func (p *CStatusServer) Init() {
	netcmd.AddCmdData(netcmd.NETDIALOK_3, p.LoginOk)

	netcmd.AddCmdData(int(netproto.CmdDefine_s2e_login), p.NetCmdLogin)

	Println("login status", p.m_strIp, p.m_iPort)
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_3)
}

func (p *CStatusServer) CmdError(conn net.Conn, data *netcmd.CmdData) {
	//断线重连
	if conn != p.m_conn {
		return
	}
	p.m_conn = nil
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_3)
	PrintfWarning("Login again status server ip=%s", p.m_strIp)
}

func (p *CStatusServer) LoginOk(conn net.Conn, data *netcmd.CmdData) error {
	Println("CStatusServer LoginOk")
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Gate
	login.ServerId = int32(g_iServerId)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CStatusServer LoginOk %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_e2s_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())
	return nil
}

func (p *CStatusServer) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		conn.Close()
		return fmt.Errorf("CStatusServer NetCmdLogin %t", err)
	}
	loginOk := &netproto.ServerLoginOk{}
	err = proto.Unmarshal(bytes, loginOk)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CStatusServer NetCmdLogin %t", err)
	}

	p.m_conn = conn

	Println("login status ok", p.m_strIp, p.m_iPort)
	return nil
}
