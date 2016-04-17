package Network

import (
	"Elev_control"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	//"time"
)

const MSGsize = 512
const masterPort = 47838
const slavePort = 47839
const alive_port = "47837"

//var lAddr *net.UDPAddr //Local address
//var bAddr *net.UDPAddr //Broadcast address

type UdpMessage struct {
	Data      Elev_control.Elevator
	Length    int //length of received data, in #bytes // N/A for sending
	Order_ID  int64
	Order     [2]int
	Btn_calls [4][2]bool
}

func Init_udp(msgToNetwork, msgFromNetwork chan UdpMessage, isMasterCh chan bool) {
	my_IP := GetLocalIP()
	fmt.Println("Lokal ip_addresse: \n", my_IP)
	peerListLocalCh := make(chan []string)
	go udpSendAlive(alive_port)
	go udpRecvAlive(alive_port, peerListLocalCh)

	msgToNetwork_master := make(chan UdpMessage, 100)
	msgToNetwork_slave := make(chan UdpMessage, 100)
	msgFromNetwork_master := make(chan UdpMessage, 100)
	msgFromNetwork_slave := make(chan UdpMessage, 100)
	go transmitMsg(masterPort, msgToNetwork_master)
	go transmitMsg(slavePort, msgToNetwork_slave)
	go receiveMsg(masterPort, msgFromNetwork_master)
	go receiveMsg(slavePort, msgFromNetwork_slave)

	go func() {
		isMaster := false
		for {
			select {
			case msg := <-msgFromNetwork_slave:
				//fmt.Println("case: msgFromNetwork_slave")
				if isMaster {
					msgFromNetwork <- msg
				}
			case msg := <-msgFromNetwork_master:
				//fmt.Println("case: msgFromNetwork_master")
				if !isMaster {
					msgFromNetwork <- msg
				}
			case msg := <-msgToNetwork:
				//fmt.Println("case: msgToNetwork")
				if isMaster {
					msgToNetwork_master <- msg
				} else {
					msgToNetwork_slave <- msg
				}
			case new_peer_list := <-peerListLocalCh:
				//fmt.Println("case: peerListLocalCh")
				//fmt.Println(new_peer_list)
				highest_IP := my_IP
				for _, IP := range new_peer_list {
					if highest_IP < IP {
						highest_IP = IP
					}
				}
				if my_IP == highest_IP {
					isMaster = true
					isMasterCh <- true
				} else {
					isMaster = false
					isMasterCh <- false
				}
				//send {isMaster, peers.filter(ip)} to user
			}
		}
	}()
}

func receiveMsg(port int, messageFromNetwork chan UdpMessage) {
	bAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(port))
	broadcastConn, err := net.ListenUDP("udp4", bAddr)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_receive_server: \nClosing connection.", r)
			broadcastConn.Close()
		}
	}()

	if err != nil {
		broadcastConn.Close()
	}

	for {
		buf := make([]byte, MSGsize)
		n, _, err := broadcastConn.ReadFromUDP(buf)
		if err != nil {
			broadcastConn.Close()
		}
		buf = buf[:n]
		var msg UdpMessage
		DecodeMessage(&msg, buf)
		messageFromNetwork <- msg
	}
}

func transmitMsg(port int, messageToNetwork chan UdpMessage) {
	bAddr, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(port))
	broadcastConn, _ := net.DialUDP("udp4", nil, bAddr)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udpTransmitServer: %s \n Closing connection.", r)
			broadcastConn.Close()
		}
	}()
	for {
		msg := <-messageToNetwork
		broadcastConn.Write(EncodeMessage(msg))
	}
}

func EncodeMessage(msg UdpMessage) []byte {
	returnMessage, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: problems encoding UdpMessage struct.\n")
		panic(err)
	}
	return returnMessage
}

