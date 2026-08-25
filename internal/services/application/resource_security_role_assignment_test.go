// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
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
			// A roles(<id>) path is the single-row fetch of the id selector and must return one
			// object, not the list envelope.
			for roleID, name := range roleNamesByID {
				if strings.Contains(req.URL.Path, "roles("+roleID+")") || strings.Contains(req.URL.RawPath, "roles%28"+roleID+"%29") {
					return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}`, roleID, name, rootBusinessID)), nil
				}
			}
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
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleAdminID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "system_user_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleAdminID),
				),
			},
			{
				// A name edit resolving to a different role is refused: moving the grant is what
				// security_role_id is for, since a name-driven replacement could adopt and then
				// destroy the old assignment under create_before_destroy.
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id       = "` + principalID + `"
					security_role_name = "MetaForm User"
				}
				`,
				ExpectError: regexp.MustCompile(`(?s)resolves to a different\s+role`),
			},
			{
				// Moving the grant by id replaces the assignment: the old association is removed
				// and the new one added.
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id   = "` + environmentID + `"
					system_user_id     = "` + principalID + `"
					security_role_id = "` + roleUserID + `"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleUserID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "system_user_id", principalID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", rootBusinessID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleUserID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "MetaForm User"),
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

	assigned := false
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			roles := ""
			if assigned {
				roles = fmt.Sprintf(`{"roleid":"%s","name":"MetaForm Global Admin","_businessunitid_value":"%s"}`, roleAdminID, rootBusinessID)
			}
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[%s]}`,
				principalID, applicationID, rootBusinessID, roles)
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
			assigned = true
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
				ImportStateId:     environmentID + "/systemusers/" + principalID + "/" + roleAdminID,
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

	assigned := false
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			roles := ""
			if assigned {
				roles = fmt.Sprintf(`{"roleid":"%s","name":"%s","_businessunitid_value":"%s"}`, roleAdminID, roleName, rootBusinessID)
			}
			body := fmt.Sprintf(`{"systemuserid":"%s","applicationid":"%s","fullname":"Example Application User","_businessunitid_value":"%s","isdisabled":false,"systemuserroles_association":[%s]}`,
				principalID, applicationID, rootBusinessID, roles)
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
			assigned = true
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
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleAdminID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", roleName),
				),
			},
			{
				// The import ID ends in the immutable role id, never the role name, so a rename
				// or a name containing "/" cannot corrupt an import.
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     environmentID + "/systemusers/" + principalID + "/" + roleAdminID,
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
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleAdminID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleAdminID),
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
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleAdminID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleAdminID),
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
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/teams/"+teamID+"/"+roleAdminID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleAdminID),
					func(_ *terraform.State) error {
						if teamRoleAssociationCalls == 0 {
							return errors.New("expected the role to be associated through teamroles_association")
						}
						return nil
					},
				),
			},
			{
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateId:     environmentID + "/teams/" + teamID + "/" + roleAdminID,
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

// Dataverse does not let access teams hold security roles, so the resource refuses them with a
// pointer at the team kinds that can.
func TestUnitSecurityRoleAssignmentResource_AccessTeam_Is_Rejected(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID = "00000000-0000-0000-0000-000000000001"
		teamID        = "00000000-0000-0000-0000-00000000000b"
	)

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/`+environmentID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/teams.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK,
				`{"teamid":"`+teamID+`","name":"Access Team","teamtype":1,"_businessunitid_value":"00000000-0000-0000-0000-000000000003"}`), nil
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
					security_role_name = "MetaForm Global Admin"
				}`,
				ExpectError: regexp.MustCompile(`access teams cannot hold security roles`),
			},
		},
	})
}

// An empty team_id must fail at plan time. Before ids were validated it satisfied the
// exactly-one-of check and then crossed onto the systemusers path, because the holder kind was
// inferred from the string being non-empty.
func TestUnitSecurityRoleAssignmentResource_Empty_TeamId_Is_Rejected(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					team_id            = ""
					security_role_name = "MetaForm Global Admin"
				}`,
				ExpectError: regexp.MustCompile(`Invalid UUID String Value`),
			},
		},
	})
}

