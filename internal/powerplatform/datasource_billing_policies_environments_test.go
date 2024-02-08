package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

//We can't test this until we don't have tenant with billing policies
// func TestAccBillingPoliciesEnvironmentsDataSource_Validate_Read(t *testing.T) {

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { TestAccPreCheck(t) },
// 		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: AcceptanceTestsProviderConfig + `
// 				data "powerplatform_billing_policies_environments" "all_pay_as_you_go_policy_envs" {
// 					billing_policy_id = "00000000-0000-0000-0000-000000000000"
// 				}`,

// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestMatchResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
// 					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.#", "0"),
// 				),
// 			},
// 		},
// 	})
// }

func TestUnitTestBillingPoliciesEnvironmentsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/datasource/environments/get_environments_for_policy.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_billing_policies_environments" "all_pay_as_you_go_policy_envs" {
					billing_policy_id = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "billing_policy_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.#", "3"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.0", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.1", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs", "environments.2", "00000000-0000-0000-0000-000000000003"),
				),
			},
		},
	})
}
