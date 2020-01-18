package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CPlayerManage struct {
}

func (p *CPlayerManage) Init() {

	netcmd.AddCmdData(int(netproto.CmdDefine_c2s_login), p.NetCmdLogin)
}

func (p *CPlayerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {

}

func (p *CPlayerManage) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	gate := g_GateServer.GetRandServer()
	if gate == nil {
		conn.Close()
		return nil
	}

	loginOk := &netproto.LoginStatus{}
	loginOk.StrIp = gate.m_strIp
	loginOk.IPort = int32(gate.m_iPort)
	res, err := proto.Marshal(loginOk)
	if err != nil {
		conn.Close()
		return fmt.Errorf("LoginStatus Marshal error")
	}

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_s2c_login))
	buffer.AddBytes(res)
	conn.Write(buffer.Data())

	conn.Close()
	return nil
}
