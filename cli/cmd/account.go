package cli

import (
	"context"

	cli "github.com/microsoft/terraform-provider-power-platform/cli"
	common "github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:       "account",
	ValidArgs: []string{"list, clear, set, get-access-token"},
	Short:     "Account management for authenticated credentials",
	Long: `
Command: terraform-provider-power-platform account
  The account command allows you to manage the authenticated accounts stored in the cache file.`,
	Example: `
  List all authenticated accounts stored in the cache file:
  terraform-provider-power-platform account --list
  
  Clear all authenticated accounts stored in the cache file:
  terraform-provider-power-platform account --clear
  
  Set the authenticated account chosing index from the list of authenticated accounts when 
  using the '--list' option (default is '1'):
  terraform-provider-power-platform account --set 1
  
  Authenticate and get an access token for the given scope:
  terraform-provider-power-platform account --get-access-token --scope https://contoso.crm4.dynamics.com/.default`,

	Run: handleAccountCommand,
}

func handleAccountCommand(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	fileCache := common.NewAuthenticationCache()

	if cmd.Flag("list").Changed {
		listAuthenticatedAccounts(ctx, cmd, fileCache)

	} else if cmd.Flag("clear").Changed {
		clearAuthenticatedAccounts(ctx, cmd, fileCache)

	} else if cmd.Flag("get-access-token").Changed {
		getAccessToken(ctx, cmd, fileCache)

	} else if cmd.Flag("set").Changed {
		setAuthenticatedAccount(ctx, cmd, fileCache)

	} else {
		cmd.Help()
	}
}

func getAccessToken(ctx context.Context, cmd *cobra.Command, fileCache *common.AuthenticationCache) {
	scope, err := cmd.Flags().GetString("scope")
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	account, err := fileCache.GetDefaultAccount(ctx)
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	authResult, err := cli.SilentLogin(ctx, []string{scope}, *account, fileCache)
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	cmd.Println(authResult[0].AccessToken)
}

func setAuthenticatedAccount(ctx context.Context, cmd *cobra.Command, fileCache *common.AuthenticationCache) {
	setFlagValue, err := cmd.Flags().GetInt32("set")
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	if setFlagValue < 1 {
		cmd.PrintErrf("Error: '--set' parameter must be greater than 0.\n")
		return
	}
	accounts, err := fileCache.GetAccounts(ctx)
	if len(accounts) == 0 {
		cmd.Println("No authenticated accounts found. Please log in first.")
		return
	}
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	if int(setFlagValue) > len(accounts) {
		cmd.PrintErrf("Error: '--set' parameter must be less than or equal to %d.\n", len(accounts))
		return
	}
	fileCache.SetDefaultAccount(ctx, accounts[setFlagValue-1])
	cmd.Printf("Authenticated account set to: %s (tenant=%s)\n\n", accounts[setFlagValue-1].PreferredUsername, accounts[setFlagValue-1].Realm)
	listAuthenticatedAccounts(ctx, cmd, fileCache)
}

func clearAuthenticatedAccounts(ctx context.Context, cmd *cobra.Command, fileCache *common.AuthenticationCache) {
	fileCache.DeleteFile(ctx)
	cmd.Println("Authenticated accounts cleared.")
}

func listAuthenticatedAccounts(ctx context.Context, cmd *cobra.Command, fileCache *common.AuthenticationCache) {
	accounts, err := fileCache.GetAccounts(ctx)
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	defaultAccount, err := fileCache.GetDefaultAccount(ctx)
	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}
	if len(accounts) == 0 {
		cmd.Println("No authenticated accounts found.")
		return
	} else {
		cmd.Println("Currently authenticated accounts:")
		for inx, account := range accounts {
			defaultString := ""
			if account.PreferredUsername == defaultAccount.PreferredUsername {
				defaultString = "(default)"
			}

			cmd.Printf("%d: %s %s (tenant=%s)\n", inx+1, defaultString, account.PreferredUsername, account.Realm)
		}
		if err != nil {
			cmd.PrintErrf("Error: %s\n", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.Flags().BoolP("list", "l", false, "List all authenticated accounts stored in the cache file.")
	accountCmd.Flags().BoolP("clear", "c", false, "Clear all authenticated accounts stored in the cache file.")
	accountCmd.Flags().Int32P("set", "s", -1, "Set the authenticated account to use.")
	accountCmd.Flags().Bool("get-access-token", false, "Authenticate and get an access token for the given scope.")
	accountCmd.Flags().String("scope", "", "The scope to get an access token for using the '--get-access-token' option.")

	accountCmd.MarkFlagsRequiredTogether("get-access-token", "scope")
}
