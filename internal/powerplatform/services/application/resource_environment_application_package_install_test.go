// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func TestAccEnvironmentApplicationPackageInstallResource_Validate_Install(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "environment" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                           = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "null_resource" "wait_60_seconds" {
					provisioner "local-exec" {
						command = "sleep 60"
					}
					depends_on = [powerplatform_environment.environment]
				}

				data "powerplatform_environment_application_packages" "application_to_install" {
					environment_id = powerplatform_environment.environment.id
					name           = "Power Platform Pipelines"
					publisher_name = "Microsoft Dynamics 365"

					depends_on = [null_resource.wait_60_seconds]
				}

				resource "powerplatform_environment_application_package_install" "development" {
					environment_id = powerplatform_environment.environment.id
  					unique_name = data.powerplatform_environment_application_packages.application_to_install.applications[0].unique_name
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment_application_package_install.development", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment_application_package_install.development", "unique_name", "msdyn_AppDeploymentAnchor"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationPackageInstallResource_Validate_Install(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Install/get_environment_%s.json", id)).String()), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Install/get_operation.json").String()), nil
		},
	)

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages/ProcessMiningAnchor/install?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages/MicrosoftFormsPro/install?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Install/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Install/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Install/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                           = "1033"
						currency_code                           = "PLN"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_environment_application_package_install" "development" {
					environment_id   = powerplatform_environment.env.id
					unique_name      = "ProcessMiningAnchor"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment_application_package_install.development", "id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment_application_package_install.development", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment_application_package_install.development", "unique_name", "ProcessMiningAnchor"),
				),
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                           = "1033"
						currency_code                           = "PLN"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_environment_application_package_install" "development" {
					environment_id   =  powerplatform_environment.env.id
					unique_name      = "MicrosoftFormsPro"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment_application_package_install.development", "id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment_application_package_install.development", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment_application_package_install.development", "unique_name", "MicrosoftFormsPro"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationPackageInstallResource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_No_Dataverse/get_operation.json").String()), nil
		},
	)

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages/ProcessMiningAnchor/install?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages/MicrosoftFormsPro/install?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name     = "displayname"
					location         = "europe"
					environment_type = "Sandbox"
				  }

				resource "powerplatform_environment_application_package_install" "development" {
					environment_id   = powerplatform_environment.env.id
					unique_name      = "ProcessMiningAnchor"
				}`,

				ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
				Check:       resource.ComposeTestCheckFunc(),
			},
		},
	})
}
