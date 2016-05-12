package sggj

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Player struct {
	Name             string
	InitPower        float64
	Power            float64
	InitIntelligence float64
	Intelligence     float64
	InitCommand      float64
	Command          float64
	Level            float64
	ExpRequired      float64
	dashPower        float64

	Army [4][5]*Soldier

	RiderLegionState int //骑兵的状态
	BuildingArray    [4]int
}

func NewPlayer(name string, power float64, interlligence float64, command float64) *Player {
	p := &Player{
		Name:         name,
		Power:        power,
		Intelligence: interlligence,
		Command:      command,
	}
	p.init()
	return p
}

func (p *Player) decipher(code string) {
	arr := strings.Split(code, ",")
	p.Name = arr[0]
	tmp, _ := strconv.Atoi(arr[1])
	p.InitPower = float64(tmp)

	tmp1, _ := strconv.Atoi(arr[2])
	p.Level = float64(tmp1)

	avgStarNum, _ := strconv.Atoi(arr[3])

	p.BuildingArray[1], _ = strconv.Atoi(arr[4])
	p.BuildingArray[3], _ = strconv.Atoi(arr[5])
	avgEquipLv, _ := strconv.Atoi(arr[6])

	fmt.Printf("======================初始化【%s】部队====================\n", arr[0])
	num := 0
	for _, soldierId := range arr[7:] {
		sId, _ := strconv.Atoi(soldierId)
		if sId > 500 {
			break
		}
		soldierInfo := SoldierData[sId]
		//fmt.Printf("soldierInfo:%v  \n", soldierInfo)
		b := p.isLegionMax(soldierInfo.Type)
		if b != -1 {
			p.Army[soldierInfo.Type][b].Set(soldierInfo.Id, soldierInfo.Name, float64(avgStarNum), float64(soldierInfo.Power), float64(soldierInfo.Intelligence), p.Level, float64(p.BuildingArray[1]), float64(avgEquipLv))
			p.Army[soldierInfo.Type][b].Reflash()

			fmt.Printf("%s ", p.Army[soldierInfo.Type][b].Name)
			num++
		}
	}
	fmt.Printf("\n=========================部队数量:%d=========================\n", num)
}

func (p *Player) isLegionMax(t int) int {
	//fmt.Printf("isLegionMax:%d  len:%d\n", t, len(p.Army[t]))
	for i := 0; i < len(p.Army[t]); i++ {
		//fmt.Printf("Army:%v \n", p.Army[t][i])
		if p.Army[t][i].Id == 0 {
			return i
		}
	}
	return -1
}

func (p *Player) isSoldierMax() bool {
	i := 0
	for i < len(p.Army) {
		if p.Army[i][4].Id == 0 {
			return false
		}
		i++
	}
	return true
}

func (p *Player) init() {
	p.Level = 1
	p.ExpRequired = p.expCalculation()

	//fmt.Printf("Player %v\n", p)
	for t := 0; t < 4; t++ { //4兵种
		for i := 0; i < 5; i++ { //每种5个士兵
			p.Army[t][i] = NewSoldier()
			p.Army[t][i].SetType(t)
		}
	}

	for i, _ := range BuildingType {
		p.BuildingArray[i] = 1
	}

}

func (p *Player) expCalculation() float64 {
	return math.Floor(float64(100*p.Level) + math.Pow(p.Level, float64(1.3)) + math.Pow(float64(1.08), p.Level))
}

func (p *Player) buildingExpCalculation(num int) int {
	return 50 + num*num*num
}

func (p *Player) IsLevelUp() bool {
	if p.ExpRequired <= 0 {
		p.ExpRequired = p.expCalculation()
		p.Level++
		p.Reflash()
		return true
	}
	return false
}

func (p *Player) getExp(a float64, b float64) float64 {
	num := 3 * a * b * math.Sqrt(p.sumStar())
	ret := math.Floor(float64(num))
	p.ExpRequired = ret
	return ret
}

func (p *Player) sumStar() float64 {
	num := 0
	for t := 0; t < 4; t++ { //4兵种
		for i := 0; i < 5; i++ {
			num += int(p.Army[t][i].Star)
		}
	}
	return float64(num)
}

func (p *Player) Reflash() {
	p.Power = math.Floor(p.InitPower*math.Sqrt(p.Level) + 8*(p.Level-1))
	p.InitIntelligence = math.Floor(p.InitIntelligence*math.Sqrt(p.Level) + 8*(p.Level-1))
	p.Command = math.Floor(0.9*p.InitCommand*math.Sqrt(p.Level) + 2*(p.Level-1) + 1)

}

func (p *Player) averageStar() float64 {
	star := 0.0
	num := 0.0
	for t := 0; t < 4; t++ {
		for n := 0; n < 5; n++ {
			if p.Army[t][n].Id != 0 {
				num++
				star += p.Army[t][n].Star
			}
		}
	}
	return star / num
}

func (p *Player) averageEquipment() float64 {
	total := 0.0
	num := 0.0
	for t := 0; t < 4; t++ {
		for n := 0; n < 5; n++ {
			if p.Army[t][n].Id != 0 {
				num++
				total += p.Army[t][n].Equipment[WEAPON].Level
				total += p.Army[t][n].Equipment[ARMOR].Level
			}
		}
	}

	return total / 2 * num
}

func (p *Player) IsAllDead() bool {
	for t := 0; t < 4; t++ {
		for n := 0; n < 5; n++ {
			//fmt.Printf("IsAllDead %d=>%d Name:%s Life:%0.2f\n", t, n, p.Army[t][n].Name, p.Army[t][n].Life)

			if p.Army[t][n].Life > 0 {
				//fmt.Printf("%s 的部队 %s (%s#%d) 活着\n", p.Name, p.Army[t][n].Name, SoldierType[t], n)
				return false
			}
		}
	}

	fmt.Printf("%s 的部队 全挂了\n", p.Name)
	return true
}
