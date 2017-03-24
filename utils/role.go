package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

// Role defines a set of credentials for the AWS API
type Role struct {
	AccessKey, SecretKey, SessionToken string
}

var translations = map[string]map[string]string{
	"envvar": {
		"AWS_ACCESS_KEY_ID":     "AccessKey",
		"AWS_SECRET_ACCESS_KEY": "SecretKey",
		"AWS_SESSION_TOKEN":     "SessionToken",
		"AWS_SECURITY_TOKEN":    "SessionToken",
	},
	"console": {
		"sessionId":    "AccessKey",
		"sessionKey":   "SecretKey",
		"sessionToken": "SessionToken",
	},
}

// NewRole creates a new Role object from provided credentials
func NewRole(creds map[string]string) (Role, error) {
	required := []string{"AccessKey", "SecretKey", "SessionToken"}
	for _, key := range required {
		elem, ok := creds[key]
		if !ok || elem == "" {
			return Role{}, fmt.Errorf("Missing required key for Role: %s", key)
		}
	}
	role := Role{
		AccessKey:    creds["AccessKey"],
		SecretKey:    creds["SecretKey"],
		SessionToken: creds["SessionToken"],
	}
	return role, nil
}

// NewRoleFromEnv creates a new Role object using credentials from the environment
func NewRoleFromEnv() (Role, error) {
	creds := make(map[string]string)
	for k, v := range translations["envvar"] {
		if creds[v] == "" {
			creds[v] = os.Getenv(k)
		}
	}
	return NewRole(creds)
}

func (r Role) toMap() map[string]string {
	return map[string]string{
		"AccessKey":    r.AccessKey,
		"SecretKey":    r.SecretKey,
		"SessionToken": r.SessionToken,
	}
}

func (r Role) translate(dictionary map[string]string) map[string]string {
	old := r.toMap()
	new := make(map[string]string)
	for k, v := range dictionary {
		new[k] = old[v]
	}
	return new
}

// ToEnvVars returns environment variables suitable for eval-ing into the shell
func (r Role) ToEnvVars() []string {
	creds := r.translate(translations["envvar"])
	var res []string
	for k, v := range creds {
		res = append(res, fmt.Sprintf("export %s=%s", k, v))
	}
	sort.Strings(res)
	return res
}

var consoleTokenURL = "https://signin.aws.amazon.com/federation"

type consoleTokenResponse struct {
	SigninToken string
}

func (r Role) toConsoleToken(lifetime int) (string, error) {
	args := []string{"Action=getSigninToken"}

	paramSession := fmt.Sprintf("SessionDuration=%d", lifetime)
	args = append(args, paramSession)

	creds := r.translate(translations["console"])
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	urlCreds := url.QueryEscape(string(jsonCreds))
	paramCreds := fmt.Sprintf("Session=%s", urlCreds)
	args = append(args, paramCreds)

	argString := strings.Join(args, "&")
	url := strings.Join([]string{consoleTokenURL, argString}, "?")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	tokenObj := consoleTokenResponse{}
	if err := json.Unmarshal(body, &tokenObj); err != nil {
		return "", err
	}

	return tokenObj.SigninToken, nil
}

// ToConsoleURL returns a console URL for the role
func (r Role) ToConsoleURL(lifetime int) (string, error) {
	consoleToken, err := r.toConsoleToken(lifetime)
	if err != nil {
		return "", nil
	}
	urlParts := []string{
		"https://signin.aws.amazon.com/federation",
		"?Action=login",
		"&Issuer=",
		"&Destination=",
		url.QueryEscape("https://console.aws.amazon.com/"),
		"&SigninToken=",
		consoleToken,
	}
	urlString := strings.Join(urlParts, "")
	return urlString, nil
}
