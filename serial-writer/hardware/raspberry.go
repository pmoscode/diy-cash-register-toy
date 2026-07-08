// Package hardware provides helper functions to detect the runtime hardware.
//
// This used to live in the shared "golang-shared-libs/utils" package, but
// there is no equivalent in "github.com/pmoscode/go-common", so it has been
// vendored locally here.
package hardware

import (
	"fmt"
	"os/exec"
	"strings"
)

func checkCommand(command string, params []string, validate string) bool {
	path, _ := exec.LookPath(command)
	cmd := exec.Command(path, params...)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	fmt.Println(string(output))

	if err != nil {
		return false
	}

	return strings.Contains(string(output), validate)
}

// IsRaspberryPiHardware tries to detect if the current process runs on a Raspberry Pi.
func IsRaspberryPiHardware() bool {
	// Test for "Hardware" content in /proc/cpuinfo
	foundSoC := checkCommand("grep", strings.Fields("Hardware /proc/cpuinfo"), "BCM2835")

	// Test for result in "uname -m"
	foundArch := checkCommand("uname", strings.Fields("-m"), "arm")

	// Test for mac address (eth0) starts with b8:27:...
	foundMAC := checkCommand("bash", []string{"-c", "ip a s eth0 | grep ether"}, "b8:27")

	return foundSoC || foundArch || foundMAC
}
