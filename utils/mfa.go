package utils

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// Mfa object encapsulates the setup of MFA for API calls
type Mfa struct {
	UseMfa  bool
	MfaCode string
}

func (m *Mfa) parseMfaFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()
	var err error
	m.UseMfa, err = flags.GetBool("mfa")
	if err != nil {
		return err
	}
	m.MfaCode, err = flags.GetString("mfacode")
	if err != nil {
		return err
	}
	return nil
}

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

// Must be called with an aws Input struct with SerialNumber and SetTokenCode fields
func (m *Mfa) configureMfa(params interface{}) error {
	if !m.UseMfa && m.MfaCode == "" {
		return nil
	}

	stype := reflect.ValueOf(params).Elem()

	serialNumber, err := API.MfaArn()
	if err != nil {
		return err
	}
	paramSerialNumber := stype.FieldByName("SerialNumber")
	if !paramSerialNumber.IsValid() {
		return fmt.Errorf("configureMfa called with struct lacking SerialNumber field")
	}
	paramSerialNumber.Set(reflect.ValueOf(&serialNumber))

	if m.MfaCode == "" {
		m.MfaCode, err = promptForMfa()
		if err != nil {
			return err
		}
	}
	paramTokenCode := stype.FieldByName("TokenCode")
	if !paramTokenCode.IsValid() {
		return fmt.Errorf("configureMfa called with struct lacking TokenCode field")
	}
	paramTokenCode.Set(reflect.ValueOf(&m.MfaCode))
	return nil
}
