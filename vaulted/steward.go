package vaulted

import (
	vendor "github.com/miquella/vaulted/lib"
)

// NoopSteward suitable for creating a vaulted Store and listing Vaults, but not for interaction
type NoopSteward struct{}

func (ns *NoopSteward) GetMFAToken(s string) (string, error) {
	return "", nil
}

func (ns *NoopSteward) GetPassword(op vendor.Operation, s string) (string, error) {
	return "", nil
}
