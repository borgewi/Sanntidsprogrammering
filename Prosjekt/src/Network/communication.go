package Network
import (
	"Elev_control"
	"fmt"
)



func Get_status_and_broadcast(receiveCh chan Elev_control.Elevator) {
	//var data []byte
	var elev Elev_control.Elevator
	toAdress := "129.241.187.143" + ":13337"
	for {
		elev = <-receiveCh
		UdpBroadcast(toAdress, elev)

		//fmt.Printf("%+v", elev)
	}
}

func Receive_status(receiveCh chan Elev_control.Elevator, timeoutCh chan int) {
	listenAddr := "129.241.187.143" + ":13337"
	go ListenToBroadcast(listenAddr, receiveCh, timeoutCh)
	var elev Elev_control.Elevator
	for {
		elev = <-receiveCh
		fmt.Printf("%+v", elev)
	}
}
