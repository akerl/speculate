package creds

import (
	"fmt"
	"strings"

	"github.com/akerl/timber/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/akerl/speculate/v2/version"
)

var logger = log.NewLogger("speculate")

// UserAgentItem defines an entry in the HTTP User Agent field
type UserAgentItem struct {
	Name, Version string
	Extra         []string
}

var speculateUserAgentItem = UserAgentItem{
	Name:    "speculate",
	Version: version.Version,
}

// Creds defines a set of AWS credentials
type Creds struct {
	AccessKey, SecretKey, SessionToken, Region string
	UserAgentItems                             []UserAgentItem
}

var namespaces = map[string]string{
	"aws":        "aws.amazon",
	"aws-us-gov": "amazonaws-us-gov",
}

func addUserAgentHandler(sess *session.Session, item UserAgentItem) {
	h := request.MakeAddToUserAgentHandler(item.Name, item.Version, item.Extra...)
	sess.Handlers.Build.PushBack(h)
}

func (c Creds) setUserAgent(sess *session.Session) {
	addUserAgentHandler(sess, speculateUserAgentItem)
	for _, x := range c.UserAgentItems {
		addUserAgentHandler(sess, x)
	}
}

// Session returns an AWS SDK session suitable for making API clients
func (c Creds) Session() (*session.Session, error) {
	logger.InfoMsg("creating new session")
	config := aws.NewConfig()
	config.WithCredentialsChainVerboseErrors(true)
	if c.AccessKey != "" {
		logger.InfoMsgf("setting session credentials to %s", c.AccessKey)
		config.WithCredentials(c.ToSdk())
	}
	if c.Region != "" {
		logger.InfoMsgf("setting session region to %s", c.Region)
		config.WithRegion(c.Region)
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return sess, err
	}
	c.setUserAgent(sess)
	return sess, nil
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
	logger.InfoMsgf("found partition: %s", pieces[1])
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
		logger.InfoMsgf("found namespace: %s", result)
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
	logger.InfoMsgf("found accountID: %s", *identity.Account)
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
	logger.InfoMsgf("found mfa arn: %s", mfaArn)
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
	logger.InfoMsgf("found user name: %s", userName)
	return userName, nil
}
