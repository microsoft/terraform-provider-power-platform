package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:       "login",
	ValidArgs: []string{"tenant", "use-device-code", "username", "password"},
	Short:     "Log in command for Microsoft Entra",
	Long: `
Command: terraform-provider-power-platform login
  By default, the login command will open a browser windows and ask you to login. If a browser is not 
  available, you can use the device code flow with the --use-device-code option or username and 
  password`,

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
  terraform-provider-power-platform login --tenant 000000000-0000-0000-0000-000000000000 --use-device-code`,

	Run: func(cmd *cobra.Command, args []string) {

		tenantId, _ := cmd.Flags().GetString("tenant")

		if cmd.Flag("use-device-code").Value.String() != "false" {

			fmt.Println("X login device code here... tenantId=" + tenantId)

		} else if cmd.Flag("username").Value.String() != "" && cmd.Flag("password").Value.String() != "" {

			fmt.Println("X login username and password here... tenantId=" + tenantId)

		} else if (cmd.Flag("username").Value.String() == "" && cmd.Flag("password").Value.String() != "") ||
			(cmd.Flag("username").Value.String() != "" && cmd.Flag("password").Value.String() == "") {

			fmt.Println("X username and password are required together... tenantId=" + tenantId)

		} else {

			println("X login interactive here... tenantId=" + tenantId)

		}
	},
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
}
