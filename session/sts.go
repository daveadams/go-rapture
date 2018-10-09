package session

import (
	"errors"
	"fmt"
	//"time"

	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/daveadams/go-rapture/config"
	"github.com/daveadams/go-rapture/log"
)

var (
	ErrBaseCredsExpired = errors.New("The base credentials are expired. Please re-initialize Rapture.")
	ErrBaseCredsInvalid = errors.New("The base credentials are missing or invalid. Please re-initialize Rapture.")
)

func (c *Credentials) stsClient() (*sts.STS, error) {
	log.Trace("session: Credentials.stsClient()")

	config := &aws.Config{Region: aws.String(config.GetConfig().Region)}
	config.Credentials = awscreds.NewStaticCredentials(c.ID, c.Secret, c.Token)
	if sess, err := awssession.NewSession(config); err != nil {
		return nil, err
	} else {
		return sts.New(sess), nil
	}
}

func (s *RaptureSession) stsSessionName() string {
	log.Trace("session: RaptureSession.stsSessionName()")

	if id := config.GetConfig().Identifier; id != "" {
		return fmt.Sprintf("rapture-%s", id)
	} else {
		return fmt.Sprintf("rapture-%s", s.ID)
	}
}

func (s *RaptureSession) AssumeRole(arn string) (*Credentials, error) {
	log.Tracef("session: RaptureSession.AssumeRole(arn='%s')", arn)

	if !s.BaseCreds.Valid() {
		return nil, ErrBaseCredsInvalid
	}

	if s.BaseCreds.NearExpiration() {
		log.Debug("Base credentials need to be renewed")
		return nil, ErrBaseCredsExpired
	}

	client, err := s.BaseCreds.stsClient()
	if err != nil {
		return nil, err
	}

	args := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(s.stsSessionName()),
		DurationSeconds: aws.Int64(config.GetConfig().SessionDuration),
	}
	if resp, err := client.AssumeRole(args); err != nil {
		return nil, err
	} else {
		return &Credentials{
			ID:        *resp.Credentials.AccessKeyId,
			Secret:    *resp.Credentials.SecretAccessKey,
			Token:     *resp.Credentials.SessionToken,
			ExpiresAt: *resp.Credentials.Expiration,
		}, nil
	}
}
