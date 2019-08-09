package main

var RaptureVersion = "2.0.0"

func CommandVersion(cmd string, args []string) int {
	shgen.Echof("rapture v%s", RaptureVersion)
	return 0
}
