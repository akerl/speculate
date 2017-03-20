package utils

import (
	"github.com/spf13/cobra"
)

// AddMfaFlags adds the flags for MFA interaction to a CLI command
func AddMfaFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	cmd.Flags().String("mfacode", "", "Code to use for MFA")
}
