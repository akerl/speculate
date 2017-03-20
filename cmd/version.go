package cmd

import (
	"fmt"

	"github.com/akerl/speculate/utils"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of speculate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", utils.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