// An association removed out of band returns 404 on the dissociation DELETE; destroy is already
// done at that point and must succeed.
func TestUnitSecurityRoleAssignmentResource_Delete_404_Is_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID  = "00000000-0000-0000-0000-000000000001"
		applicationID  = "00000000-0000-0000-0000-000000000002"
		principalID    = "00000000-0000-0000-0000-000000000008"
		rootBusinessID = "00000000-0000-0000-0000-000000000003"
		roleAdminID    = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleName       = "MetaForm Global Admin"
	)

	assigned := false
	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/`+environmentID,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			roles := ""
			if assigned {
				roles = `{"roleid":"` + roleAdminID + `","name":"` + roleName + `","_businessunitid_value":"` + rootBusinessID + `"}`
			}
			return httpmock.NewStringResponse(http.StatusOK,
				`{"systemuserid":"`+principalID+`","applicationid":"`+applicationID+`","fullname":"Example","_businessunitid_value":"`+rootBusinessID+`","isdisabled":false,"systemuserroles_association":[`+roles+`]}`), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK,
				`{"value":[{"roleid":"`+roleAdminID+`","name":"`+roleName+`","_businessunitid_value":"`+rootBusinessID+`"}]}`), nil
		})
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()`+principalID+`(%29|\))/systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			assigned = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	// the association was removed out of band: Dataverse answers the dissociation with 404
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()`+principalID+`(%29|\))/systemuserroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, `{"error":{"code":"0x80040217","message":"association does not exist"}}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "` + environmentID + `"
					system_user_id     = "` + principalID + `"
					security_role_name = "` + roleName + `"
				}`,
				Check: resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleAdminID),
			},
		},
		// the framework destroys at the end of the test; the 404 above must not fail it
	})
}

// A failed association POST never becomes successful state, even when the association was truly
// committed: the responder inserts the role into the harness's remote map exactly as a committed
// POST would, so any reintroduced read-back adoption would find it and this test would fail. The
// state is proven empty by the next plan proposing a create, without applying anything.
func TestUnitSecurityRoleAssignmentResource_Create_Ambiguous_Failure_Never_Writes_State(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		roleID      = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		principalID = "00000000-0000-0000-0000-000000000008"
	)
	h := securityRenameHarness(roleID, "Shared Role", "unused")
	// The POST commits the association into the remote map, but the response is lost in transit.
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			h.addRole(principalID, roleID)
			return nil, errors.New("connection reset by peer")
		})

	config := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		system_user_id     = "` + principalID + `"
		security_role_name = "Shared Role"
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)outcome of assigning.*is\s+unknown`),
			},
			{
				// The plan must propose a fresh create, proving nothing reached state, and the
				// re-apply must then be refused by the pre-check, proving the remote grant is
				// live: two proofs from one non-destructive step.
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("powerplatform_security_role_assignment.test", plancheck.ResourceActionCreate),
					},
				},
				ExpectError: regexp.MustCompile(`(?s)already assigned`),
			},
		},
	})
	// The grant lives on remotely, exactly one mutation happened, and Terraform never owned it.
	if h.changes != 1 {
		t.Errorf("expected exactly the committed association, got %d mutations", h.changes)
	}
	if !h.hasRole(principalID, roleID) {
		t.Error("the committed association must survive in the remote map, unowned by Terraform")
	}
}

