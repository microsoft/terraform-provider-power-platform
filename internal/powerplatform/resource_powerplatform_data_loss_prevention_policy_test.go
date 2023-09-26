package powerplatform

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

func TestUnitDataLossPreventionPolicyResource_Validate_Update(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	policyId := "00000000-0000-0000-0000-000000000000"
	policy := models.DlpPolicyModel{
		Name:             policyId,
		ETag:             "etag",
		CreatedBy:        "createdBy",
		CreatedTime:      "createdTime",
		LastModifiedBy:   "lastModifiedBy",
		LastModifiedTime: "lastModifiedTime",
	}

	steps := []resource.TestStep{
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy"
				default_connectors_classification = "Blocked"
				environment_type                  = "AllEnvironments"
				environments = []

				business_connectors = []
				non_business_connectors = []
				blocked_connectors = []
				custom_connectors_patterns = []
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "id", policyId),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "Blocked"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "AllEnvironments"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy_1"
				default_connectors_classification = "Blocked"
				environment_type                  = "AllEnvironments"
				environments = []

				business_connectors = []
				non_business_connectors = []
				blocked_connectors = []
				custom_connectors_patterns = []
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy_1"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy"
				default_connectors_classification = "General"
				environment_type                  = "OnlyEnvironments"
				environments = [
					{
						name = "00000000-0000-0000-0000-000000000000"
					}
				]

				business_connectors = []
				non_business_connectors = []
				blocked_connectors = []
				custom_connectors_patterns = []
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "General"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "OnlyEnvironments"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.#", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.0.name", "00000000-0000-0000-0000-000000000000"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy"
				default_connectors_classification = "General"
				environment_type                  = "OnlyEnvironments"
				environments = [
					{
						name = "00000000-0000-0000-0000-000000000000"
					}
				]

				business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				non_business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				blocked_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				custom_connectors_patterns = []
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "0"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "0"),

				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.#", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.default_action_rule_behavior", "Allow"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.action_rules.#", "0"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.endpoint_rules.#", "0"),

				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.#", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.default_action_rule_behavior", "Allow"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.action_rules.#", "0"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.endpoint_rules.#", "0"),

				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "0"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy"
				default_connectors_classification = "General"
				environment_type                  = "OnlyEnvironments"
				environments = [
					{
						name = "00000000-0000-0000-0000-000000000000"
					}
				]

				business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
							default_action_rule_behavior = "Allow",
							action_rules = [
							  {
								action_id = "DeleteItem_V2",
								behavior  = "Block",
							  },
							  {
								action_id = "ExecutePassThroughNativeQuery_V2",
								behavior  = "Block",
							  }
							],
							endpoint_rules = [
							  {
								order    = 1,
								behavior = "Allow",
								endpoint = "contoso.com"
							  },
							  {
								order    = 2,
								behavior = "Deny",
								endpoint = "*"
							  }
							]
						}
					])
				non_business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				blocked_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				custom_connectors_patterns = []
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.action_id", "DeleteItem_V2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.behavior", "Block"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.action_id", "ExecutePassThroughNativeQuery_V2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.behavior", "Block"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.order", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.behavior", "Allow"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.endpoint", "contoso.com"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.order", "2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.behavior", "Deny"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.endpoint", "*"),
			),
		},
		{
			Config: uniTestsProviderConfig + `
			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
				display_name                      = "Block All Policy"
				default_connectors_classification = "General"
				environment_type                  = "OnlyEnvironments"
				environments = [
					{
						name = "00000000-0000-0000-0000-000000000000"
					}
				]

				business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
							default_action_rule_behavior = "Allow",
							action_rules = [
							  {
								action_id = "DeleteItem_V2",
								behavior  = "Block",
							  },
							  {
								action_id = "ExecutePassThroughNativeQuery_V2",
								behavior  = "Block",
							  }
							],
							endpoint_rules = [
							  {
								order    = 1,
								behavior = "Allow",
								endpoint = "contoso.com"
							  },
							  {
								order    = 2,
								behavior = "Deny",
								endpoint = "*"
							  }
							]
						}
					])
				non_business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
				blocked_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							default_action_rule_behavior = "Allow",
							action_rules = [],
							endpoint_rules = [],
						}
					])
					custom_connectors_patterns = toset([
						{
						  order            = 1
						  host_url_pattern = "https://*.contoso.com"
						  data_group       = "Blocked"
						},
						{
						  order            = 2
						  host_url_pattern = "*"
						  data_group       = "Ignore"
						}
					  ])
			  }`,

			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.order", "1"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.host_url_pattern", "https://*.contoso.com"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.data_group", "Blocked"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.order", "2"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.host_url_pattern", "*"),
				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.data_group", "Ignore"),
			),
		},
	}

	clientMock.EXPECT().UpdatePolicy(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, policyToUpdate models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
		policy.DisplayName = policyToUpdate.DisplayName
		policy.DefaultConnectorsClassification = policyToUpdate.DefaultConnectorsClassification
		policy.EnvironmentType = policyToUpdate.EnvironmentType
		policy.Environments = policyToUpdate.Environments
		policy.ConnectorGroups = policyToUpdate.ConnectorGroups
		policy.CustomConnectorUrlPatternsDefinition = policyToUpdate.CustomConnectorUrlPatternsDefinition
		return &policy, nil
	}).Times(len(steps) - 1)

	clientMock.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyId string) (*models.DlpPolicyModel, error) {
		return &policy, nil
	}).AnyTimes()

	clientMock.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyToCreate models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
		policy.DisplayName = policyToCreate.DisplayName
		policy.DefaultConnectorsClassification = policyToCreate.DefaultConnectorsClassification
		policy.EnvironmentType = policyToCreate.EnvironmentType
		policy.Environments = policyToCreate.Environments
		policy.ConnectorGroups = policyToCreate.ConnectorGroups
		policy.CustomConnectorUrlPatternsDefinition = policyToCreate.CustomConnectorUrlPatternsDefinition

		return &policy, nil
	}).Times(1)

	clientMock.EXPECT().DeletePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyId string) error {
		return nil
	}).Times(1)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: steps,
	})
}

