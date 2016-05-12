package sggj

import (
	"fmt"
	"math"
)

const (
	WEAPON = 0
	ARMOR  = 1
)

type Soldier struct {
	Id               int
	Name             string
	Star             float64
	Power            float64
	InitPower        float64
	Intelligence     float64
	InitIntelligence float64
	Life             float64
	MaxLife          float64
	Type             int
	Level            float64
	ExpRequired      float64
	Equipment        [2]*Equip
}

func NewSoldier() *Soldier {
	s := &Soldier{
		Id:               0,
		Name:             "",
		Star:             0,
		Power:            0,
		InitPower:        0,
		Intelligence:     0,
		InitIntelligence: 0,
		Life:             0,
		MaxLife:          0,
		Type:             0,
		Level:            0,
		ExpRequired:      0,
	}

	s.Equipment[WEAPON] = NewEquip()
	s.Equipment[ARMOR] = NewEquip()
	return s
}

func (s *Soldier) init() {
	fmt.Println("1")
}

func (s *Soldier) Set(id int, name string, star float64, initPower float64, initIntelligence float64, level float64, equip1 float64, equip2 float64) {
	s.Id = id
	s.Name = name
	s.Star = star
	s.InitPower = initPower
	s.InitIntelligence = initIntelligence
	s.Level = level
	if equip1 == -1 {
		s.Equipment[WEAPON].Init()
		s.Equipment[ARMOR].Init()
	} else if equip1 != -2 {
		s.Equipment[WEAPON].Update(equip1, equip2)
		s.Equipment[ARMOR].Update(equip1, equip2)
	}
	//fmt.Printf("Equipment[WEAPON]:%v\n", s.Equipment[WEAPON])
	//fmt.Printf("Equipment[ARMOR]:%v\n", s.Equipment[ARMOR])
}

func (s *Soldier) Hurt(num float64) {
	s.Life = math.Max(s.Life-num, 0)
}

func (s *Soldier) Recover() {
	s.Life = s.MaxLife
}

func (s *Soldier) countPower() {
	num := 0.5*s.InitPower*s.Star*math.Sqrt(s.Level) + 2.5*s.Equipment[WEAPON].Value*s.Equipment[WEAPON].Level + 0.5*s.InitPower

	s.Power = math.Floor(num)
}

func (s *Soldier) countIntelligence() {
	num := s.InitIntelligence * s.Star * math.Sqrt(s.Level)
	s.InitIntelligence = math.Floor(num)
}

func (s *Soldier) expCalculation() float64 {
	num := 100 + math.Pow(s.Level, 1.5)*100
	return math.Floor(num)
}

func (s *Soldier) maxLifeCalculation() float64 {
	num := 5*s.InitPower*s.Star*math.Sqrt(s.Level) + 25*s.Equipment[ARMOR].Value*s.Equipment[ARMOR].Level
	return math.Floor(float64(num))
}

func (s *Soldier) getExp(num float64) bool {
	s.ExpRequired = s.ExpRequired - num
	if s.ExpRequired <= 0 {
		s.newLevel()
		return true
	}
	if s.ExpRequired > 1000000 || s.ExpRequired < -100 {
		s.ExpRequired = 2000
	}
	return false
}

func (s *Soldier) newLevel() {
	s.Level++
	s.ExpRequired = s.expCalculation()
	s.MaxLife = s.maxLifeCalculation()
	s.countPower()
	s.countIntelligence()
}

func (s *Soldier) RemoveSoldier() {
	s.Id = 0
	s.Name = ""
	s.Star = 0
	s.Power = 0
	s.InitPower = 0
	s.Intelligence = 0
	s.InitIntelligence = 0
	s.Life = 0
	s.MaxLife = 0
	s.Type = 0
	s.Level = 0
	s.ExpRequired = 0
	s.Equipment[WEAPON].Init()
	s.Equipment[ARMOR].Init()
}

func (s *Soldier) SoldierValue() float64 {
	return s.Star * math.Sqrt(s.InitPower)
}

func (s *Soldier) NewStar() {
	s.Star++
	s.Level = 1
	s.MaxLife = s.maxLifeCalculation()
	s.ExpRequired = s.expCalculation()
	s.countPower()
	s.countIntelligence()
}

func (s *Soldier) SetType(val int) {
	s.Type = val
	t1 := 0
	t2 := 1
	switch val {
	case 1:
		t1 = 2
		t2 = 3
	case 2:
		t1 = 4
		t2 = 3
	case 3:
		t1 = 5
		t2 = 3
	}
	s.Equipment[WEAPON].Type = t1
	s.Equipment[ARMOR].Type = t2
}

func (s *Soldier) Reflash() {
	s.MaxLife = s.maxLifeCalculation()
	s.Life = s.MaxLife
	s.countPower()
	s.countIntelligence()
}

func (s *Soldier) IsDead() bool {
	if s.Life == 0 {
		return true
	}
	return false
}
