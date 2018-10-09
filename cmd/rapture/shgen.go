package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/shellgen"
)

var shgen shellgen.Generator

func init() {
	shgen = shellgen.NewGenerator()
	if shgen == nil {
		fmt.Fprintf(os.Stderr, "ERROR: Your shell ('%s') is not supported.\n", filepath.Base(os.Getenv("SHELL")))
		os.Exit(1)
	}
}
