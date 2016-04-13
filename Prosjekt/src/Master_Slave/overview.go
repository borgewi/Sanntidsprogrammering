package Master_Slave

import (
	"Driver"
	"Elev_control"
	//"fmt"
	"time"
)

/*type Elevator struct {
	Floor     int
	Dir       Direction
	Requests  [Driver.NUMFLOORS][Driver.NUMBUTTONS]bool
	Behaviour ElevatorBehaviour
	Elev_ID	  int64
	//doorOpenDuration_s float
}*/

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

func setTimeStamp(btn_floor int, btn_type int) {
	btn_calls_timeStamp[btn_floor][btn_type] = Elev_control.GetActiveTime()
}

func checkTimeStamps(handleOrderAgainCh chan [2]int) {
	var errorTime int64
	errorTime = 15
	var order [2]int
	for {
		time.Sleep(1500 * time.Millisecond)
		timeNow := Elev_control.GetActiveTime()
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

func check_elevsIdleAtFloor() {
	for {
		time.Sleep(500 * time.Millisecond)
		if isMaster {
			for _, elev := range Elevators_online {
				if elev.Behaviour == Elev_control.EB_Idle || elev.Behaviour == Elev_control.EB_DoorOpen {
					All_btn_calls[elev.Floor][Elev_control.B_HallDown] = false
					btn_calls_timeStamp[elev.Floor][Elev_control.B_HallDown] = 0
					All_btn_calls[elev.Floor][Elev_control.B_HallUp] = false
					btn_calls_timeStamp[elev.Floor][Elev_control.B_HallUp] = 0
				}
			}
		}
	}
}

//All_btn_calls må oppdateres når en heis stopper i en etasje.
//Fjerne alle btn_calls som er i samme retning som den heisen.
//Eventuelt fjerne alle btn_calls i den etasjen dersom heisen har Dir D_Idle.
//Master må broadcste Elevators_online hver gang den oppdateres.
