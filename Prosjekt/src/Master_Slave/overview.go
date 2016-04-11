package Master_Slave

import (
	"Elev_control"
	"Driver"
	//"time"
	//"fmt"
)


// func oversikt skal vite hvor alle heisene er og hvilken retning de har. 
var All_elevs Elevators_online

/*type Elevator struct {
	Floor     int
	Dir       Direction
	Requests  [Driver.NUMFLOORS][Driver.NUMBUTTONS]bool
	Behaviour ElevatorBehaviour
	Elev_ID	  int64
	//doorOpenDuration_s float
}*/

type Elevators_online struct {
	Status			[]Elev_control.Elevator
	All_btn_calls 	[Driver.NUMFLOORS][Driver.NUMBUTTONS-1]bool
}

func add_elev_to_Elevators_online(new_elev Elev_control.Elevator){
	All_elevs.Status = append(All_elevs.Status, new_elev)
}
