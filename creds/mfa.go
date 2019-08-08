package creds

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MfaPrompt defines an object which recieves an Mfa ARN and returns an Mfa code
type MfaPrompt interface {
	Prompt(string) (string, error)
}

func (c Creds) handleMfa(useMfa bool, mfaCode string, mfaPrompt MfaPrompt) (string, string, error) {
	logger.InfoMsg("handling mfa options")

	// TODO: add validation for mfa code (6 digit int)
	if !useMfa && mfaCode == "" {
		logger.InfoMsg("mfa is disabled")
		return "", "", nil
	}
	mfaSerial, err := c.MfaArn()
	if err != nil {
		return "", "", err
	}
	if mfaCode != "" {
		logger.InfoMsg("mfa code already provided")
		return mfaCode, mfaSerial, nil
	}
	if mfaPrompt == nil {
		logger.InfoMsg("using default mfa prompt")
		mfaPrompt = &DefaultMfaPrompt{}
	}
	logger.InfoMsg("prompting for mfa code")
	mfaCode, err = mfaPrompt.Prompt(mfaSerial)
	if err != nil {
		return "", "", err
	}
	return mfaCode, mfaSerial, nil
}

// DefaultMfaPrompt defines the standard CLI-based MFA prompt
type DefaultMfaPrompt struct{}

// Prompt asks the user for their MFA token
func (p *DefaultMfaPrompt) Prompt(_ string) (string, error) {
	mfaReader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "MFA Code: ")
	mfa, err := mfaReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(mfa), nil
}
