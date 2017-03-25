package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// API exposes the STS API and related info
var API = api{}

type api struct {
	client   *sts.STS
	identity *sts.GetCallerIdentityOutput
}

// Client returns an API client for the STS API
func (a api) Client() sts.STS {
	if a.client == nil {
		sess := session.Must(session.NewSession(&aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		}))
		a.client = sts.New(sess)
	}
	return *a.client
}

// Identity returns the current user's information
func (a api) Identity() (*sts.GetCallerIdentityOutput, error) {
	if a.identity == nil {
		params := &sts.GetCallerIdentityInput{}
		client := a.Client()
		resp, err := client.GetCallerIdentity(params)
		if err != nil {
			return nil, err
		}
		a.identity = resp
	}
	return a.identity, nil
}

// Partition returns the AWS partition for the user
func (a api) Partition() (string, error) {
	identity, err := a.Identity()
	if err != nil {
		return "", err
	}
	pieces := strings.Split(*identity.Arn, ":")
	if len(pieces) != 6 {
		return "", fmt.Errorf("Invalid ARN provided: %s", identity.Arn)
	}
	return pieces[1], nil
}

// MfaArn returns the user's virtual MFA token ARN
func (a api) MfaArn() (string, error) {
	identity, err := a.Identity()
	if err != nil {
		return "", err
	}
	if strings.Index(*identity.Arn, ":user/") == -1 {
		return "", fmt.Errorf("Failed to parse MFA ARN for non-user: %s", identity.Arn)
	}
	mfaArn := strings.Replace(*identity.Arn, ":user/", ":mfa/", 1)
	return mfaArn, nil
}

// SessionName returns the default session name
func (a api) SessionName() (string, error) {
	identity, err := a.Identity()
	if err != nil {
		return "", err
	}
	arnChunks := strings.Split(*identity.Arn, "/")
	oldName := arnChunks[len(arnChunks)-1]
	return oldName, nil
}

// RoleArn returns the new role's ARN
func (a api) RoleArn(role, accountID string) (string, error) {
	identity, err := a.Identity()
	if err != nil {
		return "", err
	}
	partition, err := a.Partition()
	if err != nil {
		return "", err
	}
	if accountID == "" {
		accountID = *identity.Account
	}
	arn := fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, accountID, role)
	return arn, nil
}
