package Network

import (
	"Elev_control"
	"fmt"
)

func Get_status_and_broadcast(sendCh chan UdpMessage, statusCh chan Elev_control.Elevator) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	msg.Raddr = "broadcast"
	for {
		//fmt.Println("Prøver å motta status")
		elev = <-statusCh
		//fmt.Printf("%+v", elev)
		//fmt.Println("Sender status")
		msg.Data = elev
		sendCh <- msg
	}
}

func Receive_status(receiveCh chan UdpMessage) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	for {
		msg = <-receiveCh
		elev = msg.Data
		//fmt.Println("Mottar udp melding")
		Elev_control.PrintElev(elev)
		fmt.Println("")
	}
}
