package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Lifetime object encapsulates the setup of session duration
type Lifetime struct {
	LifetimeInt int64
}

func (l *Lifetime) parseLifetimeFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	l.LifetimeInt, err = flags.GetInt64("lifetime")
	if err != nil {
		return err
	}
	if l.LifetimeInt < 900 || l.LifetimeInt > 3600 {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", l.LifetimeInt)
	}
	return nil
}
