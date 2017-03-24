package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// AddMfaFlags adds the flags for MFA interaction to a CLI command
func AddMfaFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("mfa", "m", false, "Use MFA when assuming role")
	cmd.Flags().String("mfacode", "", "Code to use for MFA")
}

func promptForMfa() (string, error) {
	mfaReader := bufio.NewReader(os.Stdin)
	fmt.Print("MFA Code: ")
	mfa, err := mfaReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	mfa = strings.TrimRight(mfa, "\n")
	return mfa, nil
}
