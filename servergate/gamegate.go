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
	g_StatusServer.Init()
	g_PlayerManage.Init()
	g_GameServer.Init()

	netcmd.NewListen(g_strServerIp, g_iServerPort)
}

//断线
func CmdError(conn net.Conn, data *netcmd.CmdData) error {
	g_CentreServer.CmdError(conn, data)
	g_GameServer.CmdError(conn, data)
	g_PlayerManage.CmdError(conn, data)
	g_StatusServer.CmdError(conn, data)
	return nil
}

var g_iCreateId int = 1

func GetCreateId() int {
	g_iCreateId++
	return g_iCreateId
}

var g_CentreServer CCentreServer
var g_StatusServer CStatusServer
var g_GameServer CGameServerManage
var g_PlayerManage CPlayerManage

var g_strServerName string
var g_strServerIp string
var g_iServerPort int
var g_iServerId int

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
	status := conf.Section("status")
	if status != nil {
		g_StatusServer.m_strIp = status.Key("ip").String()
		g_StatusServer.m_iPort, _ = status.Key("port").Int()
	}

	return true
}
