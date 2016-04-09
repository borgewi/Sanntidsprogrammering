package main

import (
	"Elev_control"
	//"net"
	//"os/exec"
	//"time"
	"Driver"
	//"Network"
	"Master_Slave"
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
	receiveCh := make(chan Elev_control.Elevator)
	
	go Elev_control.Run_Elevator(receiveCh)
	go Master_Slave.Determine_Rank(receiveCh)
	for {
	}
}

