package elev_algo

include(
    "requests"
    "timer"
    "Driver"
)


// #include "con_load.h"
// #include "elevator_io_device.h"


static Elevator             elevator;
static ElevOutputDevice     outputDevice;


// static void __attribute__((constructor)) fsm_init(){
//     elevator = elevator_uninitialized();
    
//     con_load("elevator.con",
//         con_val("doorOpenDuration_s", &elevator.config.doorOpenDuration_s, "%lf")
//         con_enum("clearRequestVariant", &elevator.config.clearRequestVariant,
//             con_match(CV_All)
//             con_match(CV_InDirn)
//         )
//     )
    
//     outputDevice = elevio_getOutputDevice();
// }

func setAllLights(e Elevator){
    for(floor := 0; floor < N_FLOORS; floor++){
        for(btn := 0; btn < N_BUTTONS; btn++){
            outputDevice.requestButtonLight(floor, btn, e.Requests[floor][btn])
        }
    }
}

func fsm_onInitBetweenFloors(){
    outputDevice.motorDirection = D_Down
    elevator.Dir = D_Down
    elevator.Behaviour = EB_Moving
}

func fsm_onRequestButtonPress(btn_floor int, btn_type Button){
    switch(elevator.Behaviour){
    case EB_DoorOpen:
        if(elevator.Floor == btn_floor){
            timer_start(elevator.config.doorOpenDuration_s)
        } else {
            elevator.Requests[btn_floor][btn_type] = true
        }
        break
    case EB_Moving:
        elevator.Requests[btn_floor][btn_type] = true
        break
    case EB_Idle:
        elevator.Requests[btn_floor][btn_type] = true
        elevator.Dir = requests_chooseDirection(elevator)
        if(elevator.Dir == D_Idle){
            outputDevice.doorLight(true)
            elevator = requests_clearAtCurrentFloor(elevator)
            timer_start(elevator.config.doorOpenDuration_s)
            elevator.Behaviour = EB_DoorOpen
        } else {
            outputDevice.motorDirection(elevator.Dir)
            elevator.Behaviour = EB_Moving
        }    
        break
    }
    
    setAllLights(elevator)
}


func fsm_onFloorArrival(newFloor int){
    elevator.Floor = newFloor
    outputDevice.floorIndicator(elevator.Floor)
    switch(elevator.Behaviour){
    case EB_Moving:
        if(requests_shouldStop(elevator)){
            outputDevice.motorDirection(D_Idle)
            outputDevice.doorLight(true)
            elevator = requests_clearAtCurrentFloor(elevator)
            timer_start(elevator.config.doorOpenDuration_s)
            setAllLights(elevator)
            elevator.Behaviour = EB_DoorOpen
        }
        break
    default:
        break
    }
}


func fsm_onDoorTimeout(){
    switch(elevator.Behaviour){
    case EB_DoorOpen:
        elevator.Dir = requests_chooseDirection(elevator)
        outputDevice.doorLight(true)
        outputDevice.motorDirection(elevator.Dir)
        if(elevator.Dir == D_Idle){
            elevator.Behaviour = EB_Idle
        } else {
            elevator.Behaviour = EB_Moving
        }
        break
    default:
        break
    }
}