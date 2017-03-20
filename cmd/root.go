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
			fmt.Println("AWS Error encountered:")
			fmt.Printf("%s: %s\n", awsError.Code(), awsError.Message())
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
