package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

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

func configureMfa(a *Assumption, params *sts.AssumeRoleInput) error {
	if !a.UseMfa && a.MfaCode == "" {
		return nil
	}
	serialNumber, err := API.MfaArn()
	if err != nil {
		return err
	}
	params.SerialNumber = aws.String(serialNumber)
	if a.MfaCode == "" {
		a.MfaCode, err = promptForMfa()
		if err != nil {
			return err
		}
	}
	params.TokenCode = aws.String(a.MfaCode)
	return nil
}
