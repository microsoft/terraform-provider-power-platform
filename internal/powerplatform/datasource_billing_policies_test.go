package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccBillingPoliciesDataSource_Validate_Read(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AcceptanceTestsProviderConfig + `
				data "powerplatform_billing_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.location", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.status", regexp.MustCompile(`^(Enabled|Disabled)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.resource_group", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.subscription_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
		},
	})
}

func TestUnitTestBillingPoliciesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/licensing/test/datasource/policies/get_billing_policies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_billing_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_billing_policies.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.#", "2"),

					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.name", "billingpolicy1"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.status", "Enabled"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.resource_group", "rg-group1"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.0.billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.name", "billingpolicy2"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.status", "Enabled"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.billing_instrument.resource_group", "rg-group2"),
					resource.TestCheckResourceAttr("data.powerplatform_billing_policies.all", "billing_policies.1.billing_instrument.subscription_id", "00000000-0000-0000-0000-000000000002"),
				),
			},
		},
	})
}
