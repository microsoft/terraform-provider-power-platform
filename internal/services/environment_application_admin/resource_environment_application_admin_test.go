// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_application_admin_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentApplicationAdminResource_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	environmentId, _ := uuid.NewRandom()
	applicationId, _ := uuid.NewRandom()
	compositeId := fmt.Sprintf("%s/%s", environmentId.String(), applicationId.String())

	// Mock the POST call for enrolling the app user
	enrollUrl := fmt.Sprintf(
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/enroll?api-version=2020-10-01&environmentId=%s&appId=%s",
		environmentId.String(),
		applicationId.String(),
	)
	httpmock.RegisterResponder("POST", enrollUrl, httpmock.NewStringResponder(http.StatusOK, ""))

	// Mock the GET call to verify the app user exists
	getUrl := fmt.Sprintf(
		"https://%s.api.api.powerapps.com/api/data/v9.2/applicationusers?$filter=applicationid eq '%s'",
		environmentId.String(),
		applicationId.String(),
	)
	getResponse := fmt.Sprintf(`{
		"@odata.context": "https://%s.api.powerapps.com/api/data/v9.2/$metadata#applicationusers",
		"value": [
			{
				"applicationid": "%s"
			}
		]
	}`, environmentId.String(), applicationId.String())
	httpmock.RegisterResponder("GET", getUrl, httpmock.NewStringResponder(http.StatusOK, getResponse))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"environment_id": config.StringVariable(environmentId.String()),
					"application_id": config.StringVariable(applicationId.String()),
				},
				Config: `
				variable "environment_id" {
					type = string
				}
				
				variable "application_id" {
					type = string
				}
					
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = var.environment_id
					application_id = var.application_id
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "environment_id", environmentId.String()),
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "application_id", applicationId.String()),
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "id", compositeId),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Read_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	environmentId, _ := uuid.NewRandom()
	applicationId, _ := uuid.NewRandom()
	compositeId := fmt.Sprintf("%s/%s", environmentId.String(), applicationId.String())

	// Mock the POST call for creating the app user
	enrollUrl := fmt.Sprintf(
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/enroll?api-version=2020-10-01&environmentId=%s&appId=%s",
		environmentId.String(),
		applicationId.String(),
	)
	httpmock.RegisterResponder("POST", enrollUrl, httpmock.NewStringResponder(http.StatusOK, ""))

	// First return a valid response for the Create operation
	getUrl := fmt.Sprintf(
		"https://%s.api.api.powerapps.com/api/data/v9.2/applicationusers?$filter=applicationid eq '%s'",
		environmentId.String(),
		applicationId.String(),
	)
	getResponse := fmt.Sprintf(`{
		"@odata.context": "https://%s.api.powerapps.com/api/data/v9.2/$metadata#applicationusers",
		"value": [
			{
				"applicationid": "%s"
			}
		]
	}`, environmentId.String(), applicationId.String())
	httpmock.RegisterResponder("GET", getUrl, httpmock.NewStringResponder(http.StatusOK, getResponse))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"environment_id": config.StringVariable(environmentId.String()),
					"application_id": config.StringVariable(applicationId.String()),
				},
				Config: `
				variable "environment_id" {
					type = string
				}
				
				variable "application_id" {
					type = string
				}
					
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = var.environment_id
					application_id = var.application_id
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_admin.test", "id", compositeId),
				),
			},
			{
				ConfigVariables: config.Variables{
					"environment_id": config.StringVariable(environmentId.String()),
					"application_id": config.StringVariable(applicationId.String()),
				},
				Config: `
				variable "environment_id" {
					type = string
				}
				
				variable "application_id" {
					type = string
				}
					
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = var.environment_id
					application_id = var.application_id
				}`,
				// For the second step, simulate the user being removed
				PreConfig: func() {
					// Override the previous response with an empty result
					emptyResponse := fmt.Sprintf(`{
						"@odata.context": "https://%s.api.powerapps.com/api/data/v9.2/$metadata#applicationusers",
						"value": []
					}`, environmentId.String())
					httpmock.RegisterResponder("GET", getUrl, httpmock.NewStringResponder(http.StatusOK, emptyResponse))
				},
				// We expect the resource to be recreated since it was "removed" externally
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitEnvironmentApplicationAdminResource_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	environmentId, _ := uuid.NewRandom()
	applicationId, _ := uuid.NewRandom()
	compositeId := fmt.Sprintf("%s/%s", environmentId.String(), applicationId.String())

	// Mock the GET call for the Read operation that happens after import
	getUrl := fmt.Sprintf(
		"https://%s.api.api.powerapps.com/api/data/v9.2/applicationusers?$filter=applicationid eq '%s'",
		environmentId.String(),
		applicationId.String(),
	)
	getResponse := fmt.Sprintf(`{
		"@odata.context": "https://%s.api.powerapps.com/api/data/v9.2/$metadata#applicationusers",
		"value": [
			{
				"applicationid": "%s"
			}
		]
	}`, environmentId.String(), applicationId.String())
	httpmock.RegisterResponder("GET", getUrl, httpmock.NewStringResponder(http.StatusOK, getResponse))

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"environment_id": config.StringVariable(environmentId.String()),
					"application_id": config.StringVariable(applicationId.String()),
				},
				Config: `
				variable "environment_id" {
					type = string
				}
				
				variable "application_id" {
					type = string
				}
					
				resource "powerplatform_environment_application_admin" "test" {
					environment_id = var.environment_id
					application_id = var.application_id
				}`,
				ResourceName:      "powerplatform_environment_application_admin.test",
				ImportState:       true,
				ImportStateId:     compositeId,
				ImportStateVerify: true,
			},
		},
	})
}