func DecodeMessage(msg *UdpMessage, arr []byte) {
	json.Unmarshal(arr, msg)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

/*
channels:
	if we are master (struct { isMaster bool, slaves []addr} )
	message to network
	message from network

func Init(msgToNetwork chan UdpMessage, msgFromNetwork chan UdpMessage, isMaster chan ??){



	msgToNetwork_master := make(chan UdpMessage)
	msgToNetwork_slave := make(chan UdpMessage)
	go transmitMsg(masterPort, msgToNetwork_master)
	go transmitMsg(slavePort, msgToNetwork_slave)


	go func(){
		isMaster := false
		for {
			select {
			msg from slaves:
				if master
					send to user
				else
					ignore

			msg from master
				if master
					ignore (should not happen, can only have one master anyway)
				else
					send to user

			msg := <-msgToNetwork
				if isMaster {
					msgToNetwork_master <- msg
				} else {
					msgToNetwork_slave <- msg
				}

			new list of peers
				if ip = highest in peers
					become master
				else
					become slave
				send {isMaster, peers.filter(ip)} to user
			}
		}
	}()
}

func receiveMsg(port string, messageFromNetwork chan UdpMessage){
	set up conn for port
	for {
		receive
		decode
		shove on channel
	}
}


func transmitMsg(port string, messageToNetwork chan UdpMessage){
	set up conn for port
	for {
		read from chan
		encode
		send/write to network
	}
}



*/

/*
func InitializeConnection(localPort, broadcastPort int) (lConn, bConn *net.UDPConn, lAddr, bAddr *net.UDPAddr) {

	bAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+strconv.Itoa(broadcastPort))

	tempConn, err := net.DialUDP("udp4", nil, bAddr)
	defer tempConn.Close()

	tempAddr := tempConn.LocalAddr()
	lAddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
	lAddr.Port = localPort

	localConn, err := net.ListenUDP("udp4", lAddr)
	if err != nil {
		fmt.Println(err)
	}

	broadcastConn, err := net.ListenUDP("udp", bAddr)

	if err != nil {
		localConn.Close()
	}

	return localConn, broadcastConn, lAddr, bAddr
}

func UdpTransmitInit(localSendPort, broadcastSendPort int, sendCh chan UdpMessage) {
	lConn, bConn, _, bAddr := InitializeConnection(localSendPort, broadcastSendPort)
	go udpTransmitServer(lConn, bConn, sendCh, bAddr)
}

func UdpReceiveInit(localListenPort, broadcastListenPort int, receiveCh chan UdpMessage) {

	lConn, bConn, lAddr, _ := InitializeConnection(localListenPort, broadcastListenPort)
	go udpReceiveServer(lConn, bConn, receiveCh, lAddr)
}

func udpTransmitServer(lConn, bConn *net.UDPConn, sendCh chan UdpMessage, bAddr *net.UDPAddr) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udpTransmitServer: %s \n Closing connection.", r)
			lConn.Close()
			bConn.Close()
		}
	}()

	var err error
	var n int

	for {
		msg := <-sendCh
		if msg.Raddr == "broadcast" {
			n, err = lConn.WriteToUDP(EncodeMessage(msg.Data), bAddr)

		} else {
			rAddr, err := net.ResolveUDPAddr("udp", msg.Raddr)
			if err != nil {
				fmt.Println("Error: udpTransmitServer: could not resolve raddr\n")
				panic(err)
			}

			n, err = lConn.WriteToUDP(EncodeMessage(msg.Data), rAddr)
		}
		if err != nil || n < 0 {
			fmt.Println("Error: udp_transmit_server: writing\n")
			panic(err)
		}
	}
}

func udpReceiveServer(lConn, bConn *net.UDPConn, receiveCh chan UdpMessage, lAddr *net.UDPAddr) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_receive_server: \nClosing connection.", r)
			lConn.Close()
			bConn.Close()
		}
	}()

	bConnRcvCh := make(chan UdpMessage)
	lConnRcvCh := make(chan UdpMessage)
	timeoutCh := make(chan UdpMessage)

	go udpConnectionReader(lConn, lConnRcvCh, timeoutCh)
	go udpConnectionReader(bConn, bConnRcvCh, timeoutCh)

	for {
		select {
		case buf := <-bConnRcvCh:
			receiveCh <- buf

		case buf := <-lConnRcvCh:
			receiveCh <- buf

		case buf := <-timeoutCh:
			receiveCh <- buf

		}
	}
}

func udpConnectionReader(conn *net.UDPConn, rcvCh, timeoutCh chan UdpMessage) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udpConnectionReader: \nClosing connection.", r)
			conn.Close()
		}
	}()

	for {
		buf := make([]byte, MSGsize)

		conn.SetDeadline(time.Now().Add(2000 * time.Millisecond))
		n, rAddr, err := conn.ReadFromUDP(buf)

		buf = buf[:n]

		//if err != nil || n < 0 { // make timeout specific
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				timeoutCh <- UdpMessage{Length: -1}
				panic(err)
			}
		}

		var TempData Elev_control.Elevator
		DecodeMessage(&TempData, buf)

		rcvCh <- UdpMessage{Raddr: rAddr.String(), Data: TempData, Length: n}
	}
}

*/
