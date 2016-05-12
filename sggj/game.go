package sggj

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	//"time"
)

type GameManage struct {
	actionTargetNumber   int //行动玩家
	actionLegionNumber   int //行动部队
	defenderLegionNumber int //防御部队
	attackerNumber       int //攻击士兵
	defenderNumber       int //防御士兵
	Cfg                  *Config
	battleArray          [2]*Player
	randomValue          [3]float64
	battleFunction       func()
	fightStatus          int
	bout                 int
	b                    int
}

var RandomCoefficientDic = map[int]float64{
	3:  0.5,
	4:  0.56,
	5:  0.61,
	6:  0.67,
	7:  0.72,
	8:  0.78,
	9:  0.83,
	10: 0.89,
	11: 0.94,
	12: 1,
	13: 1.06,
	14: 1.11,
	15: 1.17,
	16: 1.22,
	17: 1.28,
	18: 1.33,
	19: 1.39,
	20: 1.44,
	21: 1.5,
}

func NewGame(a int, b int) *GameManage {
	gm := new(GameManage)

	gm.battleArray[0] = NewPlayer("A", 10, 10, 10)
	gm.battleArray[1] = NewPlayer("B", 10, 10, 10)

	//code0 := "孙权,30,185,19,185,185,185,323,334,12,70,335,19,81,33,192,65,263,27,341,342,342,342,343,344,344,344,5285"
	//code1 := "袁绍,29,175,18,175,175,175,241,237,60,26,214,219,267,260,332,13,10,152,333,341,341,341,342,344,344,344,5508"

	code0 := PlayerData[a].Code
	code1 := PlayerData[b].Code
	gm.battleArray[0].decipher(code0)
	gm.battleArray[1].decipher(code1)

	return gm
}

func (gm *GameManage) Start() {
	debug("cfg SoldierData %v\n", SoldierData)
	debug("cfg EquipData %v\n", EquipData)
	debug("cfg PlayerData %v\n", PlayerData)

	debug("battleArray 0 %s \n", gm.battleArray[0].Name)
	debug("battleArray 1 %s \n", gm.battleArray[1].Name)

	//	for {
	//		gm.playerFighting()
	//	}

	fmt.Println("战斗开始")
	gm.battleArray[0].RiderLegionState = 0
	gm.battleArray[0].dashPower = 0

	gm.battleArray[1].RiderLegionState = 0
	gm.battleArray[1].dashPower = 0

	gm.actionTargetNumber = 0
	gm.actionLegionNumber = 0

	gm.attackerNumber = 0
	gm.defenderNumber = 0
	debug("Start defenderNumber#%d \n", gm.defenderNumber)
	gm.switchBattleState(BATTLE_STATE_FIGHTING)
	gm.fightStatus = 0

	gm.bout = 0
	gm.b = 0
	for {
		show("轮换#%d=>", gm.bout)
		gm.playerFighting()
		debug("\n")
		gm.bout++
		if gm.fightStatus == 1 {
			break
		}
	}

	show("战斗结束 总回合数#%d\n", gm.bout)
}

func (gm *GameManage) battleStateFighting() {
	debug("[%d] 开始计算战斗\n", gm.bout)
	if gm.battleArray[0].IsAllDead() {
		fmt.Printf("%s 的部队 获得胜利\n", gm.battleArray[1].Name)
		gm.fightStatus = 1
	} else if gm.battleArray[1].IsAllDead() {
		fmt.Printf("%s 的部队 获得胜利\n", gm.battleArray[0].Name)
		gm.fightStatus = 1
	} else {
		for {
			debug("battleStateFighting 行动方:%d 行动部队:%d\n", gm.actionTargetNumber, gm.actionLegionNumber)

			if gm.isAllDead(gm.actionTargetNumber, gm.actionLegionNumber) {
				gm.getNextActionTarget()
				continue
			}
			break
		}
		debug("2===========================开始行动方:%s\n", SoldierType[gm.actionLegionNumber])

		switch gm.actionLegionNumber { //检查行动兵种为哪一种
		case 0: //刀兵行动
			gm.switchSaberLegionTarget()
		case 1: //枪兵行动
			gm.switchLancerLegionTarget()
		case 2: //弓兵行动
			gm.switchArcherLegionTarget()
		case 3: //骑兵行动
			gm.switchRiderLegionTarget()
		}

		if gm.actionLegionNumber != 3 || gm.battleArray[gm.actionTargetNumber].RiderLegionState == RIDER_STATE_FAITING {
			if gm.defenderLegionNumber != -1 {
				//fmt.Print(" 攻击防御士兵重置 ")
				gm.switchBattleState(BATTLE_STATE_LOCALFIGHTIONG)
				gm.attackerNumber = 0
				gm.defenderNumber = 0
				debug("battleStateFighting defenderNumber#%d 部队:%d \n", gm.defenderNumber, gm.defenderLegionNumber)
			} else {
				//fmt.Print(" 防御方兵种全挂了 ")
				gm.getNextActionTarget()
			}
		} else if gm.battleArray[gm.actionTargetNumber].RiderLegionState == RIDER_STATE_OUTFLANK || gm.battleArray[gm.actionTargetNumber].RiderLegionState == RIDER_STATE_ASSAULT {
			//fmt.Print(" 状态[迂回|冲锋] ")
			gm.getNextActionTarget()
		}
		show("\n")
		debug("战斗状态执行结束")
	}
}

