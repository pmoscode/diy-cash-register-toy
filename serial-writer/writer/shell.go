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
	cmd := exec.Command("bash", "-c", "echo -e", message, ">", s.InterfaceName)
	err := cmd.Run()

	var exitCode = "0"
	if werr, ok := err.(*exec.ExitError); ok {
		if s := werr.Error(); s != "0" {
			exitCode = s
		}
	}
	log.Println("command executed with exit code ", exitCode)
}
