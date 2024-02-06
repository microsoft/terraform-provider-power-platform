package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	cli "github.com/microsoft/terraform-provider-power-platform/cli/cmd"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform"
)

// Generate the provider document.
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()

	if len(flag.Args()) == 0 {
		serveOpts := providerserver.ServeOpts{
			Debug:   debug,
			Address: "registry.terraform.io/microsoft/power-platform",
		}

		err := providerserver.Serve(ctx, powerplatform.NewPowerPlatformProvider(ctx), serveOpts)

		if err != nil {
			log.Fatalf("Error serving provider: %s", err)
		}
	} else if debug {
		fmt.Println("To use CLI, run `terraform-provider-power-platform --help`")
		fmt.Println()
		fmt.Println()

	} else {
		cli.Execute()
	}

}
