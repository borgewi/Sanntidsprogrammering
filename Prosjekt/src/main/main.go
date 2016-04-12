package main

import (
	"Elev_control"
	//"net"
	//"os/exec"
	//"time"
	"Driver"
	//"Network"
	"Master_Slave"
	//"fmt"
)

const (
	NUMBUTTONS  int = 3
	NUMFLOORS   int = 4
	BUTTON_DOWN     = 0 + iota
	BUTTON_UP
	BUTTON_COMMAND
)

func main() {
	Driver.ElevInit()
	localStatusCh := make(chan Elev_control.Elevator)
	sendBtnCallsCh := make(chan [2]int)

	go Elev_control.Run_Elevator(localStatusCh, sendBtnCallsCh)
	go Master_Slave.Princess(localStatusCh, sendBtnCallsCh)
	for {
	}
}
