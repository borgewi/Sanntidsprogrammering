package Master_Slave

import (
	"Elev_control"
)

//skal vi legge til current_order i elevator struct? Nei, tenker jeg.
//Returns index of the elevator that should handle the specific btn_call. Elementer i Elevators_online må ikke ha endret seg!
//Dersom den returnerer -1: vent et halvt sekund og prøv på nytt. Eventuelt gi oppdrag til første heis som får Dir D_Idle
func cost_function(btn_floor int, btn_type Elev_control.Button, elevs_online []Elev_control.Elevator) int {
	i_best_time := -1
	best_time := 100
	var time_to_handle int
	var floors_between int

	for i, elev := range elevs_online {
		if elev.Error == true {
			continue
		}
		floors_between = 0
		time_to_handle = 0
		if elev.Floor == btn_floor {
			if elev.Behaviour == Elev_control.EB_Idle || elev.Behaviour == Elev_control.EB_DoorOpen {
				return i
			}
		}
		switch elev.Dir {
		case Elev_control.D_Down:
			if btn_type == Elev_control.B_HallUp {
				floors_between += elev.Floor + btn_floor
			} else { //B_HallDown
				if elev.Floor <= btn_floor { //Vanskelig å regne ut tid, siden så langt unna
					floors_between = 10
					break
				}
				floors_between += elev.Floor - btn_floor
			}
		case Elev_control.D_Idle:
			if elev.Floor == btn_floor {
				return i
			} else if elev.Floor > btn_floor {
				floors_between += elev.Floor - btn_floor
			} else { //elev.Floor < btn_floor
				floors_between += btn_floor - elev.Floor
			}
		case Elev_control.D_Up:
			if btn_type == Elev_control.B_HallDown {
				floors_between += 6 - elev.Floor - btn_floor
			} else {
				if elev.Floor >= btn_floor { //Vanskelig å regne ut tid, siden så langt unna
					floors_between = 10
					break
				}
				floors_between += btn_floor - elev.Floor
			}
		}

		time_to_handle = calculate_time(floors_between)
		if time_to_handle < best_time {
			i_best_time = i
			best_time = time_to_handle
		}
	}
	//Få med elev.Behaviour

	return i_best_time
}

func calculate_time(floors_between int) int {
	time_between_floors := 1
	door_open_time := 1
	return floors_between*(time_between_floors+door_open_time) - door_open_time
}

//Ha en funksjon som kjører cost_function på alle nåværende btn_calls dersom en heis
//får en error (antas som død) og Elevs_online oppdateres.
//Normalt kjøres bare cost_function når en ny bestilling kommer inn.
