package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	//"time"
)

const masterPort = 47838
const slavePort = 84620

var isMaster bool

//const extern_Addr = "129.241.187.255" + ":13337"

func Princess(localStatusCh chan Elev_control.Elevator){
	master_elev := <-localStatusCh
	update_Elevators_online(master_elev)
	msgToNetwork := make(chan Network.UdpMessage)
	msgFromNetwork := make(chan Network.UdpMessage)
	updateElevsCh := make(chan Elev_control.Elevator)
	isMasterCh := make(chan bool)
	sendOrderCh := make(chan [2]int)
	sendBtnCallsCh := make(chan [4][2]bool)

	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	go Network.MH_HandleIncomingMsg(msgFromNetwork, updateElevsCh)
	go Network.MH_HandleOutgoingMsg(msgToNetwork, localStatusCh, updateElevsCh, sendOrderCh, sendBtnCallsCh)
	go update_All_elevs(updateElevsCh)
	for{
		isMaster = <- isMasterCh
		delete_All_elevs()
		if isMaster{
			fmt.Println("                        Er master")
			Network.UpdateMasterStatus(isMaster)
			test_send_order(msgToNetwork)
		} else{
			fmt.Println("                        Er slave")
			Network.UpdateMasterStatus(isMaster)
		}
	}

	//Init All_elevs
	/*master_elev := <-statusCh
	msg := <- receiveCh
	slave_elev := msg.Data
	add_elev_to_Elevators_online(master_elev)
	add_elev_to_Elevators_online(slave_elev)
	print_All_elevs_status()
	Receive_msg(receiveCh chan UdpMessage)
	*/
}

func update_All_elevs(updateElevsCh chan Elev_control.Elevator){
	for{
		elev := <- updateElevsCh
		//elev := status_msg.Data
		update_Elevators_online(elev)
		print_All_elevs_status()
	}
}

func test_send_order(msgToNetwork chan Network.UdpMessage){
	var new_order [2]int
	new_order[0] = 2
	new_order[1] = 1
	order_ID := Elevators_online[0].Elev_ID
	Network.MH_Master_send_order(order_ID, new_order, msgToNetwork)
}

func print_All_elevs_status(){
	fmt.Println("\n")
	fmt.Printf("%+v",Elevators_online)
	fmt.Println("\n")
}