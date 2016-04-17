package main

import (
	"Elev_control"
	//"net"
	"os/exec"
	//"time"
	"Driver"
	"Master_Slave"
	"fmt"
)

const (
	NUMBUTTONS  int = 3
	NUMFLOORS   int = 4
	BUTTON_DOWN     = 0 + iota
	BUTTON_UP
	BUTTON_COMMAND
)

func main() {
	//go Backup()
	Driver.ElevInit()
	localStatusCh := make(chan Elev_control.Elevator, 100)
	sendBtnCallCh := make(chan [NUMBUTTONS - 1]int, 100)
	errorCh := make(chan int)
	receiveAllBtnCallsCh := make(chan [NUMFLOORS][NUMBUTTONS - 1]bool, 100)
	setLights_setExtBtnsCh := make(chan [4][2]bool, 100)

	go Elev_control.Run_Elevator(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, setLights_setExtBtnsCh, errorCh)
	go Master_Slave.Princess(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, setLights_setExtBtnsCh, errorCh)
	checkForError(errorCh)
	callBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	callBackup.Run()
}

func checkForError(errorCh chan int) {
	var err int
	for {
		err = <-errorCh
		if err == Elev_control.ERR_MOTORSTOP {
			fmt.Println("Error har oppstått. Har vært i EB_Moving for lenge. err = ", err)
			for {
				err = <-errorCh
				if err == Elev_control.ERR_NO_ERROR {
					fmt.Println("Error er fikset")
					break
				} else if err == Elev_control.ERR_NO_ELEVS_OPERABLE {
					fmt.Println("Ingen heiser er operatible.\nStarter program på nytt på ny terminal")
					return
				}
			}
		} else if err == Elev_control.ERR_NO_ELEVS_OPERABLE {
			fmt.Println("Ingen heiser er operatible.\nStarter program på nytt på ny terminal")
			return
		}
	}
}
