// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccEnvironmentsResource_Validate_Update(t *testing.T) {
	envName := fmt.Sprintf("orgtest%d", rand.Intn(100000))
	envNameNew := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + envName + `"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  =  "` + envName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", envName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", envName),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://"+envName+".crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + envName + `"
					domain									  =  "` + envNameNew + `"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", envName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", envNameNew),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://"+envNameNew+".crm4.dynamics.com/"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create(t *testing.T) {
	envName := fmt.Sprintf("orgtest%d", rand.Intn(100000))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "` + envName + `"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  =  "` + envName + `"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", envName),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://"+envName+".crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", envName),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "templates", regexp.MustCompile(`D365_FinOps_Finance$`)),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "template_metadata", regexp.MustCompile(`{"PostProvisioningPackages": [{ "applicationUniqueName": "msdyn_FinanceAndOperationsProvisioningAppAnchor",\n "parameters": "DevToolsEnabled=true\|DemoDataEnabled=true"\n }\n ]\n }`)),
					// resource.TestMatchResourceAttr("powerplatform_environment.development", "linked_app_url", regexp.MustCompile(`\.operations\.dynamics\.com$`)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	envIdResponseInx := -1
	envIdResponseArray := []string{"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
		"00000000-0000-0000-0000-000000000003",
		"00000000-0000-0000-0000-000000000004",
		"00000000-0000-0000-0000-000000000005"}

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_transactioncurrencies_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			envIdResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_lifecycle_%s.json", envIdResponseArray[envIdResponseInx])).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/([a-z]+)/environmentLanguages\?api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/([a-z]+)/environmentCurrencies\?api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000001"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
	
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000002"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "EUR"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000003"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "EUR"
					environment_type                          = "Trial"
					domain									  = "00000000-0000-0000-0000-000000000004"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Trial"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1031"
					currency_code                             = "EUR"
					environment_type                          = "Trial"
					domain									  = "00000000-0000-0000-0000-000000000005"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000005"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1031"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
		},
	})

}

func TestUnitEnvironmentsResource_Validate_Create_And_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	getLifecycleResponseInx := 0
	patchResponseInx := 0

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_lifecycle_%d.json", getLifecycleResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			getLifecycleResponseInx++
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_environment_%d.json", patchResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_environments_%d.json", patchResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000001"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example123"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000001"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_With_Billing_Policy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_With_Billing_Policy/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_With_Billing_Policy/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_With_Billing_Policy/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000002"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Update_With_Billing_Policy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	getResponseInx := 0
	patchResponseInx := 0

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Update_With_Billing_Policy/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			getResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Update_With_Billing_Policy/get_environment_%d.json", getResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Update_With_Billing_Policy/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001/environments/add?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000002/environments/add?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001/environments/remove?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000002/environments/remove?api-version=2022-03-01-preview",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2022-05-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Update_With_Billing_Policy/get_environments_%d.json", patchResponseInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000002"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					billing_policy_id                         = ""
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "billing_policy_id", ""),
				),
			},
		},
	})
}
func TestUnitEnvironmentsResource_Validate_Create_With_D365_Template(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_With_D365_Template/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_With_D365_Template/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_With_D365_Template/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_currencies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
					billing_policy_id                         = "00000000-0000-0000-0000-000000000002"
					templates = ["D365_FinOps_Finance"]
  					template_metadata = "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("powerplatform_environment.development", "templates.*", "D365_FinOps_Finance"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "template_metadata", "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"),
				),
			},
		},
	})
}