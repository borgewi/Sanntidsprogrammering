package Elev_control

import (
	"Driver"
	"fmt"
	"time"
)

func Run_Elevator(statusCh chan Elevator) {
	//var elevator2 Elevator

	//Init elev_state
	if Driver.ElevGetFloorSensorSignal() == -1 {
		fsm_onInitBetweenFloors()
	}
	elevator_uninitialized()
	go send_status(statusCh)
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
			fsm_onDoorTimeout()
			timer_stop()
		}
		time.Sleep(25 * time.Millisecond) //Hardkoding
	}
}

func send_status(statusCh chan Elevator) {
	for {
		time.Sleep(1000 * time.Millisecond)
		//var data []byte
		//data = json.Marshal(elevator)
		statusCh <- elevator
	}
}

func PrintElev(elev Elevator) {
	fmt.Println("")
	fmt.Println("Floor: ", elev.Floor)
	fmt.Println("Direction: ", elev.Dir)
	for f := elev.Floor + 1; f < 4; f++ {
		fmt.Printf("%+v", elev.Requests[f])
		fmt.Println("")
	}
}
