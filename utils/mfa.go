package utils

import (
	"github.com/spf13/cobra"
)

func AddMfaFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	cmd.Flags().String("mfacode", "", "Code to use for MFA")
}
