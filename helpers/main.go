package helpers

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// PrintAwsError attempts to prettyprint AWS errors if possible
func PrintAwsError(err error) {
	if awsError, ok := err.(awserr.Error); ok {
		fmt.Fprintln(os.Stderr, "AWS Error encountered:")
		fmt.Fprintf(os.Stderr, "%s: %s\n", awsError.Code(), awsError.Message())
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}
