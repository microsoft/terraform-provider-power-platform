//go:build tools

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	// mocks generation
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/microsoft/terraform-provider-power-platform/tools/quickstartgen"
)