// A definitive rejection of the association POST fails immediately: the principal is never read
// back, which the harness's own read counter proves, and nothing reaches state.
func TestUnitSecurityRoleAssignmentResource_Create_Definitive_Failure_Fails_Immediately(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		roleID      = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		principalID = "00000000-0000-0000-0000-000000000008"
	)
	h := securityRenameHarness(roleID, "Shared Role", "unused")
	readsAtFailure := -1
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref$`,
		func(_ *http.Request) (*http.Response, error) {
			readsAtFailure = h.principalReads
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"error":{"code":"0x80048306","message":"principal is disabled"}}`), nil
		})

	config := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		system_user_id     = "` + principalID + `"
		security_role_name = "Shared Role"
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)Failed to assign security\s+role.*principal is disabled`),
			},
			{
				// The plan must propose a fresh create, proving nothing reached state; the
				// re-apply then fails on the same definitive rejection.
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("powerplatform_security_role_assignment.test", plancheck.ResourceActionCreate),
					},
				},
				ExpectError: regexp.MustCompile(`(?s)Failed to assign security\s+role.*principal is disabled`),
			},
		},
	})
	if readsAtFailure < 0 {
		t.Fatal("the POST responder never fired")
	}
	if h.principalReads != readsAtFailure {
		t.Errorf("a definitive failure must not read the principal back, got %d reads after the POST", h.principalReads-readsAtFailure)
	}
}

// A role selected by security_role_id needs no name resolution: the role row is fetched by id, its
// business unit becomes the computed business_unit_id, and its name is filled in as a courtesy.
func TestUnitSecurityRoleAssignmentResource_Create_By_Role_Id(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID = "00000000-0000-0000-0000-000000000001"
		applicationID = "00000000-0000-0000-0000-000000000007"
		principalID   = "00000000-0000-0000-0000-000000000008"
		rootBusiness  = "00000000-0000-0000-0000-000000000003"
		roleID        = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	)
	assigned := false

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			roles := "[]"
			if assigned {
				roles = `[{"roleid":"` + roleID + `","name":"MetaForm Global Admin","_businessunitid_value":"` + rootBusiness + `"}]`
			}
			body := `{"systemuserid":"` + principalID + `","applicationid":"` + applicationID + `","fullname":"Example Application User","_businessunitid_value":"` + rootBusiness + `","isdisabled":false,"systemuserroles_association":` + roles + `}`
			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+body+`]}`), nil
		})
	// The single-row fetch must be registered before any roles.* pattern so it matches first.
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles(%28|\()`+roleID+`(%29|\)).*`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"roleid":"`+roleID+`","name":"MetaForm Global Admin","_businessunitid_value":"`+rootBusiness+`"}`), nil
		})
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000008%29/systemuserroles_association/\$ref$`,
		func(_ *http.Request) (*http.Response, error) {
			assigned = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers(%28|\()00000000-0000-0000-0000-000000000008(%29|\))/systemuserroles_association/\$ref.*`,
		func(_ *http.Request) (*http.Response, error) {
			assigned = false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id   = "` + environmentID + `"
					system_user_id   = "` + principalID + `"
					security_role_id = "` + roleID + `"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "id", environmentID+"/systemusers/"+principalID+"/"+roleID),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "MetaForm Global Admin"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", rootBusiness),
				),
			},
		},
	})
}

// A configured business unit must agree with the business unit the role row actually belongs to.
func TestUnitSecurityRoleAssignmentResource_Refuses_Business_Unit_Mismatch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		environmentID = "00000000-0000-0000-0000-000000000001"
		applicationID = "00000000-0000-0000-0000-000000000007"
		principalID   = "00000000-0000-0000-0000-000000000008"
		rootBusiness  = "00000000-0000-0000-0000-000000000003"
		otherBusiness = "00000000-0000-0000-0000-000000000004"
		roleID        = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	)

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			body := `{"systemuserid":"` + principalID + `","applicationid":"` + applicationID + `","fullname":"Example Application User","_businessunitid_value":"` + rootBusiness + `","isdisabled":false,"systemuserroles_association":[]}`
			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+body+`]}`), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles(%28|\()`+roleID+`(%29|\)).*`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"roleid":"`+roleID+`","name":"MetaForm Global Admin","_businessunitid_value":"`+rootBusiness+`"}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id   = "` + environmentID + `"
					system_user_id   = "` + principalID + `"
					security_role_id = "` + roleID + `"
					business_unit_id = "` + otherBusiness + `"
				}`,
				ExpectError: regexp.MustCompile(`(?s)belongs to business unit`),
			},
		},
	})
}

// The role is selected by exactly one of name or id, from both directions.
func TestUnitSecurityRoleAssignmentResource_Role_Selector_Exactly_One(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	base := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id = "00000000-0000-0000-0000-000000000001"
		system_user_id = "00000000-0000-0000-0000-000000000008"
	`
	cases := []struct {
		name   string
		config string
	}{
		{"both", base + `
			security_role_name = "MetaForm Global Admin"
			security_role_id   = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		}`},
		{"neither", base + `}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination`),
					},
				},
			})
		})
	}
}

// securityRenameHarness wires principal-generic mocks for an environment with one assignable role
// whose display name can change between steps. Any systemuser id is served, so replacement tests
// can move the holder. The returned function switches the catalogue to the renamed form, and the
// counter tracks every association add or removal.
// secHarness owns the mock environment's remote truth: the association map, the mutation count
// and the principal-read count, so tests assert against the same store production reads.
type secHarness struct {
	changes        int
	principalReads int
	assigned       map[string][]string
	flip           func()
	roleID         string
}

