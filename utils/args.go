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
func RoleArn(role string, accountID string) (string, error) {
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

// SessionName returns a new session name based on user input or identity
func SessionName(sessionName string) (string, error) {
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
