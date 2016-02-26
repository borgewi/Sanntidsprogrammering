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

func ListenToBroadcast(listenAddress string, receiveCh, timeoutCh chan int) {

	ServerAddress, err := net.ResolveUDPAddr("udp", listenAddress)
	CheckError(err)

	ServerConnection, err := net.ListenUDP("udp", ServerAddress)
	CheckError(err)

	defer ServerConnection.Close()

	var count int64

	mainloop:
	for {
		buffer := make([]byte, 64)
		ServerConnection.SetDeadline(time.Now().Add(3000 * time.Millisecond))
			n, _, err := ServerConnection.ReadFromUDP(buffer)
			CheckError(err)
			count, _ = binary.Varint(buffer[0:n])
			c := int(count)
			receiveCh <- c

			if err != nil {
				timeoutCh <- 1
				break mainloop
			}

			
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

			//time.Sleep(time.Millisecond * 1000)
		default:

		}
	}
}

func Backup() {
	receiveCh := make(chan int)
	timeoutCh := make(chan int)
	countCh := make(chan int)
	latestValue := 0
	var p *int = &latestValue
	localAddr := "129.241.187.255" + ":13337"

	go ListenToBroadcast(localAddr, receiveCh, timeoutCh)

	for {
		select {
		case count := <-receiveCh: //Finnes allerede en primal som kjører
			if count>0{
				*p = count
			}
		case <-timeoutCh: //Finner ikke primal, start backup og counter
			go Counter(countCh, *p)
			go UdpBroadcast(localAddr, countCh)
			callBackup := exec.Command("gnome-terminal", "-x",  "sh", "-c", "go run main.go")
			callBackup.Run()
			break
		}
	}
}

func Counter(countCh chan int, latestValue int) {
	for {
		fmt.Println(latestValue)
		latestValue++
		countCh <- latestValue
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	Backup()
}
