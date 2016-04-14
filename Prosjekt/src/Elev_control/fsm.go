package Elev_control

import (
	"Driver"
	"fmt"
	"time"
	//"encoding/json"
)

//Gjøre om denne til å styres av interne kommandoer OG når all_btn_calls mottas.
func setAllLights() {
	for floor := 0; floor < Driver.NUMFLOORS; floor++ {
		for btn := 0; btn < Driver.NUMBUTTONS; btn++ {
			if btn == int(B_Cab) {
				Driver.ElevSetButtonLight(btn, floor, elevator.Requests[floor][btn])
			} else {
				Driver.ElevSetButtonLight(btn, floor, allExtBtns[floor][btn])
			}
		}
	}
}

func updateAllExtLights(receiveAllBtnCallsCh chan [4][2]bool) {
	for {
		allExtBtns = <-receiveAllBtnCallsCh
		fmt.Println("Mottar allExtBtns fra receiveAllBtnCallsCh")
		setAllLights()
	}
}

func fsm_onInitBetweenFloors() {
	Driver.ElevSetMotorDirection(int(D_Down))
	elevator.Dir = D_Down
	elevator.Behaviour = EB_Moving
	lastFloorTime = GetActiveTime()
	for Driver.ElevGetFloorSensorSignal() == -1 {
	}
	Driver.ElevSetMotorDirection(int(D_Idle))
}

func fsm_onRequestButtonPress(btn_floor int, btn_type Button, sendBtnCallCh chan [2]int) {
	switch btn_type {
	case B_Cab:
		fsm_onNewActiveRequest(btn_floor, btn_type)
	case B_HallDown:
		fsm_SendNewOrderToMaster(btn_floor, btn_type, sendBtnCallCh)
	case B_HallUp:
		fsm_SendNewOrderToMaster(btn_floor, btn_type, sendBtnCallCh)
	}
}

func fsm_onNewActiveRequest(btn_floor int, btn_type Button) {
	switch elevator.Behaviour {
	case EB_DoorOpen:
		if elevator.Floor == btn_floor {
			timer_start(3000 * time.Millisecond)
		} else {
			elevator.Requests[btn_floor][btn_type] = true
			//fsm_setTimeStamp(btn_floor,btn_type)
		}
		break
	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = true
		//fsm_setTimeStamp(btn_floor,btn_type)
		break
	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = true
		//fsm_setTimeStamp(btn_floor,btn_type)
		elevator.Dir = requests_chooseDirection(elevator)
		if elevator.Dir == D_Idle {
			Driver.ElevSetDoorLight(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			timer_start(3000 * time.Millisecond)
			elevator.Behaviour = EB_DoorOpen
		} else {
			Driver.ElevSetMotorDirection(int(elevator.Dir))
			elevator.Behaviour = EB_Moving
			lastFloorTime = GetActiveTime()
		}
		break
	}

	setAllLights()
}

func fsm_SendNewOrderToMaster(btn_floor int, btn_type Button, sendBtnCallCh chan [2]int) {
	var newBtnCall [2]int
	newBtnCall[0] = btn_floor
	newBtnCall[1] = int(btn_type)
	sendBtnCallCh <- newBtnCall
}

func fsm_onFloorArrival(newFloor int) {
	elevator.Floor = newFloor
	Driver.ElevSetFloorIndicator(newFloor)
	switch elevator.Behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			Driver.ElevSetMotorDirection(int(D_Idle))
			Driver.ElevSetDoorLight(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			timer_start(3000 * time.Millisecond)
			setAllLights()
			elevator.Behaviour = EB_DoorOpen
			//fsm_deleteTimeStamp(newFloor)
			lastFloorTime = GetActiveTime()
		}
		break
	default:
		break
	}
}

func fsm_onDoorTimeout() {
	switch elevator.Behaviour {
	case EB_DoorOpen:
		elevator.Dir = requests_chooseDirection(elevator)
		Driver.ElevSetMotorDirection(int(elevator.Dir))
		Driver.ElevSetDoorLight(false)
		if elevator.Dir == D_Idle {
			elevator.Behaviour = EB_Idle
		} else {
			elevator.Behaviour = EB_Moving
			lastFloorTime = GetActiveTime()
		}
		break
	default:
		break
	}
}

func fsm_elevatorUninitialized() {
	elevator.Dir = D_Idle
	elevator.Behaviour = EB_Idle
	elevator.Floor = Driver.ElevGetFloorSensorSignal()
	elevator.Elev_ID = GetActiveTime()
	for f := 0; f < Driver.NUMFLOORS; f++ {
		for b := 0; b < Driver.NUMBUTTONS; b++ {
			elevator.Requests[f][b] = false
		}
	}
}

func Fsm_addOrder(Order [2]int, Order_ID int64) {
	if Order_ID == elevator.Elev_ID {
		fsm_onNewActiveRequest(Order[0], Button(Order[1]))
	} else {
		fmt.Println("Feil Order_ID")
	}
}

/*
func fsm_setTimeStamp(btn_floor int,btn_type Button){
	if requests_timeStamp[btn_floor][btn_type] == 0{
		requests_timeStamp[btn_floor][btn_type] = GetActiveTime()
	}
}

func fsm_deleteTimeStamp(newFloor int){
	switch(elevator.Dir){
	case D_Down:
		requests_timeStamp[newFloor][B_HallDown] = 0
		requests_timeStamp[newFloor][B_Cab] = 0
		if newFloor == 0{
			requests_timeStamp[newFloor][B_HallUp] = 0
		}
	case D_Up:
		requests_timeStamp[newFloor][B_HallUp] = 0
		requests_timeStamp[newFloor][B_Cab] = 0
		if newFloor == Driver.NUMFLOORS-1{
			requests_timeStamp[newFloor][B_HallDown] = 0
		}
	case D_Idle:
		requests_timeStamp[newFloor][B_HallDown] = 0
		requests_timeStamp[newFloor][B_HallUp] = 0
		requests_timeStamp[newFloor][B_Cab] = 0
	}
}

*/
