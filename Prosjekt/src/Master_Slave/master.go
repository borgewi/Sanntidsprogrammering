package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	//"time"
)

const masterPort = 47838
const slavePort = 84620

//const extern_Addr = "129.241.187.255" + ":13337"

func Princess(localStatusCh chan Elev_control.Elevator){
	master_elev := <-localStatusCh
	add_elev_to_Elevators_online(master_elev)
	msgToNetwork := make(chan Network.UdpMessage)
	msgFromNetwork := make(chan Network.UdpMessage)
	isMasterCh := make(chan bool)
	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	//test_send_order(msgToNetwork)
	go Network.MH_Get_status_and_broadcast(msgToNetwork, localStatusCh)
	go Network.MH_HandleIncomingMsg(msgFromNetwork)
	var isMaster bool
	for{
		isMaster = <- isMasterCh
		if isMaster{
			fmt.Println("                        Er master")
			test_send_order(msgToNetwork)
		} else{
			fmt.Println("                        Er slave")
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


func test_send_order(msgToNetwork chan Network.UdpMessage){
	var new_order [2]int
	new_order[0] = 2
	new_order[1] = 1
	order_ID := All_elevs.Status[0].Elev_ID
	Network.MH_Master_send_order(order_ID, new_order, msgToNetwork)
}

func print_All_elevs_status(){
	fmt.Println("\n\n\n")
	fmt.Printf("%+v",All_elevs)
	fmt.Println("\n\n\n")
}