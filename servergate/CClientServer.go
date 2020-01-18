package main

import (
	"fmt"
	"net"
	"netcmd"

	netproto "proto"

	"github.com/golang/protobuf/proto"
)

type CPlayerInfo struct {
	m_iId  int
	m_bTmp bool

	m_strName string
	m_conn    net.Conn
}

func (p *CPlayerInfo) Write(bytes []byte) {
	if p.m_conn != nil {
		p.m_conn.Write(bytes)
	}
}

type CPlayerManage struct {
	m_oPlayerList     map[net.Conn]*CPlayerInfo
	m_oPlayerListId   map[int]*CPlayerInfo
	m_oPlayerListName map[string]*CPlayerInfo

	m_oTmpPlayerLogin     map[int]*CPlayerInfo
	m_oTmpPlayerLoginName map[string]*CPlayerInfo
}

func (p *CPlayerManage) Init() {
	p.m_oPlayerList = make(map[net.Conn]*CPlayerInfo)
	p.m_oPlayerListId = make(map[int]*CPlayerInfo)
	p.m_oPlayerListName = make(map[string]*CPlayerInfo)

	p.m_oTmpPlayerLogin = make(map[int]*CPlayerInfo)
	p.m_oTmpPlayerLoginName = make(map[string]*CPlayerInfo)
	netcmd.AddCmdData(int(netproto.CmdDefine_c2r_login), p.NetCmdLogin)
}

func (p *CPlayerManage) CmdError(conn net.Conn, data *netcmd.CmdData) {
	player, ok := p.m_oPlayerList[conn]
	if ok {
		Printf("gate server %d out line\n", player.m_iId)
		if player.m_bTmp {
			delete(p.m_oTmpPlayerLogin, player.m_iId)
			delete(p.m_oTmpPlayerLoginName, player.m_strName)
		} else {
			delete(p.m_oPlayerListId, player.m_iId)
			delete(p.m_oPlayerListName, player.m_strName)
		}
		delete(p.m_oPlayerList, conn)
	}
}
func (p *CPlayerManage) FindPlayer(id int) *CPlayerInfo {
	player, ok := p.m_oPlayerListId[id]
	if ok {
		return player
	}
	return nil
}
func (p *CPlayerManage) FindTmpPlayer(id int) *CPlayerInfo {
	player, ok := p.m_oTmpPlayerLogin[id]
	if ok {
		return player
	}
	return nil
}
func (p *CPlayerManage) ReplacePlayer(tmpid int, userid int) bool {
	player, ok := p.m_oTmpPlayerLogin[tmpid]
	if !ok {
		return false
	}
	delete(p.m_oTmpPlayerLogin, tmpid)
	delete(p.m_oTmpPlayerLoginName, player.m_strName)

	player.m_iId = userid
	p.m_oPlayerListId[userid] = player
	p.m_oPlayerListName[player.m_strName] = player

	return true
}

func (p *CPlayerManage) NetCmdLogin(conn net.Conn, data *netcmd.CmdData) error {
	bytes, err := data.GetBytes()
	if err != nil {
		return fmt.Errorf("CPlayerManage NetCmdLogin %t", err)
	}

	login := &netproto.Login{}
	err = proto.Unmarshal(bytes, login)
	if err != nil {
		return fmt.Errorf("CPlayerManage NetCmdLogin %t", err)
	}

	_, ok := p.m_oTmpPlayerLoginName[login.PlayerName]
	_, ok1 := p.m_oPlayerListName[login.PlayerName]
	if ok || ok1 {
		//已经在登入了
		loginerror := &netproto.LoginError{}
		loginerror.ErrId = 1

		res, err := proto.Marshal(loginerror)
		if err != nil {
			conn.Close()
			return fmt.Errorf("CStatusServer LoginOk %t", err)
		}

		var buffer netcmd.CmdData
		buffer.AddCmdID(int(netproto.CmdDefine_r2c_login_error))
		buffer.AddBytes(res)
		conn.Write(buffer.Data())
		conn.Close()
		return nil
	}

	Println("client login name=", login.PlayerName)
	player := &CPlayerInfo{}
	player.m_iId = GetCreateId()
	player.m_bTmp = false
	player.m_strName = login.PlayerName
	player.m_conn = conn

	p.m_oPlayerList[conn] = player
	p.m_oTmpPlayerLogin[player.m_iId] = player
	p.m_oTmpPlayerLoginName[player.m_strName] = player

	var buffer netcmd.CmdData
	buffer.AddCmdID(int(netproto.CmdDefine_e2t_player_login))
	buffer.AddInt(player.m_iId)
	buffer.AddBytes(bytes)
	g_CentreServer.Write(buffer.Data())

	return nil
}
