package sggj

const (
	UNKNOWN = byte(iota)
	OFF_FREE
	OFF_RAID
	OFF_PROT
	ON_FREE
	ON_PROT
)

const (
	PLAYER_STATE_WAITING     = 0 //玩家等待
	PLAYER_STATE_FIGHTING    = 1 //玩家战斗
	PLAYER_STATE_BUILDING    = 2 //玩家建筑
	PLAYER_STATE_FARMING     = 3 //玩家种田
	PLAYER_STATE_DEFANCE     = 4 //玩家反抗
	PLAYER_STATE_CUTTING     = 5 //玩家伐木
	PLAYER_STATE_MANUFACTURE = 6 //玩家制造
	PLAYER_STATE_TOUR        = 7 //玩家转转
)

const (
	BATTLE_STATE_TRAP           = 0 //战斗初始
	BATTLE_STATE_FIGHTING       = 1 //进入战斗
	BATTLE_STATE_RUNAWAY        = 2 //逃跑
	BATTLE_STATE_GET_ITEM       = 3 //获取物品
	BATTLE_STATE_SURRENDER      = 4 //搜刮
	BATTLE_STATE_RECOVER        = 5 //恢复
	BATTLE_STATE_LOCALFIGHTIONG = 6 //攻击
	BATTLE_STATE_AOEFIGHTING    = 7 //范围攻击
)

const (
	RIDER_STATE_OUTFLANK = 0 //迂回
	RIDER_STATE_ASSAULT  = 1 //冲锋
	RIDER_STATE_BREAK    = 2 //破阵
	RIDER_STATE_FAITING  = 3 //离开
)

var SoldierType = map[int]string{
	0: "刀盾兵",
	1: "长枪兵",
	2: "弓箭兵",
	3: "骑兵",
}

var BuildingType = map[int]string{
	0: "主城",
	1: "铁匠铺",
	2: "农场",
	3: "陷阱",
}

//[当前骑兵部队状态]
var RiderLegionState = [4][4]float64{
	{0.4, 0.6, 0, 0},   //迂回
	{0, 0, 1, 0},       //冲锋
	{0.3, 0, 0.3, 0.4}, //破阵
	{0.5, 0, 0, 0.5},   //离开
}

//[行动部队][目标部队]
var ArmsPlus = [4][4]float64{
	{0.8, 1.2, 1.5, 0.8},
	{0.5, 1, 1.5, 2},
	{0.4, 1, 1, 1},
	{0.5, 0.8, 0.8, 0.5},
}
