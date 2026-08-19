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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitEnvironmentApplicationUserSecurityRoleAssignmentResource_CreateReplaceDelete(t *testing.T) {
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
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "system_user_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "role_id", roleAdminID),
				),
			},
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm User"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/MetaForm User"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "system_user_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "role_id", roleUserID),
				),
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserSecurityRoleAssignmentResource_Import(t *testing.T) {
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
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
			},
			{
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     environmentID + "/systemusers/" + principalID + "/MetaForm Global Admin",
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserSecurityRoleAssignmentResource_ImportRoleNameWithSlash(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleName       = "MetaForm Admin/Delegate"
	)

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}]}`,
				principalID, applicationID, rootBusinessID, roleAdminID, roleName, rootBusinessID)
			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[%s]}`, body)), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}]}`,
				roleAdminID, roleName, rootBusinessID)
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
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "` + roleName + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleName),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", roleName),
				),
			},
			{
				// The import ID contains "/" inside the security role name; SplitN must keep
				// the role name's trailing segments intact.
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     environmentID + "/systemusers/" + principalID + "/" + roleName,
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserSecurityRoleAssignmentResource_Read_RoleDeletedOutOfBand(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	)

	roleDeleted := false
	roleAssigned := false

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			rolePayload := ""
			if roleAssigned && !roleDeleted {
				rolePayload = fmt.Sprintf(`{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}`, roleAdminID, rootBusinessID)
			}
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[%s]}`,
				principalID, applicationID, rootBusinessID, rolePayload)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			if roleDeleted {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`, roleAdminID, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			roleAssigned = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the role assignment successfully.
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "role_id", roleAdminID),
				),
			},
			{
				// Step 2: refresh when the security role has been deleted out-of-band.
				// ResolveSecurityRoleNames surfaces ErrObjectNotFound -> the resource is
				// removed from state instead of failing the refresh; the follow-up plan
				// recreates it.
				PreConfig:          func() { roleDeleted = true },
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitEnvironmentApplicationUserSecurityRoleAssignmentResource_Validate_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	)

	environmentDeleted := false
	roleAssigned := false

	// The application client's getEnvironment(). Once the parent environment has been
	// deleted out-of-band (step 2), it returns 404 so an ErrObjectNotFound surfaces.
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			if environmentDeleted {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/application_admin/Read_ParentDeleted/get_environment.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			rolePayload := ""
			if roleAssigned {
				rolePayload = fmt.Sprintf(`{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}`, roleAdminID, rootBusinessID)
			}
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[%s]}`,
				principalID, applicationID, rootBusinessID, rolePayload)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			body := fmt.Sprintf(`{"value":[{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}]}`, roleAdminID, rootBusinessID)
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			roleAssigned = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the role assignment successfully.
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm Global Admin"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "role_id", roleAdminID),
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

// A security role can be assigned to a team, which lives in its own Dataverse table with its own
// role association, so the resource must address teams rather than systemusers.
func TestUnitSecurityRoleAssignmentResource_Team(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		teamID         = "00000000-0000-0000-0000-00000000000a"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleName       = "MetaForm Global Admin"
	)

	assignedRoleIDs := []string{}
	teamRoleAssociationCalls := 0

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/`+environmentID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/teams.*`,
		func(req *http.Request) (*http.Response, error) {
			rolePayload := make([]map[string]string, 0, len(assignedRoleIDs))
			for _, roleID := range assignedRoleIDs {
				rolePayload = append(rolePayload, map[string]string{
					"roleid":                roleID,
					"name":                  roleName,
					"_businessunitid_value": rootBusinessID,
				})
			}
			body, err := json.Marshal(map[string]any{
				"teamid":                teamID,
				"name":                  "Example Team",
				"_businessunitid_value": rootBusinessID,
				"teamroles_association": rolePayload,
			})
			if err != nil {
				return nil, err
			}
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(
				`{"value":[{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}]}`, roleAdminID, roleName, rootBusinessID)), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/teams(%28|\()`+teamID+`(%29|\))/teamroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			teamRoleAssociationCalls++
			if !slices.Contains(assignedRoleIDs, roleAdminID) {
				assignedRoleIDs = append(assignedRoleIDs, roleAdminID)
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/teams(%28|\()`+teamID+`(%29|\))/teamroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			assignedRoleIDs = []string{}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					team_id            = "` + teamID + `"
					security_role_name = "` + roleName + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "team_id", teamID),
					resource.TestCheckNoResourceAttr("powerplatform_security_role_assignment.test", "system_user_id"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/teams/"+teamID+"/"+roleName),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "role_id", roleAdminID),
					func(_ *terraform.State) error {
						if teamRoleAssociationCalls == 0 {
							return fmt.Errorf("expected the role to be associated through teamroles_association")
						}
						return nil
					},
				),
			},
			{
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateId:     environmentID + "/teams/" + teamID + "/" + roleName,
				ImportStateVerify: true,
			},
		},
	})
}

// Exactly one principal must be named.
func TestUnitSecurityRoleAssignmentResource_Principal_Is_ExactlyOneOf(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "neither" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					security_role_name = "MetaForm Global Admin"
				}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
			{
				Config: `
				resource "powerplatform_security_role_assignment" "both" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					system_user_id     = "00000000-0000-0000-0000-000000000008"
					team_id            = "00000000-0000-0000-0000-00000000000a"
					security_role_name = "MetaForm Global Admin"
				}`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}
