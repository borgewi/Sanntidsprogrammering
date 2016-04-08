package Elev_control

import (
	"Driver"
	"fmt"
	"time"
)

func Merry_go_around(recieveCh chan Elevator) {
	//var elevator2 Elevator

	//Init elev_state
	if Driver.ElevGetFloorSensorSignal() == -1 {
		fsm_onInitBetweenFloors()
	}
	elevator_uninitialized()
	//go send_status(recieveCh)
	running := true
	var prev_button [Driver.NUMFLOORS][Driver.NUMBUTTONS]int
	var prev_floor int
	prev_floor = Driver.ElevGetFloorSensorSignal()
	for running {
		// Request button
		for f := 0; f < Driver.NUMFLOORS; f++ {
			for b := 0; b < Driver.NUMBUTTONS; b++ {
				v := Driver.ElevGetButtonSignal(b, f)
				if v&int(v) != prev_button[f][b] {
					fsm_onRequestButtonPress(f, Button(b))
				}
				prev_button[f][b] = v
			}
		}
		// Floor sensor
		f := Driver.ElevGetFloorSensorSignal()
		if f != -1 {
			if f != prev_floor {
				fsm_onFloorArrival(f)
			}
		}
		prev_floor = f
		// Timer
		if timer_timedOut() {
			fmt.Println("timer stoppet")
			fsm_onDoorTimeout()
			timer_stop()
		}
		time.Sleep(25 * time.Millisecond) //Hardkoding
	}
}