//攻击
func (gm *GameManage) battleStateLocalFighting() {

	debug("攻击方:%d 攻击兵种:%d 攻击士兵:%d 防御兵种:%d 防御士兵:%d\n", gm.actionTargetNumber, gm.actionLegionNumber, gm.attackerNumber, gm.defenderLegionNumber, gm.defenderNumber)

	soldier1 := new(Soldier)
	soldier2 := new(Soldier)

	for {
		gm.attackerNumber++
		if gm.attackerNumber == 5 {
			gm.attackerNumber = 0
		}
		debug("1.切换攻击方的士兵#%d \n", gm.attackerNumber)
		soldier1 = gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][gm.attackerNumber]
		debug("soldier1 部队：%v\n", soldier1)
		debug("soldier1 部队：%s\n", soldier1.Name)
		if soldier1.Life > 0 {
			break
		}
	}

	for {
		gm.defenderNumber++
		if gm.defenderNumber == 5 {
			gm.defenderNumber = 0
		}
		debug("battleStateLocalFighting defenderNumber#%d \n", gm.defenderNumber)

		debug("1.切换防御方的士兵#%d \n", gm.defenderNumber)
		soldier2 = gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber][gm.defenderNumber]
		debug("soldier2 部队：%v\n", soldier2)
		debug("soldier2 部队：%s Life：%0.2f,defenderLegionNumber:%d  gm.defenderNumber:%d\n", soldier2.Name, soldier2.Life, gm.defenderLegionNumber, gm.defenderNumber)
		if soldier2.Life > 0 {
			break
		}

	}

	gm.getRandomValue()

	hurtNum := soldier1.Power * ArmsPlus[gm.actionLegionNumber][gm.defenderLegionNumber] * gm.randomCoefficient() * math.Pow(1.05, gm.battleArray[gm.actionTargetNumber].InitPower-15)

	debug("====  %0.2f, %0.2f, %0.2f, %0.2f\n", soldier1.Power, ArmsPlus[gm.actionLegionNumber][gm.defenderLegionNumber], gm.randomCoefficient(), math.Pow(1.05, gm.battleArray[gm.actionTargetNumber].InitPower-15))
	soldier2.Hurt(hurtNum)
	show("%s(%s)[%s#%d] 攻击 %s(%s)[%s#%d] (%d)造成了 %0.2f 点伤害.\n", gm.battleArray[gm.actionTargetNumber].Name, soldier1.Name, SoldierType[gm.actionLegionNumber], gm.attackerNumber, gm.battleArray[1-gm.actionTargetNumber].Name, soldier2.Name, SoldierType[gm.defenderLegionNumber], gm.defenderNumber, int(soldier2.Life), hurtNum)

	if gm.isAllDead(1-gm.actionTargetNumber, gm.defenderLegionNumber) || gm.isActionOver() {
		debug("=======================【行动结束 切换目标】===========================\n")
		gm.switchBattleState(BATTLE_STATE_FIGHTING)
		gm.getNextActionTarget()
	}
}

func (gm *GameManage) randomCoefficient() float64 {
	return RandomCoefficientDic[gm.sumRandomValue()]
}

func (gm *GameManage) getRandomValue() {
	for i := 0; i < len(gm.randomValue); i++ {
		gm.randomValue[i] = float64(rand.Intn(7))
	}

	debug("getRandomValue %v\n", gm.randomValue)
}

