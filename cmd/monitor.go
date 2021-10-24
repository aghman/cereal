package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"

	"github.com/aghman/cereal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(monitorCmd)
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Starts monitoring of serial port",
	Long:  `Starts monitoring of serial port`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

type DataLineHandler func(string) error

var portToUse string
var runningConfig config.CerealConfig

func run() {

	location := viper.Get("location")
	runningConfig = *config.NewCerealConfig(location.(map[string]interface{}))
	fmt.Println(runningConfig)
	mode := &serial.Mode{
		BaudRate: 115200,
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
		fmt.Printf("%v", stringData)
		if strings.Contains(stringData, "\n") {

			newLineSplitParts := strings.Split(stringData, "\n")
			if len(newLineSplitParts) == 2 {
				lineOfData += newLineSplitParts[0]
				HandleSerialOutput(lineOfData)
				lineOfData = ""
			}
		} else {
			lineOfData += stringData
		}

	}

}
func HandleSerialOutput(dataLine string) error {
	if dataLine == "" {
		return fmt.Errorf("dataline is empty. %s", dataLine)
	}
	dataFile, err := os.Create(runningConfig.Location.OutputFile)
	if err != nil {
		return fmt.Errorf("cannot create file: %v", err)
	}
	defer dataFile.Close()
	_, err = dataFile.WriteString(dataLine)
	if err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	dataFile.Sync()

	return nil
}
