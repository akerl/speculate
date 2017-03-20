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
