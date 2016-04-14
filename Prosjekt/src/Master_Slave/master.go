package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	"time"
)

var isMaster bool

//const extern_Addr = "129.241.187.255" + ":13337"

func Princess(localStatusCh chan Elev_control.Elevator, sendBtnCallCh chan [2]int, receiveAllBtnCallsCh chan [4][2]bool, errorCh chan int) {
	master_elev := <-localStatusCh
	update_Elevators_online(master_elev)
	msgToNetwork := make(chan Network.UdpMessage, 5)
	msgFromNetwork := make(chan Network.UdpMessage, 5)
	updateElevsCh := make(chan Elev_control.Elevator)
	isMasterCh := make(chan bool)
	sendOrderCh := make(chan Network.UdpMessage)
	receiveBtnCallCh := make(chan [2]int, 5)
	handleOrderAgainCh := make(chan [2]int)

	Network.Init_udp(msgToNetwork, msgFromNetwork, isMasterCh)
	go Network.MH_HandleIncomingMsg(msgFromNetwork, updateElevsCh, receiveBtnCallCh, receiveAllBtnCallsCh)
	go Network.MH_HandleOutgoingMsg(msgToNetwork, sendOrderCh, localStatusCh, updateElevsCh, sendBtnCallCh, receiveBtnCallCh)
	go update_btnCall_run_costFunction(receiveBtnCallCh, handleOrderAgainCh, sendOrderCh, receiveAllBtnCallsCh)
	go receive_allBtnCalls(receiveAllBtnCallsCh)
	go distribute_allBtnCalls_Master(sendOrderCh, receiveAllBtnCallsCh)
	go update_All_elevs(updateElevsCh)
	go checkForError(errorCh)
	go check_elevsIdleAtFloor()
	for {
		isMaster = <-isMasterCh
		delete_All_elevs()
		if isMaster {
			fmt.Println("                        Er master")
			Network.MH_UpdateMasterStatus(isMaster)
			runCost_AllUnfinishedOrders(handleOrderAgainCh)
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

func update_btnCall_run_costFunction(receiveBtnCallCh, handleOrderAgainCh chan [2]int, sendOrderCh chan Network.UdpMessage, receiveAllBtnCallsCh chan [4][2]bool) {
	go checkTimeStamps(handleOrderAgainCh)
	var oldCall bool
	var call [2]int
	for {
		if isMaster {
			oldCall = false
			select {
			case call = <-receiveBtnCallCh:
				break
			case call = <-handleOrderAgainCh:
				oldCall = true
				fmt.Println("oldCall: ", oldCall)
			}
			if update_btnCalls(call) || oldCall { //Hvis det er en ny ordre
				//OBS!!: Index verdi kan være -1. Må lage funksjonalitet for dette senere.
				index_elev := cost_function(call[0], Elev_control.Button(call[1]))
				for index_elev == -1 {
					fmt.Println("Fant ingen heiser lett tilgjengelig. Dette skal ikke skje!")
					time.Sleep(500 * time.Millisecond)
					index_elev = cost_function(call[0], Elev_control.Button(call[1]))
				}
				elev_ID := Elevators_online[index_elev].Elev_ID
				fmt.Printf("%+v", elev_ID)
				Network.MH_send_new_order(elev_ID, call, sendOrderCh)
				Network.MH_broadcast_all_btn_calls(All_btn_calls, sendOrderCh)
				receiveAllBtnCallsCh <- All_btn_calls
			}
		}
	}
}

func runCost_AllUnfinishedOrders(handleOrderAgainCh chan [2]int) {
	var oldOrder [2]int
	for i, k := range All_btn_calls {
		for j, call := range k {
			fmt.Println(i, j, call)
			if call {
				oldOrder[0] = i
				oldOrder[1] = j
				handleOrderAgainCh <- oldOrder
			}
		}
	}
}

func receive_allBtnCalls(receiveAllBtnCallsCh chan [4][2]bool) {
	for {
		All_btn_calls = <-receiveAllBtnCallsCh
		//fmt.Printf("%+v", All_btn_calls)
	}
}

func distribute_allBtnCalls_Master(sendOrderCh chan Network.UdpMessage, receiveAllBtnCallsCh chan [4][2]bool) {
	for {
		time.Sleep(500 * time.Millisecond)
		if isMaster {
			Network.MH_broadcast_all_btn_calls(All_btn_calls, sendOrderCh)
			receiveAllBtnCallsCh <- All_btn_calls
		}
	}
}

func checkForError(errorCh chan int) {
	var err int
	for {
		time.Sleep(1 * time.Second)
		err = <-errorCh
		if err == 1 {
			fmt.Println("Error has occured")
		}
		//Master: Fjern denne heisen fra Elevators_online inntil den er operatibel igjen.
		//Kjør cost_function på nytt for All_btn_calls
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