func (gm *GameManage) sumRandomValue() int {
	num := 0
	for _, n := range gm.randomValue {
		num += int(n)
	}
	return num
}

func (gm *GameManage) playerFighting() {
	debug("battleFunction %v\n", &gm.battleFunction)
	gm.battleFunction()
}

func (gm *GameManage) switchBattleState(t int) {
	debug("switchBattleState %d\n", t)
	switch t {

	case BATTLE_STATE_TRAP:
		gm.battleFunction = gm.battleStateTrap

	case BATTLE_STATE_FIGHTING:
		gm.battleFunction = gm.battleStateFighting

	case BATTLE_STATE_RUNAWAY:
		gm.battleFunction = gm.battleStateRunaway

	case BATTLE_STATE_GET_ITEM:
		gm.battleFunction = gm.battleStateGetItem

	case BATTLE_STATE_SURRENDER:
		gm.battleFunction = gm.battleStateSurrender

	case BATTLE_STATE_RECOVER:
		gm.battleFunction = gm.battleStateRecover

	case BATTLE_STATE_LOCALFIGHTIONG:
		gm.battleFunction = gm.battleStateLocalFighting

	case BATTLE_STATE_AOEFIGHTING:
		gm.battleFunction = gm.battleStateAoeFighting

	}
	//	fmt.Printf("battleFunction %v\n", gm.battleFunction)
}

func (gm *GameManage) battleStateTrap() {

}

func (gm *GameManage) battleStateRunaway() {

}

func (gm *GameManage) battleStateGetItem() {

}

func (gm *GameManage) battleStateSurrender() {

}

func (gm *GameManage) battleStateRecover() {

}

func (gm *GameManage) isActionOver() bool {

	n := gm.attackerNumber
	debug("isActionOver %d\n", gm.attackerNumber+1)
	for {
		n++
		if n >= len(gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber]) {
			break
		}
		if gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][n].Life > 0 {
			return false
		}

	}
	return true

}

//寻找刀盾兵目标
func (gm *GameManage) switchSaberLegionTarget() {
	debug("行动方:%d 寻找刀盾兵目标\n", gm.actionTargetNumber)
	target := 1 - gm.actionTargetNumber //目标
	gm.defenderLegionNumber = -1        //表示该兵种全挂了
	if !gm.isAllDead(target, 0) {
		gm.defenderLegionNumber = 0 //防御方为刀兵
	} else if !gm.isAllDead(target, 1) {
		gm.defenderLegionNumber = 1 //防御方为枪兵
	} else if !gm.isAllDead(target, 2) {
		gm.defenderLegionNumber = 2 //防御方为弓兵
	} else if !gm.isAllDead(target, 3) &&
		gm.battleArray[target].RiderLegionState == RIDER_STATE_FAITING ||
		gm.battleArray[target].RiderLegionState == RIDER_STATE_BREAK { //骑兵破阵状态
		gm.defenderLegionNumber = 3 //防御方为骑兵
	}
	debug("行动方:刀盾兵=>防御方兵种:%d \n", gm.defenderLegionNumber)
}

//寻找长枪兵目标
func (gm *GameManage) switchLancerLegionTarget() {
	debug("行动方:%d 寻找长枪兵目标\n", gm.actionTargetNumber)
	target := 1 - gm.actionTargetNumber
	gm.defenderLegionNumber = -1 //表示改兵种全挂了

	if !gm.isAllDead(target, 3) && (gm.battleArray[target].RiderLegionState == RIDER_STATE_FAITING || gm.battleArray[target].RiderLegionState == RIDER_STATE_BREAK) {
		gm.defenderLegionNumber = 3 //防御方骑兵未全死亡 状态为离开和破阵  防御方为骑兵
	} else if !gm.isAllDead(target, 0) {
		gm.defenderLegionNumber = 0 //防御方为刀兵
	} else if !gm.isAllDead(target, 1) {
		gm.defenderLegionNumber = 1 //防御方为枪兵
	} else if !gm.isAllDead(target, 2) {
		gm.defenderLegionNumber = 2 //防御方为弓兵
	}
	debug("行动方:长枪兵=>防御方兵种:%d \n", gm.defenderLegionNumber)
}

