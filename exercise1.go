package main

import (
	"fmt"
	"time"
)

i int = 0;

func plus(){
	for j := 0; j < 10; j++ {
		i += 1
	}
}

func minus(){
	for j := 0; j < 10; j++ {
		i -= 1
	}
}

func main(){
	time.Sleep(100*time.Millisecond)
}