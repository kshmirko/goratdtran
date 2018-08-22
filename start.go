package main

import (
	_ "fmt"
	"github.com/kshmirko/rt3v1/rtcode"
	"os"
	"path/filepath"
)

func main() {

	v := rtcode.New()
	v.SetSizeDistrib(0.1, 1.0, -3.5, 101)
	v.SetGalbedo(0.0)
	v.DoCalc()
	a, _ := v.UnmarshalData()
	_, _, _ = a.DumpDownwardRadiation(true)

	//Поиск и удаление всех файов с индикатриссой рассеяния
	files, err := filepath.Glob(".scat*")

	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
	//fmt.Println(mu, I, Q)
}
