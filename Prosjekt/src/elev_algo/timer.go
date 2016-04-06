package elev_algo

import (
	"time"
)


var timerEndTime time.time
var timerActive bool

func timer_start(duration float){
    timerEndTime    = time.now().Add(duration* time.second)
    timerActive     = true
}

func timer_stop(){
    timerActive = false
}

func timer_timedOut() bool{
    return (timerActive  &&  time.now().After(timerEndTime))
}