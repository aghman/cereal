package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DataLineHandler func(string) error

var influxClient influxdb2.Client
var writeAPI api.WriteAPI
var portToUse string

func main() {

	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "hEwSGkJMW1Nua6jfNL3q63IlUB2hgUWjfcorFQsw9cwUbKSepzwbZUgLgj3uSAz2oQXxHMcra61gWf2PT1DBgA=="
	const bucket = "garden"
	const org = "home"

	influxClient = influxdb2.NewClient("http://es:8086", token)
	// always close client at the end
	defer influxClient.Close()
	writeAPI = influxClient.WriteAPI(org, bucket)

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
	/* Read a chunk of bytes. Convert to string.
	 * If we have a newline, it indicates that the data entry is complete.
	 * If there is no newline, append to lineOfData.
	 * If there is a newline, and only one element is in the split array,
	 * append to lineOfData and submit to the handling method.
	 * If there is more than one element in the split array,
	 * take second part (or more) and add to new instance of lineOfData.
	 */

	for { //loop until we die
		n, err := serialPort.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 { // serial port is gone
			fmt.Println("\nEOF")
			break
		}

		stringData := string(buff[:n])
		//fmt.Printf("%v", stringData)
		if strings.Contains(stringData, "\n") {

			newLineSplitParts := strings.Split(stringData, "\n")
			if len(newLineSplitParts) == 2 {
				lineOfData += newLineSplitParts[0]
				// send to handler - TODO
				HandleParticleData(lineOfData)
				lineOfData = ""
			}
		} else {
			lineOfData += stringData
		}

	}

}

const expectedParticleDataElements = 9

type PanasonicSNGCJA5Data struct {
	PM_1_0 float64
	PM_2_5 float64
	PM_10  float64
	PC_0_5 int64
	PC_1   int64
	PC_2_5 int64
	PC_5   int64
	PC_7_5 int64
	PC_10  int64
}

/* HandleParticleData processes a serial line of data from the Panasonic SN-GCJA5 Particle Sensor
 * Data format is as follows, split by commas:
 * 1) PM1.0
 * 2) PM2.5
 * 3) PM10
 * 4) Particulate Count 0.5
 * 5) Particulate Count 1.0
 * 6) Particulate Count 2.5
 * 7) Particulate Count 5.0
 * 8) Particulate Count 7.5
 * 9) Particulate Count 10.0
 */

func HandleParticleData(dataLine string) error {
	fmt.Printf("Handling %s", dataLine)
	if dataLine == "" {
		return fmt.Errorf("dataline is empty. %s", dataLine)
	}
	dataComponents := strings.Split(dataLine, ",")
	if len(dataComponents) != expectedParticleDataElements {
		return fmt.Errorf("unexpected data format. Expected %d, received %d ", expectedParticleDataElements, len(dataComponents))
	}
	var currentData PanasonicSNGCJA5Data
	currentData.PM_1_0, _ = strconv.ParseFloat(dataComponents[0], 64)
	currentData.PM_2_5, _ = strconv.ParseFloat(dataComponents[1], 64)
	currentData.PM_10, _ = strconv.ParseFloat(dataComponents[2], 64)
	currentData.PC_0_5, _ = strconv.ParseInt(dataComponents[3], 0, 0)
	currentData.PC_1, _ = strconv.ParseInt(dataComponents[4], 0, 0)
	currentData.PC_2_5, _ = strconv.ParseInt(dataComponents[5], 0, 0)
	currentData.PC_5, _ = strconv.ParseInt(dataComponents[6], 0, 0)
	currentData.PC_7_5, _ = strconv.ParseInt(dataComponents[7], 0, 0)
	currentData.PC_10, _ = strconv.ParseInt(dataComponents[8], 0, 0)
	//fmt.Print(currentData)
	// write line protocol
	writeAPI.WriteRecord(fmt.Sprintf("pm1.0,unit=ugm3 sensor=%f", currentData.PM_1_0))
	writeAPI.WriteRecord(fmt.Sprintf("pm2.5,unit=ugm3 sensor=%f", currentData.PM_2_5))
	writeAPI.WriteRecord(fmt.Sprintf("pm10,unit=ugm3 sensor=%f", currentData.PM_10))
	writeAPI.WriteRecord(fmt.Sprintf("pc0.5,unit=particles count=%d", currentData.PC_0_5))
	writeAPI.WriteRecord(fmt.Sprintf("pc1.0,unit=particles count=%d", currentData.PC_1))
	writeAPI.WriteRecord(fmt.Sprintf("pc2.5,unit=particles count=%d", currentData.PC_2_5))
	writeAPI.WriteRecord(fmt.Sprintf("pc5.0,unit=particles count=%d", currentData.PC_5))
	writeAPI.WriteRecord(fmt.Sprintf("pc7.5,unit=particles count=%d", currentData.PC_7_5))
	writeAPI.WriteRecord(fmt.Sprintf("pc10,unit=particles count=%d", currentData.PC_10))
	// Flush writes
	writeAPI.Flush()

	return nil
}
