package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	"time"
)

const masterPort = 47838
const slavePort = 84620

var isMaster bool

//const extern_Addr = "129.241.187.255" + ":13337"

func Princess(localStatusCh chan Elev_control.Elevator, sendBtnCallsCh chan [2]int) {
	master_elev := <-localStatusCh
	update_Elevators_online(master_elev)
	msgToNetwork := make(chan Network.UdpMessage)
	msgFromNetwork := make(chan Network.UdpMessage)
	updateElevsCh := make(chan Elev_control.Elevator)
	isMasterCh := make(chan bool)
	sendOrderCh := make(chan Network.UdpMessage)
	receiveBtnCallsCh := make(chan [2]int)

	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	go Network.MH_HandleIncomingMsg(msgFromNetwork, updateElevsCh, receiveBtnCallsCh)
	go Network.MH_HandleOutgoingMsg(msgToNetwork, sendOrderCh, localStatusCh, updateElevsCh, sendBtnCallsCh, receiveBtnCallsCh)
	go update_btnCall_run_costFunction(receiveBtnCallsCh, sendOrderCh)
	go update_All_elevs(updateElevsCh)
	for {
		isMaster = <-isMasterCh
		delete_All_elevs()
		if isMaster {
			fmt.Println("                        Er master")
			Network.MH_UpdateMasterStatus(isMaster)
		} else {
			fmt.Println("                        Er slave")
			Network.MH_UpdateMasterStatus(isMaster)
		}
	}
}

func update_All_elevs(updateElevsCh chan Elev_control.Elevator) {
	go print_All_elevs_status()
	for {
		elev := <-updateElevsCh
		update_Elevators_online(elev)
	}
}

func update_btnCall_run_costFunction(receiveBtnCallsCh chan [2]int, sendOrderCh chan Network.UdpMessage) {
	for {
		newCall := <-receiveBtnCallsCh
		if update_btnCalls(newCall) { //Hvis det er en ny ordre
			//OBS!!: Index verdi kan være -1. Må lage funksjonalitet for dette senere.
			index_elev := cost_function(newCall[0], Elev_control.Button(newCall[1]))
			for index_elev == -1 {
				fmt.Println("Fant ingen heiser lett tilgjengelig. Prøver på nytt")
				time.Sleep(500 * time.Millisecond)
				index_elev = cost_function(newCall[0], Elev_control.Button(newCall[1]))
			}
			elev_ID := Elevators_online[index_elev].Elev_ID
			fmt.Printf("%+v", elev_ID)
			Network.MH_send_new_order(elev_ID, newCall, sendOrderCh)
		}
	}
}

func print_All_elevs_status() {
	for {
		time.Sleep(3000 * time.Millisecond)
		fmt.Println("\n")
		fmt.Printf("%+v", Elevators_online)
		fmt.Println("\n")
	}
}
