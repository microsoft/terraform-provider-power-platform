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

func TestAccUserResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				//lintignore:AT004
				Config: TestsProviderConfig + `

				terraform {
					required_providers {
					  azuread = {
						source = "hashicorp/azuread"
					  }
					  random = {
						source = "hashicorp/random"
					  }
					}
				}

				provider "azuread" {
					use_cli = true
				}
				  
				data "azuread_domains" "aad_domains" {
					only_initial = true
				}
				  
				locals {
					domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
				}
				  
				resource "random_password" "passwords" {
					length           = 16
					special          = true
					override_special = "_%@"
				}
				  
				resource "azuread_user" "test_user" {
					user_principal_name = "TestAccUserResource_Validate_Create@${local.domain_name}"
					display_name        = "TestAccUserResource_Validate_Create"
					mail_nickname       = "TestAccUserResource_Validate_Create"
					password            = random_password.passwords.result
					usage_location      = "US"
				}
				  
				resource "powerplatform_environment" "dataverse_user_example" {
					display_name      = "TestAccUserResource_Validate_Create"
					location          = "europe"
					environment_type  = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}
				  
				resource "powerplatform_user" "new_user" {
					environment_id = powerplatform_environment.dataverse_user_example.id
					security_roles = [
					  "e0d2794e-82f3-e811-a951-000d3a1bcf17", // bot author
					]
					aad_id         = azuread_user.test_user.id
					disable_delete = false
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_user.new_user", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_user.new_user", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.#", "1"),
					resource.TestMatchResourceAttr("powerplatform_user.new_user", "aad_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "first_name", "#"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "last_name", "TestAccUserResource_Validate_Create"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "disable_delete", "false"),

					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.0", "e0d2794e-82f3-e811-a951-000d3a1bcf17"),
				),
			},
		},
	})
}

func TestUnitUserResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addUser?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29&%24filter=azureactivedirectoryobjectid+eq+00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_systemusers.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29/systemuserroles_association/$ref",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_systemuser_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					security_roles = [
					  "d58407f2-48d5-e711-a82c-000d3a37c848",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = false
				}`,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("powerplatform_user.new_user", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "aad_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "user_principal_name", "jdoe@contoso.onmicrosoft.com"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "first_name", "#"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "last_name", "John Doe"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "disable_delete", "false"),

					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.0", "d58407f2-48d5-e711-a82c-000d3a37c848"),
				),
			},
		},
	})

}

func TestUnitUserResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)/addUser\?api-version=2023-06-01`),
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create_And_Force_Recreate/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create_And_Force_Recreate/get_environment_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)\.crm4\.dynamics\.com/api/data/v9\.2/systemusers\?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29&%24filter=azureactivedirectoryobjectid\+eq\+00000000-0000-0000-0000-000000000002`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create_And_Force_Recreate/get_systemusers.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)\.crm4\.dynamics\.com/api/data/v9\.2/systemusers%2800000000-0000-0000-0000-000000000002%29/systemuserroles_association/\$ref`),
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)\.crm4.dynamics\.com/api/data/v9\.2/systemusers%2800000000-0000-0000-0000-000000000002%29\?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create_And_Force_Recreate/get_systemuser_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					security_roles = [
					  "d58407f2-48d5-e711-a82c-000d3a37c848",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = false
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "environment_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000002"
					security_roles = [
					  "d58407f2-48d5-e711-a82c-000d3a37c848",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = false
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "environment_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitUserResource_Validate_Create_And_Force_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	var systemUserGetInx = 0
	var systemUsersGetInx = 0

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addUser?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create_And_Update/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29&%24filter=azureactivedirectoryobjectid+eq+00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			systemUsersGetInx++
			url := fmt.Sprintf("services/authorization/tests/resource/user/Validate_Create_And_Update/get_systemusers_%d.json", systemUsersGetInx)

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(url).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29/systemuserroles_association/$ref",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29",
		func(req *http.Request) (*http.Response, error) {
			systemUserGetInx++

			url := fmt.Sprintf("services/authorization/tests/resource/user/Validate_Create_And_Update/get_systemuser_00000000-0000-0000-0000-000000000002_%d.json", systemUserGetInx)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(url).String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29/systemuserroles_association/$ref?%24id=https%3A%2F%2F00000000-0000-0000-0000-000000000001.crm4.dynamics.com%2Fapi%2Fdata%2Fv9.2%2Froles%2800000000-0000-0000-0000-000000000001%29",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					security_roles = [
					  "00000000-0000-0000-0000-000000000001",
					  "00000000-0000-0000-0000-000000000002",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = false
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.1", "00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					security_roles = [
					  "00000000-0000-0000-0000-000000000002",
					  "00000000-0000-0000-0000-000000000003",
					  "00000000-0000-0000-0000-000000000004",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = false
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.#", "3"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.0", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.1", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}

func TestUnitUserResource_Validate_Disable_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001/addUser?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29&%24filter=azureactivedirectoryobjectid+eq+00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_systemusers.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29/systemuserroles_association/$ref",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29?%24expand=systemuserroles_association%28%24select%3Droleid%2Cname%2Cismanaged%2C_businessunitid_value%29",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/authorization/tests/resource/user/Validate_Create/get_systemuser_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/systemusers%2800000000-0000-0000-0000-000000000002%29",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusNoContent, "")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_user" "new_user" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					security_roles = [
					  "d58407f2-48d5-e711-a82c-000d3a37c848",
					]
					aad_id         = "00000000-0000-0000-0000-000000000002"
					disable_delete = true
				}`,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("powerplatform_user.new_user", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "aad_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "user_principal_name", "jdoe@contoso.onmicrosoft.com"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "first_name", "#"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "last_name", "John Doe"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "disable_delete", "true"),

					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_user.new_user", "security_roles.0", "d58407f2-48d5-e711-a82c-000d3a37c848"),
				),
			},
		},
	})

}
