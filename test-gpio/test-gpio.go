package main

import (
	log "github.com/pmoscode/go-common/logging"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

var pinNumber = "18"
var sleepSeconds = 2

var logger *log.Logger

func main() {
	logger = log.NewLogger()

	if _, err := host.Init(); err != nil {
		logger.Fatal(err, "Could not initialize host")
	}

	pin := gpioreg.ByName("P1_" + pinNumber)
	high := true

	for {
		if high {
			err := pin.Out(gpio.High)
			if err != nil {
				logger.Error(err, "Error setting pin high")
				return
			}
		} else {
			err := pin.Out(gpio.Low)
			if err != nil {
				logger.Error(err, "Error setting pin low")
				return
			}
		}
		high = !high
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}
}
