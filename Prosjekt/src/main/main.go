package main


import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"time"
	"Driver"
)


func main() {
	Driver.ElevInit()
	error := 0
	for error < 1{
		check_buttons()
	}
}

func button_pressed(floor int,elev_button_type_t button) {
	Driver.ElevSetMotorDirection(1)
}

func check_buttons(){
    for sfloor := 0;sfloor < 3;sfloor ++ {
        if(Driver.elev_get_button_signal(BUTTON_CALL_UP, sfloor) == 1){
            button_pressed(sfloor,BUTTON_CALL_UP);
        }
    }
    for sfloor := 1;sfloor < 4;sfloor ++ {
        if(Driver.elev_get_button_signal(BUTTON_CALL_DOWN, sfloor) == 1){
            button_pressed(sfloor,BUTTON_CALL_DOWN);
        }
    }
    for sfloor := 0;sfloor < 4;sfloor ++ {
        if(Driver.elev_get_button_signal(BUTTON_COMMAND, sfloor) == 1){
            button_pressed(sfloor,BUTTON_COMMAND);
        }
    }
}
