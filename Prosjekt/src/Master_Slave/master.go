package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	"time"
)

const masterPort = 13337
const slavePort = 13338

//const extern_Addr = "129.241.187.255" + ":13337"

func Master(statusCh chan Elev_control.Elevator){
	sendCh := make(chan Network.UdpMessage)
	receiveCh := make(chan Network.UdpMessage)

	Network.UdpTransmitInit(masterPort, masterPort, sendCh)
	Network.UdpReceiveInit(slavePort, slavePort, receiveCh)

	wait_for_slave(statusCh, receiveCh, sendCh)
	test_send_order(sendCh)
	//start backup
}

func Slave(statusCh chan Elev_control.Elevator, receiveMasterCh chan Network.UdpMessage) { //motta ordrer på bestemt port. Sende til master på masters IP-addr og bestemt port.
	sendCh := make(chan Network.UdpMessage)

	Network.UdpTransmitInit(slavePort, slavePort, sendCh)
	go Network.Receive_msg(receiveMasterCh)
	go Network.Get_status_and_broadcast(sendCh, statusCh)
}

func Determine_Rank(statusCh chan Elev_control.Elevator) {
	receiveMasterCh := make(chan Network.UdpMessage)

	Network.UdpReceiveInit(masterPort, masterPort, receiveMasterCh)

	if checkTimeout(receiveMasterCh) {
		fmt.Println("Starting master")
		go Master(statusCh)
		//go Backup
	} else {
		fmt.Println("Starting slave")
		go Slave(statusCh, receiveMasterCh)
	}
}

func wait_for_slave(statusCh chan Elev_control.Elevator, receiveCh, sendCh chan Network.UdpMessage) { 
	go Network.Receive_msg(receiveCh)
	go Network.Get_status_and_broadcast(sendCh, statusCh)
	for {
		if checkTimeout(receiveCh) {
			time.Sleep(1000 * time.Millisecond)
			Network.UdpReceiveInit(slavePort, slavePort, receiveCh)
			fmt.Println("timeout")
		} else {
			break
		}
	}
	fmt.Println("Slave er i live")

	//Init All_elevs
	master_elev := <-statusCh
	msg := <- receiveCh
	slave_elev := msg.Data
	add_elev_to_Elevators_online(master_elev)
	add_elev_to_Elevators_online(slave_elev)
	print_All_elevs_status()
}

func checkTimeout(receiveCh chan Network.UdpMessage) bool {
	msg := <-receiveCh
	if msg.Length == -1 {
		return true
	}
	return false
}


func test_send_order(sendCh chan Network.UdpMessage){
	var new_order [2]int
	new_order[0] = 2
	new_order[1] = 1
	order_ID := All_elevs.Status[1].Elev_ID
	Network.Master_send_order(order_ID, new_order, sendCh)
}


func print_All_elevs_status(){
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("%+v",All_elevs)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}