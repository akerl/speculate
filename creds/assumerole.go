package creds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sts"
)

// AssumeRoleOptions defines the available parameters for assuming roles
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

// AssumeRole executes an AWS role assumption
func (c Creds) AssumeRole(options AssumeRoleOptions) (Creds, error) {
	logger.InfoMsg("assuming role")
	logger.DebugMsgf("assumerole parameters: %+v", options)

	err := c.assumeRolePreflight(&options)
	if err != nil {
		return Creds{}, err
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
	logger.InfoMsgf("generated target arn: %s", arn)

	params := &sts.AssumeRoleInput{
		RoleArn: &arn,
	}

	if options.Policy != "" {
		params.Policy = &options.Policy
	}

	if options.Lifetime != 0 {
		params.DurationSeconds = &options.Lifetime
	}

	if options.SessionName != "" {
		params.RoleSessionName = &options.SessionName
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
	logger.InfoMsg("running assumerole api call")
	resp, err := client.AssumeRole(params)
	if err != nil {
		return Creds{}, err
	}

	newCreds, err := NewFromStsSdk(resp.Credentials)
	newCreds.Region = c.Region
	return newCreds, err
}

func (c Creds) assumeRolePreflight(options *AssumeRoleOptions) error {
	var err error

	if options.RoleName == "" {
		return fmt.Errorf("role name cannot be empty")
	}

	if options.AccountID == "" {
		options.AccountID, err = c.AccountID()
		if err != nil {
			return err
		}
	}

	if options.SessionName == "" {
		options.SessionName, err = c.UserName()
		if err != nil {
			return err
		}
	}

	options.Lifetime, err = validateLifetime(options.Lifetime, AssumeRoleLifetimeLimits)
	return err
}
