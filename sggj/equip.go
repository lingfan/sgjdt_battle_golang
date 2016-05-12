package sggj

import (
	"math"
)

type Equip struct {
	Id    int
	Type  int
	Value float64
	Name  string
	Level float64
}

//铁匠铺等级,装备等级
func NewEquip() *Equip {
	e := new(Equip)
	return e
}

func (e *Equip) Init() {
	e.Id = 0
	e.Value = 0
	e.Name = "空名字"
	e.Level = 0
}

func (e *Equip) newLevel() {
	e.Level++
}

func (e *Equip) Update(smithyLv float64, level float64) {
	e.Id = int(math.Max(float64(smithyLv/5+1), float64(40)))
	e.Level = level
	e.Name = EquipData[e.Id].Name
	e.Value = EquipData[e.Id].Value
}
