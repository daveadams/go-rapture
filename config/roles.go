package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type RoleMap map[string]string

func RolesFilename() string {
	return filepath.Join(ConfigDir(), "aliases.json")
}

func LoadRoles() (RoleMap, error) {
	rv := RoleMap{}
	fn := RolesFilename()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no roles file, return empty map
		return rv, nil
	} else {
		bytes, err := ioutil.ReadFile(fn)
		if err != nil {
			return rv, err
		}
		err = json.Unmarshal(bytes, &rv)
		if err != nil {
			return rv, err
		}
	}
	return rv, nil
}

func WriteRoles(roleMap RoleMap) error {
	bytes, err := json.Marshal(roleMap)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(RolesFilename(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
