package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/akerl/speculate/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			fmt.Fprintln(os.Stderr, "AWS Error encountered:")
			fmt.Fprintf(os.Stderr, "%s: %s\n", awsError.Code(), awsError.Message())
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
