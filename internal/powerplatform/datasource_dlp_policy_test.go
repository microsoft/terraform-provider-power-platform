package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestUnitDlpPolicyDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dlp_policy/tests/datasource/Validate_Read/get_policies.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dlp_policy/tests/datasource/Validate_Read/get_policy_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/00000000-0000-0000-0000-000000000002`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/dlp_policy/tests/datasource/Validate_Read/get_policy_00000000-0000-0000-0000-000000000002.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_data_loss_prevention_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					//Verify returned count
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.#", "2"),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.blocked_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.non_business_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.created_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.created_time", "2023-10-02T07:38:50.3269899Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.default_connectors_classification", "General"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.display_name", "a1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.last_modified_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.last_modified_time", "2023-10-02T07:38:50.3269899Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.environment_type", "AllEnvironments"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.environments.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.data_group", "Ignore"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.host_url_pattern", "*"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.order", "1"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.blocked_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.business_connectors.#", "4"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.id", "/providers/Microsoft.PowerApps/apis/shared_office365users"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.default_action_rule_behavior", ""),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.action_rules.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.default_action_rule_behavior", "Block"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.#", "13"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.0.action_id", "CreateFile_V2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.0.behavior", "Allow"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.behavior", "Deny"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.endpoint", "*"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.order", "1"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.created_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.created_time", "2023-10-02T07:38:56.6864176Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.default_connectors_classification", "General"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.display_name", "a2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.last_modified_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.last_modified_time", "2023-10-02T07:56:43.9700369Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environment_type", "ExceptEnvironments"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environments.0.name", "be0eb809-e58a-ec1b-8fce-ea40b0e53442"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.data_group", "NonBusiness"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.host_url_pattern", "http://aaa.com"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.order", "1"),
				),
			},
		},
	})
}
