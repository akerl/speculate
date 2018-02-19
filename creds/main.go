package creds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Creds defines a set of AWS credentials
type Creds struct {
	AccessKey, SecretKey, SessionToken, Region string
}

// New initializes credentials from a map
func New(argCreds map[string]string) (Creds, error) {
	required := []string{"AccessKey", "SecretKey", "SessionToken"}
	for _, key := range required {
		elem, ok := argCreds[key]
		if !ok || elem == "" {
			return Creds{}, fmt.Errorf("Missing required key for Creds: %s", key)
		}
	}
	c := Creds{
		AccessKey:    argCreds["AccessKey"],
		SecretKey:    argCreds["SecretKey"],
		SessionToken: argCreds["SessionToken"],
		Region:       argCreds["Region"],
	}
	return c, nil
}

// NewFromStsSdk initializes a credential object from an AWS SDK Credentials object
func NewFromStsSdk(stsCreds *sts.Credentials) (Creds, error) {
	return New(map[string]string{
		"AccessKey":    *stsCreds.AccessKeyId,
		"SecretKey":    *stsCreds.SecretAccessKey,
		"SessionToken": *stsCreds.SessionToken,
	})
}

// NewFromEnv initializes credentials from the environment variables
func NewFromEnv() (Creds, error) {
	// TODO: Handle region env vars here
	envCreds := make(map[string]string)
	for k, v := range Translations["envvar"] {
		if envCreds[v] == "" {
			envCreds[v] = os.Getenv(k)
		}
	}
	return New(envCreds)
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

// ToSdk returns an AWS SDK Credentials object
func (c *Creds) ToSdk() *credentials.Credentials {
	return credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)
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
	namespace, err := c.namespace()
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
	namespace, err := c.namespace()
	if err != nil {
		return "", err
	}
	baseURL := fmt.Sprintf(consoleTokenURL, namespace)
	url := strings.Join([]string{baseURL, "/oauth?Action=logout"}, "")
	return url, nil
}

// Client returns an AWS STS client for these creds
func (c Creds) Client() *sts.STS {
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	if c.AccessKey != "" {
		config.WithCredentials(c.ToSdk())
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *config,
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sts.New(sess)
}

func (c Creds) identity() (*sts.GetCallerIdentityOutput, error) {
	params := &sts.GetCallerIdentityInput{}
	client := c.Client()
	return client.GetCallerIdentity(params)
}

func (c Creds) partition() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	pieces := strings.Split(*identity.Arn, ":")
	return pieces[1], nil
}

func (c Creds) namespace() (string, error) {
	partition, err := c.partition()
	if err != nil {
		return "", err
	}
	result, ok := namespaces[partition]
	if ok {
		return result, nil
	}
	return "", fmt.Errorf("Unknown partition: %s", partition)
}

var namespaces = map[string]string{
	"aws":        "aws.amazon",
	"aws-us-gov": "amazonaws-us-gov",
}

// MfaArn returns the user's virtual MFA token ARN
func (c Creds) MfaArn() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	if strings.Index(*identity.Arn, ":user/") == -1 {
		return "", fmt.Errorf("Failed to parse MFA ARN for non-user: %s", *identity.Arn)
	}
	mfaArn := strings.Replace(*identity.Arn, ":user/", ":mfa/", 1)
	return mfaArn, nil
}

// SessionName returns the default session name
func (c Creds) SessionName() (string, error) {
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	arnChunks := strings.Split(*identity.Arn, "/")
	oldName := arnChunks[len(arnChunks)-1]
	return oldName, nil
}

// NextRoleArn returns the new role's ARN
func (c Creds) NextRoleArn(role, accountID string) (string, error) {
	if role == "" {
		return "", fmt.Errorf("Role name cannot be empty")
	}
	identity, err := c.identity()
	if err != nil {
		return "", err
	}
	partition, err := c.partition()
	if err != nil {
		return "", err
	}
	if accountID == "" {
		accountID = *identity.Account
	}
	arn := fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, accountID, role)
	return arn, nil
}
