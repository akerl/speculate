package executors

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/akerl/speculate/creds"
)

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
	if val < 900 || val > 3600 {
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
	// TODO: Check if valid ARN
	m.mfaSerial = val
	return nil
}

// SetMfaCode sets the OTP for MFA
func (m *Mfa) SetMfaCode(val string) error {
	// TODO: Check if valid code
	m.mfaCode = val
	return nil
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
		c, err := creds.Creds{}
		if err != nil {
			return "", err
		}
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
		m.mfaCode = strings.TrimRight(mfa, "\n")
	}
	return m.mfaCode, nil
}
