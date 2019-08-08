package creds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	consoleTokenURL = "https://signin.%s.com"
)

// Translations defines common mappings for credential variables
var Translations = map[string]map[string]string{
	"envvar": {
		"AWS_ACCESS_KEY_ID":     "AccessKey",
		"AWS_SECRET_ACCESS_KEY": "SecretKey",
		"AWS_SESSION_TOKEN":     "SessionToken",
		"AWS_SECURITY_TOKEN":    "SessionToken",
		"AWS_DEFAULT_REGION":    "Region",
	},
	"console": {
		"sessionId":    "AccessKey",
		"sessionKey":   "SecretKey",
		"sessionToken": "SessionToken",
	},
}

// Translate converts credentials based on a map of field names
func (c Creds) Translate(dictionary map[string]string) map[string]string {
	old := c.ToMap()
	new := make(map[string]string)
	for k, v := range dictionary {
		new[k] = old[v]
	}
	return new
}

// ToMap returns the credentials as a map of field names to strings
func (c Creds) ToMap() map[string]string {
	return map[string]string{
		"AccessKey":    c.AccessKey,
		"SecretKey":    c.SecretKey,
		"SessionToken": c.SessionToken,
		"Region":       c.Region,
	}
}

// ToSdk returns an AWS SDK Credentials object
func (c *Creds) ToSdk() *credentials.Credentials {
	return credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)
}

// ToEnvVars returns environment variables suitable for eval-ing into the shell
func (c Creds) ToEnvVars() []string {
	envCreds := c.Translate(Translations["envvar"])
	var res []string
	for k, v := range envCreds {
		if v != "" {
			res = append(res, fmt.Sprintf("export %s=%s", k, v))
		}
	}
	sort.Strings(res)
	return res
}

type consoleTokenResponse struct {
	SigninToken string
}

func (c Creds) toConsoleToken() (string, error) {
	args := []string{"?Action=getSigninToken"}

	consoleCreds := c.Translate(Translations["console"])
	jsonCreds, err := json.Marshal(consoleCreds)
	if err != nil {
		return "", err
	}
	urlCreds := url.QueryEscape(string(jsonCreds))
	paramCreds := fmt.Sprintf("Session=%s", urlCreds)
	args = append(args, paramCreds)

	argString := strings.Join(args, "&")
	namespace, err := c.namespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	url := strings.Join([]string{baseURL, "/federation", argString}, "")

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tokenObj := consoleTokenResponse{}
	if err := json.Unmarshal(body, &tokenObj); err != nil {
		return "", err
	}

	return tokenObj.SigninToken, nil
}

// ToConsoleURL returns a console URL for the role
func (c Creds) ToConsoleURL() (string, error) {
	return c.ToCustomConsoleURL("")
}

// ToCustomConsoleURL returns a console URL with a custom path
func (c Creds) ToCustomConsoleURL(dest string) (string, error) {
	consoleToken, err := c.toConsoleToken()
	if err != nil {
		return "", err
	}
	namespace, err := c.namespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	var targetURL string
	if c.Region != "" {
		targetURL = fmt.Sprintf("https://%s.console.%s.com/%s", c.Region, namespace, dest)
	} else {
		targetURL = fmt.Sprintf("https://console.%s.com/%s", namespace, dest)
	}
	urlParts := []string{
		baseURL,
		"/federation",
		"?Action=login",
		"&Issuer=",
		"&Destination=",
		url.QueryEscape(targetURL),
		"&SigninToken=",
		consoleToken,
	}
	urlString := strings.Join(urlParts, "")
	return urlString, nil
}

// ToSignoutURL returns a signout URL for the console
func (c Creds) ToSignoutURL() (string, error) {
	namespace, err := c.namespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	url := strings.Join([]string{baseURL, "/oauth?Action=logout"}, "")
	return url, nil
}