func (h *secHarness) hasRole(principal, role string) bool {
	return slices.Contains(h.assigned[principal], role)
}

// addRole inserts an association into the remote map directly, the way a concurrent caller or a
// committed-but-unacknowledged POST would, and counts it as a mutation.
func (h *secHarness) addRole(principal, role string) {
	if !slices.Contains(h.assigned[principal], role) {
		h.assigned[principal] = append(h.assigned[principal], role)
		h.changes++
	}
}

func securityRenameHarness(roleID, oldName, newName string) *secHarness {
	const (
		environmentID = "00000000-0000-0000-0000-000000000001"
		applicationID = "00000000-0000-0000-0000-000000000007"
		rootBusiness  = "00000000-0000-0000-0000-000000000003"
	)
	currentName := oldName
	h := &secHarness{assigned: map[string][]string{}, roleID: roleID}
	guidRe := regexp.MustCompile(`systemusers(?:%28|\()([0-9a-f-]+)(?:%29|\))`)
	principalOf := func(raw string) string {
		m := guidRe.FindStringSubmatch(raw)
		if len(m) == 2 {
			return m[1]
		}
		return ""
	}

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			principal := principalOf(req.URL.RawPath + req.URL.Path)
			if principal == "" {
				principal = "00000000-0000-0000-0000-000000000008"
			} else {
				h.principalReads++
			}
			roles := make([]string, 0, len(h.assigned[principal]))
			for _, id := range h.assigned[principal] {
				name := currentName
				if id != roleID {
					name = "A Different Role"
				}
				roles = append(roles, `{"roleid":"`+id+`","name":"`+name+`","_businessunitid_value":"`+rootBusiness+`"}`)
			}
			body := `{"systemuserid":"` + principal + `","applicationid":"` + applicationID + `","fullname":"Example Application User","_businessunitid_value":"` + rootBusiness + `","isdisabled":false,"systemuserroles_association":[` + strings.Join(roles, ",") + `]}`
			if strings.Contains(req.URL.Path, "/systemusers(") || strings.Contains(req.URL.RawPath, "/systemusers%28") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+body+`]}`), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles(%28|\()`+roleID+`(%29|\)).*`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{"roleid":"`+roleID+`","name":"`+currentName+`","_businessunitid_value":"`+rootBusiness+`"}`), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(_ *http.Request) (*http.Response, error) {
			body := `{"value":[
				{"roleid":"` + roleID + `","name":"` + currentName + `","_businessunitid_value":"` + rootBusiness + `"},
				{"roleid":"99999999-0000-0000-0000-000000000009","name":"A Different Role","_businessunitid_value":"` + rootBusiness + `"}
			]}`
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			principal := principalOf(req.URL.RawPath + req.URL.Path)
			var payload map[string]string
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				return nil, err
			}
			for _, candidate := range []string{roleID, "99999999-0000-0000-0000-000000000009"} {
				if strings.Contains(payload["@odata.id"], candidate) {
					h.addRole(principal, candidate)
				}
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			principal := principalOf(req.URL.RawPath + req.URL.Path)
			for _, candidate := range []string{roleID, "99999999-0000-0000-0000-000000000009"} {
				if strings.Contains(req.URL.RawQuery, candidate) && slices.Contains(h.assigned[principal], candidate) {
					h.assigned[principal] = slices.DeleteFunc(h.assigned[principal], func(id string) bool { return id == candidate })
					h.changes++
					return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
				}
			}
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	h.flip = func() { currentName = newName }
	return h
}

// A renamed Dataverse role is followed in place: the association is neither removed nor re-added,
// so create_before_destroy has no adopt-then-destroy window to fall into.
func TestUnitSecurityRoleAssignmentResource_Rename_Same_Role_Updates_In_Place(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	h := securityRenameHarness(roleID, "Old Role Name", "New Role Name")

	config := func(name string) string {
		return `
		resource "powerplatform_security_role_assignment" "test" {
			environment_id     = "00000000-0000-0000-0000-000000000001"
			system_user_id     = "00000000-0000-0000-0000-000000000008"
			security_role_name = "` + name + `"
			timeouts = {
				update = "5m"
			}
		}`
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("Old Role Name"),
				Check: func(_ *terraform.State) error {
					h.flip()
					return nil
				},
			},
			{
				Config: config("New Role Name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "New Role Name"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleID),
					func(_ *terraform.State) error {
						if h.changes != 1 {
							return fmt.Errorf("a same-role rename must not touch the association; expected only the create, got %d changes", h.changes)
						}
						return nil
					},
				),
			},
		},
	})
}

