package main

import (
	"flag"
	"fmt"
	"strconv"
)

func setupCommandLine() (*string, *string, *int, *string, *int, *bool) {
	clkPin := flag.String("clk-pin", "16", "Clock pin on RPI by gpio block pin number")
	dataPin := flag.String("data-pin", "18", "Data pin on RPI by gpio block pin number")
	delay := flag.Int("delay", 2, "Delay between digit transition in milliseconds")
	number := flag.String("number", "", "Set only one number (0-9) and exit")
	position := flag.Int("position", 1, "Only considered whe 'number' is set")
	debug := flag.Bool("debug", false, "Debug or not debug??")
	flag.Parse()

	return clkPin, dataPin, delay, number, position, debug
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func main() {
	clkPin, dataPin, delay, number, position, debug := setupCommandLine()
	fmt.Printf("Starting VDF controller with Clock-Pin: %s and Data-Pin %s\n", *clkPin, *dataPin)

	vfdController := ControllerFactory(*clkPin, *dataPin, *delay, *debug)

	if *number != "" {
		if len(*number) > 1 {
			numStr := reverse(*number)
			for true {
				for pos, char := range numStr {
					letter := string(char)
					num, _ := strconv.Atoi(letter)
					//fmt.Printf("Printing %d at position %d\n", num, pos+*position)
					vfdController.SetNumber(num, pos+*position)
				}
			}
		} else {
			num, _ := strconv.Atoi(*number)
			vfdController.SetNumber(num, *position)
		}
	} else {
		fmt.Println("MQTT subscribe not implemented yet... exiting...")
		//for true {
		//	vfdController.SetNegation(5)
		//	vfdController.SetNumber(2, 4)
		//	vfdController.SetNumberWithDot(9, 3)
		//	vfdController.SetNumber(9, 2)
		//	vfdController.SetNumber(5, 1)
		//}
	}
}
