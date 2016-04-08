package main

import (
	"Elev_control"
	"fmt"
	//"net"
	//"os/exec"
	//"time"
	"Driver"
	"Network"
)

const (
	NUMBUTTONS  int = 3
	NUMFLOORS   int = 4
	BUTTON_DOWN     = 0 + iota
	BUTTON_UP
	BUTTON_COMMAND
)

func main() {
	receiveCh := make(chan Elev_control.Elevator)
	timeoutCh := make(chan int)
	Driver.ElevInit()
	go Elev_control.Merry_go_around(receiveCh)
	//go get_status_and_broadcast(receiveCh)
	go receive_status(receiveCh, timeoutCh)
	for {
	}
}

func get_status_and_broadcast(receiveCh chan Elev_control.Elevator) {
	//var data []byte
	var elev Elev_control.Elevator
	toAdress := "129.241.187.161" + ":13337"
	for {
		elev = <-receiveCh
		Network.UdpBroadcast(toAdress, elev)

		//fmt.Printf("%+v", elev)
	}
}

func receive_status(receiveCh chan Elev_control.Elevator, timeoutCh chan int) {
	listenAdress := "129.241.187.161" + ":13337"
	go Network.ListenToBroadcast(listenAdress, receiveCh, timeoutCh)
	var elev Elev_control.Elevator
	for {

		elev = <-receiveCh
		fmt.Printf("%+v", elev)
	}
}

// func button_pressed(floor int, button int) {
// 	Driver.ElevSetMotorDirection(1)
// 	fmt.Println("drive")
// }

// func check_buttons() {
// 	for sfloor := 0; sfloor < 3; sfloor++ {
// 		if Driver.ElevGetButtonSignal(BUTTON_UP, sfloor) == 1 {
// 			button_pressed(sfloor, BUTTON_UP)
// 			fmt.Println("opp")
// 		}
// 	}
// 	for sfloor := 1; sfloor < 4; sfloor++ {
// 		if Driver.ElevGetButtonSignal(BUTTON_DOWN, sfloor) == 1 {
// 			button_pressed(sfloor, BUTTON_DOWN)
// 		}
// 	}
// 	for sfloor := 0; sfloor < 4; sfloor++ {
// 		if Driver.ElevGetButtonSignal(BUTTON_COMMAND, sfloor) == 1 {
// 			button_pressed(sfloor, BUTTON_COMMAND)
// 		}
// 	}
// }
