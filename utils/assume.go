package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func AddAssumeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account ID to assume role on (defaults to source account")
	AddMfaFlags(cmd)
}

func AssumeRole(role string, flags *pflag.FlagSet) Role {
	new_role := Role{Access_key: "test"}
	return new_role
}
