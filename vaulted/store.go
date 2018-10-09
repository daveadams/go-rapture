package vaulted

import (
	vendor "github.com/miquella/vaulted/lib"
)

func New() vendor.Store {
	return vendor.New(&NoopSteward{})
}
