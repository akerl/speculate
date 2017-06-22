package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// Signin describes the parameters to perform GetSigninToken
type Signin struct {
	Lifetime
	Mfa
}

// ParseFlags for signin object
func (s *Signin) ParseFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	for _, key := range []string{"account", "session", "policy"} {
		val, err := flags.GetString(key)
		if err != nil {
			return err
		}
		if val != "" {
			return fmt.Errorf("%s is not valid without a role name", key)
		}
	}
	err = s.parseLifetimeFlags(cmd)
	if err != nil {
		return err
	}
	err = s.parseMfaFlags(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Execute actions the signin object with creds from the environment
func (s *Signin) Execute() (Creds, error) {
	return s.ExecuteWithCreds(Creds{})
}

// ExecuteWithCreds actions the signin object with the provided creds
func (s *Signin) ExecuteWithCreds(c Creds) (Creds, error) {
	newCreds := Creds{}

	if s.LifetimeInt == 0 {
		s.LifetimeInt = 3600
	}

	params := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(s.LifetimeInt),
	}
	if err := s.configureMfa(params); err != nil {
		return newCreds, err
	}

	client := API.ClientWithCreds(c)
	resp, err := client.GetSessionToken(params)
	if err != nil {
		return newCreds, err
	}

	err = newCreds.NewFromStsSdk(resp.Credentials)
	return newCreds, err
}
