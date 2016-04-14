package Elev_control

import (
	"Driver"
	"fmt"
	"time"
)

type Elevator struct {
	Floor     int
	Dir       Direction
	Requests  [Driver.NUMFLOORS][Driver.NUMBUTTONS]bool
	Behaviour ElevatorBehaviour
	Elev_ID   int64
	//doorOpenDuration_s float
}

//var requests_timeStamp [4][3]int64 //Setter timeStamp når en ordre aktiveres.

var (
	elevator      Elevator
	lastFloorTime int64
)

func Run_Elevator(localStatusCh chan Elevator, sendBtnCallsCh chan [2]int, errorCh chan int) {
	//var elevator2 Elevator

	//Init elev_state
	if Driver.ElevGetFloorSensorSignal() == -1 {
		fsm_onInitBetweenFloors()
	}
	fsm_elevatorUninitialized()
	fmt.Printf("%+v", elevator.Elev_ID)
	fmt.Println("")
	go send_status(localStatusCh)

	running := true
	var prev_button [Driver.NUMFLOORS][Driver.NUMBUTTONS]int
	var prev_floor int
	prev_floor = Driver.ElevGetFloorSensorSignal()

	go checkElevMoving(errorCh)

	for running {
		// Request button
		for f := 0; f < Driver.NUMFLOORS; f++ {
			for b := 0; b < Driver.NUMBUTTONS; b++ {
				v := Driver.ElevGetButtonSignal(b, f)
				if v&int(v) != prev_button[f][b] {
					fsm_onRequestButtonPress(f, Button(b), sendBtnCallsCh)
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

func send_status(localStatusCh chan Elevator) {
	for {
		time.Sleep(1000 * time.Millisecond)
		localStatusCh <- elevator
	}
}

func checkElevMoving(errorCh chan int) {
	var errorTime int64
	var timeNow int64
	errorTime = 6
	for {
		if elevator.Dir != D_Idle {
			timeNow = GetActiveTime()
			if lastFloorTime-timeNow > errorTime {
				errorCh <- 1
			}
		}
	}
}

/*func review_timeStamps(errorCh chan int){
	var errorTime int64
	errorTime =
	for{
		timeNow := getActiveTime()
		for _,timeStamp in range requests_timeStamp{
			if
		}
	}
}
*/
func PrintElev(elev Elevator) {
	fmt.Println("")
	fmt.Println("Floor: ", elev.Floor)
	fmt.Println("Direction: ", elev.Dir)
	for f := 0; f < 4; f++ {
		fmt.Printf("%+v", elev.Requests[f])
		fmt.Println("")
	}
}
