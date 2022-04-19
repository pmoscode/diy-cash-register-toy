package main

import (
	"flag"
	"fmt"
)

func setupCommandLine() (*string, *string, *int, *bool) {
	clkPin := flag.String("clk-pin", "16", "Clock pin on RPI by gpio block pin number")
	dataPin := flag.String("data-pin", "18", "Data pin on RPI by gpio block pin number")
	delay := flag.Int("delay", 2, "Delay between digit transition in milliseconds")
	debug := flag.Bool("debug", false, "Debug or not debug??")
	flag.Parse()

	return clkPin, dataPin, delay, debug
}

func main() {
	clkPin, dataPin, delay, debug := setupCommandLine()
	fmt.Printf("Starting VDF controller with Clock-Pin: %s and Data-Pin %s", *clkPin, *dataPin)

	vfdController := ControllerFactory(*clkPin, *dataPin, *delay, *debug)

	for true {
		vfdController.SetNegation(5)
		vfdController.SetNumber(2, 4)
		vfdController.SetNumberWithDot(9, 3)
		vfdController.SetNumber(9, 2)
		vfdController.SetNumber(5, 1)
	}
}