//寻找弓箭兵目标
func (gm *GameManage) switchArcherLegionTarget() {
	debug("行动方:%d 寻找弓箭兵目标 \n", gm.actionTargetNumber)
	target := 1 - gm.actionTargetNumber

	for {
		legion := rand.Intn(4)
		gm.defenderLegionNumber = legion
		debug("3=================isAllDead:%d=>%d \n", target, legion)
		if !gm.isAllDead(target, legion) {
			break
		}
	}
	debug("行动方:弓箭兵=>防御方兵种:%d \n", gm.defenderLegionNumber)
}

//寻找骑兵目标
func (gm *GameManage) switchRiderLegionTarget() {
	debug("行动方:%d  寻找骑兵目标 \n", gm.actionTargetNumber)
	gm.switchRiderLegionState()
	debug("行动方:%d  骑兵状态:%d \n", gm.actionTargetNumber, gm.battleArray[gm.actionTargetNumber].RiderLegionState)
	switch gm.battleArray[gm.actionTargetNumber].RiderLegionState {
	case RIDER_STATE_OUTFLANK:
		gm.getRandomValue()
		gm.battleArray[gm.actionTargetNumber].dashPower += 2 * float64(gm.sumRandomValue())
		show("%s的%s正在【迂回】！冲刺优势增加了%d。\n", gm.battleArray[gm.actionTargetNumber].Name, SoldierType[gm.actionLegionNumber], 2*gm.sumRandomValue())
	case RIDER_STATE_ASSAULT:
		gm.getRandomValue()
		gm.battleArray[gm.actionTargetNumber].dashPower += 3 * float64(gm.sumRandomValue())
		show("%s的%s发起了【冲锋】！冲刺优势增加了%d。\n", gm.battleArray[gm.actionTargetNumber].Name, SoldierType[gm.actionLegionNumber], 3*gm.sumRandomValue())
	case RIDER_STATE_BREAK:
		gm.switchArcherLegionTarget()
		show("%s的%s冲入了敌人的%s中！【破阵】\n", gm.battleArray[gm.actionTargetNumber].Name, SoldierType[gm.actionLegionNumber], SoldierType[gm.defenderLegionNumber])
		gm.getRandomValue()
		gm.switchBattleState(BATTLE_STATE_AOEFIGHTING)
		gm.attackerNumber = 0
		gm.defenderNumber = 0
		debug("switchRiderLegionTarget defenderNumber#%d \n", gm.defenderNumber)
	case RIDER_STATE_FAITING:
		gm.switchArcherLegionTarget()
		gm.battleArray[gm.actionTargetNumber].dashPower = 0

	}
	debug("骑兵=>防御方:%d \n", gm.defenderLegionNumber)
}

//骑兵状态更新
func (gm *GameManage) switchRiderLegionState() {

	n := gm.battleArray[gm.actionTargetNumber].RiderLegionState
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := rand.Float64()
	debug("switchRiderLegionState==rnd:%0.2f \n", rnd)
	debug("switchRiderLegionState==n:%d  => %v \n", n, RiderLegionState[n])

	state := 0
	for _, v := range RiderLegionState[n] {
		rnd -= v
		if rnd <= 0 {
			break
		}
		state++
	}

	debug("switchRiderLegionState==state:%d \n", state)

	gm.battleArray[gm.actionTargetNumber].RiderLegionState = state
}

func (gm *GameManage) isAllDead(target int, legion int) bool {
	debug("=================isAllDead: 目标:%d  部队%d Start\n", target, legion)
	for _, v := range gm.battleArray[target].Army[legion] {
		if !v.IsDead() {
			return false
		}
		debug("=================v.Life:%d \n", v.Life)
	}
	debug("=================isAllDead: 目标:%d  部队%d AllDead\n", target, legion)

	return true
}

func (gm *GameManage) Close() {

}

//获取下一个攻击目标
func (gm *GameManage) getNextActionTarget() {
	show("%d开始===========================行动方:%d 行动兵种:%d===========================\n", gm.b, gm.actionTargetNumber, gm.actionLegionNumber)
	if gm.actionTargetNumber == 0 { //切换目标
		gm.actionTargetNumber = 1
	} else {
		gm.actionTargetNumber = 0
		if gm.actionLegionNumber == 3 { //骑兵为最后一个行动部队
			gm.actionLegionNumber = 0 //切换为刀盾兵 行动
		} else {
			gm.actionLegionNumber++ //下一个 行动兵种
		}
	}
	gm.b++
	debug("结束===========================行动方:%d 行动兵种:%d===========================\n", gm.actionTargetNumber, gm.actionLegionNumber)
}

