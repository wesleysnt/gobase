package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gobase",
	Short: "GoBase — a Go project template",
	Long:  `A batteries-included Go project with CLI, HTTP API, database migrations, and JWT auth.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
