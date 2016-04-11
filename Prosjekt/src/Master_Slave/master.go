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
	msgToNetwork := make(chan Network.UdpMessage)
	msgFromNetwork := make(chan Network.UdpMessage)
	isMasterCh := make(chan bool)
	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	//test_send_order(sendCh)
	go Network.Get_status_and_broadcast(msgToNetwork, localStatusCh)




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


func test_send_order(sendCh chan Network.UdpMessage){
	var new_order [2]int
	new_order[0] = 2
	new_order[1] = 1
	order_ID := All_elevs.Status[1].Elev_ID
	Network.Master_send_order(order_ID, new_order, sendCh)
}

func print_All_elevs_status(){
	fmt.Println("\n\n\n")
	fmt.Printf("%+v",All_elevs)
	fmt.Println("\n\n\n")
}