package main

import (
	"os"

	"github.com/daveadams/go-rapture/log"
	vaulted "github.com/miquella/vaulted/lib"
)

func PrintWhoami() int {
	log.Trace("[main] PrintWhoami()")

	creds := vaulted.AWSCredentials{
		ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Token:  os.Getenv("AWS_SESSION_TOKEN"),
	}

	if !creds.Valid() {
		shgen.ErrEcho("ERROR: No AWS credentials were found in your environment")
		return 1
	}

	arn, err := creds.GetCallerIdentity()
	if err != nil {
		shgen.ErrEcho("ERROR: Could not determine AWS identity. Are your credentials loaded?")
		return 1
	}

	shgen.Echo(arn.String())

	return 0
}

func CommandWhoami(cmd string, args []string) int {
	log.Tracef("[main] CommandWhoami(cmd='%s', args=%s)", cmd, args)
	return PrintWhoami()
}
