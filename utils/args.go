package utils

import (
	"fmt"
	"os"
)

func RoleParse(args []string) string {
	if len(args) < 1 {
		fmt.Printf("No role name provided\n")
		os.Exit(1)
	}
	return args[0]
}
