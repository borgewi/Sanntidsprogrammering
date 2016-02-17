package main
 
import (
    "fmt"
    "net"
    "os"
    "strings"
    "strconv"
    "encoding/binary"
    "time"
)


func CheckError( err error ) {
    if err  != nil {
        fmt.Println( "Error: " , err )
        os.Exit(0)
    }
}

//New udp initialize function, which uses resolveUDPAddr, DialUDP, and ListenUDP

func udpBroadcast( toAdress string , activationTime int64 ) {
    ServerAddress,err := net.ResolveUDPAddr( "udp" , toAdress  )
    CheckError( err )

    Connection, err := net.DialUDP( "udp", nil, ServerAddress ) //Returns type UDPConn which we can send msg to
    CheckError( err )

    defer Connection.Close() //Closes the connection after udpBroadcast functioncall

    for {

        buffer := make( []byte, 1024 )
        binary.PutVarint( buffer , activationTime )

        _,err := Connection.Write( buffer )
        CheckError( err )

        time.Sleep( time.Millisecond * 1000 ) //Will we be needing much sleeptime?
    }
}

func listenToBroadcast( fromPort string ) {

    ServerAddress,err := net.ResolveUDPAddr( "udp" , ":" + fromPort)
    CheckError(err)

    ServerConnection, err := net.ListenUDP( "udp" , ServerAddress )
    CheckError(err)
    
    defer ServerConnection.Close()
    
    buffer := make( []byte, 1024 )
 
    for {
        n, senderAddress ,err := ServerConnection.ReadFromUDP( buffer ) //Might want senderAddr
        CheckError( err )

        receivedActivationTime, _ := binary.Varint(buffer[0:n])
        fmt.Println( "Received IP adress" ,  senderAddress.IP.String() , "and activation time" ,receivedActivationTime)
        //fmt.Println(serverConn.receiveFrom()) 
    }
    //receivedId := elevatorId{ activationTime: strconv.Atoi(receivedStrings[0]) , elevatorIndex: strconv.Atoi(receivedStrings[1]) }
    //return receivedId
}

func main() {
    //****************************
    referenceTime := time.Unix(1451606400,0) //Unix 2016 01/01/00
    currentTime := time.Now() //returns local unix time

    var duration time.Duration = currentTime.Sub( referenceTime )
    var activationTime = int64(duration)
    //****************************

    //udpBroadcast( "129.241.187.255:30007" , activationTime )
}

//ch <- v    // Send v to channel ch.
//v := <-ch  // Receive from ch, and
//           // assign value to v.
