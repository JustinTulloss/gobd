package main

import (
	"flag"
	"fmt"

	"github.com/JustinTulloss/gobd"
)

func main() {
	flag.Parse()
	fmt.Println("Hello! Testing things out now.")
	obd, err := gobd.NewOBD(flag.Arg(0))
	if err != nil {
		panic(err.Error())
	}
	defer obd.Close()
	fmt.Println("Finding the available pids")
	err = obd.SendCommand([]byte("010C"))
	if err != nil {
		panic(err.Error())
	}
	resp, err := obd.ReadResult()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(resp))
}
