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
	Println("start Centre Server Ok")
	for true {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func InitCmd() {
	netcmd.AddCmdData(netcmd.NETERROR, CmdError)
	g_GameServer.Init()
	g_GateServer.Init()
	g_DbServer.Init()
	g_PlayerManage.Init()

	netcmd.NewListen(g_strServerIp, g_iServerPort)
}

//断线
func CmdError(conn net.Conn, data *netcmd.CmdData) error {
	g_GameServer.CmdError(conn, data)
	g_GateServer.CmdError(conn, data)
	g_DbServer.CmdError(conn, data)
	return nil
}

var g_GameServer CGameServerManage
var g_GateServer CGateServerManage
var g_DbServer CDbServer

var g_PlayerManage CPlayerManage

var g_strServerName string
var g_strServerIp string
var g_iServerPort int

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
		g_strServerIp = setup.Key("ip").String()
		g_iServerPort, _ = setup.Key("port").Int()

		g_bPrintLog, _ = setup.Key("log").Bool()
	}

	db := conf.Section("db")
	if db != nil {
		g_strDbIp = db.Key("ip").String()
		g_iDbPort, _ = db.Key("port").Int()
	}

	return true
}
