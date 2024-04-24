package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	tokens "tinybird-cli/cmd/tokens/actions"
)

var File string

var RootCmd = &cobra.Command{
	Use:   "tokens",
	Short: "The Tokens API allows you to list, create, update or delete your Tinybird Auth Tokens",
	Args:  cobra.ExactArgs(1), // action
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tokens called")
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "Path to file with tokens")
	err := RootCmd.MarkPersistentFlagRequired("file")
	if err != nil {
		fmt.Println("Error setting up flags:", err)
		return
	}

	RootCmd.AddCommand(tokens.PutToken)
}
