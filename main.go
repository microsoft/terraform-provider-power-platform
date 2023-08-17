package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform"
)

// Generate the provider document.
//go:generate tfplugindocs generate

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()
	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "github.com/microsoft/powerplatform",
	}

	err := providerserver.Serve(ctx, powerplatform.NewPowerPlatformProvider(), serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
