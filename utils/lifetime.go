package utils

import (
	"github.com/spf13/cobra"
)

// Lifetime object encapsulates the setup of session duration
type Lifetime struct {
	lifetime int64
}

func (l *Lifetime) parseLifetimeFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	l.lifetime, err = flags.GetInt64("lifetime")
	if err != nil {
		return err
	}
	if l.lifetime < 900 || l.lifetime > 3600 {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", a.lifetime)
	}
	return nil
}