func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

	policy := models.DlpPolicyModel{}
	policy.Name = "00000000-0000-0000-0000-000000000000"
	policy.ETag = "etag"
	policy.CreatedBy = "createdBy"
	policy.CreatedTime = "createdTime"
	policy.LastModifiedBy = "lastModifiedBy"
	policy.LastModifiedTime = "lastModifiedTime"

	clientMock.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyID string) (*models.DlpPolicyModel, error) {
		return &policy, nil
	}).AnyTimes()

	clientMock.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyToCreate models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
		policy.DisplayName = policyToCreate.DisplayName
		policy.DefaultConnectorsClassification = policyToCreate.DefaultConnectorsClassification
		policy.EnvironmentType = policyToCreate.EnvironmentType
		policy.Environments = policyToCreate.Environments
		policy.ConnectorGroups = policyToCreate.ConnectorGroups
		policy.CustomConnectorUrlPatternsDefinition = policyToCreate.CustomConnectorUrlPatternsDefinition

		return &policy, nil
	}).Times(1)

	clientMock.EXPECT().DeletePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyID string) error {
		return nil
	}).Times(1)

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
		},
		Steps: []resource.TestStep{
			{
				Config: uniTestsProviderConfig + `
				resource "powerplatform_data_loss_prevention_policy" "my_policy" {
					display_name                      = "Block All Policy"
					default_connectors_classification = "Blocked"
					environment_type                  = "OnlyEnvironments"
					environments = [
						{
							name = "00000000-0000-0000-0000-000000000000"
						}
					]

					business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
							default_action_rule_behavior = "Allow",
							action_rules = [
							  {
								action_id = "DeleteItem_V2",
								behavior  = "Block",
							  },
							  {
								action_id = "ExecutePassThroughNativeQuery_V2",
								behavior  = "Block",
							  }
							],
							endpoint_rules = [
							  {
								order    = 1,
								behavior = "Allow",
								endpoint = "contoso.com"
							  },
							  {
								order    = 2,
								behavior = "Deny",
								endpoint = "*"
							  }
							]
						  }
					])
					non_business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							default_action_rule_behavior = "Allow",
							action_rules                 = [],
							endpoint_rules               = []
						},
					])
					blocked_connectors      = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							default_action_rule_behavior = "Allow",
							action_rules                 = []
							endpoint_rules               = []
						  },
					])
					custom_connectors_patterns = toset([
					  {
						order            = 1
						host_url_pattern = "https://*.contoso.com"
						data_group       = "Blocked"
					  },
					  {
						order            = 2
						host_url_pattern = "*"
						data_group       = "Ignore"
					  }
					])
				  }`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "id", policy.Name),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "Blocked"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "OnlyEnvironments"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.0.name", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.action_id", "DeleteItem_V2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.behavior", "Block"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.action_id", "ExecutePassThroughNativeQuery_V2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.behavior", "Block"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.order", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.endpoint", "contoso.com"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.order", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.behavior", "Deny"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.endpoint", "*"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.default_action_rule_behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.action_rules.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.default_action_rule_behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.action_rules.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.order", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.host_url_pattern", "https://*.contoso.com"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.data_group", "Blocked"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.order", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.host_url_pattern", "*"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.data_group", "Ignore"),
				),
			},
		},
	})
}
