package executors

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/akerl/speculate/creds"

	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	mfaCodeRegexString   = `^\d{6}$`
	mfaArnRegexString    = `^arn:aws(?:-gov)?:iam::\d+:mfa/[\w+=,.@-]+$`
	iamEntityRegexString = `^[\w+=,.@-]+$`
	accountIDRegexString = `^\d{12}$`
)

var mfaCodeRegex = regexp.MustCompile(mfaCodeRegexString)
var mfaArnRegex = regexp.MustCompile(mfaArnRegexString)
var iamEntityRegex = regexp.MustCompile(iamEntityRegexString)
var accountIDRegex = regexp.MustCompile(accountIDRegexString)

// Executor defines the interface for requesting a new set of AWS creds
type Executor interface {
	Execute() (creds.Creds, error)
	ExecuteWithCreds(creds.Creds) (creds.Creds, error)
	SetAccountID(string) error
	SetRoleName(string) error
	SetSessionName(string) error
	SetPolicy(string) error
	SetLifetime(int64) error
	SetMfa(bool) error
	SetMfaSerial(string) error
	SetMfaCode(string) error
	GetAccountID() (string, error)
	GetRoleName() (string, error)
	GetSessionName() (string, error)
	GetPolicy() (string, error)
	GetLifetime() (int64, error)
	GetMfa() (bool, error)
	GetMfaSerial() (string, error)
	GetMfaCode() (string, error)
}

// Lifetime object encapsulates the setup of session duration
type Lifetime struct {
	lifetimeInt int64
}

// SetLifetime allows setting the credential lifespan
func (l *Lifetime) SetLifetime(val int64) error {
	if val != 0 && (val < 900 || val > 3600) {
		return fmt.Errorf("Lifetime must be between 900 and 3600: %d", val)
	}
	l.lifetimeInt = val
	return nil
}

// GetLifetime returns the lifetime of the executor
func (l *Lifetime) GetLifetime() (int64, error) {
	if l.lifetimeInt == 0 {
		l.lifetimeInt = 3600
	}
	return l.lifetimeInt, nil
}

// Mfa object encapsulates the setup of MFA for API calls
type Mfa struct {
	useMfa    bool
	mfaSerial string
	mfaCode   string
}

// SetMfa sets whether MFA is used
func (m *Mfa) SetMfa(val bool) error {
	m.useMfa = val
	return nil
}

// SetMfaSerial sets the ARN of the MFA device
func (m *Mfa) SetMfaSerial(val string) error {
	if val == "" || mfaArnRegex.MatchString(val) {
		m.mfaSerial = val
		return nil
	}
	return fmt.Errorf("MFA Serial is malformed: %s", val)
}

// SetMfaCode sets the OTP for MFA
func (m *Mfa) SetMfaCode(val string) error {
	if val == "" || mfaCodeRegex.MatchString(val) {
		m.mfaCode = val
		return nil
	}
	return fmt.Errorf("MFA Code is malformed: %s", val)
}

// GetMfa returns if MFA will be used
func (m *Mfa) GetMfa() (bool, error) {
	if !m.useMfa && m.mfaCode != "" {
		m.useMfa = true
	}
	return m.useMfa, nil
}

// GetMfaSerial returns the ARN of the MFA device
func (m *Mfa) GetMfaSerial() (string, error) {
	if m.mfaSerial == "" {
		c := creds.Creds{}
		var err error
		m.mfaSerial, err = c.MfaArn()
		if err != nil {
			return "", err
		}
	}
	return m.mfaSerial, nil
}

// GetMfaCode returns the OTP to use
func (m *Mfa) GetMfaCode() (string, error) {
	if m.mfaCode == "" {
		mfaReader := bufio.NewReader(os.Stdin)
		fmt.Fprint(os.Stderr, "MFA Code: ")
		mfa, err := mfaReader.ReadString('\n')
		if err != nil {
			return "", err
		}
		m.mfaCode = strings.TrimSpace(mfa)
	}
	return m.mfaCode, nil
}

func (m *Mfa) configureMfa(paramsIface interface{}) error {
	useMfa, err := m.GetMfa()
	if err != nil {
		return err
	}
	if !useMfa {
		return nil
	}

	mfaCode, err := m.GetMfaCode()
	if err != nil {
		return err
	}
	mfaSerial, err := m.GetMfaSerial()
	if err != nil {
		return err
	}

	switch params := paramsIface.(type) {
	case *sts.AssumeRoleInput:
		params.TokenCode = &mfaCode
		params.SerialNumber = &mfaSerial
	case *sts.GetSessionTokenInput:
		params.TokenCode = &mfaCode
		params.SerialNumber = &mfaSerial
	default:
		return fmt.Errorf("Expected AssumeRoleInput or GetSessionTokenInput, received %T", params)
	}
	return nil
}
