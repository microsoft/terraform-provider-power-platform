package cli

import (
	"context"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	cli "github.com/microsoft/terraform-provider-power-platform/cli"
	common "github.com/microsoft/terraform-provider-power-platform/common"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:       "login",
	ValidArgs: []string{"tenant", "use-device-code", "username", "password"},
	Short:     "Log in command for Microsoft Entra",
	Long: `
Command: terraform-provider-power-platform login
  By default, the login command will open a browser windows and ask you to login. If a browser
  is not available, you can use the device code flow with the --use-device-code option or 
  username and password`,

	Example: `
  Log in interactively:
  terraform-provider-power-platform login
  
  Log in with a user name and password. This does not work with accounts that have 
  multi-factor authentication (MFA) enabled"
  terraform-provider-power-platform login --username johndoe@contoso.com --password p@ssw0rd
 
  Login using device code. This is useful if you are not able to use the default 
  browser flow:
  terraform-provider-power-platform login --use-device-code
  
  For all login options you can also specify the tenant ID.
  terraform-provider-power-platform login --tenant 00000000-0000-0000-0000-000000000000
  terraform-provider-power-platform login --tenant 00000000-0000-0000-0000-000000000000 
  --username johndoe@contoso.com --password p@ssw0rd
  terraform-provider-power-platform login --tenant 000000000-0000-0000-0000-000000000000 
  --use-device-code`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		fileCache := common.NewAuthenticationCache()
		tenantId, _ := cmd.Flags().GetString("tenant")
		err := error(nil)
		authResults := []public.AuthResult{}

		if cmd.Flag("use-device-code").Changed {
			authResults, err = cli.DeviceCodeLogin(ctx, &tenantId, constants.REQUIRED_SCOPES, fileCache)

		} else if cmd.Flag("username").Changed && cmd.Flag("password").Changed {
			username := cmd.Flag("username").Value.String()
			pass := cmd.Flag("password").Value.String()
			authResults, err = cli.UserPassLogin(ctx, &tenantId, &username, &pass, constants.REQUIRED_SCOPES, fileCache)

		} else {
			authResults, err = cli.InteractiveLogin(ctx, &tenantId, constants.REQUIRED_SCOPES, fileCache)
		}

		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
		}

		err = setDefaultAccountIfNoneSet(ctx, authResults[0].Account, fileCache)
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
		}

	},
}

func setDefaultAccountIfNoneSet(ctx context.Context, account public.Account, fileCache *common.AuthenticationCache) error {
	defaultAccount, err := fileCache.GetDefaultAccount(ctx)
	if err != nil {
		return err
	}

	if defaultAccount == nil {
		err := fileCache.SetDefaultAccount(ctx, account)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("password", "p", "", `Password to login. Flag 'username' is required 
if you use this option.`)
	loginCmd.Flags().StringP("tenant", "t", "", "The Microsoft Entra ID.")
	loginCmd.Flags().BoolP("use-device-code", "d", false, `Flow based on device code. Use this if you are not able to use 
the default browser flow.`)
	loginCmd.Flags().StringP("username", "u", "", `User name to login. Use if format: 
'username@domain.onmicrosoft.com' or 'username@domain'. 
Flag 'password' is required if you use this option.`)

	loginCmd.MarkFlagsRequiredTogether("username", "password")
}
