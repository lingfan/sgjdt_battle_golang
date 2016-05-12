package sggj

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CfgSoldier struct {
	Id           int
	Name         string
	Power        int
	Intelligence int
	MaxStar      int
	Type         int
}
type CfgEquip struct {
	Id    int
	Name  string
	Value float64
}
type CfgPlayer struct {
	Id   int
	Code string
}

type Config struct {
	Soldier *map[int]CfgSoldier
	Equip   *map[int]CfgEquip
	Player  *map[int]CfgEquip
}

func NewConfig() *Config {
	c := new(Config)
	return c
}

var SoldierData = map[int]CfgSoldier{}

var EquipData = map[int]CfgEquip{}

var PlayerData = map[int]CfgPlayer{}

func (cfg *Config) Init(path string) {

	file1 := path + "data/soldier.csv"
	err1 := cfg.parseSoldier(file1)
	if err1 != nil {
		fmt.Println(err1)
	}

	file2 := path + "data/equipment.csv"
	err2 := cfg.parseEquip(file2)
	if err2 != nil {
		fmt.Println(err2)
	}

	file3 := path + "data/player.csv"
	err3 := cfg.parsePlayer(file3)
	if err3 != nil {
		fmt.Println(err3)
	}
}

func (cfg *Config) parseSoldier(file string) error {
	//soldiers := make(map[int]CfgSoldier)
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return err
		}

		power, err := strconv.Atoi(line[2])
		if err != nil {
			return err
		}

		intelligence, err := strconv.Atoi(line[3])
		if err != nil {
			return err
		}

		star, err := strconv.Atoi(line[4])
		if err != nil {
			return err
		}

		t, err := strconv.Atoi(line[5])
		if err != nil {
			return err
		}

		SoldierData[id] = CfgSoldier{
			Id:           id,
			Name:         line[1],
			Power:        power,
			Intelligence: intelligence,
			MaxStar:      star,
			Type:         t,
		}

		//fmt.Printf("Soldier %v\n", soldiers[id])
	}
	//cfg.Soldier = &soldiers
	return nil
}

func (cfg *Config) parseEquip(file string) error {
	//equip := make(map[int]CfgEquip)
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return err
		}

		val, err := strconv.Atoi(line[2])
		if err != nil {
			return err
		}

		EquipData[id] = CfgEquip{
			Id:    id,
			Name:  line[1],
			Value: float64(val),
		}
	}

	//cfg.Equip = &equip
	return nil
}

func (cfg *Config) parsePlayer(file string) error {
	//player := make(map[int]CfgPlayer)
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return err
		}

		PlayerData[id] = CfgPlayer{
			Id:   id,
			Code: strings.Join(line[1:], ","),
		}
	}

	//cfg.Equip = &player
	return nil
}
