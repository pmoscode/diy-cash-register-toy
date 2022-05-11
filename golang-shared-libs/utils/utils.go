package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func checkCommand(command string, params []string, validate string) bool {
	isRaspi := true

	path, _ := exec.LookPath(command)
	cmd := exec.Command(path, params...)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	fmt.Println(string(output))

	if err != nil {
		isRaspi = false
	} else {
		if !strings.Contains(string(output), validate) {
			isRaspi = false
		}
	}

	return isRaspi
}

func IsRaspberryPiHardware() bool {
	// Test for "Hardware" content in /proc/cpuinfo
	foundSoC := checkCommand("grep", strings.Fields("Hardware /proc/cpuinfo"), "BCM2835")

	// Test for result in "uname -m"
	foundArch := checkCommand("uname", strings.Fields("-m"), "arm")

	// Test for mac address (eth0) starts with b8:27:...
	foundMAC := checkCommand("bash", []string{"-c", "ip a s eth0 | grep ether"}, "b8:27")

	return foundSoC || foundArch || foundMAC
}
