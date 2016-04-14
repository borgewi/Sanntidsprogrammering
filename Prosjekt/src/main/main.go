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
	sendBtnCallCh := make(chan [2]int, 5)
	errorCh := make(chan int)
	receiveAllBtnCallsCh := make(chan [4][2]bool, 5)

	go Elev_control.Run_Elevator(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, errorCh)
	go Master_Slave.Princess(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, errorCh)
	for {
	}
}
