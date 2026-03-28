// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentGitIntegrationResource_ValidateConfig_AllowsUnknownProjectName(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Config: `
resource "powerplatform_environment_git_integration" "seed" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  git_provider      = "AzureDevOps"
  scope             = "Solution"
  organization_name = "example-org"
  project_name      = "example-project"
  repository_name   = "example-repo"
}

resource "powerplatform_environment_git_integration" "test" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  git_provider      = "AzureDevOps"
  scope             = "Solution"
  organization_name = "example-org"
  project_name      = powerplatform_environment_git_integration.seed.id
  repository_name   = "example-repo"
}
`,
			},
		},
	})
}
