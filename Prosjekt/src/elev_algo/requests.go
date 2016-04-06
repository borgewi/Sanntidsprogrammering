package elev_algo

include(
    "Driver"
    "time"
)

type Direction struct {
    D_Down -1 +iota
    D_Idle
    D_Up
}

type Button struct { 
    B_HallUp 0 + iota
    B_HallDown
    B_Cab
}

type ElevatorBehaviour struct {
    EB_Idle 0 + iota
    EB_DoorOpen
    EB_Moving
} 

//type ClearRequestVariant struct {
    // Assume everyone waiting for the elevator gets on the elevator, even if 
    // they will be traveling in the "wrong" direction for a while
//    CV_All 0 + iota
    
    // Assume that only those that want to travel in the current direction 
    // enter the elevator, and keep waiting outside otherwise
//    CV_InDirn
//}

type Elevator struct{
    Floor int
    Dir Direction
    Requests[N_FLOORS][N_BUTTONS] bool
    Behaviour ElevatorBehaviour
    type config struct {
        clearRequestVariant ClearRequestVariant
        doorOpenDuration_s float
    }
}


func requests_above(e Elevator) int{
    for(f := e.floor+1; f < N_FLOORS; f++){
        for(btn := 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}

func requests_below(e Elevator) int{
    for(f := 0; f < e.floor; f++){
        for(btn := 0; btn < N_BUTTONS; btn++){
            if(e.requests[f][btn]){
                return 1;
            }
        }
    }
    return 0;
}

func requests_chooseDirection(e Elevator) Direction{
    switch(e.Dir){
    case D_Up:
        if requests_above(e){
            return D_Up
        } else if requests_below(e){
            return D_Down
        } else{
            return D_Idle
        }
    case D_Down:
        if requests_above(e){
            return D_Up
        } else if requests_below(e){
            return D_Down
        } else{
            return D_Idle
        }
    case D_Idle: // there should only be one request in this case. Checking up or down first is arbitrary.
        if requests_above(e){
            return D_Up
        } else if requests_below(e){
            return D_Down
        } else{
            return D_Idle
        }
    default:
        return D_Idle;
    }
}

func requests_shouldStop(e Elevator) bool{
    switch(e.Dir){
    case D_Down:
        return e.requests[e.floor][B_HallDown] || e.requests[e.floor][B_Cab] || !requests_below(e)
    case D_Up:
        return e.requests[e.floor][B_HallUp] || e.requests[e.floor][B_Cab] || !requests_above(e)
    case D_Stop:
        return true
    default:
        return true
    }
}


func requests_clearAtCurrentFloor(e Elevator) Elevator{
    e.requests[e.floor][B_Cab] = false;
    switch(e.Dir){
    case D_Up:
        e.requests[e.floor][B_HallUp] = false
        if(!requests_above(e)){
            e.requests[e.floor][B_HallDown] = false
        }
        break
        
    case D_Down:
        e.requests[e.floor][B_HallDown] = false
        if(!requests_below(e)){
            e.requests[e.floor][B_HallUp] = false
        }
        break
        
    case D_Idle:
    default:
        e.requests[e.floor][B_HallUp] = false
        e.requests[e.floor][B_HallDown] = false
        break
    }
    
    return e;
}