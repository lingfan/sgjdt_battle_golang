package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"./sggj"
)

var ArgA = flag.Int("a", 10, "攻击方")
var ArgB = flag.Int("b", 20, "防御方")

func main() {
	fmt.Println("GAME Start")
	path := GetCurrPath()

	flag.Parse()

	cfg := sggj.NewConfig()
	cfg.Init(path)

	g01 := sggj.NewGame(*ArgA, *ArgB)
	//	g.Cfg = cfg
	g01.Start()

	fmt.Println("GAME END")
}

func GetCurrPath() string {

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}
