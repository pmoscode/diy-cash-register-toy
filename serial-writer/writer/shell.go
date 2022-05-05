package writer

import (
	"log"
	"os/exec"
)

type Shell struct {
	InterfaceName string
}

func (s Shell) Connect() {}

func (s Shell) Disconnect() {}

func (s Shell) Write(message string) {
	cmd := exec.Command("echo", "-e", message, ">", s.InterfaceName)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
