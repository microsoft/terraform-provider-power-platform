// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccApplicationsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name      = "TestAccApplicationsDataSource_Validate_Read"
					location          = "europe"
					environment_type  = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_environment_application_packages" "all_applications" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "id", regexp.MustCompile(`^[\dA-Fa-f]{8}-[\dA-Fa-f]{4}-[\dA-Fa-f]{4}-[\dA-Fa-f]{4}-[\dA-Fa-f]{12}_\d+$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.unique_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.version", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.description", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.publisher_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.publisher_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.learn_more_url", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.state", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_visibility", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitApplicationsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_Read/get_applications.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_Read/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {

			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/application/tests/datasource/environment_application_packages/Validate_Read/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_Read/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "PLN"
						domain                                    = "00000000-0000-0000-0000-000000000001"
						security_group_id                         = "00000000-0000-0000-0000-000000000000"
					}
				}

				data "powerplatform_environment_application_packages" "all_applications" {
					environment_id = powerplatform_environment.env.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "id", "00000000-0000-0000-0000-000000000001_2"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_id", "4bbd5362-21f6-47a8-bcd9-e2a75e8242ef"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.application_name", "Dynamics 365 Customer Voice"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.0.unique_name", "MicrosoftFormsPro"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.application_id", "f50a3059-435a-401b-a7ee-1bca67da5657"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.application_name", "Intelligent Order Management Portal"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_application_packages.all_applications", "applications.1.unique_name", "msdyn_IOMOrderReturnsPortalAnchor"),
				),
			},
		},
	})
}

func TestUnitApplicationsDataSource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_No_Dataverse/get_applications.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {

			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/application/tests/datasource/environment_application_packages/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/application/tests/datasource/environment_application_packages/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}

				data "powerplatform_environment_application_packages" "all_applications" {
					environment_id = powerplatform_environment.env.id
				}`,
				ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
