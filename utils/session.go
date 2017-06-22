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
	// TODO: implement caching of client and identity
}

// Client returns an API client for the STS API
func (a api) Client() *sts.STS {
	return a.client(Creds{})
}

// ClientWithCreds returns an API client using the provided creds
func (a api) ClientWithCreds(c Creds) *sts.STS {
	return a.client(c)
}

func (a api) client(c Creds) *sts.STS {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	if c.AccessKey != "" {
		config.WithCredentials(c.ToSdk())
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sts.New(sess)
}

// Identity returns the current user's information
func (a api) Identity() (*sts.GetCallerIdentityOutput, error) {
	params := &sts.GetCallerIdentityInput{}
	client := a.Client()
	return client.GetCallerIdentity(params)
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
