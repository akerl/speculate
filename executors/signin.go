package executors

import (
	"fmt"

	"github.com/akerl/speculate/creds"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Signin describes the parameters to perform GetSigninToken
type Signin struct {
	Lifetime
	Mfa
}

// Execute actions the signin object with creds from the environment
func (s *Signin) Execute() (creds.Creds, error) {
	return s.ExecuteWithCreds(creds.Creds{})
}

// ExecuteWithCreds actions the signin object with the provided creds
func (s *Signin) ExecuteWithCreds(c creds.Creds) (creds.Creds, error) {
	newCreds := creds.Creds{}

	lifetime, err := s.GetLifetime()
	if err != nil {
		return newCreds, err
	}

	params := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(lifetime),
	}

	err = s.configureMfa(params)
	if err != nil {
		return newCreds, err
	}

	client := c.Client()
	resp, err := client.GetSessionToken(params)
	if err != nil {
		return newCreds, err
	}

	newCreds, err = creds.NewFromStsSdk(resp.Credentials)
	return newCreds, err
}

// SetAccountID sets the target account ID
func (s *Signin) SetAccountID(_ string) error {
	return fmt.Errorf("accountID is invalid for GetSessionToken")
}

// SetRoleName sets the target role name
func (s *Signin) SetRoleName(_ string) error {
	return fmt.Errorf("roleName is invalid for GetSessionToken")
}

// SetSessionName sets the target session name
func (s *Signin) SetSessionName(_ string) error {
	return fmt.Errorf("sessionName is invalid for GetSessionToken")
}

// SetPolicy sets the new IAM policy
func (s *Signin) SetPolicy(_ string) error {
	return fmt.Errorf("policy is invalid for GetSessionToken")
}

// GetAccountID gets the target account ID
func (s *Signin) GetAccountID() (string, error) {
	return "", fmt.Errorf("accountID is invalid for GetSessionToken")
}

// GetRoleName gets the target role name
func (s *Signin) GetRoleName() (string, error) {
	return "", fmt.Errorf("roleName is invalid for GetSessionToken")
}

// GetSessionName gets the target session name
func (s *Signin) GetSessionName() (string, error) {
	return "", fmt.Errorf("sessionName is invalid for GetSessionToken")
}

// GetPolicy gets the new IAM policy
func (s *Signin) GetPolicy() (string, error) {
	return "", fmt.Errorf("policy is invalid for GetSessionToken")
}
