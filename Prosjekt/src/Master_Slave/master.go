package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
)

const masterPort = 13337
const slavePort = 13338

//const extern_Addr = "129.241.187.255" + ":13337"

func Master(statusCh chan Elev_control.Elevator) { // broadcaste ordrer til alle på bestemt port. Motta fra slaver på egen IP-addr på bestemt port
	sendCh := make(chan Network.UdpMessage)
	receiveCh := make(chan Network.UdpMessage)

	Network.UdpTransmitInit(masterPort, masterPort, sendCh)
	Network.UdpReceiveInit(slavePort, slavePort, receiveCh)

	go Network.Receive_status(receiveCh)
	go Network.Get_status_and_broadcast(sendCh, statusCh)
	for {
		if checkTimeout(receiveCh) {
			Network.UdpReceiveInit(slavePort, slavePort, receiveCh)
			fmt.Println("timeout")
		} else {
			break
		}
	}
	fmt.Println("Kommet ut")
	//start backup
}

func Slave(statusCh chan Elev_control.Elevator, receiveMasterCh chan Network.UdpMessage) { //motta ordrer på bestemt port. Sende til master på masters IP-addr og bestemt port.
	sendCh := make(chan Network.UdpMessage)

	Network.UdpTransmitInit(slavePort, slavePort, sendCh)
	go Network.Receive_status(receiveMasterCh)
	Network.Get_status_and_broadcast(sendCh, statusCh)
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

func checkTimeout(receiveCh chan Network.UdpMessage) bool {
	msg := <-receiveCh
	if msg.Length == -1 {
		return true
	}
	return false
}
