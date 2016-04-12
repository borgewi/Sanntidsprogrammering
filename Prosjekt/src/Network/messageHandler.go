package Network

import (
	"Elev_control"
	"fmt"
)

var isMaster bool

func MH_HandleOutgoingMsg(msgToNetwork chan UdpMessage, localStatusCh, updateElevsCh chan Elev_control.Elevator, sendOrderCh chan [2]int, sendBtnCallsCh chan [4][2]bool) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	for {

		if isMaster{
			select{
			case elev = <-localStatusCh:
				updateElevsCh <- elev
			case order := <- sendOrderCh:
				msg.Order = order 
				msgToNetwork <- msg
			case btn_calls := <- sendBtnCallsCh:
				msg.Btn_calls = btn_calls
				msgToNetwork <- msg
			}
		} else{
			select{
				case elev = <-localStatusCh:
					msg.Data = elev
					msgToNetwork <- msg		
				case btn_calls := <- sendBtnCallsCh:
					msg.Btn_calls = btn_calls
					msgToNetwork <- msg
			}
		}
	}
}

func MH_HandleIncomingMsg(msgFromNetwork chan UdpMessage, updateElevsCh chan Elev_control.Elevator) {
	var elev Elev_control.Elevator
	var msg UdpMessage
	for {
		msg = <-msgFromNetwork
		elev = msg.Data
		switch(msg.Order_ID){
		case 0:
			updateElevsCh <- msg.Data
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

func MH_Master_send_order(ID int64, new_order[2] int, msgToNetwork chan UdpMessage){
	//new_order = Elev_control.Elevator.Requests
	var msg UdpMessage
	//var elev Elev_control.Elevator
	msg.Order_ID = ID
	msg.Order = new_order
	//msg.Data = elev
	msgToNetwork <- msg
	fmt.Println("Order sent")
}

func UpdateMasterStatus(isMasterFrom_Master_Slave  bool){
	isMaster = isMasterFrom_Master_Slave
}