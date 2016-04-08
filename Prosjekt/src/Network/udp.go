package Network

import (
	"encoding/binary/json"
	"fmt"
	"net"
	"time"
	"Elev_control"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func ListenToBroadcast(listenAddress string, receiveCh chan Elevator, timeoutCh chan int) {

	ServerAddress, err := net.ResolveUDPAddr("udp", listenAddress)
	CheckError(err)

	ServerConnection, err := net.ListenUDP("udp", ServerAddress)
	CheckError(err)

	defer ServerConnection.Close()

	var elev Elevator

	mainloop:
	for {
		buffer := make([]byte, 64)
		ServerConnection.SetDeadline(time.Now().Add(3000 * time.Millisecond))
			n, _, err := ServerConnection.ReadFromUDP(buffer)
			CheckError(err)
			err := json.Unmarshal(buffer, &elev)
			receiveCh <- elev

			if err != nil {
				timeoutCh <- 1
				break mainloop
			}
	}
}


func UdpBroadcast(toAdress string, elev Elevator) {
	ServerAddress, err := net.ResolveUDPAddr("udp", toAdress)
	CheckError(err)

	Connection, err := net.DialUDP("udp", nil, ServerAddress) //Returns type UDPConn which we can send msg to
	CheckError(err)

	defer Connection.Close() //Closes the connection after udpBroadcast functioncall
	buffer := make([]byte, 64)
	buffer, err := json.Marshal(elev)
	_, err := Connection.Write(buffer)
	CheckError(err)
	//time.Sleep(time.Millisecond * 1000)
	}
}