package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var sess = session.Must(session.NewSession(&aws.Config{
	CredentialsChainVerboseErrors: aws.Bool(true),
}))
var StsSession = sts.New(sess)

var stsIdentityCache = make(map[string]string)

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
