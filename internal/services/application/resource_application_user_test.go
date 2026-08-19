// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentApplicationUserResource_CreateDelete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_data_role_assignment.example",
				Config: `
				resource "azuread_application_registration" "example_app" {
					display_name = "TestAccAdminManagementApplicationResource Application"
				}

				resource "azuread_service_principal" "example_sp" {
					client_id = azuread_application_registration.example_app.client_id
				}
		
				resource "powerplatform_environment" "env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_application_user" "application_user" {
					environment_id = powerplatform_environment.env.id
					application_id = azuread_application_registration.example_app.client_id
				}

				resource "powerplatform_data_role_assignment" "example" {
					environment_id     = powerplatform_environment.env.id
					principal_id       = powerplatform_application_user.application_user.id
					security_role_name = "Basic User"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerplatform_data_role_assignment.example", "id"),
					resource.TestCheckResourceAttrSet("powerplatform_data_role_assignment.example", "environment_id"),
					resource.TestCheckResourceAttrSet("powerplatform_data_role_assignment.example", "principal_id"),
					resource.TestCheckResourceAttr("powerplatform_data_role_assignment.example", "security_role_name", "Basic User"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_CreateDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
	)

	userDeleted := false
	userCreated := false
	userDisabled := true

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			userCreated = true
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", systemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if userDeleted || !userCreated {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":%t,"systemuserroles_association":[]}`,
				systemUserID, applicationID, rootBusinessID, userDisabled)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			userDisabled = false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			userDeleted = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "system_user_id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "disabled", "false"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_UpdateDisabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
	)

	userDisabled := false

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", systemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":%t,"systemuserroles_association":[]}`,
				systemUserID, applicationID, rootBusinessID, userDisabled)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			userDisabled = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "disabled", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
					disabled       = true
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "disabled", "true"),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_RecreatesDeletedUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID   = "00000000-0000-0000-0000-000000000001"
		applicationID   = "00000000-0000-0000-0000-000000000002"
		firstSystemUser = "00000000-0000-0000-0000-000000000008"
		nextSystemUser  = "00000000-0000-0000-0000-000000000009"
		rootBusinessID  = "00000000-0000-0000-0000-000000000003"
	)

	createCount := 0
	currentSystemUserID := firstSystemUser
	userDeleted := false

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			createCount++
			userDeleted = false
			if createCount > 1 {
				currentSystemUserID = nextSystemUser
			}

			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", currentSystemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if createCount == 0 {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			deletedState := 0
			if userDeleted {
				deletedState = 1
			}
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"deletedstate":%d,"systemuserroles_association":[]}`,
				currentSystemUserID, applicationID, rootBusinessID, deletedState)

			if strings.Contains(req.URL.Path, "/systemusers(") || strings.Contains(req.URL.RawPath, "/systemusers%28") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			userDeleted = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000009(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", firstSystemUser),
				),
			},
			{
				PreConfig: func() {
					userDeleted = true
				},
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", nextSystemUser),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
	)

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", systemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[]}`,
				systemUserID, applicationID, businessUnitID)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
			},
			{
				ResourceName:      "powerplatform_application_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     environmentID + "/" + applicationID,
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_Read_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					application_id = "00000000-0000-0000-0000-000000000002"
				}
				`,
				ExpectError: regexp.MustCompile(".*"),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_CreateAdoptsExistingAndReconcilesDisabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
	)
	postAttempts := 0
	patchAttempts := 0
	userDisabled := true

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			postAttempts++
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"message":"A record with matching applicationid already exists."}}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			row := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Existing Application User","_businessunitid_value":"%s","isdisabled":%t,"deletedstate":0,"systemuserroles_association":[]}`,
				systemUserID, applicationID, businessUnitID, userDisabled)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, row), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, row)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			patchAttempts++
			userDisabled = false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "system_user_id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "business_unit_id", businessUnitID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "disabled", "false"),
				),
			},
		},
	})
	if postAttempts != 0 {
		t.Fatalf("existing application user should be adopted without POST; observed %d create attempt(s)", postAttempts)
	}
	if patchAttempts < 1 {
		t.Fatal("adopted disabled state was not reconciled")
	}
}

