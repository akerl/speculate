package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var sess = session.Must(session.NewSession(&aws.Config{
	CredentialsChainVerboseErrors: aws.Bool(true),
}))

// StsSession is a shared API client for talking to STS
var StsSession = sts.New(sess)

var stsIdentityCache = make(map[string]string)

// StsIdentity is a shared cache for the origin AWS identity
func StsIdentity() (map[string]string, error) {
	if stsIdentityCache["Account"] != "" {
		return stsIdentityCache, nil
	}
	var params *sts.GetCallerIdentityInput
	resp, err := StsSession.GetCallerIdentity(params)
	if err != nil {
		return stsIdentityCache, err
	}
	stsIdentityCache["Account"] = *resp.Account
	stsIdentityCache["Arn"] = *resp.Arn
	stsIdentityCache["UserId"] = *resp.UserId
	return stsIdentityCache, nil
}
