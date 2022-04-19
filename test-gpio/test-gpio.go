package main

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

var pinNumber = "18"
var sleepSeconds = 2

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	pin := gpioreg.ByName("P1_" + pinNumber)
	high := true

	for {
		if high {
			err := pin.Out(gpio.High)
			if err != nil {
				fmt.Print("Error: ", err)
				return
			}
		} else {
			err := pin.Out(gpio.Low)
			if err != nil {
				fmt.Print("Error: ", err)
				return
			}
		}
		high = !high
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}
