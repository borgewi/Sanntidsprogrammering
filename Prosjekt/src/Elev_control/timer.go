package Elev_algo

import (
	"time"
)


var timerEndTime time.Time
var timerActive bool

func timer_start(duration time.Duration){
    timerEndTime    = time.Now().Add(duration)
    timerActive     = true
}

func timer_stop(){
    timerActive = false
}

func timer_timedOut() bool{
    return (timerActive  &&  time.Now().After(timerEndTime))
}