// A name edit that resolves to a different role is refused instead of replacing the grant.
func TestUnitSecurityRoleAssignmentResource_Rename_To_Different_Role_Fails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	_ = securityRenameHarness(roleID, "Old Role Name", "unused")

	config := func(name string) string {
		return `
		resource "powerplatform_security_role_assignment" "test" {
			environment_id     = "00000000-0000-0000-0000-000000000001"
			system_user_id     = "00000000-0000-0000-0000-000000000008"
			security_role_name = "` + name + `"
		}`
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("Old Role Name")},
			{
				Config:      config("A Different Role"),
				ExpectError: regexp.MustCompile(`(?s)resolves to a different\s+role`),
			},
		},
	})
}

// Terraform plans a replacement's new instance from the configuration alone, so a holder change on
// a name-selected assignment cannot carry the stored role id and is refused with the id handed
// over.
func TestUnitSecurityRoleAssignmentResource_Holder_Replacement_Requires_Explicit_Id(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	h := securityRenameHarness(roleID, "Shared Role", "unused")

	config := func(principal string) string {
		return `
		resource "powerplatform_security_role_assignment" "test" {
			environment_id     = "00000000-0000-0000-0000-000000000001"
			system_user_id     = "` + principal + `"
			security_role_name = "Shared Role"
		}`
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("00000000-0000-0000-0000-000000000008")},
			{
				Config:      config("00000000-0000-0000-0000-000000000009"),
				ExpectError: regexp.MustCompile(`(?s)requires the explicit role\s+id`),
			},
			{
				// The refused plan must not have dissociated anything, and the replacement runs
				// with the id the error handed over, granting exactly the anchored role to the
				// new holder.
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id   = "00000000-0000-0000-0000-000000000001"
					system_user_id   = "00000000-0000-0000-0000-000000000009"
					security_role_id = "` + roleID + `"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "system_user_id", "00000000-0000-0000-0000-000000000009"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleID),
					func(_ *terraform.State) error {
						// create P8 (1), then the replacement: dissociate P8 (2), associate P9 (3)
						if h.changes != 3 {
							return fmt.Errorf("expected exactly the create and the explicit-id replacement, got %d association changes", h.changes)
						}
						return nil
					},
				),
			},
		},
	})
}

// A name-selected assignment cannot move to another environment: the stored role id does not exist
// there and the name could resolve to a different role.
func TestUnitSecurityRoleAssignmentResource_Environment_Move_By_Name_Fails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	_ = securityRenameHarness(roleID, "Shared Role", "unused")

	config := func(env string) string {
		return `
		resource "powerplatform_security_role_assignment" "test" {
			environment_id     = "` + env + `"
			system_user_id     = "00000000-0000-0000-0000-000000000008"
			security_role_name = "Shared Role"
		}`
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("00000000-0000-0000-0000-000000000001")},
			{
				Config:      config("00000000-0000-0000-0000-000000000002"),
				ExpectError: regexp.MustCompile(`(?s)cannot move to another\s+environment or business unit`),
			},
		},
	})
}

