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

	netcmd.AddCmdData(int(netproto.CmdDefine_t2g_login), p.NetCmdLogin)

	Println("login centre", p.m_strIp, p.m_iPort)
	netcmd.CmdDial(p.m_strIp, p.m_iPort, netcmd.NETDIALOK_1)
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
	login := &netproto.ServerLogin{}
	login.ServerType = netproto.ServerType_Type_Game
	login.ServerId = int32(g_iServerId)
	login.ServerIp = g_strServerIp
	login.ServerPort = int32(g_iServerPort)
	login.ServerPWD = "123456"
	res, err := proto.Marshal(login)
	if err != nil {
		conn.Close()
		return fmt.Errorf("CCentreServer LoginOk %t", err)
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_g2t_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())
	return nil
}

func (p *CCentreServer) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	Println("CCentreServer NetCmdLogin")
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

	p.m_conn = conn

	Println("login centre ok", p.m_strIp, p.m_iPort)
	return nil
}
