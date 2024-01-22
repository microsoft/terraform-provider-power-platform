package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/microsoft/terraform-provider-power-platform/cli"
	"github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform"
)

// Generate the provider document.
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	tenantId := flag.String("tenantid", "", "tenant id")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	isListAccountsMode := flag.Bool("list-accounts", false, "List accounts in cache")
	isGetTokenMode := flag.Bool("get-token", false, "Get token for given --scope in 'http://..../.default' format")
	scope := flag.String("scope", "", "The scope in 'http://..../.default' format")

	flag.Parse()
	ctx := context.Background()

	if *tenantId == "" {
		log.Default().Println("Working in provider mode")
		actAsProvider(debug, ctx)
	} else {
		actAsClient(ctx, tenantId, username, password, scope, isListAccountsMode, isGetTokenMode)
	}
}

func actAsClient(ctx context.Context, tenantId, username, password, scope *string, isListAccountsMode, isGetTokenMode *bool) {
	cache := &common.AuthenticationCache{}

	if *isGetTokenMode {
		cli.TokenMode(ctx, tenantId, username, password, scope, cache, isListAccountsMode)

	} else if *isListAccountsMode {
		cli.ListAccountMode(ctx, tenantId, cache)

	} else {
		cli.LoginMode(ctx, tenantId, username, password, scope, cache)

	}

}

func actAsProvider(debug bool, ctx context.Context) {
	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/microsoft/power-platform",
	}

	err := providerserver.Serve(ctx, powerplatform.NewPowerPlatformProvider(), serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
