package cli

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "terraform-provider-power-platform",
	Short: "CLI for Power Platform Terraform Provider",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		var debug bool

		if len(args) > 0 && (args[0] == "debug" || args[0] == "--debug") {
			debug = true
		}

		log.Default().Println("Working in provider mode")
		serveOpts := providerserver.ServeOpts{
			Debug:   debug,
			Address: "registry.terraform.io/microsoft/power-platform",
		}

		ctx := context.Background()
		err := providerserver.Serve(ctx, powerplatform.NewPowerPlatformProvider(), serveOpts)

		if err != nil {
			log.Fatalf("Error serving provider: %s", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}
