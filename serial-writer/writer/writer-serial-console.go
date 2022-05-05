package writer

import (
	"github.com/tarm/serial"
	"log"
)

type SerialConsole struct {
	InterfaceName     string
	InterfaceBaudRate int
	port              *serial.Port
}

func (c SerialConsole) Connect() {
	config := &serial.Config{
		Name: c.InterfaceName,
		Baud: c.InterfaceBaudRate,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	c.port = port
}

func (c *SerialConsole) Write(message string) {
	if c.port == nil {
		log.Fatalln("SerialConsole was not setup properly. Use 'connect' method to do this.")
	} else {
		_, err := c.port.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
	}
}
