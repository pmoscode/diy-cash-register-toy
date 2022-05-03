package main

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"time"
)

type Controller struct {
	clkPin  gpio.PinOut
	dataPin gpio.PinOut
	delay   int
	debug   bool
}

func (c Controller) SetNumber(number int, position int) {
	c.set(number, position, false)
}

func (c Controller) SetNumberWithDot(number int, position int) {
	c.set(number, position, true)
}

func (c Controller) SetNegation(position int) {

	var blocks [3]byte

	blocks[0] = 0b00000000
	blocks[1] = 0b00000000
	blocks[2] = 0b00000000

	c.processPosition(&blocks, position)

	blocks[2] = blocks[2] | 0b01000000

	c.sendData(blocks)
}

func (c Controller) digitalWrite(pin gpio.PinOut, level gpio.Level) {
	err := pin.Out(level)
	if err != nil {
		fmt.Printf("error caugth: %t \n", err)
		return
	}
}

func (c Controller) delayNanoseconds(value int) {
	time.Sleep(time.Duration(value) * time.Nanosecond)
}

func (c Controller) delayMicroseconds(value int) {
	time.Sleep(time.Duration(value) * time.Microsecond)
}

func (c Controller) delayMilliseconds(value int) {
	time.Sleep(time.Duration(value) * time.Millisecond)
}

func (c Controller) writeCharacter(data byte, mask byte, lastDelay int) {
	if c.debug {
		fmt.Printf("Writing data: %08b \n", data)
	}

	var dataTemp byte = 170
	var maskTemp byte = 1

	dataTemp = data

	for maskTemp = mask; maskTemp > 0; maskTemp <<= 1 {
		c.digitalWrite(c.clkPin, gpio.Low)
		c.delayMicroseconds(5)
		if dataTemp&maskTemp > 0 {
			c.digitalWrite(c.dataPin, gpio.High)
			if c.debug {
				fmt.Printf("Pin %s -> %s \n", c.dataPin, "HIGH")
			}
		} else {
			c.digitalWrite(c.dataPin, gpio.Low)
			if c.debug {
				fmt.Printf("Pin %s -> %s \n", c.dataPin, "LOW")
			}
		}
		c.digitalWrite(c.clkPin, gpio.High)
		c.delayMicroseconds(lastDelay)
	}

	if c.debug {
		fmt.Println()
	}
}

func (c Controller) writeCharacterFull(data byte) {
	c.writeCharacter(data, 0b00000001, 5)
}

func (c Controller) writeCharacterHalf(data byte) {
	c.writeCharacter(data, 0b00010000, 1)
}

func (c Controller) sendData(blocks [3]byte) {
	if c.debug {
		fmt.Println("Send data: ", blocks)
		fmt.Print("Send data (binary): ")
		for _, n := range blocks {
			fmt.Printf("%08b ", n)
		}
		fmt.Println()
	}

	c.delayMicroseconds(1)
	c.writeCharacterFull(blocks[0])
	c.writeCharacterFull(blocks[1])
	c.writeCharacterHalf(blocks[2])
	c.delayMicroseconds(1)
}

func (c Controller) processNumber(blocks *[3]byte, number int) {
	switch number {
	case 0:
		blocks[0] = blocks[0] | 0b01111000
		blocks[2] = blocks[2] | 0b10100000
		break
	case 1:
		blocks[0] = blocks[0] | 0b00110000
		blocks[2] = blocks[2] | 0b00000000
		break
	case 2:
		blocks[0] = blocks[0] | 0b01011000
		blocks[2] = blocks[2] | 0b01100000
		break
	case 3:
		blocks[0] = blocks[0] | 0b01111000
		blocks[2] = blocks[2] | 0b01000000
		break
	case 4:
		blocks[0] = blocks[0] | 0b00110000
		blocks[2] = blocks[2] | 0b11000000
		break
	case 5:
		blocks[0] = blocks[0] | 0b01101000
		blocks[2] = blocks[2] | 0b11000000
		break
	case 6:
		blocks[0] = blocks[0] | 0b01101000
		blocks[2] = blocks[2] | 0b11100000
		break
	case 7:
		blocks[0] = blocks[0] | 0b00111000
		blocks[2] = blocks[2] | 0b00000000
		break
	case 8:
		blocks[0] = blocks[0] | 0b01111000
		blocks[2] = blocks[2] | 0b11100000
		break
	case 9:
		blocks[0] = blocks[0] | 0b01111000
		blocks[2] = blocks[2] | 0b11000000
		break

	default:
		break
	}
}

func (c Controller) processPosition(blocks *[3]byte, position int) {
	switch position {
	case 1:
		blocks[0] = blocks[0] | 0b10000000
		break
	case 2:
		blocks[1] = blocks[1] | 0b00000001
		break
	case 3:
		blocks[1] = blocks[1] | 0b00000010
		break
	case 4:
		blocks[1] = blocks[1] | 0b00000100
		break
	case 5:
		blocks[1] = blocks[1] | 0b00001000
		break
	case 6:
		blocks[1] = blocks[1] | 0b00010000
		break
	case 7:
		blocks[1] = blocks[1] | 0b00100000
		break
	case 8:
		blocks[1] = blocks[1] | 0b01000000
		break
	case 9:
		blocks[1] = blocks[1] | 0b10000000
		break
	case 10:
		blocks[0] = blocks[0] | 0b00000010
		break

	default:
		break
	}
}

func (c Controller) set(number int, position int, dot bool) {
	if c.debug {
		fmt.Printf("Printing number %d on position %d with dot %t \n", number, position, dot)
	}

	var blocks [3]byte

	blocks[0] = 0b00000000
	blocks[1] = 0b00000000
	blocks[2] = 0b00000000

	c.processNumber(&blocks, number)
	c.processPosition(&blocks, position)

	if dot {
		blocks[2] = blocks[2] | 0b00010000
	}

	c.sendData(blocks)
	c.delayMilliseconds(c.delay)
}

func ControllerFactory(clkPin string, dataPin string, delay int, debug bool) Controller {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	return Controller{
		clkPin:  gpioreg.ByName("P1_" + clkPin),
		dataPin: gpioreg.ByName("P1_" + dataPin),
		delay:   delay,
		debug:   debug,
	}
}
