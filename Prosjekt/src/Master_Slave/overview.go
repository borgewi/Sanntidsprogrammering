package Master_Slave

import (
	"Driver"
	"Elev_control"
	"fmt"
	//"sync"
	"time"
)

//var Fuckit = sync.Mutex{}
var Elevators_online []Elev_control.Elevator
var All_btn_calls [Driver.NUMFLOORS][Driver.NUMBUTTONS - 1]bool
var btn_calls_timeStamp [Driver.NUMFLOORS][Driver.NUMBUTTONS - 1]int64

func update_Elevators_online(curr_elev Elev_control.Elevator) {
	for i, elev := range Elevators_online {
		if elev.Elev_ID == curr_elev.Elev_ID {
			Elevators_online[i] = curr_elev
			return
		}
	}
	Elevators_online = append(Elevators_online, curr_elev)
}

func delete_All_elevs() {
	Elevators_online = Elevators_online[:1]
}

func update_btnCalls(newCall [2]int) bool {
	if All_btn_calls[newCall[0]][newCall[1]] {
		return false
	}
	All_btn_calls[newCall[0]][newCall[1]] = true
	setTimeStamp(newCall[0], newCall[1])
	return true
}

func getElevators_Online() []Elev_control.Elevator {
	return Elevators_online
}

func setElevators_Online(elevs_online []Elev_control.Elevator) {
	Elevators_online = elevs_online
}

func get_All_btn_calls() [4][2]bool {
	return All_btn_calls
}

func setAll_btn_calls(btn_calls [4][2]bool) {
	All_btn_calls = btn_calls
}

/*func getOrSet_ElevatorsOnline() {
	for {
		select {
		case <-getToUpdateCh:
			ToUpdateCh <- Elevators_online
		case <-getToCostFunctionCh:
			ToCostFunctionCh <- Elevators_online
		case Elevators_online = <-Set_ElevatorsOnlineCh:
			break
		}
	}
}*/

func setTimeStamp(btn_floor int, btn_type int) {
	btn_calls_timeStamp[btn_floor][btn_type] = time.Now().Unix()
	//fmt.Println("Satte timeStamp")
}

func checkTimeStamps(handleOrderAgainCh chan [2]int) {
	var errorTime int64
	var timeNow int64
	errorTime = 10
	var order [2]int
	for {
		if !isMaster {
			time.Sleep(10000 * time.Millisecond)
		} else {
			time.Sleep(1500 * time.Millisecond)
			timeNow = time.Now().Unix()
			fmt.Println("")
			fmt.Printf("%+v", btn_calls_timeStamp)
			fmt.Println("")
			for i, k := range btn_calls_timeStamp {
				for j, timeStamp := range k {
					if timeStamp != 0 {
						if timeNow-timeStamp > errorTime {
							order[0] = i
							order[1] = j
							handleOrderAgainCh <- order
							setTimeStamp(i, j)
						}
					}
				}
			}
		}
	}
}

func check_elevsIdleAtFloor() {
	time.Sleep(500 * time.Millisecond)
	if isMaster {
		//Fuckit.Lock()
		for _, elev := range Elevators_online {
			if elev.Behaviour == Elev_control.EB_Idle || elev.Behaviour == Elev_control.EB_DoorOpen {
				All_btn_calls[elev.Floor][Elev_control.B_HallDown] = false
				btn_calls_timeStamp[elev.Floor][Elev_control.B_HallDown] = 0
				All_btn_calls[elev.Floor][Elev_control.B_HallUp] = false
				btn_calls_timeStamp[elev.Floor][Elev_control.B_HallUp] = 0
			}
		}
		//Fuckit.Unlock()
	}
}
