package creds

import (
	"github.com/aws/aws-sdk-go/service/sts"
)

type GetSessionTokenOptions struct {
	Lifetime  int64
	UseMfa    bool
	MfaCode   string
	MfaPrompt MfaPrompt
}

func (c Creds) GetSessionToken(options GetSessionTokenOptions) (Creds, error) {
	// TODO: add validation for lifetime (between 900 and 3600 or 0)
	params := &sts.GetSessionTokenInput{
		DurationSeconds: &options.Lifetime,
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
	resp, err := client.GetSessionToken(params)
	if err != nil {
		return Creds{}, err
	}

	newCreds, err := NewFromStsSdk(resp.Credentials)
	return newCreds, err
}
