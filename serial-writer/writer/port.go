package writer

import (
	"github.com/tarm/serial"
	"log"
)

type Port struct {
	InterfaceName     string
	InterfaceBaudRate int
	port              *serial.Port
}

func (c *Port) Connect() {
	log.Println("Setting up Serial with interface ", c.InterfaceName, " and at baud rate ", c.InterfaceBaudRate)

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

func (c *Port) Write(message string) {
	if c.port == nil {
		log.Fatalln("Port was not setup properly. Use 'connect' method to do this.")
	} else {
		log.Println("Sending '", message, "' to Port...")
		writtenBytes, err := c.port.Write([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("... with bytes written: ", writtenBytes)
	}
}

func (c *Port) Disconnect() {
	c.port.Close()
}
