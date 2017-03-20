package utils

import (
	"fmt"
	"strings"
)

func RoleNameParse(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("No role name provided")
	}
	return args[0], nil
}

func RoleArn(role string, account_id string) (string, error) {
	if account_id == "" {
		identity, err := StsIdentity()
		if err != nil {
			return "", err
		}
		account_id = identity["Account"]
	}
	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", account_id, role)
	return arn, nil
}

func SessionName(session_name string) (string, error) {
	if session_name != "" {
		return session_name, nil
	}
	identity, err := StsIdentity()
	if err != nil {
		return "", err
	}
	arn_chunks := strings.Split(identity["Arn"], "/")
	old_name := arn_chunks[len(arn_chunks)-1]
	return old_name, nil
}
