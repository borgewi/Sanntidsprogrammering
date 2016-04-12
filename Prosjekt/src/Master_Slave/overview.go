package Master_Slave

import (
	"Elev_control"
	"Driver"
	//"time"
	//"fmt"
)


// func oversikt skal vite hvor alle heisene er og hvilken retning de har. 

/*type Elevator struct {
	Floor     int
	Dir       Direction
	Requests  [Driver.NUMFLOORS][Driver.NUMBUTTONS]bool
	Behaviour ElevatorBehaviour
	Elev_ID	  int64
	//doorOpenDuration_s float
}*/

/*type Elevators_online struct {
	Status			[]Elev_control.Elevator
	All_btn_calls 	[Driver.NUMFLOORS][Driver.NUMBUTTONS-1]bool
}*/


var Elevators_online 	[]Elev_control.Elevator
var All_btn_calls 		[Driver.NUMFLOORS][Driver.NUMBUTTONS-1]bool


func update_Elevators_online(curr_elev Elev_control.Elevator){
	for i,elev := range Elevators_online{
		if elev.Elev_ID == curr_elev.Elev_ID{
			Elevators_online[i] = curr_elev
			return
		}
	}
	Elevators_online = append(Elevators_online, curr_elev)
}

func delete_All_elevs(){
	Elevators_online  = Elevators_online[:1]
	
	//Elevators_online = Elevators_online
	//Funker dette?
}

//All_btn_calls må oppdateres når en heis stopper i en etasje.
//Fjerne alle btn_calls som er i samme retning som den heisen.
//Eventuelt fjerne alle btn_calls i den etasjen dersom heisen har Dir D_Idle