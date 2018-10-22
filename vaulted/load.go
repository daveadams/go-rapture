package vaulted

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/daveadams/go-rapture/log"
)

func LoadVault(name string) (map[string]string, error) {
	log.Tracef("vaulted: LoadVault(name='%s')", name)

	rv := map[string]string{}

	out, err := exec.Command("vaulted", "env", "--format", "json", name).Output()
	if err != nil {
		return nil, fmt.Errorf("Could not load variables from Vault: %s", err)
	}

	err = json.Unmarshal(out, &rv)
	if err != nil {
		return nil, err
	}

	return rv, nil
}