// The same rule covers an explicit business unit change, which re-scopes the name resolution.
func TestUnitSecurityRoleAssignmentResource_Business_Unit_Move_By_Name_Fails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	_ = securityRenameHarness(roleID, "Shared Role", "unused")

	base := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		system_user_id     = "00000000-0000-0000-0000-000000000008"
		security_role_name = "Shared Role"
	`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: base + `}`},
			{
				Config:      base + `	business_unit_id = "00000000-0000-0000-0000-000000000004"` + "\n\t}",
				ExpectError: regexp.MustCompile(`(?s)cannot move to another\s+environment or business unit`),
			},
		},
	})
}

// A name edit combined with a holder replacement is refused: the rename must be applied first.
func TestUnitSecurityRoleAssignmentResource_Rename_With_Replacement_Fails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	h := securityRenameHarness(roleID, "Shared Role", "Renamed Role")

	config := func(principal, name string) string {
		return `
		resource "powerplatform_security_role_assignment" "test" {
			environment_id     = "00000000-0000-0000-0000-000000000001"
			system_user_id     = "` + principal + `"
			security_role_name = "` + name + `"
		}`
	}

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("00000000-0000-0000-0000-000000000008", "Shared Role"),
				Check: func(_ *terraform.State) error {
					h.flip()
					return nil
				},
			},
			{
				Config:      config("00000000-0000-0000-0000-000000000009", "Renamed Role"),
				ExpectError: regexp.MustCompile(`(?s)cannot change in the same apply as a\s+replacement`),
			},
		},
	})
}

// An unknown name is still name-selected: supplied through an unresolved expression it must not be
// mistaken for the id selector, or a holder replacement would slip past the refusal.
func TestUnitSecurityRoleAssignmentResource_Unknown_Name_Replacement_Refused(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	h := securityRenameHarness(roleID, "Shared Role", "unused")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					system_user_id     = "00000000-0000-0000-0000-000000000008"
					security_role_name = "Shared Role"
				}`,
			},
			{
				Config: `
				resource "terraform_data" "name" {
					input = "Shared Role"
				}
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					system_user_id     = "00000000-0000-0000-0000-000000000009"
					security_role_name = terraform_data.name.output
				}`,
				ExpectError: regexp.MustCompile(`(?s)requires the explicit role\s+id`),
			},
			{
				Config: `
				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					system_user_id     = "00000000-0000-0000-0000-000000000008"
					security_role_name = "Shared Role"
				}`,
				Check: func(_ *terraform.State) error {
					if h.changes != 1 {
						return fmt.Errorf("the refused plan must not touch associations, got %d changes", h.changes)
					}
					return nil
				},
			},
		},
	})
}

