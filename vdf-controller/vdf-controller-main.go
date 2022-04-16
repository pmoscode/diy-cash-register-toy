package main

import (
	"flag"
	"log"
	"periph.io/x/host/v3"
)

func setupCommandLine() (*string, *string, *int) {
	clkPin := flag.String("clk-pin", "16", "Clock pin on RPI by gpio block pin number")
	dataPin := flag.String("data-pin", "18", "Data pin on RPI by gpio block pin number")
	delay := flag.Int("delay", 2, "Delay between digit transition in milliseconds")
	flag.Parse()

	return clkPin, dataPin, delay
}

func main() {
	clkPin, dataPin, delay := setupCommandLine()

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	controller := Controller{
		clkPin:  *clkPin,
		dataPin: *dataPin,
		delay:   *delay,
	}

	for true {
		controller.SetNegation(5)
		controller.SetNumber(2, 4)
		controller.SetNumberWithDot(9, 3)
		controller.SetNumber(9, 2)
		controller.SetNumber(5, 1)
	}
}
