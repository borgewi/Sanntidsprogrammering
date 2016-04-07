package main


import (
	//"encoding/binary"
	"fmt"
    "Elev_algo"
	//"net"
	//"os/exec"
	//"time"
	"Driver"
)


const(
    NUMBUTTONS int = 3
    NUMFLOORS int = 4
    BUTTON_DOWN = 0 + iota
    BUTTON_UP 
    BUTTON_COMMAND 
)


func main() {
    receiveCh := make(chan []byte)
    //timeoutCh := make(chan int)

	Driver.ElevInit()
    go Elev_algo.Merry_go_around(receiveCh)
    go get_status_and_broadcast(receiveCh)
    for{}
}


func get_status_and_broadcast(receiveCh chan []byte){
    var data []byte
    var elev Elev_algo.Elevator
    for{
        data = <-receiveCh
        json.Unmarshal(data, &elev)
        fmt.Printf("%+v", elev)
    }
}












func button_pressed(floor int,button int) {
	Driver.ElevSetMotorDirection(1)
    fmt.Println("drive")
}

func check_buttons(){
    for sfloor := 0;sfloor < 3;sfloor ++ {
        if(Driver.ElevGetButtonSignal(BUTTON_UP, sfloor) == 1){
            button_pressed(sfloor,BUTTON_UP);
            fmt.Println("opp")
        }
    }
    for sfloor := 1;sfloor < 4;sfloor ++ {
        if(Driver.ElevGetButtonSignal(BUTTON_DOWN, sfloor) == 1){
            button_pressed(sfloor,BUTTON_DOWN);
        }
    }
    for sfloor := 0;sfloor < 4;sfloor ++ {
        if(Driver.ElevGetButtonSignal(BUTTON_COMMAND, sfloor) == 1){
            button_pressed(sfloor,BUTTON_COMMAND);
        }
    }
}
