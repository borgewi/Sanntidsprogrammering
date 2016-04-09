package Master_Slave

import (
	"Network"
	"fmt"
	"Elev_control"
)

const localAddr = "129.241.187.143" + ":13337"
//const extern_Addr = "129.241.187.146" + ":13337"

func Master(receiveCh chan Elev_control.Elevator, timeoutCh chan int){
	go Network.Receive_status(receiveCh, timeoutCh)
	for{
		select{
		case <- timeoutCh:
			go Network.ListenToBroadcast(localAddr, receiveCh, timeoutCh)
		}
	}
	//start backup
}

func Slave(receiveCh chan Elev_control.Elevator){
	Network.Get_status_and_broadcast(receiveCh)
}




func Determine_Rank(receiveCh chan Elev_control.Elevator) {
	timeoutCh := make(chan int)

	go Network.ListenToBroadcast(localAddr, receiveCh, timeoutCh)

	select {
	case <-receiveCh: //Finnes allerede en primal som kjører
		fmt.Println("Starting slave")
		go Slave(receiveCh)
	case <-timeoutCh: //Finner ikke primal, start backup og counter
		fmt.Println("Starting master")
		go Master(receiveCh, timeoutCh)
		//go Backup
	}
}