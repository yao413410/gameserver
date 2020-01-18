package main

import (
	"net"
	"netcmd"
	"time"

	"gopkg.in/ini.v1"
)

var g_bPrintLog bool = false

func main() {
	LoadSetUp()
	InitCmd()
	Println("start Gate Server Ok")
	for true {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
func InitCmd() {
	netcmd.AddCmdData(netcmd.NETERROR, CmdError)

	g_CentreServer.Init()
	g_GateServer.Init()
	g_DbServer.Init()
	netcmd.NewListen(g_strServerIp, g_iServerPort)
}

//断线
func CmdError(conn net.Conn, data *netcmd.CmdData) error {
	g_CentreServer.CmdError(conn, data)
	g_GateServer.CmdError(conn, data)
	g_DbServer.CmdError(conn, data)
	return nil
}

var g_CentreServer CCentreServer
var g_GateServer CGateServerManage
var g_DbServer CDbServer

var g_strServerName string
var g_strServerIp string
var g_iServerPort int
var g_iServerId int

var g_strDbIp string
var g_iDbPort int

func LoadSetUp() bool {
	conf, err := ini.Load("config/setup.ini")
	if err != nil {
		PrintfWarning("try load config file[setup.ini] error[%s]\n", err.Error())
		return false
	}
	setup := conf.Section("setup")
	if setup != nil {
		g_strServerName = setup.Key("name").String()
		g_iServerId, _ = setup.Key("id").Int()
		g_strServerIp = setup.Key("ip").String()
		g_iServerPort, _ = setup.Key("port").Int()

		g_bPrintLog, _ = setup.Key("log").Bool()
	}
	centre := conf.Section("centre")
	if centre != nil {
		g_CentreServer.m_strIp = centre.Key("ip").String()
		g_CentreServer.m_iPort, _ = centre.Key("port").Int()
	}

	db := conf.Section("db")
	if db != nil {
		g_strDbIp = db.Key("ip").String()
		g_iDbPort, _ = db.Key("port").Int()
	}

	return true
}
