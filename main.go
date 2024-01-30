package main

import (
	cli "github.com/microsoft/terraform-provider-power-platform/cli/cmd"
)

// Generate the provider document.
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

func main() {
	/*
		exe login [tenant=<>]
		exe login [tenant=<>] use-device-code
		exe login [tenant=<>] username=<> password=<>
		exe account list
		exe account clear
		exe account set
		exe account get-access-token
	*/

	cli.Execute()

	//var debug bool

	//flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")

	// login := flag.Bool("login", false, "login to power platform")
	// useDeviceCode := flag.Bool("use-device-code", false, "use device code to login")
	// username := flag.String("username", "", "username for login")
	// password := flag.String("password", "", "password for login")
	// tenant := flag.String("tenant", "", "optional tenant id for login")

	// account := flag.Bool("account", false, "account management")
	// listAccounts := flag.Bool("list", false, "list accounts")
	// clearAccounts := flag.Bool("clear", false, "clear accounts")
	// setAccount := flag.Bool("set", false, "set account")
	// getAccessToken := flag.Bool("get-access-token", false, "get access token")

	//tenantId := flag.String("tenantid", "", "tenant id")
	//username := flag.String("username", "", "username")
	// password := flag.String("password", "", "password")
	// isListAccountsMode := flag.Bool("list-accounts", false, "List accounts in cache")
	// isGetTokenMode := flag.Bool("get-token", false, "Get token for given --scope in 'http://..../.default' format")
	// scope := flag.String("scope", "", "The scope in 'http://..../.default' format")

	// flag.Parse()
	// ctx := context.Background()
	// cache := &common.AuthenticationCache{}

	// if *login {
	// 	cli.Login(ctx, tenant, username, password, useDeviceCode, cache)

	// } else if *account {
	// 	cli.Account()
	// } else {
	// 	log.Default().Println("Working in provider mode")
	// 	actAsProvider(debug, ctx)
	// }

	// if *tenantId == "" {
	// 	log.Default().Println("Working in provider mode")
	// 	actAsProvider(debug, ctx)
	// } else {
	// 	actAsClient(ctx, tenantId, username, password, scope, isListAccountsMode, isGetTokenMode)
	// }
}

// func actAsClient(ctx context.Context, tenantId, username, password, scope *string, isListAccountsMode, isGetTokenMode *bool) {
// 	cache := &common.AuthenticationCache{}

// 	if *isGetTokenMode {
// 		cli.TokenMode(ctx, tenantId, username, password, scope, cache, isListAccountsMode)

// 	} else if *isListAccountsMode {
// 		cli.ListAccountMode(ctx, tenantId, cache)

// 	} else {
// 		cli.LoginMode(ctx, tenantId, username, password, scope, cache)

// 	}

// }

// func actAsProvider(debug bool, ctx context.Context) {
// 	serveOpts := providerserver.ServeOpts{
// 		Debug:   debug,
// 		Address: "registry.terraform.io/microsoft/power-platform",
// 	}

// 	err := providerserver.Serve(ctx, powerplatform.NewPowerPlatformProvider(), serveOpts)

// 	if err != nil {
// 		log.Fatalf("Error serving provider: %s", err)
// 	}
// }
