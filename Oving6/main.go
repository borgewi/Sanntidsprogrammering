package Network

import "main"
import "fmt"
import "udp"

func Backup(){
	receiveCh := make(chan int)
	timeoutCh := make(chan int)
	countCh := make(chan int)
	latestValue := 0

	address = "localhost:25000"
	go ListenToBroadcast(address, receiveCh, timeoutCh chan int)

	for{
		select{
			case count := <- receiveCh: //Finnes allerede en primal som kjører
				if(latestValue > 0) {latestValue = count}  
			case <- timeoutCh: //Finner ikke primal, start backup og counter
				go Counter(countCh, latestValue)
				go UdpBroadcast(address, receiveCh,timeoutCh)
				exec.Command("gnome-terminal","-x", "sh", "-c", "go run main.go")
		}
	}
}


func Counter(countCh chan int, latestValue int){
	for{
		latestValue++
		countCh <- latestValue
		time.Sleep(500*time.Millisecond)
	}
}

func main() {
	go Backup()
}

