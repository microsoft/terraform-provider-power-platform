package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:       "account",
	ValidArgs: []string{"list, clear, set, get-access-token"},
	Short:     "Account management for authenticated credentials",
	Long:      ``,
	Example:   ``,

	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flag("list").Value.String() == "true" {
			fmt.Println("X account list")
		} else if cmd.Flag("clear").Value.String() == "true" {
			fmt.Println("X account clear")
		} else if cmd.Flag("get-access-token").Value.String() == "true" {
			fmt.Println("X account get-access-token")
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.Flags().BoolP("list", "l", false, "List all authenticated accounts stored in the cache file.")
	accountCmd.Flags().BoolP("clear", "c", false, "Clear all authenticated accounts stored in the cache file.")
	//accountCmd.Flags().BoolP("set", "s", false, "Set the authenticated account to use.")
	accountCmd.Flags().BoolP("get-access-token", "g", false, "Authenticate and get an access token for the given scope.")
}
