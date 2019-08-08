package creds

import (
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetSessionTokenOptions defines the available parameters for session tokens
type GetSessionTokenOptions struct {
	Lifetime  int64
	UseMfa    bool
	MfaCode   string
	MfaPrompt MfaPrompt
}

// GetSessionToken executes an AWS session token request
func (c Creds) GetSessionToken(options GetSessionTokenOptions) (Creds, error) {
	logger.InfoMsg("getting session token")
	logger.DebugMsg(fmt.Sprintf("getsessiontoken parameters: %+v", options))

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
	logger.InfoMsg("running getsessiontoken api call")
	resp, err := client.GetSessionToken(params)
	if err != nil {
		return Creds{}, err
	}

	newCreds, err := NewFromStsSdk(resp.Credentials)
	return newCreds, err
}
