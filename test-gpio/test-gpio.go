package main

import (
	log "gitlab.com/pmoscode/golang-shared-libs/logging"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

var pinNumber = "18"
var sleepSeconds = 2

var logger *log.Logger

func main() {
	logger = log.NewLogger("log.log")

	if _, err := host.Init(); err != nil {
		logger.Fatal(err)
	}

	pin := gpioreg.ByName("P1_" + pinNumber)
	high := true

	for {
		if high {
			err := pin.Out(gpio.High)
			if err != nil {
				logger.Print("Error: ", err)
				return
			}
		} else {
			err := pin.Out(gpio.Low)
			if err != nil {
				logger.Print("Error: ", err)
				return
			}
		}
		high = !high
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}