func TestUnitEnvironmentApplicationUserResource_CreateFailsOnAmbiguousActiveUsers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
	)
	postAttempts := 0

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[
				{"systemuserid":"00000000-0000-0000-0000-000000000008","applicationid":"%s","_businessunitid_value":"%s","isdisabled":false,"deletedstate":0,"systemuserroles_association":[]},
				{"systemuserid":"00000000-0000-0000-0000-000000000009","applicationid":"%s","_businessunitid_value":"%s","isdisabled":false,"deletedstate":0,"systemuserroles_association":[]}
			]}`, applicationID, businessUnitID, applicationID, businessUnitID)), nil
		})
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			postAttempts++
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `
			resource "powerplatform_application_user" "test" {
				environment_id = "` + environmentID + `"
				application_id = "` + applicationID + `"
			}
			`,
			ExpectError: regexp.MustCompile(`(?s)multiple active application\s+users found.*00000000-0000-0000-0000-000000000008.*00000000-0000-0000-0000-000000000009`),
		}},
	})
	if postAttempts != 0 {
		t.Fatalf("ambiguous application users must fail before POST; observed %d create attempt(s)", postAttempts)
	}
}

func TestUnitEnvironmentApplicationUserResource_AdoptionFailureNeverDeletesExistingUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
	)
	deleteAttempts := 0

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			row := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","_businessunitid_value":"%s","isdisabled":true,"deletedstate":0,"systemuserroles_association":[]}`,
				systemUserID, applicationID, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, row)), nil
		})
	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"message":"reconciliation failed"}}`), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			deleteAttempts++
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `
			resource "powerplatform_application_user" "test" {
				environment_id = "` + environmentID + `"
				application_id = "` + applicationID + `"
			}
			`,
			ExpectError: regexp.MustCompile(`(?s)Failed to reconcile disabled state for adopted application user.*pre-existing application user was not deleted`),
		}},
	})
	if deleteAttempts != 0 {
		t.Fatalf("failed adoption must never delete the pre-existing user; observed %d delete attempt(s)", deleteAttempts)
	}
}

func TestUnitEnvironmentApplicationUserResource_CreatePurgesDeletedConflict(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		deletedUserID  = "00000000-0000-0000-0000-000000000010"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
	)

	createAttempts := 0
	deleteAttempts := 0
	conflictingDeletedUserPurged := false

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			createAttempts++
			if createAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusPreconditionFailed, `{"error":{"message":"A record with matching key values already exists."}}`), nil
			}

			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", systemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if !conflictingDeletedUserPurged {
				body := fmt.Sprintf(`{"value":[{"systemuserid":"%s","applicationid":"%s","fullname":"Deleted Application User","_businessunitid_value":"%s","isdisabled":true,"deletedstate":1,"systemuserroles_association":[]}]}`,
					deletedUserID, applicationID, businessUnitID)
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"deletedstate":0,"systemuserroles_association":[]}`,
				systemUserID, applicationID, businessUnitID)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000010(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			deleteAttempts++
			if deleteAttempts >= 1 {
				conflictingDeletedUserPurged = true
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "disabled", "false"),
				),
			},
		},
	})
}

func TestUnitApplicationUserResource_Validate_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
	)

	environmentDeleted := false
	userCreated := false

	// The application client's getEnvironment(). Once the parent environment has been
	// deleted out-of-band (step 2), it returns 404 so an ErrObjectNotFound surfaces.
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			if environmentDeleted {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/application_admin/Read_ParentDeleted/get_environment.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/businessunits.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"businessunitid":"%s","name":"root"}]}`, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers$`,
		func(req *http.Request) (*http.Response, error) {
			userCreated = true
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%s)", systemUserID))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if !userCreated {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[]}`,
				systemUserID, applicationID, rootBusinessID)
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the application user successfully.
				Config: `
				resource "powerplatform_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_application_user.test", "system_user_id", systemUserID),
				),
			},
			{
				// Step 2: refresh when the parent environment has been deleted out-of-band.
				// getEnvironment returns 404 -> ErrObjectNotFound -> the resource is removed
				// from state instead of failing the refresh; the follow-up plan recreates it.
				PreConfig:          func() { environmentDeleted = true },
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
