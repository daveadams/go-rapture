package session

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/daveadams/go-rapture/log"
)

const (
	CacheDirMode  = 0755
	CacheFileMode = 0600

	BaseCredentialsArn = ":base:"
)

type CachedCredentials struct {
	RoleArn        string          `json:"role_arn"`
	EncryptedCreds string          `json:"credentials"`
	ExpiresAt      time.Time       `json:"expires_at,omitempty"`
	Creds          *Credentials    `json:"-"`
	sess           *RaptureSession `json:"-"`
}

// fetch cached credentials if they are valid and not yet expired
// otherwise, try to assume the role
// base credentials are stored with the sentinel role arn of BaseCredentialsArn
func (s *RaptureSession) CredentialsForRole(arn string) (*CachedCredentials, error) {
	return s.GetCredentialsForRole(arn, false)
}

// fetch cached credentials if they are valid and not yet expired, and force_refresh is not on
// otherwise, try to assume the role
// base credentials are stored with the sentinel role arn of BaseCredentialsArn
func (s *RaptureSession) GetCredentialsForRole(arn string, forceRefresh bool) (*CachedCredentials, error) {
	log.Tracef("session: RaptureSession.GetCredentialsForRole(arn='%s', forceRefresh='%t')", arn, forceRefresh)

	rv := &CachedCredentials{
		RoleArn: arn,
		sess:    s,
	}

	if ok, err := rv.load(); ok {
		log.Debug("Loaded credentials from cache successfully")
		if !rv.Creds.NearExpiration() && !forceRefresh {
			log.Debug("These credentials are not expired!")
			return rv, nil
		}

		log.Debug("These credentials are expired or will soon, and need to be renewed!")
		if arn == BaseCredentialsArn {
			log.Debug("But they are base credentials, so we have to hand that back to the UI.")
			return nil, ErrBaseCredsExpired
		}
	} else {
		if err != nil {
			log.Debugf("ERROR: Failed to load credentials because: %s", err)
		} else {
			log.Debug("No cached credentials exist")
		}

		// if we are trying to load base credentials, give up now
		if arn == BaseCredentialsArn {
			return nil, err
		}
	}

	// attempt to assume the role using the base credentials
	if val, err := s.AssumeRole(arn); err != nil {
		return nil, err
	} else {
		rv.Creds = val

		if err := rv.save(); err != nil {
			return rv, err
		} else {
			return rv, nil
		}
	}
}

// generate filename from sha256 sum of the session salt and role ARN
func (cc *CachedCredentials) Filename() string {
	log.Trace("session: CachedCredentials.Filename()")

	hash := sha256.Sum256([]byte(cc.sess.Salt + cc.RoleArn))
	bytes := make([]byte, len(hash))
	copy(bytes, hash[:])

	return filepath.Join(cc.sess.CacheDir(), base64.RawURLEncoding.EncodeToString(bytes))
}

// encrypt creds in memory
func (cc *CachedCredentials) encrypt() error {
	log.Trace("session: CachedCredentials.encrypt()")

	if cc.Creds == nil {
		return fmt.Errorf("No credentials loaded for this role")
	}

	if enc, err := cc.sess.EncryptCredentials(cc.Creds); err != nil {
		return err
	} else {
		cc.EncryptedCreds = enc
		return nil
	}
}

// decrypt creds in memory
func (cc *CachedCredentials) decrypt() error {
	log.Trace("session: CachedCredentials.decrypt()")

	if cc.EncryptedCreds == "" {
		return fmt.Errorf("No encrypted credentials loaded for this role")
	}

	if dc, err := cc.sess.DecryptCredentials(cc.EncryptedCreds); err != nil {
		return err
	} else {
		cc.Creds = dc
		return nil
	}
}

// encrypt and save cached creds to disk
func (cc *CachedCredentials) save() error {
	log.Trace("session: CachedCredentials.save()")

	if err := cc.encrypt(); err != nil {
		return err
	}

	if cc.Creds != nil {
		cc.ExpiresAt = cc.Creds.ExpiresAt
	}

	fn := cc.Filename()
	if err := os.MkdirAll(filepath.Dir(fn), CacheDirMode); err != nil {
		return err
	}

	bytes, err := json.Marshal(cc)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fn, bytes, CacheFileMode); err != nil {
		return err
	}

	return nil
}

// decrypt and load from disk
func (cc *CachedCredentials) load() (bool, error) {
	log.Trace("session: CachedCredentials.load()")

	fn := cc.Filename()

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// no cache file, return false, but no error
		return false, nil
	} else {
		bytes, err := ioutil.ReadFile(fn)
		if err != nil {
			return false, err
		}

		// unmarshal into this object
		err = json.Unmarshal(bytes, cc)
		if err != nil {
			return false, err
		}
	}

	// finally, try to decrypt the credentials
	if err := cc.decrypt(); err != nil {
		return false, err
	}

	return true, nil
}
