package config

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	d := os.Getenv("RAPTURE_CONF_DIR")
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".rapture")
	}
	return d
}

func CacheDir() string {
	d, _ := os.UserCacheDir()
	if d == "" {
		d = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return filepath.Join(d, "rapture")
}

func SessionsCacheDir() string {
	return filepath.Join(CacheDir(), "sessions")
}
