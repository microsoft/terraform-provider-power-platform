// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

// Generate the provider document.
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()

	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/microsoft/power-platform",
	}

	err := providerserver.Serve(ctx, provider.NewPowerPlatformProvider(ctx), serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
