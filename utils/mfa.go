package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// Mfa object encapsulates the setup of MFA for API calls
type Mfa struct {
	useMfa  bool
	mfaCode string
}

func (m *Mfa) parseMfaFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	m.useMfa, err = flags.GetBool("mfa")
	if err != nil {
		return err
	}
	m.mfaCode, err = flags.GetString("mfacode")
	if err != nil {
		return err
	}
	return nil
}

func promptForMfa() (string, error) {
	mfaReader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "MFA Code: ")
	mfa, err := mfaReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	mfa = strings.TrimRight(mfa, "\n")
	return mfa, nil
}

func (m *Mfa) configureMfa(params *sts.AssumeRoleInput) error {
	if !m.useMfa && m.mfaCode == "" {
		return nil
	}
	serialNumber, err := API.MfaArn()
	if err != nil {
		return err
	}
	params.SerialNumber = aws.String(serialNumber)
	if m.mfaCode == "" {
		m.mfaCode, err = promptForMfa()
		if err != nil {
			return err
		}
	}
	params.TokenCode = aws.String(m.mfaCode)
	return nil
}
