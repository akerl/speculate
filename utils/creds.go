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

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

// CredsExecutor defines the interface for requesting a new set of AWS creds
type CredsExecutor interface {
	ParseFlags(*cobra.Command) error
	Execute() (Creds, error)
	ExecuteWithCreds(Creds) (Creds, error)
}

// Creds defines a set of AWS credentials
type Creds struct {
	AccessKey, SecretKey, SessionToken, Region string
}

// New initializes credentials from a map
func (c *Creds) New(argCreds map[string]string) error {
	required := []string{"AccessKey", "SecretKey", "SessionToken"}
	for _, key := range required {
		elem, ok := argCreds[key]
		if !ok || elem == "" {
			return fmt.Errorf("Missing required key for Creds: %s", key)
		}
	}
	c.AccessKey = argCreds["AccessKey"]
	c.SecretKey = argCreds["SecretKey"]
	c.SessionToken = argCreds["SessionToken"]
	return nil
}

// NewFromStsSdk initializes a credential object from an AWS SDK Credentials object
func (c *Creds) NewFromStsSdk(stsCreds *sts.Credentials) error {
	return c.New(map[string]string{
		"AccessKey":    *stsCreds.AccessKeyId,
		"SecretKey":    *stsCreds.SecretAccessKey,
		"SessionToken": *stsCreds.SessionToken,
	})
}

// ToSdk returns an AWS SDK Credentials object
func (c *Creds) ToSdk() *credentials.Credentials {
	return credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)
}

// NewFromEnv initializes credentials from the environment variables
func (c *Creds) NewFromEnv() error {
	envCreds := make(map[string]string)
	for k, v := range Translations["envvar"] {
		if envCreds[v] == "" {
			envCreds[v] = os.Getenv(k)
		}
	}
	return c.New(envCreds)
}

// Translations defines common mappings for credential variables
var Translations = map[string]map[string]string{
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

func (c Creds) toMap() map[string]string {
	return map[string]string{
		"AccessKey":    c.AccessKey,
		"SecretKey":    c.SecretKey,
		"SessionToken": c.SessionToken,
	}
}

// Translate converts credentials based on a map of field names
func (c Creds) Translate(dictionary map[string]string) map[string]string {
	old := c.toMap()
	new := make(map[string]string)
	for k, v := range dictionary {
		new[k] = old[v]
	}
	return new
}

// ToEnvVars returns environment variables suitable for eval-ing into the shell
func (c Creds) ToEnvVars() []string {
	envCreds := c.Translate(Translations["envvar"])
	var res []string
	for k, v := range envCreds {
		res = append(res, fmt.Sprintf("export %s=%s", k, v))
	}
	sort.Strings(res)
	return res
}

var consoleTokenURL = "https://signin.%s.com"

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
	namespace, err := getNamespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	url := strings.Join([]string{baseURL, "/federation", argString}, "")

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
func (c Creds) ToConsoleURL() (string, error) {
	consoleToken, err := c.toConsoleToken()
	if err != nil {
		return "", err
	}
	namespace, err := getNamespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	targetURL := fmt.Sprintf("https://console.%s.com/", namespace)
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
	namespace, err := getNamespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	url := strings.Join([]string{baseURL, "/oauth?Action=logout"}, "")
	return url, nil
}

var namespaces = map[string]string{
	"aws":     "aws.amazon",
	"aws-gov": "amazonaws-us-gov",
}

func getNamespace() (string, error) {
	partition, err := API.Partition()
	if err != nil {
		return "", err
	}
	result, ok := namespaces[partition]
	if ok {
		return result, nil
	}
	return "", fmt.Errorf("Unknown partition: %s", partition)
}
