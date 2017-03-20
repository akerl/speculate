package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
)

func RoleNameParse(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("No role name provided")
	}
	return args[0], nil
}

func RoleArn(role string) (string, error) {
	var params *sts.GetCallerIdentityInput
	resp, err := StsSession.GetCallerIdentity(params)
	if err != nil {
		return "", err
	}
	accountid := *resp.Account
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountid, role)
	return arn, nil
}
