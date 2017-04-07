package utils

import (
	"github.com/spf13/cobra"
)

// Signin describes the parameters to perform GetSigninToken
type Signin struct {
	lifetime int64
	Mfa
}

// ParseFlags for signin object
func (s *Signin) ParseFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	for _, key := range []string{"account", "session"} {
		val, err := flags.GetString(key)
		if err != nil {
			return err
		}
		if val != "" {
			return fmt.Errorf("%s is not valid without a role name", key)
		}
	}
	s.lifetime, err = flags.GetInt64("lifetime")
	if err != nil {
		return err
	}
	if a.lifetime < 900 || a.lifetime > 3600 {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", a.lifetime)
	}
	a.policy, err = flags.GetString("policy")
	if err != nil {
		return err
	}
	err := s.parseMfaFlags(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Execute actions the signin object
func (s *Signin) Execute() (Creds, error) {
	return Creds{}, nil
}
