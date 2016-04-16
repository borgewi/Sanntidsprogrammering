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

func updateAllExtLights(receiveAllBtnCallsCh chan [4][2]bool, allExtBtns [4][2]bool) {
	for {
		allExtBtns = <-receiveAllBtnCallsCh
		setAllLights()
	}
}

func fsm_onInitBetweenFloors() {
	Driver.ElevSetMotorDirection(int(D_Down))
	elevator.Dir = D_Down
	elevator.Behaviour = EB_Moving
	//lastFloorTime = time.Now().Unix()
	for Driver.ElevGetFloorSensorSignal() == -1 {
	}
	Driver.ElevSetMotorDirection(int(D_Idle))
}

func fsm_onRequestButtonPress(btn_floor int, btn_type Button, sendBtnCallCh chan [2]int) bool {
	switch btn_type {
	case B_Cab:
		fsm_onNewActiveRequest(btn_floor, btn_type)
		return true
	case B_HallDown:
		fsm_SendNewOrderToMaster(btn_floor, btn_type, sendBtnCallCh)
		return false
	case B_HallUp:
		fsm_SendNewOrderToMaster(btn_floor, btn_type, sendBtnCallCh)
		return false
	}
	return false
}

func fsm_onNewActiveRequest(btn_floor int, btn_type Button) {
	fmt.Println("\nKommer inn i fsm_onani\nFloor: ", btn_floor, " Button: ", btn_type)
	switch elevator.Behaviour {
	case EB_DoorOpen:
		fmt.Println("EB_DoorOpen")
		if elevator.Floor == btn_floor {
			timer_start(3000 * time.Millisecond)
		} else {
			elevator.Requests[btn_floor][btn_type] = true
		}
		break
	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = true
		break
	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = true
		//fmt.Println("Kjører requests_chooseDirection")
		elevator.Dir = requests_chooseDirection(elevator)
		//fmt.Println("Kommer ut av requests_chooseDirection")
		if elevator.Dir == D_Idle {
			Driver.ElevSetDoorLight(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			//timer_start(3000 * time.Millisecond)
			elevator.Behaviour = EB_DoorOpen
		} else {
			Driver.ElevSetMotorDirection(int(elevator.Dir))
			elevator.Behaviour = EB_Moving
			//lastFloorTime = time.Now().Unix()
		}
		break
	}
	setAllLights()
	fmt.Println("Kommer ut av fsm_onani")
}

func fsm_SendNewOrderToMaster(btn_floor int, btn_type Button, sendBtnCallCh chan [2]int) {
	var newBtnCall [2]int
	newBtnCall[0] = btn_floor
	newBtnCall[1] = int(btn_type)
	sendBtnCallCh <- newBtnCall
}

func Fsm_addOrder(Order [2]int, Order_ID int64) {

	if Order_ID == elevator.Elev_ID {
		fsm_onNewActiveRequest(Order[0], Button(Order[1]))
		fmt.Println("")
		fmt.Println(" Mottat ordre og kjører fsm_onNewActiveRequest")
	} else {
		fmt.Println("Feil Order_ID")
	}
}

func fsm_onFloorArrival(newFloor int) {
	elevator.Floor = newFloor
	Driver.ElevSetFloorIndicator(newFloor)
	//fsm_checkExtRequestsStillActive() //Skal vi ha med denne?
	switch elevator.Behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			Driver.ElevSetMotorDirection(int(D_Idle))
			Driver.ElevSetDoorLight(true)
			elevator = requests_clearAtCurrentFloor(elevator)
			timer_start(3000 * time.Millisecond)
			setAllLights()
			elevator.Behaviour = EB_DoorOpen
			//lastFloorTime = time.Now().Unix()
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
			//lastFloorTime = time.Now().Unix()
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

//Skal vi ha med denne?
/*func fsm_checkExtRequestsStillActive() {
	for i, k := range elevator.Requests {
		for j, _ := range k {
			if j != 2 && elevator.Requests[i][j] {
				if !allExtBtns[i][j] {
					//Bestillingen er allerede utført av noen andre
					elevator.Requests[i][j] = false
				}
			}
		}
	}
}
*/
