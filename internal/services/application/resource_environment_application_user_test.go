// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
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
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleUserID     = "31c0083e-67f6-f011-8407-000d3a7a0cab"
	)

	roleNamesByID := map[string]string{
		roleAdminID: "MetaForm Global Admin",
		roleUserID:  "MetaForm User",
	}
	assignedRoleIDs := []string{}
	userDeleted := false
	userCreated := false

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

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"%s","_businessunitid_value":"%s"},{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}]}`,
				roleAdminID, roleNamesByID[roleAdminID], rootBusinessID,
				roleUserID, roleNamesByID[roleUserID], rootBusinessID,
			)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			var payload map[string]string
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				return nil, err
			}
			for roleID := range roleNamesByID {
				if strings.Contains(payload["@odata.id"], roleID) && !slices.Contains(assignedRoleIDs, roleID) {
					assignedRoleIDs = append(assignedRoleIDs, roleID)
				}
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))/systemuserroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			roleID := ""
			for candidate := range roleNamesByID {
				if strings.Contains(req.URL.RawQuery, candidate) {
					roleID = candidate
					break
				}
			}
			if roleID == "" {
				return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
			}
			assignedRoleIDs = slices.DeleteFunc(assignedRoleIDs, func(id string) bool { return id == roleID })
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if userDeleted || !userCreated {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}

			rolePayload := make([]map[string]string, 0, len(assignedRoleIDs))
			for _, roleID := range assignedRoleIDs {
				rolePayload = append(rolePayload, map[string]string{
					"roleid":                roleID,
					"name":                  roleNamesByID[roleID],
					"_businessunitid_value": rootBusinessID,
				})
			}

			responsePayload := map[string]any{
				"systemuserid":                systemUserID,
				"applicationid":               applicationID,
				"fullname":                    "Example Application User",
				"_businessunitid_value":       rootBusinessID,
				"isdisabled":                  false,
				"systemuserroles_association": rolePayload,
			}

			bodyBytes, err := json.Marshal(responsePayload)
			if err != nil {
				return nil, err
			}

			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, string(bodyBytes)), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, string(bodyBytes))), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
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
				resource "powerplatform_environment_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
					security_roles = ["MetaForm Global Admin"]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "id", environmentID+"/"+applicationID),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "system_user_id", systemUserID),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "security_roles.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.0.name", "MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.0.role_id", roleAdminID),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
					security_roles = ["MetaForm Global Admin", "MetaForm User"]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "security_roles.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.0.name", "MetaForm User"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.0.role_id", roleUserID),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.1.name", "MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_environment_application_user.test", "resolved_security_roles.1.role_id", roleAdminID),
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
		roleID         = "7d0690d3-6af6-f011-8407-000d3a7a035d"
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

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`,
				roleID, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.Path, "/systemusers("+systemUserID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+systemUserID+"%29") {
				body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`,
					systemUserID, applicationID, businessUnitID, roleID, businessUnitID,
				)
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}

			body := fmt.Sprintf(`{"value":[{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}]}`,
				systemUserID, applicationID, businessUnitID, roleID, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
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
				resource "powerplatform_environment_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
					security_roles = ["MetaForm Global Admin"]
				}
				`,
			},
			{
				ResourceName:      "powerplatform_environment_application_user.test",
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
				resource "powerplatform_environment_application_user" "test" {
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
				resource "powerplatform_environment_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
				}
				`,
				ExpectError: regexp.MustCompile(`(?s)already exists.*cannot adopt existing application users during create`),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserResource_CreateRoleResolutionFailureRollsBack(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		systemUserID   = "00000000-0000-0000-0000-000000000008"
		businessUnitID = "00000000-0000-0000-0000-000000000003"
		roleID1        = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleID2        = "31c0083e-67f6-f011-8407-000d3a7a0cab"
	)

	deactivateCalls := 0
	deleteCalls := 0

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

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"},{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`,
				roleID1, businessUnitID, roleID2, businessUnitID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://test-env.crm.dynamics.com/api/data/v9.0/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			deactivateCalls++
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))$`,
		func(req *http.Request) (*http.Response, error) {
			deleteCalls++
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_application_user" "test" {
					environment_id = "` + environmentID + `"
					application_id = "` + applicationID + `"
					security_roles = ["MetaForm Global Admin"]
				}
				`,
				ExpectError: regexp.MustCompile(`(?s)ambiguous.*Rollback succeeded`),
			},
		},
	})

	if deactivateCalls != 1 {
		t.Fatalf("expected 1 deactivate call during rollback, got %d", deactivateCalls)
	}
	if deleteCalls != 1 {
		t.Fatalf("expected 1 delete call during rollback, got %d", deleteCalls)
	}
}
