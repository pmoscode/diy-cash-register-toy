package writer

import (
	"github.com/tarm/serial"
	"log"
)

type Console struct {
}

func (c *Console) Write(message string) {
	// setup serial
	config := &serial.Config{
		Name: "COM5",
		Baud: 115200,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()
}
