package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:           "speculate",
	Short:         "Tool for assuming roles in AWS accounts",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if aws_error, ok := err.(awserr.Error); ok {
			fmt.Println("AWS Error encountered:")
			fmt.Printf("%s: %s\n", aws_error.Code(), aws_error.Message())
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