//范围攻击
func (gm *GameManage) battleStateAoeFighting() {

	soldier1 := new(Soldier)
	soldier2 := new(Soldier)

	debug("battleStateAoeFighting attackerNumber#%d \n", gm.attackerNumber)
	debug("battleStateAoeFighting defenderNumber#%d \n", gm.defenderNumber)
	for gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][gm.attackerNumber].Life == 0 {
		gm.attackerNumber++
		debug("battleStateAoeFighting attackerNumber#%d,%d,%0.2f \n", gm.actionLegionNumber, gm.attackerNumber, gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][gm.attackerNumber].Life)
	}

	//	for _, v := range gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber] {
	//		fmt.Printf("battleStateAoeFighting actionLegionNumber#%d,%d,%0.2f \n", gm.actionLegionNumber, gm.attackerNumber, v.Life)
	//		if v.Life == 0 {
	//			gm.attackerNumber++
	//			fmt.Printf("2.切换攻击方的士兵#%d \n", gm.attackerNumber)
	//		}
	//	}

	soldier1 = gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][gm.attackerNumber]

	for gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber][gm.defenderNumber].Life == 0 {
		gm.defenderNumber++
		debug("battleStateAoeFighting defenderLegionNumber#%d,%d,%0.2f \n", gm.defenderLegionNumber, gm.defenderNumber, gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber][gm.defenderNumber].Life)
	}

	//	for _, v := range gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber] {
	//		fmt.Printf("battleStateAoeFighting defenderLegionNumber#%d,%d,%0.2f \n", gm.defenderLegionNumber, gm.defenderNumber, v.Life)
	//		if v.Life == 0 {
	//			gm.defenderNumber++
	//			fmt.Printf("battleStateAoeFighting defenderNumber#%d \n", gm.defenderNumber)
	//			fmt.Printf("4.切换防御方的士兵#%d \n", gm.defenderNumber)
	//		}
	//	}

	soldier2 = gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber][gm.defenderNumber]

	gm.getRandomValue()

	hurtNum := soldier1.Power * ArmsPlus[gm.actionLegionNumber][gm.defenderLegionNumber] * gm.randomCoefficient() * math.Pow(1.05, gm.battleArray[gm.actionTargetNumber].InitPower-15)

	hurtNum *= 1 + gm.battleArray[gm.actionTargetNumber].dashPower/100

	show("%s 冲刺攻击 %s! 造成了 %0.2f 点伤害。\n", soldier1.Name, soldier2.Name, hurtNum)

	soldier2.Hurt(hurtNum)

	if gm.isAllDead(1-gm.actionTargetNumber, gm.defenderLegionNumber) {
		gm.switchBattleState(BATTLE_STATE_FIGHTING)
		gm.getNextActionTarget()
	} else {
		for {

			gm.defenderNumber++
			debug("battleStateAoeFighting defenderNumber#%d \n", gm.defenderNumber)
			debug("2.切换防御方的士兵#%d => %d \n", gm.defenderNumber, len(gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber]))
			if gm.defenderNumber >= len(gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber])-1 {
				gm.defenderNumber = 0
				debug("battleStateAoeFighting defenderNumber#%d \n", gm.defenderNumber)
				debug("3.切换防御方的士兵#%d \n", gm.defenderNumber)
				for {
					gm.attackerNumber++
					debug("3.切换攻击方的士兵#%d \n", gm.attackerNumber)
					if gm.attackerNumber >= len(gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber])-1 {
						break
					}
					if gm.battleArray[gm.actionTargetNumber].Army[gm.actionLegionNumber][gm.attackerNumber].Life <= 0 {
						continue
					}
				}
				gm.switchBattleState(BATTLE_STATE_FIGHTING)
				gm.getNextActionTarget()
				break
			}

			if gm.battleArray[1-gm.actionTargetNumber].Army[gm.defenderLegionNumber][gm.defenderNumber].Life > 0 {
				break
			}
		}
	}

}

func debug(format string, a ...interface{}) (n int, err error) {
	//return fmt.Fprintf(os.Stdout, format, a...)
	return 0, nil
}

func show(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, format, a...)
}
