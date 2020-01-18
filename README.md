# gameserver
网络游戏架构gate-centre-game

centre 中心服务器,链接db
status 分配GIT
gate 网关服务器 链接status,centre,game
game 游戏逻辑服务器 链接centre,db
db 数据库服务器,链接db
