package main

import (
	_ "fmt"
	"github.com/kshmirko/rt3v1/rtcode"
)

func main() {

	v := rtcode.New()
	v.DoCalc()
	a, _ := v.UnmarshalData()
	_, _, _ = a.DumpDownwardRadiation(true)
	//fmt.Println(mu, I, Q)
}
