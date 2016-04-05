package elev_algo

import (
	"time"
)

func get_wall_time() float{
    struct timeval time
    gettimeofday(&time, NULL)
    return (double)time.tv_sec + (double)time.tv_usec * .000001
}


var timerEndTime time.time
var timerActive bool

func timer_start(duration float){
    timerEndTime    = time.now() + duration;
    timerActive     = 1;
}

func timer_stop(){
    timerActive = 0
}

func timer_timedOut() int{
    return (timerActive  &&  get_wall_time() > timerEndTime)
}



