package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/log"
)

type RaptureConfig struct {
	Region          string `json:"region,omitempty"`
	Quiet           bool   `json:"quiet,omitempty"`
	Identifier      string `json:"identifier,omitempty"`
	SessionDuration int64  `json:"session_duration,omitempty"`
	DefaultVault    string `json:"default_vault,omitempty"`
}

var config *RaptureConfig

func ConfigFilename() string {
	log.Trace("config: ConfigFilename()")
	return filepath.Join(ConfigDir(), "config.json")
}

func DefaultConfig() *RaptureConfig {
	log.Trace("config: DefaultConfig()")
	return &RaptureConfig{
		Region:          "us-east-1",
		Quiet:           false,
		Identifier:      os.Getenv("USER"),
		SessionDuration: 3600,
		DefaultVault:    "default",
	}
}

func LoadConfig() (*RaptureConfig, error) {
	log.Trace("config: LoadConfig()")

	if config != nil {
		return config, nil
	}

	config := DefaultConfig()
	fn := ConfigFilename()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no roles file, return defaults
		return config, nil
	} else {
		bytes, err := ioutil.ReadFile(fn)
		if err != nil {
			return config, err
		}
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

// just return a config even if there's an error
func GetConfig() *RaptureConfig {
	log.Trace("config: GetConfig()")

	if c, err := LoadConfig(); err != nil {
		log.Debugf("Failed to load config: %s", err)
		return DefaultConfig()
	} else {
		return c
	}
}
