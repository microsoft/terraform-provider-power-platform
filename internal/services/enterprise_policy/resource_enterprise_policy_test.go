// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitTestEnterpisePolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Create/get_environment_%s.json", id)).String()), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_operation.json").String()), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/enterprisePolicies/NetworkInjection/link?api-version=2019-10-01`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/enterprisePolicies/NetworkInjection/unlink?api-version=2019-10-01`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_enterprise_policy" "network_injection" {
						environment_id = "00000000-0000-0000-0000-000000000001"
						system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
						policy_type    = "NetworkInjection"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_enterprise_policy.network_injection", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_enterprise_policy.network_injection", "system_id", "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_enterprise_policy.network_injection", "policy_type", "NetworkInjection"),
					resource.TestCheckResourceAttr("powerplatform_enterprise_policy.network_injection", "id", "00000000-0000-0000-0000-000000000001_networkinjection")),
			},
		},
	})
}
