package vaulted

import (
	"os/exec"
)

func Installed() bool {
	err := exec.Command("which", "vaulted").Run()
	if err != nil {
		return false
	}
	return true
}
