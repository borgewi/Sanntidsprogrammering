package Master_Slave

import (
	"Elev_control"
)


//skal vi legge til current_order i elevator struct?
//Returns index of the elevator that should handle the specific btn_call
//Dersom den returnerer -1: vent et halvt sekund og prøv på nytt. Eventuelt gi oppdrag til første heis som får Dir D_Idle
func cost_function(btn_floor int, btn_type Elev_control.Button) int{
	num_elevs := len(Elevators_online)
	i_best_time := -1
	best_time := 100
	time_to_handle int
	floors_between int

	for i,elev := range Elevators_online{
		floors_between = 0
		time_to_handle = 0
		switch(elev.Dir){
		case D_Down:
			if btn_type == B_HallUp{
				floors_between += elev.floor + btn_floor
			} else {
				if elev.floor <= btn_floor{ //Vanskelig å regne ut tid, siden så langt unna
					floors_between = 100
					break
				}
				floors_between += elev.floor - btn_floor	
			}
		case D_Idle:
			if elev.floor == btn_floor{
				return i
			} else if elev.floor > btn_floor{
				floors_between += elev.floor - btn_floor
			} else {
				floors_between += btn_floor - elev.floor
			}
		case D_Up:
			if btn_type == B_HallDown{
				floors_between += 6 - elev.floor - btn_floor
				} else {
					if elev.floor >= btn_floor{ //Vanskelig å regne ut tid, siden så langt unna
						floors_between = 100
						break
					}
					floors_between += btn_floor - elev.floor
				}
			}
		}


		time_to_handle = calculate_time(floors_between)
		if time_to_handle < best_time{
			i_best_time = i
		}
	}
	//Få med elev.Behaviour
	
	return i_best_time
}

func calculate_time(floors_between int) int{
	time_between_floors := 1
	door_open_time := 1
	return floors_between*(time_between_floors+door_open_time)-door_open_time
}



//Ha en funksjon som kjører cost_function på alle nåværende btn_calls dersom en heis
//får en error (antas som død) og Elevs_online oppdateres.
//Normalt kjøres bare cost_function når en ny bestilling kommer inn.