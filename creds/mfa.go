package creds

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	mfaCodeRegexString = `^\d{6}$`
)

var mfaCodeRegex = regexp.MustCompile(mfaCodeRegexString)

// MfaPrompt defines an object which recieves an Mfa ARN and returns an Mfa code
type MfaPrompt interface {
	Prompt(string) (string, error)
}

// revive:disable-next-line:flag-parameter
func (c Creds) handleMfa(useMfa bool, mfaCode string, mfaPrompt MfaPrompt) (string, string, error) {
	logger.InfoMsg("handling mfa options")

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
		err := validateMfa(mfaCode)
		if err != nil {
			return "", "", err
		}

		return mfaCode, mfaSerial, nil
	}
	if mfaPrompt == nil {
		logger.InfoMsg("using default mfa prompt")
		// revive:disable-next-line:modifies-parameter
		mfaPrompt = &DefaultMfaPrompt{}
	}
	logger.InfoMsg("prompting for mfa code")
	// revive:disable-next-line:modifies-parameter
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
	trimmed := strings.TrimSpace(mfa)
	err = validateMfa(trimmed)
	if err != nil {
		return "", err
	}
	return trimmed, nil
}

func validateMfa(code string) error {
	logger.InfoMsg("validating mfa")
	if mfaCodeRegex.MatchString(code) {
		return nil
	}
	return fmt.Errorf("provided mfa code does not match the necessary format")
}

// MultiMfaPrompt allows a slice of sequential backends to check for Mfa
type MultiMfaPrompt struct {
	Backends []MfaPrompt
}

// Prompt iterates through the backends to find an Mfa code
func (m *MultiMfaPrompt) Prompt(arn string) (string, error) {
	logger.InfoMsgf("looking up %s in multi mfa prompt", arn)
	for index, item := range m.Backends {
		res, err := item.Prompt(arn)
		if err == nil {
			return res, nil
		}
		logger.InfoMsg("failed backend lookup %d with %s", index, err)
	}
	return "", fmt.Errorf("all backends failed to return mfa code")
}
