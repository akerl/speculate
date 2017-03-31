package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "speculate",
	Short:         "Tool for assuming roles in AWS accounts",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute function is the entrypoint for the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			fmt.Fprintln(os.Stderr, "AWS Error encountered:")
			fmt.Fprintf(os.Stderr, "%s: %s\n", awsError.Code(), awsError.Message())
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
