// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentApplicationUserRoleAssignmentResource_CreateReplaceDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleUserID     = "31c0083e-67f6-f011-8407-000d3a7a0cab"
	)

	roleNamesByID := map[string]string{
		roleAdminID: "MetaForm Global Admin",
		roleUserID:  "MetaForm User",
	}
	assignedRoleIDs := []string{}

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			rolePayload := make([]map[string]string, 0, len(assignedRoleIDs))
			for _, roleID := range assignedRoleIDs {
				rolePayload = append(rolePayload, map[string]string{
					"roleid":                roleID,
					"name":                  roleNamesByID[roleID],
					"_businessunitid_value": rootBusinessID,
				})
			}

			responsePayload := map[string]any{
				"systemuserid":                principalID,
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

			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, string(bodyBytes)), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, string(bodyBytes))), nil
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

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					principal_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentID+"/"+principalID+"/MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "principal_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "role_id", roleAdminID),
				),
			},
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					principal_id       = "` + principalID + `"
					security_role_name = "MetaForm User"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "id", environmentID+"/"+principalID+"/MetaForm User"),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "principal_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_role_assignment.test", "role_id", roleUserID),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserRoleAssignmentResource_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	)

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`,
				principalID, applicationID, rootBusinessID, roleAdminID, rootBusinessID)
			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`,
				roleAdminID, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))/systemuserroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					principal_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
			},
			{
				ResourceName:      "powerplatform_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     environmentID + "/" + principalID + "/MetaForm Global Admin",
			},
		},
	})
}
