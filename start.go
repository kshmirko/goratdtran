package main

import (
	"fmt"
	"github.com/kshmirko/rt3v1/rtcode"
)

func main() {

	v := rtcode.New()
	a, _ := v.UnmarshalData()
	//v.DoCalc()

	fmt.Println(a.Z, len(*a.Mu))

}
