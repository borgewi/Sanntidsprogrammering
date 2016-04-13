package Master_Slave

import (
	"Driver"
	"Elev_control"
	"time"
	//"fmt"
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
	return true
}

func check_elevsIdleAtFloor(){
	for{
		time.Sleep(500 * time.Millisecond)
		if isMaster{
			for _,elev := range Elevators_online{
				if elev.Behaviour == Elev_control.EB_Idle{
					All_btn_calls[elev.Floor][0] = false
					All_btn_calls[elev.Floor][1] = false
				} else if elev.Behaviour == Elev_control.EB_DoorOpen{
					switch(elev.Dir){
					case Elev_control.D_Idle:
						All_btn_calls[elev.Floor][0] = false
						All_btn_calls[elev.Floor][1] = false
					case Elev_control.D_Down:
						All_btn_calls[elev.Floor][0] = false
					case Elev_control.D_Up:
						All_btn_calls[elev.Floor][1] = false
					}
				}
			}
		}
	}
}


//All_btn_calls må oppdateres når en heis stopper i en etasje.
//Fjerne alle btn_calls som er i samme retning som den heisen.
//Eventuelt fjerne alle btn_calls i den etasjen dersom heisen har Dir D_Idle.
//Master må broadcste Elevators_online hver gang den oppdateres.
