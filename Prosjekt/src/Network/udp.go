package Network

import (
	"Elev_control"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const MSGsize = 512

var lAddr *net.UDPAddr //Local address
var bAddr *net.UDPAddr //Broadcast address

type UdpMessage struct {
	Raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	Data   Elev_control.Elevator
	Length int //length of received data, in #bytes // N/A for sending
}

func InitializeConnection(localPort, broadcastPort int) (lConn, bConn *net.UDPConn, lAddr, bAddr *net.UDPAddr) {

	bAddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:"+strconv.Itoa(broadcastPort))

	tempConn, err := net.DialUDP("udp4", nil, bAddr)
	defer tempConn.Close()

	tempAddr := tempConn.LocalAddr()
	lAddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
	lAddr.Port = localPort

	localConn, err := net.ListenUDP("udp4", lAddr)

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

func EncodeMessage(e Elev_control.Elevator) []byte {
	returnMessage, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("Error: problems encoding elevator struct.\n")
		panic(err)
	}
	return returnMessage
}

func DecodeMessage(Msg *Elev_control.Elevator, arr []byte) {
	json.Unmarshal(arr, Msg)
}

/*
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func ListenToBroadcast(listenAddress string, receiveUdpCh chan Elev_control.Elevator, timeoutCh chan int) {
	fmt.Println("Running ListenToBroadcast")
	ServerAddress, err := net.ResolveUDPAddr("udp", listenAddress)
	CheckError(err)

	ServerConnection, err := net.ListenUDP("udp", ServerAddress)
	CheckError(err)

	defer ServerConnection.Close()

	var elev Elev_control.Elevator

mainloop:
	for {
		buffer := make([]byte, 1024)
		ServerConnection.SetDeadline(time.Now().Add(3000 * time.Millisecond))
		_, _, err := ServerConnection.ReadFromUDP(buffer)
		CheckError(err)
		DecodeMessage(buffer, &elev)
		if err != nil {
			timeoutCh <- 1
			break mainloop
		}
		receiveUdpCh <- elev
	}
}

func UdpBroadcast(toAdress string, elev Elev_control.Elevator) {
	ServerAddress, err := net.ResolveUDPAddr("udp", toAdress)
	CheckError(err)

	Connection, err := net.DialUDP("udp", nil, ServerAddress) //Returns type UDPConn which we can send msg to
	CheckError(err)

	defer Connection.Close() //Closes the connection after udpBroadcast functioncall
	//buffer := make([]byte, 1024)
	buffer := EncodeMessage(elev)
	_, err = Connection.Write(buffer)
	CheckError(err)
	fmt.Println("Sender status")
	//time.Sleep(time.Millisecond * 1000)

}
*/
