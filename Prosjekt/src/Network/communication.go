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

func Receive_msg(receiveCh chan UdpMessage) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	for {
		msg = <-receiveCh
		elev = msg.Data
		switch(msg.Order_ID){
		case 0:
			Elev_control.PrintElev(elev)
			fmt.Println("")
			break
		default:
			Elev_control.Fsm_addOrder(msg.Order, msg.Order_ID)
			fmt.Println("Ordre mottas til: ",msg.Order_ID)
		}
		//fmt.Println("Mottar udp melding")
	}
}

func Master_send_order(ID int64, new_order[2] int, sendCh chan UdpMessage){
	//new_order = Elev_control.Elevator.Requests
	var msg UdpMessage
	msg.Order_ID = ID
	msg.Order = new_order
	sendCh <- msg
}