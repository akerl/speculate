package creds

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/sts"
)

// New initializes credentials from a map
func New(argCreds map[string]string) (Creds, error) {
	logger.InfoMsg(fmt.Sprintf(
		"creating new creds from map with access key %s",
		argCreds["AccessKey"],
	))
	required := []string{"AccessKey", "SecretKey", "SessionToken"}
	for _, key := range required {
		elem, ok := argCreds[key]
		if !ok || elem == "" {
			return Creds{}, fmt.Errorf("missing required key for Creds: %s", key)
		}
	}
	c := Creds{
		AccessKey:    argCreds["AccessKey"],
		SecretKey:    argCreds["SecretKey"],
		SessionToken: argCreds["SessionToken"],
		Region:       argCreds["Region"],
	}
	return c, nil
}

// NewFromStsSdk initializes a credential object from an AWS SDK Credentials object
func NewFromStsSdk(stsCreds *sts.Credentials) (Creds, error) {
	logger.InfoMsg("creating new creds from SDK")
	return New(map[string]string{
		"AccessKey":    *stsCreds.AccessKeyId,
		"SecretKey":    *stsCreds.SecretAccessKey,
		"SessionToken": *stsCreds.SessionToken,
	})
}

// NewFromEnv initializes credentials from the environment variables
func NewFromEnv() (Creds, error) {
	logger.InfoMsg("creating new creds from env")
	envCreds := make(map[string]string)
	for k, v := range Translations["envvar"] {
		if envCreds[v] == "" {
			envCreds[v] = os.Getenv(k)
		}
	}
	return New(envCreds)
}
