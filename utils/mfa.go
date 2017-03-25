package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
