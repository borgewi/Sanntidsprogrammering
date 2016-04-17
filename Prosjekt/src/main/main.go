package main

import (
	"Elev_control"
	//"net"
	//"os/exec"
	//"time"
	"Driver"
	//"Network"
	"Master_Slave"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	//"strings"
)

const (
	NUMBUTTONS  int = 3
	NUMFLOORS   int = 4
	BUTTON_DOWN     = 0 + iota
	BUTTON_UP
	BUTTON_COMMAND
)

func main() {
	//go Backup()
	Driver.ElevInit()
	localStatusCh := make(chan Elev_control.Elevator, 100)
	sendBtnCallCh := make(chan [NUMBUTTONS - 1]int, 100)
	errorCh := make(chan int)
	receiveAllBtnCallsCh := make(chan [NUMFLOORS][NUMBUTTONS - 1]bool, 100)
	setLights_setExtBtnsCh := make(chan [4][2]bool, 100)

	go Elev_control.Run_Elevator(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, setLights_setExtBtnsCh, errorCh)
	go Master_Slave.Princess(localStatusCh, sendBtnCallCh, receiveAllBtnCallsCh, setLights_setExtBtnsCh, errorCh)
	var err int
	for {
		err = <-errorCh
		if err == 1 {
			fmt.Println("Error har oppstått. Har vært i EB_Moving for lenge. err = ", err)
			for {
				err = <-errorCh
				if err == 0 {
					fmt.Println("Error er fikset")
					break
				}
			}
		}
	}
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func writeLines(lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	//writer := bufio.NewWriter(file)
	for _, item := range lines {
		fmt.Println("				item")
		_, err := file.WriteString(item + "\n")
		file.Write([]byte(item))
		if err != nil {
			fmt.Println("debug")
			fmt.Println(err)
			break
		}
	}
	/*content := strings.Join(lines, "\n")
	  _, err = writer.WriteString(content)*/
	return
}

func Backup() {
	lines, err := readLines("/home/student/Desktop/BorgOsk/internalOrders.txt")
	if err != nil {
		fmt.Println("Error: %s\n", err)
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	//array := []string{"7.0", "8.5", "9.1"}
	err = writeLines(lines, "/home/student/Desktop/BorgOsk/internalOrders.txt")
	fmt.Println(err)
}
