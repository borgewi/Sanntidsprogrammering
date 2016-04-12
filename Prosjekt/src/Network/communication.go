package Network

import (
	"Elev_control"
	"fmt"
)

func Get_status_and_broadcast(msgToNetwork chan UdpMessage, localStatusCh chan Elev_control.Elevator) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	msg.Raddr = "broadcast"
	for {
		fmt.Println("Prøver å motta status")
		elev = <-localStatusCh
		fmt.Println("Status mottatt fra Elev_control")
		//fmt.Printf("%+v", elev)
		//fmt.Println("Sender status")
		msg.Data = elev
		msgToNetwork <- msg
		fmt.Printf("%+v",msg)
		fmt.Println("")
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
			fmt.Println("Ordre mottas til: ",msg.Order_ID)
			Elev_control.Fsm_addOrder(msg.Order, msg.Order_ID)
		}
		//fmt.Println("Mottar udp melding")
	}
}

func Master_send_order(ID int64, new_order[2] int, msgToNetwork chan UdpMessage){
	//new_order = Elev_control.Elevator.Requests
	var msg UdpMessage
	//var elev Elev_control.Elevator
	msg.Raddr = "broadcast"
	msg.Order_ID = ID
	msg.Order = new_order
	//msg.Data = elev
	msgToNetwork <- msg
	fmt.Println("ssssssssssssssssssssssssssssssssssssssssssssssssssssssssss")
}