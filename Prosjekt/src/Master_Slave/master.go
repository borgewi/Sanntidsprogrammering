package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	//"sync"
	"time"
)

var isMaster bool

func Princess(localStatusCh chan Elev_control.Elevator, sendBtnCallCh chan [2]int, receiveAllBtnCallsCh, setLights_setExtBtnsCh chan [4][2]bool, errorCh chan int) {
	isMaster = true
	Network.MH_UpdateMasterStatus(isMaster)
	master_elev := <-localStatusCh
	update_Elevators_online(master_elev)
	msgToNetwork := make(chan Network.UdpMessage, 100)
	msgFromNetwork := make(chan Network.UdpMessage, 100)
	updateElevsCh := make(chan Elev_control.Elevator, 100)
	isMasterCh := make(chan bool)
	sendOrderCh := make(chan Network.UdpMessage, 100)
	receiveBtnCallCh := make(chan [2]int, 100)
	handleOrderAgainCh := make(chan [2]int, 100)
	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	go Network.MH_HandleIncomingMsg(msgFromNetwork, updateElevsCh, receiveBtnCallCh, receiveAllBtnCallsCh)
	go Network.MH_HandleOutgoingMsg(msgToNetwork, updateElevsCh, sendOrderCh, localStatusCh, sendBtnCallCh, receiveBtnCallCh)
	var call [2]int
	for {
		select {
		case call = <-receiveBtnCallCh:
			fmt.Println("Kommer til receiveBtnCallCh")
			if update_btnCalls(call) {
				fmt.Println("UPDATEBTNCALLS ")
				temp_Elevators_online := getElevators_Online()
				index_elev := cost_function(call[0], Elev_control.Button(call[1]), temp_Elevators_online)
				if index_elev == -1 {
					fmt.Println("Sender inn 2 i errorCh")
					errorCh <- Elev_control.ERR_NO_ELEVS_OPERABLE
					break
				}
				elev_ID := temp_Elevators_online[index_elev].Elev_ID
				fmt.Printf("ElevID fra receiveBtnCallCh: ")
				fmt.Printf("%+v", elev_ID)
				Network.MH_send_new_order(elev_ID, call, sendOrderCh)
				fmt.Println("Ferdig med updatebtncalls")
			}
		case call = <-handleOrderAgainCh:
			fmt.Println("Kommer til handleOrderAgainCh")
			temp_Elevators_online := getElevators_Online()
			index_elev := cost_function(call[0], Elev_control.Button(call[1]), temp_Elevators_online)
			if index_elev == -1 {
				errorCh <- Elev_control.ERR_NO_ELEVS_OPERABLE
				break
			}
			elev_ID := temp_Elevators_online[index_elev].Elev_ID
			fmt.Printf("%+v", elev_ID)
			Network.MH_send_new_order(elev_ID, call, sendOrderCh)
		case allCalls := <-receiveAllBtnCallsCh:
			fmt.Println("Kommer inn i receiveAllBtnCallsCh")
			setAll_btn_calls(allCalls)
			fmt.Printf("%+v", allCalls)
			setLights_setExtBtnsCh <- allCalls
		case elev := <-updateElevsCh:
			update_All_elevs(elev)
			//fmt.Println("HAR LAGT TIL NY ELEVATOR!!!", elev, "\n")
		case isMaster = <-isMasterCh:
			delete_All_elevs()
			if isMaster {
				runCost_AllUnfinishedOrders(handleOrderAgainCh)
				fmt.Println("                        				Er master")

			} else {
				fmt.Println("                        				Er slave")
			}
			Network.MH_UpdateMasterStatus(isMaster)

		}

		check_elevsIdleAtFloor()
		temp_All_btn_calls := get_All_btn_calls()
		if isMaster {
			//fmt.Println("Kommer før MH_broadcast_all_btn_calls")
			Network.MH_broadcast_all_btn_calls(temp_All_btn_calls, sendOrderCh)
			//fmt.Println("Er det analen?")
			setLights_setExtBtnsCh <- temp_All_btn_calls
			//fmt.Println("Sendt til setLights_setExtBtnsCh")
			checkTimeStamps(handleOrderAgainCh)
		}
	}
}

func runCost_AllUnfinishedOrders(handleOrderAgainCh chan [2]int) {
	var oldOrder [2]int
	time.Sleep(1 * time.Second)
	for i, k := range all_btn_calls {
		for j, call := range k {
			//fmt.Println(i, j, call)
			if call {
				oldOrder[0] = i
				oldOrder[1] = j
				handleOrderAgainCh <- oldOrder
			}
		}
	}
}
