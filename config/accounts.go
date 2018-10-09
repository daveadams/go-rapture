package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AccountMap map[string]string

func AccountsFilename() string {
	return filepath.Join(ConfigDir(), "accounts.json")
}

func LoadAccounts() (AccountMap, error) {
	rv := AccountMap{}
	fn := AccountsFilename()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no accounts file, return empty map
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

func WriteAccounts(acctMap AccountMap) error {
	bytes, err := json.Marshal(acctMap)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(AccountsFilename(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
