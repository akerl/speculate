package utils

import (
	"fmt"

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
	for _, key := range []string{"account", "session"} {
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

// Execute actions the signin object
func (s *Signin) Execute() (Creds, error) {
	return Creds{}, nil
}
