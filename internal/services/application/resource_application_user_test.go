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
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentApplicationUserResource_CreateUpdateDelete(t *testing.T) {
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
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"deletedstate":%d,"systemuserroles_association":[]}`,
				currentSystemUserID, applicationID, rootBusinessID, boolToDeletedState(userDeleted))

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

func boolToDeletedState(deleted bool) int {
	if deleted {
		return 1
	}

	return 0
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

func TestUnitEnvironmentApplicationUserResource_CreateAlreadyExists(t *testing.T) {
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
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"message":"A record with matching applicationid already exists."}}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"systemuserid":"%s","applicationid":"%s","fullname":"Existing Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[]}]}`,
				systemUserID, applicationID, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
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
				ExpectError: regexp.MustCompile(`(?s)already exists.*cannot adopt existing application users during create`),
			},
		},
	})
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
