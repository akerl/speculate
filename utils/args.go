package utils

import (
	"fmt"
	"strings"
)

// RoleNameParse returns the role name or an error
func RoleNameParse(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("No role name provided")
	}
	return args[0], nil
}

// RoleArn returns the new role ARN
func roleArn(role string, accountID string) (string, error) {
	if accountID == "" {
		identity, err := StsIdentity()
		if err != nil {
			return "", err
		}
		accountID = identity["Account"]
	}
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, role)
	return arn, nil
}

// mfaArn converts a User's ARN into their MFA ARN
func mfaArn() (string, error) {
	identity, err := StsIdentity()
	if err != nil {
		return "", err
	}
	userArn := identity["Arn"]
	if strings.Index(userArn, ":user/") == -1 {
		return "", fmt.Errorf("Failed to parse MFA ARN for non-user: %s", userArn)
	}
	mfaArn := strings.Replace(userArn, ":user/", ":mfa/", 1)
	return mfaArn, nil
}

// SessionName returns a new session name based on user input or identity
func sessionName(sessionName string) (string, error) {
	if sessionName != "" {
		return sessionName, nil
	}
	identity, err := StsIdentity()
	if err != nil {
		return "", err
	}
	arnChunks := strings.Split(identity["Arn"], "/")
	oldName := arnChunks[len(arnChunks)-1]
	return oldName, nil
}