// An unknown business unit is an unprovable cross-catalogue move and must fail before anything is
// dissociated, not on the refresh after the damage.
func TestUnitSecurityRoleAssignmentResource_Unknown_Business_Unit_Move_Refused(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	h := securityRenameHarness(roleID, "Shared Role", "unused")

	base := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		system_user_id     = "00000000-0000-0000-0000-000000000008"
		security_role_name = "Shared Role"
	`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: base + `}`},
			{
				Config: `
				resource "terraform_data" "bu" {
					input = "00000000-0000-0000-0000-000000000003"
				}
				` + base + `	business_unit_id = terraform_data.bu.output` + "\n\t}",
				ExpectError: regexp.MustCompile(`(?s)cannot move to another\s+environment or business unit`),
			},
			{
				Config: base + `}`,
				Check: func(_ *terraform.State) error {
					if h.changes != 1 {
						return fmt.Errorf("the refused plan must not touch associations, got %d changes", h.changes)
					}
					return nil
				},
			},
		},
	})
}

// A forced recreate carries no attribute change to detect, so it resolves the name afresh like the
// create it is, and the resolved id lands visibly in state.
func TestUnitSecurityRoleAssignmentResource_Taint_Reresolves_Name(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const (
		roleA        = "7d0690d3-6af6-f011-8407-000d3a7a035d"
		roleB        = "99999999-0000-0000-0000-000000000009"
		rootBusiness = "00000000-0000-0000-0000-000000000003"
		principalID  = "00000000-0000-0000-0000-000000000008"
	)
	// After the flip, the display name belongs to a different role row.
	flipped := false
	assigned := []string{}

	httpmock.RegisterResponder("GET", `=~^https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/application_admin/Create/get_environment.json").String()), nil
		})
	roleName := func(id string) string {
		if (id == roleA) != flipped {
			return "Shared Role"
		}
		return "Some Other Name"
	}
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*`,
		func(req *http.Request) (*http.Response, error) {
			roles := make([]string, 0, len(assigned))
			for _, id := range assigned {
				roles = append(roles, `{"roleid":"`+id+`","name":"`+roleName(id)+`","_businessunitid_value":"`+rootBusiness+`"}`)
			}
			body := `{"systemuserid":"` + principalID + `","applicationid":"00000000-0000-0000-0000-000000000007","fullname":"Example Application User","_businessunitid_value":"` + rootBusiness + `","isdisabled":false,"systemuserroles_association":[` + strings.Join(roles, ",") + `]}`
			if strings.Contains(req.URL.Path, "/systemusers("+principalID+")") || strings.Contains(req.URL.RawPath, "/systemusers%28"+principalID+"%29") {
				return httpmock.NewStringResponse(http.StatusOK, body), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, `{"value":[`+body+`]}`), nil
		})
	httpmock.RegisterResponder("GET", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/roles.*`,
		func(_ *http.Request) (*http.Response, error) {
			body := `{"value":[
				{"roleid":"` + roleA + `","name":"` + roleName(roleA) + `","_businessunitid_value":"` + rootBusiness + `"},
				{"roleid":"` + roleB + `","name":"` + roleName(roleB) + `","_businessunitid_value":"` + rootBusiness + `"}
			]}`
			return httpmock.NewStringResponse(http.StatusOK, body), nil
		})
	httpmock.RegisterResponder("POST", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			var payload map[string]string
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				return nil, err
			}
			for _, candidate := range []string{roleA, roleB} {
				if strings.Contains(payload["@odata.id"], candidate) && !slices.Contains(assigned, candidate) {
					assigned = append(assigned, candidate)
				}
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})
	httpmock.RegisterResponder("DELETE", `=~^https://test-env.crm.dynamics.com/api/data/v9.2/systemusers.*systemuserroles_association/\$ref.*`,
		func(req *http.Request) (*http.Response, error) {
			for _, candidate := range []string{roleA, roleB} {
				if strings.Contains(req.URL.RawQuery, candidate) {
					assigned = slices.DeleteFunc(assigned, func(id string) bool { return id == candidate })
				}
			}
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	config := `
	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		system_user_id     = "` + principalID + `"
		security_role_name = "Shared Role"
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleA),
					func(_ *terraform.State) error {
						flipped = true
						return nil
					},
				),
			},
			{
				Config: config,
				Taint:  []string{"powerplatform_security_role_assignment.test"},
				Check:  resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", roleB),
			},
		},
	})
}

// The scenario that rules out automatic adoption: a forced recreate under create_before_destroy
// plans create first, and an adopting create would record the deposed instance's own association,
// which the deposed delete would then dissociate. With adoption removed the create leg refuses
// instead, before anything is dissociated, and the grant survives.
func TestUnitSecurityRoleAssignmentResource_CBD_Taint_Same_Tuple_Fails_Safely(t *testing.T) {
	const roleID = "7d0690d3-6af6-f011-8407-000d3a7a035d"
	selectors := []struct {
		name     string
		roleLine string
	}{
		{"by_name", `security_role_name = "Shared Role"`},
		{"by_id", `security_role_id = "` + roleID + `"`},
	}
	for _, tc := range selectors {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			h := securityRenameHarness(roleID, "Shared Role", "unused")

			config := `
			resource "powerplatform_security_role_assignment" "test" {
				environment_id = "00000000-0000-0000-0000-000000000001"
				system_user_id = "00000000-0000-0000-0000-000000000008"
				` + tc.roleLine + `
				lifecycle {
					create_before_destroy = true
				}
			}`

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: config},
					{
						Config:      config,
						Taint:       []string{"powerplatform_security_role_assignment.test"},
						ExpectError: regexp.MustCompile(`(?s)already assigned`),
					},
					{
						// Before the suite's cleanup destroy runs, prove the refusal touched
						// nothing: exactly one mutation so far, and the original association is
						// still live. A count alone cannot tell the safe pair (create plus
						// cleanup) from an unsafe one (create plus an illicit removal).
						PreConfig: func() {
							if h.changes != 1 {
								t.Errorf("only the original create may have touched associations before cleanup, got %d changes", h.changes)
							}
							if !h.hasRole("00000000-0000-0000-0000-000000000008", roleID) {
								t.Error("the original association must still exist before cleanup")
							}
						},
						Config:  config,
						Destroy: true,
					},
				},
			})
			// After the explicit destroy: exactly the create and the cleanup dissociation.
			if h.changes != 2 {
				t.Errorf("expected exactly the create and the final cleanup, got %d association changes", h.changes)
			}
		})
	}
}

// The name selector resolves "Basic User", a role every Dataverse environment ships with, inside the
// holder's business unit and pins the id it found. The import step proves the composite id
// round-trips: everything Read rebuilds from it matches what Create wrote.
func TestAccSecurityRoleAssignmentResource_Validate_Create_By_Name(t *testing.T) {
	config := `
	resource "azuread_application_registration" "test_app" {
		display_name = "` + mocks.TestName() + `"
	}

	resource "azuread_service_principal" "test_sp" {
		client_id = azuread_application_registration.test_app.client_id
	}

	resource "powerplatform_environment" "test_env" {
		display_name     = "` + mocks.TestName() + `"
		location         = "europe"
		environment_type = "Sandbox"
		dataverse = {
			language_code     = "1033"
			currency_code     = "USD"
			security_group_id = "00000000-0000-0000-0000-000000000000"
		}
	}

	resource "powerplatform_application_user" "test_user" {
		environment_id = powerplatform_environment.test_env.id
		application_id = azuread_service_principal.test_sp.client_id
	}

	resource "powerplatform_security_role_assignment" "test" {
		environment_id     = powerplatform_environment.test_env.id
		system_user_id     = powerplatform_application_user.test_user.system_user_id
		security_role_name = "Basic User"
	}`

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
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "environment_id", "powerplatform_environment.test_env", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "system_user_id", "powerplatform_application_user.test_user", "system_user_id"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "Basic User"),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "id", regexp.MustCompile(`(?i)^[0-9a-f-]{36}/systemusers/[0-9a-f-]{36}/[0-9a-f-]{36}$`)),
					resource.TestCheckNoResourceAttr("powerplatform_security_role_assignment.test", "team_id"),
				),
			},
			{
				Config:            config,
				ResourceName:      "powerplatform_security_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// The id selector pins the role directly, so the name and business unit are filled in from the role
// row rather than resolved from a name.
func TestAccSecurityRoleAssignmentResource_Validate_Create_By_Role_Id(t *testing.T) {
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
				Config: `
				resource "azuread_application_registration" "test_app" {
					display_name = "` + mocks.TestName() + `"
				}

				resource "azuread_service_principal" "test_sp" {
					client_id = azuread_application_registration.test_app.client_id
				}

				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_application_user" "test_user" {
					environment_id = powerplatform_environment.test_env.id
					application_id = azuread_service_principal.test_sp.client_id
				}

				data "powerplatform_security_roles" "all" {
					environment_id   = powerplatform_environment.test_env.id
					business_unit_id = powerplatform_application_user.test_user.business_unit_id
				}

				resource "powerplatform_security_role_assignment" "test" {
					environment_id   = powerplatform_environment.test_env.id
					system_user_id   = powerplatform_application_user.test_user.system_user_id
					security_role_id = one([
						for role in data.powerplatform_security_roles.all.security_roles :
						role.role_id if role.name == "Basic User"
					])
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "environment_id", "powerplatform_environment.test_env", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "business_unit_id", "powerplatform_application_user.test_user", "business_unit_id"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "Basic User"),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "id", regexp.MustCompile(`(?i)^[0-9a-f-]{36}/systemusers/[0-9a-f-]{36}/[0-9a-f-]{36}$`)),
				),
			},
		},
	})
}

// A team holds roles through its own association table, so the same configuration shape has to land
// on `teams` rather than `systemusers`.
func TestAccSecurityRoleAssignmentResource_Validate_Create_Team(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "europe"
					environment_type = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "test_team" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "team"
					columns = {
						name        = "` + mocks.TestName() + `"
						description = "Owner team for the security role assignment acceptance test"
					}
				}

				resource "powerplatform_security_role_assignment" "test" {
					environment_id     = powerplatform_environment.test_env.id
					team_id            = powerplatform_data_record.test_team.id
					security_role_name = "Basic User"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "environment_id", "powerplatform_environment.test_env", "id"),
					resource.TestCheckResourceAttrPair("powerplatform_security_role_assignment.test", "team_id", "powerplatform_data_record.test_team", "id"),
					resource.TestCheckResourceAttr("powerplatform_security_role_assignment.test", "security_role_name", "Basic User"),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "security_role_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "business_unit_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_security_role_assignment.test", "id", regexp.MustCompile(`(?i)^[0-9a-f-]{36}/teams/[0-9a-f-]{36}/[0-9a-f-]{36}$`)),
					resource.TestCheckNoResourceAttr("powerplatform_security_role_assignment.test", "system_user_id"),
				),
			},
		},
	})
}
