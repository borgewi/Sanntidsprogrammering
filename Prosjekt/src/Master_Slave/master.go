package Master_Slave

import (
	"Elev_control"
	"Network"
	"fmt"
	//"sync"
	"time"
)

var isMaster bool

//const extern_Addr = "129.241.187.255" + ":13337"

func Princess(localStatusCh chan Elev_control.Elevator, sendBtnCallCh chan [2]int, receiveAllBtnCallsCh, setLights_setExtBtnsCh chan [4][2]bool, errorCh chan int) {
	fmt.Println("Comment")
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
	go Network.MH_HandleOutgoingMsg(msgToNetwork, sendOrderCh, localStatusCh, updateElevsCh, sendBtnCallCh, receiveBtnCallCh)
	go checkForError(errorCh)
	fmt.Println("Kommer hit")
	wasMaster := false

	go func() {
		for {
			isMaster = <-isMasterCh
		}
	}()
	var call [2]int
	isMaster = true
	Network.MH_UpdateMasterStatus(isMaster)
	for {
		select {
		case call = <-receiveBtnCallCh:
			fmt.Println("Kommer til receiveBtnCallCh")
			if update_btnCalls(call) {
				temp_Elevators_online := getElevators_Online()
				index_elev := cost_function(call[0], Elev_control.Button(call[1]), temp_Elevators_online)
				elev_ID := temp_Elevators_online[index_elev].Elev_ID
				fmt.Printf("%+v", elev_ID)
				Network.MH_send_new_order(elev_ID, call, sendOrderCh)
				fmt.Println("Ferdig med updatebtncalls")
			}
		case call = <-handleOrderAgainCh:
			fmt.Println("Kommer handleOrderAgainCh")
			temp_Elevators_online := getElevators_Online()
			index_elev := cost_function(call[0], Elev_control.Button(call[1]), temp_Elevators_online)
			elev_ID := temp_Elevators_online[index_elev].Elev_ID
			fmt.Printf("%+v", elev_ID)
			Network.MH_send_new_order(elev_ID, call, sendOrderCh)
		case allCalls := <-receiveAllBtnCallsCh:
			fmt.Println("Kommer receiveAllBtnCallsCh")
			setAll_btn_calls(allCalls)
			setLights_setExtBtnsCh <- allCalls
		case elev := <-updateElevsCh:
			update_All_elevs(elev)
			fmt.Println("Kommer updateElevsCh")
		}
		check_elevsIdleAtFloor()
		temp_All_btn_calls := get_All_btn_calls()
		if isMaster {
			fmt.Println("Kommer før MH_broadcast_all_btn_calls")
			Network.MH_broadcast_all_btn_calls(temp_All_btn_calls, sendOrderCh)
			fmt.Println("Er det analen?")
			setLights_setExtBtnsCh <- temp_All_btn_calls
			fmt.Println("Sendt til setLights_setExtBtnsCh")
			checkTimeStamps(handleOrderAgainCh)
		}
		if wasMaster && !isMaster { //Blir Slave
			delete_All_elevs()
			fmt.Println("                        Er slave")
			Network.MH_UpdateMasterStatus(isMaster)
			wasMaster = false
		} else if !wasMaster && isMaster { //Blir Master
			delete_All_elevs()
			fmt.Println("                        Er master")
			Network.MH_UpdateMasterStatus(isMaster)
			runCost_AllUnfinishedOrders(handleOrderAgainCh)
			wasMaster = true
		}
	}
}

func update_All_elevs(elev Elev_control.Elevator) {
	print_All_elevs_status()
	update_Elevators_online(elev)
}

//receiveAllBtnCallsCh <- All_btn_calls
//, receiveAllBtnCallsCh chan [4][2]bool

//mottar knappetrykk og kjører kostfunk på dem
/*func runCostfunctionOnBtnCalls(receiveBtnCallCh, handleOrderAgainCh chan [2]int, sendOrderCh chan Network.UdpMessage) {
	//go checkTimeStamps(handleOrderAgainCh)
	var oldCall bool
	var call [2]int
	if isMaster {
		oldCall = false
		select {
		case call = <-receiveBtnCallCh:
			break
		case call = <-handleOrderAgainCh:
			oldCall = true
			fmt.Println("oldCall: ", oldCall)
		}
		//Fuckit.Lock()
		if update_btnCalls(call) || oldCall { //Hvis det er en ny ordre
			elevs_online := Elevators_online
			index_elev := cost_function(call[0], Elev_control.Button(call[1]), elevs_online)
			elev_ID := Elevators_online[index_elev].Elev_ID
			fmt.Printf("%+v", elev_ID)
			Network.MH_send_new_order(elev_ID, call, sendOrderCh)
		}
		//Fuckit.Unlock()
	}
}*/

func runCost_AllUnfinishedOrders(handleOrderAgainCh chan [2]int) {
	var oldOrder [2]int
	//Fuckit.Lock()
	//defer Fuckit.Unlock()
	for i, k := range all_btn_calls {
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
	fmt.Println("\n")
	fmt.Printf("%+v", elevators_online)
	fmt.Println("\n")
}
