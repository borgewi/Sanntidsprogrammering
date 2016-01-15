package main

import (
	"fmt"
	"time"
)

var i int = 0 

func plus(ic chan int){
	for j := 0; j < 10; j++ {
	  ic <- 1
		i += 1
		<- ic
	}
}

func minus(ic chan int){
	for j := 0; j < 10; j++ {
	  ic <- 1
		i -= 1
		<- ic
	}
}

func main(){
  ic := make(chan int, 1)
  go plus(ic)
	go minus(ic)
	time.Sleep(100*time.Millisecond)
	fmt.Println("i: ",i)
}
