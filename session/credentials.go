package session

import (
	"os"
	"time"

	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/shellgen"
)

type Credentials struct {
	ID        string    `json:"id"`
	Secret    string    `json:"secret"`
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

const (
	// five minutes is too close for comfort
	ExpirationWindow = time.Duration(5) * time.Minute
)

func ReadCredentialsFromEnvironment() *Credentials {
	log.Trace("session: ReadCredentialsFromEnvironment()")

	rv := &Credentials{
		ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Token:  os.Getenv("AWS_SESSION_TOKEN"),
	}

	if rv.Token == "" {
		rv.Token = os.Getenv("AWS_SECURITY_TOKEN")
	}

	if val, ok := os.LookupEnv("VAULTED_ENV_EXPIRATION"); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			rv.ExpiresAt = t
		} else {
			log.Debugf("WARNING: Could not parse '%s' as a time: %s", val, err)
		}
	} else {
		log.Debug("INFO: No expiration time was found in the environment")
	}

	return rv
}

// exports these credentials into the environment using the provided shell generator
func (c *Credentials) ExportToEnvironment(shgen shellgen.Generator) {
	log.Trace("session: Credentials.ExportToEnvironment(<shellgen.Generator>)")

	shgen.Export("AWS_ACCESS_KEY_ID", c.ID)
	shgen.Export("AWS_SECRET_ACCESS_KEY", c.Secret)
	if c.Token == "" {
		shgen.Unset("AWS_SESSION_TOKEN")
		shgen.Unset("AWS_SECURITY_TOKEN")
	} else {
		shgen.Export("AWS_SESSION_TOKEN", c.Token)
		shgen.Export("AWS_SECURITY_TOKEN", c.Token)
	}

	if c.ExpiresAt.IsZero() {
		shgen.Unset(AssumedRoleExpirationEnvVar)
	} else {
		shgen.Export(AssumedRoleExpirationEnvVar, c.ExpiresAt.UTC().Format(time.RFC3339))
	}
}

// checks if the credentials are valid/not expired
func (c *Credentials) Valid() bool {
	log.Trace("session: Credentials.Valid()")
	return c != nil && c.ID != "" && c.Secret != ""
}

// checks if the credentials are near expiration
func (c *Credentials) NearExpiration() bool {
	log.Trace("session: Credentials.NearExpiration()")
	if c.ExpiresAt.IsZero() {
		log.Debug("ExpiresAt is not set for these credentials")
		// if we don't know the expiration time, assume they don't expire
		// which is correct for static credentials
		return false
	}

	log.Debugf("These credentials will expire at %s", c.ExpiresAt.String())

	return time.Now().Add(ExpirationWindow).After(c.ExpiresAt)
}
