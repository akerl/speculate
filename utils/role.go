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

func (r Role) ToEnvVars() []string {
	creds := r.translate(translations["envvar"])
	var res []string
	for k, v := range creds {
		res = append(res, fmt.Sprintf("export %s=%s", k, v))
	}
	sort.Strings(res)
	return res
}

var base_console_token_url = "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session="

type console_token_response struct {
	SigninToken string
}

func (r Role) toConsoleToken() (string, error) {
	creds := r.translate(translations["console"])

	json_creds, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	url_creds := url.QueryEscape(string(json_creds))
	url := strings.Join([]string{base_console_token_url, url_creds}, "")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	token_obj := console_token_response{}
	if err := json.Unmarshal(body, &token_obj); err != nil {
		return "", err
	}

	return token_obj.SigninToken, nil
}

func (r Role) ToConsoleUrl() (string, error) {
	console_token, err := r.toConsoleToken()
	if err != nil {
		return "", nil
	}
	url_parts := []string{
		"https://signin.aws.amazon.com/federation",
		"?Action=login",
		"&Issuer=",
		"&Destination=",
		url.QueryEscape("https://console.aws.amazon.com/"),
		"&SigninToken=",
		console_token,
	}
	url_string := strings.Join(url_parts, "")
	return url_string, nil
}
