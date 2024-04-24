package cmd

import (
	"os"
	"tinybird-cli/cmd/tokens"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tinybird-cli",
	Short: "A simple CLI tool to manage some bulk operations for TinyBird API",
}

var AdminToken string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tokens.RootCmd)
	rootCmd.PersistentFlags().StringVarP(&AdminToken, "admin-token", "a", "", "TinyBird admin token (required)")
	err := rootCmd.MarkPersistentFlagRequired("admin-token")
	if err != nil {
		return
	}
}
