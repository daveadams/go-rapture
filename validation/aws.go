package validation

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws/arn"
)

var (
	validAwsAccountId = regexp.MustCompile(`^[0-9]{12}$`)
)

func IsValidAwsAccountId(s string) bool {
	return validAwsAccountId.MatchString(s)
}

func IsValidIamRoleArn(s string) bool {
	if _, err := arn.Parse(s); err != nil {
		return false
	} else {
		return true
	}
}
