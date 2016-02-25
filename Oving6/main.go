package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func ListenToBroadcast(fromPort string, receiveCh, timeoutCh chan int) {

	ServerAddress, err := net.ResolveUDPAddr("udp", ":"+fromPort)
	CheckError(err)

	ServerConnection, err := net.ListenUDP("udp", ServerAddress)
	CheckError(err)

	defer ServerConnection.Close()

	var count int64

	//mainloop:
	for {
		buffer := make([]byte, 64)
		ServerConnection.SetDeadline(time.Now().Add(1000 * time.Millisecond))
		n, _, err := ServerConnection.ReadFromUDP(buffer)
		CheckError(err)
		count, _ = binary.Varint(buffer[0:n])

		if err != nil {
			timeoutCh <- 1
			break //mainloop
		}

		receiveCh <- int(count)
	}
}

func UdpBroadcast(toAdress string, countCh chan int) {
	ServerAddress, err := net.ResolveUDPAddr("udp", toAdress)
	CheckError(err)

	Connection, err := net.DialUDP("udp", nil, ServerAddress) //Returns type UDPConn which we can send msg to
	CheckError(err)

	defer Connection.Close() //Closes the connection after udpBroadcast functioncall

	for {
		select {
		case count := <-countCh:
			buffer := make([]byte, 64)
			binary.PutVarint(buffer, int64(count))
			_, err := Connection.Write(buffer)
			CheckError(err)

			time.Sleep(time.Millisecond * 100)
		default:

		}
	}
}

func Backup() {
	receiveCh := make(chan int)
	timeoutCh := make(chan int)
	countCh := make(chan int)
	latestValue := 0
	localAddr := "129.241.187.161"
	port := "25000"

	go ListenToBroadcast(port, receiveCh, timeoutCh)

	for {
		select {
		case count := <-receiveCh: //Finnes allerede en primal som kjører
			if latestValue > 0 {
				latestValue = count
			}
		case <-timeoutCh: //Finner ikke primal, start backup og counter
			go Counter(countCh, latestValue)
			go UdpBroadcast(localAddr, countCh)
			exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			fmt.Println(latestValue)
		}
	}
}

func Counter(countCh chan int, latestValue int) {
	for {
		latestValue++
		countCh <- latestValue
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	Backup()
}
