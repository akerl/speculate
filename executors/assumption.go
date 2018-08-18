package executors

import (
	"fmt"

	"github.com/akerl/speculate/creds"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Assumption describes the parameters that result in a Role
type Assumption struct {
	roleName    string
	accountID   string
	sessionName string
	policy      string
	Lifetime
	Mfa
}

// Execute actions a role assumption object with creds from the environment
func (a *Assumption) Execute() (creds.Creds, error) {
	return a.ExecuteWithCreds(creds.Creds{})
}

// ExecuteWithCreds actions a role assumption with provided creds
func (a *Assumption) ExecuteWithCreds(c creds.Creds) (creds.Creds, error) {
	newCreds := creds.Creds{}

	roleName, err := a.GetRoleName()
	if err != nil {
		return newCreds, err
	}

	accountID, err := a.GetAccountID()
	if err != nil {
		return newCreds, err
	}

	arn, err := c.NextRoleArn(roleName, accountID)
	if err != nil {
		return newCreds, err
	}
	sessionName, err := a.GetSessionName()
	if err != nil {
		return newCreds, err
	}
	lifetime, err := a.GetLifetime()
	if err != nil {
		return newCreds, err
	}

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(lifetime),
	}

	err = a.configureMfa(params)
	if err != nil {
		return newCreds, err
	}

	policy, err := a.GetPolicy()
	if err != nil {
		return newCreds, err
	}
	if policy != "" {
		params.Policy = aws.String(policy)
	}

	client := c.Client()
	resp, err := client.AssumeRole(params)
	if err != nil {
		return newCreds, err
	}

	newCreds, err = creds.NewFromStsSdk(resp.Credentials)
	return newCreds, err
}

// SetAccountID sets the target account ID
func (a *Assumption) SetAccountID(val string) error {
	if val == "" || accountIDRegex.MatchString(val) {
		a.accountID = val
		return nil
	}
	return fmt.Errorf("Account ID is malformed: %s", val)
}

// SetRoleName sets the target role name
func (a *Assumption) SetRoleName(val string) error {
	if val == "" || iamEntityRegex.MatchString(val) {
		a.roleName = val
		return nil
	}
	return fmt.Errorf("Role name is malformed: %s", val)
}

// SetSessionName sets the target session name
func (a *Assumption) SetSessionName(val string) error {
	if val == "" || iamEntityRegex.MatchString(val) {
		a.sessionName = val
		return nil
	}
	return fmt.Errorf("Session name is malformed: %s", val)
}

// SetPolicy sets the new IAM policy
func (a *Assumption) SetPolicy(val string) error {
	a.policy = val
	return nil
}

// GetAccountID gets the target account ID
func (a *Assumption) GetAccountID() (string, error) {
	if a.accountID == "" {
		c := creds.Creds{}
		var err error
		a.accountID, err = c.AccountID()
		if err != nil {
			return "", err
		}
	}
	return a.accountID, nil
}

// GetRoleName gets the target role name
func (a *Assumption) GetRoleName() (string, error) {
	return a.roleName, nil
}

// GetSessionName gets the target session name
func (a *Assumption) GetSessionName() (string, error) {
	if a.sessionName == "" {
		c := creds.Creds{}
		var err error
		a.sessionName, err = c.SessionName()
		if err != nil {
			return "", err
		}
	}
	return a.sessionName, nil
}

// GetPolicy gets the new IAM policy
func (a *Assumption) GetPolicy() (string, error) {
	return a.policy, nil
}
