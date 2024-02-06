package main

import (
	cli "github.com/microsoft/terraform-provider-power-platform/cli/cmd"
)

// Generate the provider document.
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

func main() {
	cli.Execute()
}
