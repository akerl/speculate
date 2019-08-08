package creds

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Creds defines a set of AWS credentials
type Creds struct {
	AccessKey, SecretKey, SessionToken, Region string
}

var namespaces = map[string]string{
	"aws":        "aws.amazon",
	"aws-us-gov": "amazonaws-us-gov",
}

// Client returns an AWS STS client for these creds
func (c Creds) Client() (*sts.STS, error) {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	if c.AccessKey != "" {
		config.WithCredentials(c.ToSdk())
	}
	if c.Region != "" {
		config.WithRegion(c.Region)
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}
	return sts.New(sess), nil
}

func (c Creds) identity() (*sts.GetCallerIdentityOutput, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	return client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (c Creds) partition() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	pieces := strings.Split(*identity.Arn, ":")
	return pieces[1], nil
}

func (c Creds) namespace() (string, error) {
	partition, err := c.partition()
	if err != nil {
		return "", err
	}
	result, ok := namespaces[partition]
	if ok {
		return result, nil
	}
	return "", fmt.Errorf("unknown partition: %s", partition)
}

// AccountID returns the user's account ID
func (c Creds) AccountID() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	return *identity.Account, nil
}

// MfaArn returns the user's virtual MFA token ARN
func (c Creds) MfaArn() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	if !strings.Contains(*identity.Arn, ":user/") {
		return "", fmt.Errorf("failed to parse MFA ARN for non-user: %s", *identity.Arn)
	}
	mfaArn := strings.Replace(*identity.Arn, ":user/", ":mfa/", 1)
	return mfaArn, nil
}

// SessionName returns the default session name
func (c Creds) SessionName() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	arnChunks := strings.Split(*identity.Arn, "/")
	oldName := arnChunks[len(arnChunks)-1]
	return oldName, nil
}
