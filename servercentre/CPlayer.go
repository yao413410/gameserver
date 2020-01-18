package main

type CPlayerInfo struct {
	m_iId         int
	m_strName     string
	m_strRoleName string
}

type CPlayerManage struct {
	m_oPlayerList     map[int]*CPlayerInfo
	m_oPlayerListName map[string]*CPlayerInfo
}

func (p *CPlayerManage) Init() {
	p.m_oPlayerList = make(map[int]*CPlayerInfo)
	p.m_oPlayerListName = make(map[string]*CPlayerInfo)
}

func (p *CPlayerManage) AddPlayer(player *CPlayerInfo) {
	if player == nil {
		return
	}
	p.m_oPlayerList[player.m_iId] = player
	p.m_oPlayerListName[player.m_strName] = player
}

func (p *CPlayerManage) FindPlayer(id int) *CPlayerInfo {
	player, ok := p.m_oPlayerList[id]
	if ok {
		return player
	}
	return nil
}
func (p *CPlayerManage) FindPlayerName(name string) *CPlayerInfo {
	player, ok := p.m_oPlayerListName[name]
	if ok {
		return player
	}
	return nil
}
