package main

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"time"
)

type Controller struct {
	clkPin  string
	dataPin string
	delay   int
}

func (c Controller) Init(clkPin string, dataPin string, delay int) {
	c.dataPin = dataPin
	c.clkPin = clkPin
	c.delay = delay
}

func (c Controller) digitalWrite(pin string, level gpio.Level) {
	err := gpioreg.ByName(pin).Out(level)
	if err != nil {
		return
	}
}

func (c Controller) delayMicroseconds(mseconds int) {
	time.Sleep(time.Duration(mseconds) * time.Microsecond)
}

func (c Controller) delayMilliseconds(mseconds int) {
	time.Sleep(time.Duration(mseconds) * time.Millisecond)
}

func (c Controller) writeCharacter(data byte, mask byte, lastDelay int) {
	var maskTemp byte
	for maskTemp = mask; maskTemp > 0; maskTemp <<= 1 {
		c.digitalWrite(c.clkPin, gpio.Low)
		time.Sleep(5 * time.Microsecond)
		if data&maskTemp > 0 {
			c.digitalWrite(c.dataPin, gpio.High)
		} else {
			c.digitalWrite(c.dataPin, gpio.Low)
		}
		c.digitalWrite(c.clkPin, gpio.High)
		c.delayMicroseconds(lastDelay)
	}
}

func (c Controller) writeCharacterFull(data byte) {
	c.writeCharacter(data, 0b00000001, 5)
}

func (c Controller) writeCharacterHalf(data byte) {
	c.writeCharacter(data, 0b00010000, 1)
}

func (c Controller) sendData(blocks [3]byte) {
	c.delayMicroseconds(1)
	c.writeCharacterFull(blocks[0])
	c.writeCharacterFull(blocks[1])
	c.writeCharacterHalf(blocks[2])
	c.delayMicroseconds(1)
}

func (c Controller) processNumber(blocks [3]byte, number int) {
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

func (c Controller) processPosition(blocks [3]byte, position int) {
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

func (c Controller) SetNumber(number int, position int) {
	c.set(number, position, false)
}

func (c Controller) SetNumberWithDot(number int, position int) {
	c.set(number, position, true)
}

func (c Controller) set(number int, position int, dot bool) {
	var blocks [3]byte

	blocks[0] = 0b00000000
	blocks[1] = 0b00000000
	blocks[2] = 0b00000000

	c.processNumber(blocks, number)
	c.processPosition(blocks, position)

	if dot {
		blocks[2] = blocks[2] | 0b00010000
	}

	c.sendData(blocks)
	c.delayMilliseconds(c.delay)
}

func (c Controller) SetNegation(position int) {

	var blocks [3]byte

	blocks[0] = 0b00000000
	blocks[1] = 0b00000000
	blocks[2] = 0b00000000

	c.processPosition(blocks, position)

	blocks[2] = blocks[2] | 0b01000000

	c.sendData(blocks)
}
