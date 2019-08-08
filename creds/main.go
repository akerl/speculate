package creds

import (
	"fmt"
	"strings"

	"github.com/akerl/timber/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var logger = log.NewLogger("speculate")

// Creds defines a set of AWS credentials
type Creds struct {
	AccessKey, SecretKey, SessionToken, Region string
}

var namespaces = map[string]string{
	"aws":        "aws.amazon",
	"aws-us-gov": "amazonaws-us-gov",
}

// Session returns an AWS SDK session suitable for making API clients
func (c Creds) Session() (*session.Session, error) {
	logger.InfoMsg("creating new session")
	config := aws.NewConfig()
	config.WithCredentialsChainVerboseErrors(true)
	if c.AccessKey != "" {
		logger.InfoMsg(fmt.Sprintf("setting session credentials to %s", c.AccessKey))
		config.WithCredentials(c.ToSdk())
	}
	if c.Region != "" {
		logger.InfoMsg(fmt.Sprintf("setting session region to %s", c.Region))
		config.WithRegion(c.Region)
	}
	return session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	})
}

// Client returns an AWS STS client for these creds
func (c Creds) Client() (*sts.STS, error) {
	logger.InfoMsg("creating new STS client")
	session, err := c.Session()
	if err != nil {
		return nil, err
	}
	return sts.New(session), nil
}

func (c Creds) identity() (*sts.GetCallerIdentityOutput, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	logger.InfoMsg("looking up identity")
	return client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

func (c Creds) partition() (string, error) {
	logger.InfoMsg("looking up partition")
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	pieces := strings.Split(*identity.Arn, ":")
	logger.InfoMsg(fmt.Sprintf("found partition: %s", pieces[1]))
	return pieces[1], nil
}

func (c Creds) namespace() (string, error) {
	logger.InfoMsg("looking up namespace")
	partition, err := c.partition()
	if err != nil {
		return "", err
	}
	result, ok := namespaces[partition]
	if ok {
		logger.InfoMsg(fmt.Sprintf("found namespace: %s", result))
		return result, nil
	}
	return "", fmt.Errorf("unknown partition: %s", partition)
}

// AccountID returns the user's account ID
func (c Creds) AccountID() (string, error) {
	logger.InfoMsg("looking up accountID")
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	logger.InfoMsg(fmt.Sprintf("found accountID: %s", *identity.Account))
	return *identity.Account, nil
}

// MfaArn returns the user's virtual MFA token ARN
func (c Creds) MfaArn() (string, error) {
	logger.InfoMsg("looking up mfa arn")
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	if !strings.Contains(*identity.Arn, ":user/") {
		return "", fmt.Errorf("failed to parse MFA ARN for non-user: %s", *identity.Arn)
	}
	mfaArn := strings.Replace(*identity.Arn, ":user/", ":mfa/", 1)
	logger.InfoMsg(fmt.Sprintf("found mfa arn: %s", mfaArn))
	return mfaArn, nil
}

// UserName returns the current user name
func (c Creds) UserName() (string, error) {
	logger.InfoMsg("looking up current user name")
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	arnChunks := strings.Split(*identity.Arn, "/")
	userName := arnChunks[len(arnChunks)-1]
	logger.InfoMsg(fmt.Sprintf("found user name: %s", userName))
	return userName, nil
}
