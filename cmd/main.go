package main

import (
	"fmt"
	"log"
	"os"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {

	var portToUse string = ""
	mode := &serial.Mode{
		BaudRate: 9600,
	}

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	for _, port := range ports {
		//fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("Found USB Port at %s, using for serial scanning. \n", port.Name)
			portToUse = port.Name
		}
	}

	buff := make([]byte, 100)
	var lineOfData string = ""

	if portToUse == "" {
		fmt.Println("Unable to find a serial port. Exiting...")
		os.Exit(-1)
	}
	serialPort, err := serial.Open(portToUse, mode)
	if err != nil {
		log.Fatal(err)
	}
	// Read a chunk of bytes. Convert to string.
	// If we have a newline, it indicates that the data entry is complete.
	// If there is no newline, append to lineOfData.
	// If there is a newline, and only one element is in the split array,
	// append to lineOfData and submit to the handling method.
	// If there is more than one element in the split array,
	// take second part (or more) and add to new instance of lineOfData.
	for { //loop until we die
		n, err := serialPort.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}

}
