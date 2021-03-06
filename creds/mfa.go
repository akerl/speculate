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

// WritableMfaPrompt defines an MFA Prompt which can store a new secret
type WritableMfaPrompt interface {
	Store(string, string) error
	RetryText(string) string
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
type DefaultMfaPrompt struct {
	PromptTextFunc func(string) string
}

func defaultPromptTextFunc(_ string) string {
	return "MFA Code: "
}

// Prompt asks the user for their MFA token
func (p *DefaultMfaPrompt) Prompt(arn string) (string, error) {
	mfaReader := bufio.NewReader(os.Stdin)
	pf := p.PromptTextFunc
	if pf == nil {
		pf = defaultPromptTextFunc
	}
	fmt.Fprint(os.Stderr, pf(arn))
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
		logger.InfoMsgf("failed backend lookup %d with %s", index, err)
	}
	return "", fmt.Errorf("all backends failed to return mfa code")
}

// Store attempts to store the Mfa in a backend
func (m *MultiMfaPrompt) Store(arn, seed string) error {
	logger.InfoMsgf("attempting to store %s in multi mfa object", arn)
	writer := m.getWriter()
	if writer == nil {
		return fmt.Errorf("no writable mfa backends found")
	}
	return writer.Store(arn, seed)
}

// RetryText returns helper text for retrying a failed mfa storage
func (m *MultiMfaPrompt) RetryText(arn string) string {
	writer := m.getWriter()
	if writer == nil {
		return ""
	}
	return writer.RetryText(arn)
}

func (m *MultiMfaPrompt) getWriter() WritableMfaPrompt {
	for index, item := range m.Backends {
		writer, ok := item.(WritableMfaPrompt)
		if ok {
			logger.InfoMsgf("found mfa writer with index %d", index)
			return writer
		}
	}
	return nil
}
