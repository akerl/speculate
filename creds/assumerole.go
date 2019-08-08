package creds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
)

type AssumeRoleOptions struct {
	RoleName    string
	AccountID   string
	SessionName string
	Policy      string
	Lifetime    int64
	UseMfa      bool
	MfaCode     string
	MfaPrompt   MfaPrompt
}

func (c Creds) AssumeRole(options AssumeRoleOptions) (Creds, error) {
	if options.RoleName == "" {
		return Creds{}, fmt.Errorf("role name cannot be empty")
	}

	if options.AccountID == "" {
		identity, err := c.identity()
		if err != nil {
			return Creds{}, err
		}
		options.AccountID = *identity.Account
	}

	partition, err := c.partition()
	if err != nil {
		return Creds{}, err
	}
	arn := fmt.Sprintf(
		"arn:%s:iam::%s:role/%s",
		partition,
		options.AccountID,
		options.RoleName,
	)

	// TODO: add validation for lifetime (between 900 and 3600 or 0)
	params := &sts.AssumeRoleInput{
		RoleArn:         &arn,
		RoleSessionName: &options.SessionName,
		DurationSeconds: &options.Lifetime,
		Policy:          &options.Policy,
	}

	tokenCode, serialNumber, err := c.handleMfa(
		options.UseMfa,
		options.MfaCode,
		options.MfaPrompt,
	)
	if err != nil {
		return Creds{}, err
	}
	if tokenCode != "" {
		params.TokenCode = &tokenCode
		params.SerialNumber = &serialNumber
	}

	client, err := c.Client()
	if err != nil {
		return Creds{}, err
	}
	resp, err := client.AssumeRole(params)
	if err != nil {
		return Creds{}, err
	}

	newCreds, err := NewFromStsSdk(resp.Credentials)
	return newCreds, err
}
