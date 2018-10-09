package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
	"github.com/daveadams/go-rapture/shellgen"
	"github.com/google/uuid"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KeySize   = 32
	NonceSize = 24
	SaltSize  = 16

	IDEnvVar   = "_rapture_session_id"
	KeyEnvVar  = "_rapture_session_key"
	SaltEnvVar = "_rapture_session_salt"

	AssumedRoleArnEnvVar   = "RAPTURE_ASSUMED_ROLE_ARN"
	AssumedRoleAliasEnvVar = "RAPTURE_ROLE"
)

type RaptureSession struct {
	ID               string
	key              []byte
	Salt             string
	BaseCreds        *Credentials
	AssumedRoleArn   string
	AssumedRoleAlias string
}

var (
	currentSession *RaptureSession
)

func newSessionID() string {
	log.Trace("session: newSessionID()")
	return uuid.New().String()
}

func isValidSessionID(id string) bool {
	log.Tracef("session: isValidSessionID(id='%s')", id)
	_, err := uuid.Parse(id)
	return err == nil
}

// generate a new encryption key for this session
func newSessionKey() ([]byte, error) {
	log.Trace("session: newSessionKey()")

	bytes := make([]byte, KeySize)
	if _, err := rand.Read(bytes); err != nil {
		return bytes, err
	}
	return bytes, nil
}

func parseEncodedKey(s string) ([]byte, error) {
	log.Tracef("session: parseEncodedKey(s='%s')", s)

	if bytes, err := base64.StdEncoding.DecodeString(s); err != nil {
		return bytes, err
	} else {
		return bytes, nil
	}
}

func isValidSessionKey(key []byte) bool {
	log.Tracef("session: isValidSessionKey(key=0x%s)", hex.EncodeToString(key))
	return len(key) == KeySize
}

// create a session-specific salt for encoding filenames for cached credentials
func newSessionSalt() (string, error) {
	log.Trace("session: newSessionSalt()")

	bytes := make([]byte, SaltSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func isValidSessionSalt(s string) bool {
	log.Tracef("session: isValidSessionSalt(s='%s')", s)
	return len(s) == base64.StdEncoding.EncodedLen(SaltSize)
}

// loads the current session from env vars if possible, or creates a new one if not
func CurrentSession() (*RaptureSession, bool, error) {
	log.Trace("session: CurrentSession()")

	if currentSession != nil {
		return currentSession, false, nil
	}

	new := false

	id, ok := os.LookupEnv(IDEnvVar)
	if !ok || id == "" {
		id = newSessionID()
		new = true
	} else {
		if !isValidSessionID(id) {
			return nil, false, fmt.Errorf("Your session ID is invalid. Please run 'rapture reset'.")
		}
	}

	var key []byte
	encodedKey, ok := os.LookupEnv(KeyEnvVar)
	if !ok || encodedKey == "" || new == true {
		if val, err := newSessionKey(); err != nil {
			return nil, false, fmt.Errorf("Could not generate a new session key: %s", err)
		} else {
			key = val
			new = true
		}
	} else {
		if val, err := parseEncodedKey(encodedKey); err != nil || !isValidSessionKey(val) {
			return nil, false, fmt.Errorf("Your session key is invalid. Please run 'rapture reset'.")
		} else {
			key = val
		}
	}

	salt, ok := os.LookupEnv(SaltEnvVar)
	if !ok || salt == "" || new == true {
		if val, err := newSessionSalt(); err != nil {
			return nil, false, fmt.Errorf("Could not generate a new session salt: %s", err)
		} else {
			salt = val
			new = true
		}
	} else {
		if !isValidSessionSalt(salt) {
			return nil, false, fmt.Errorf("Your session salt is invalid. Please run 'rapture reset'.")
		}
	}

	rv := &RaptureSession{
		ID:   id,
		key:  key,
		Salt: salt,
	}

	if new {
		rv.BaseCreds = ReadCredentialsFromEnvironment()
	} else {
		if err := rv.LoadBaseCredentials(); err != nil {
			return nil, new, err
		}
	}

	rv.AssumedRoleArn = os.Getenv(AssumedRoleArnEnvVar)
	rv.AssumedRoleAlias = os.Getenv(AssumedRoleAliasEnvVar)

	return rv, new, nil
}

func (s *RaptureSession) LoadBaseCredentials() error {
	log.Trace("session: RaptureSession.LoadBaseCredentials()")

	if bc, err := s.CredentialsForRole(BaseCredentialsArn); err != nil {
		return err
	} else {
		s.BaseCreds = bc.Creds
		return nil
	}
}

func (s *RaptureSession) SaveBaseCredentials() error {
	log.Trace("session: RaptureSession.SaveBaseCredentials()")

	cc := CachedCredentials{
		RoleArn: BaseCredentialsArn,
		sess:    s,
		Creds:   s.BaseCreds,
	}

	if err := cc.save(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *RaptureSession) EncodedKey() string {
	log.Trace("session: RaptureSession.EncodedKey()")
	return base64.StdEncoding.EncodeToString(s.key)
}

func (s *RaptureSession) Key() [KeySize]byte {
	log.Trace("session: RaptureSession.Key()")
	var rv [KeySize]byte
	copy(rv[:], s.key)
	return rv
}

func (s *RaptureSession) Save(g shellgen.Generator) error {
	log.Trace("session: RaptureSession.Save(<shellgen.Generator>)")
	g.Set(IDEnvVar, s.ID)
	g.Set(KeyEnvVar, s.EncodedKey())
	g.Set(SaltEnvVar, s.Salt)

	if s.AssumedRoleArn == "" {
		g.Unset(AssumedRoleArnEnvVar)
	} else {
		g.Export(AssumedRoleArnEnvVar, s.AssumedRoleArn)
	}

	if s.AssumedRoleAlias == "" {
		g.Unset(AssumedRoleAliasEnvVar)
	} else {
		g.Export(AssumedRoleArnEnvVar, s.AssumedRoleArn)
	}

	if err := s.SaveBaseCredentials(); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *RaptureSession) CacheDir() string {
	log.Trace("session: RaptureSession.CacheDir()")
	return filepath.Join(config.SessionsCacheDir(), s.ID)
}

// Decrypts a base64-encoded string prefixed with its own nonce as generated by Credentials.Encrypt
func (s *RaptureSession) DecryptCredentials(enc string) (*Credentials, error) {
	log.Tracef("session: RaptureSession.DecryptCredentials(enc='%s')", enc)

	key := s.Key()

	ciphertext, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	// nonce is the first 24 bytes of the ciphertext
	var nonce [24]byte
	copy(nonce[:], ciphertext[:24])

	plaintext, ok := secretbox.Open(nil, ciphertext[24:], &nonce, &key)
	if !ok {
		return nil, fmt.Errorf("Could not decrypt credentials!")
	}

	rv := Credentials{}
	err = json.Unmarshal(plaintext, &rv)
	if err != nil {
		return nil, err
	}

	return &rv, nil
}

// Converts the Credentials object to a base64-encoded encrypted string
func (s *RaptureSession) EncryptCredentials(c *Credentials) (string, error) {
	log.Trace("session: RaptureSession.EncryptCredentials(<*Credentials>)")

	message, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	key := s.Key()

	// use a new nonce for each encryption
	var nonce [NonceSize]byte
	rand.Read(nonce[:])

	// prefix the ciphertext with the nonce
	prefix := make([]byte, NonceSize)
	copy(prefix, nonce[:])

	ciphertext := secretbox.Seal(prefix, message, &nonce, &key)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
