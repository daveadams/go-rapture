package validation

import "regexp"

var (
	validAwsAccountId = regexp.MustCompile(`^[0-9]{12}$`)
	validIamRoleArn   = regexp.MustCompile(`^arn:aws:iam::[0-9]{12}:role/.+$`)
)

func IsValidAwsAccountId(s string) bool {
	return validAwsAccountId.MatchString(s)
}

func IsValidIamRoleArn(s string) bool {
	return validIamRoleArn.MatchString(s)
}
