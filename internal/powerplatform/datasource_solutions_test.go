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

func TestUnitSolutionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource_solutions_test/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource_solutions_test/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"@odata.etag": "W/\"1874400\"",
						"installedon": "2023-10-10T08:09:56Z",
						"solutionpackageversion": "9.0",
						"_configurationpageid_value": null,
						"solutionid": "70edca66-e4c2-4384-92e0-4300465c1894",
						"modifiedon": "2023-10-10T08:09:58Z",
						"uniquename": "ProductivityToolsAnchor",
						"isapimanaged": false,
						"_publisherid_value": "9bb8ab98-18b3-4766-9a2e-243d39523107",
						"ismanaged": true,
						"isvisible": true,
						"thumbprint": null,
						"pinpointpublisherid": null,
						"version": "9.2.1.1020",
						"_modifiedonbehalfby_value": null,
						"_parentsolutionid_value": null,
						"pinpointassetid": null,
						"pinpointsolutionid": null,
						"friendlyname": "ProductivityTools",
						"_organizationid_value": "11afca7f-025d-ee11-a382-000d3a25be4d",
						"versionnumber": 1874400,
						"templatesuffix": null,
						"upgradeinfo": null,
						"_createdonbehalfby_value": null,
						"_modifiedby_value": "f3134d74-515d-ee11-be6f-000d3aaae21d",
						"createdon": "2023-10-10T08:09:56Z",
						"updatedon": null,
						"description": "Productivity Tools for Dynamics 365 apps",
						"solutiontype": null,
						"pinpointsolutiondefaultlocale": null,
						"_createdby_value": "f3134d74-515d-ee11-be6f-000d3aaae21d",
						"publisherid": {
							"@odata.etag": "W/\"1919045\"",
							"address2_line1": null,
							"address1_county": null,
							"pinpointpublisherdefaultlocale": null,
							"address2_utcoffset": null,
							"address2_fax": null,
							"modifiedon": "2023-10-11T00:45:55Z",
							"entityimage_url": null,
							"address1_line1": null,
							"address1_name": null,
							"uniquename": "microsoftdynamics",
							"address1_postalcode": null,
							"address2_line3": null,
							"address1_addressid": null,
							"publisherid": "9bb8ab98-18b3-4766-9a2e-243d39523107",
							"address1_line3": null,
							"address2_name": null,
							"address1_utcoffset": null,
							"address2_city": null,
							"pinpointpublisherid": null,
							"address2_county": null,
							"emailaddress": null,
							"address2_postofficebox": null,
							"address1_stateorprovince": null,
							"address2_telephone3": null,
							"address2_addresstypecode": null,
							"address2_telephone2": null,
							"address2_telephone1": null,
							"address2_shippingmethodcode": null,
							"_modifiedonbehalfby_value": null,
							"isreadonly": true,
							"entityimage_timestamp": null,
							"address2_stateorprovince": null,
							"address1_latitude": null,
							"address1_longitude": null,
							"customizationoptionvalueprefix": 19235,
							"address2_latitude": null,
							"friendlyname": "microsoftdynamics",
							"address1_line2": null,
							"supportingwebsiteurl": "http://crm.dynamics.com",
							"address2_postalcode": null,
							"address2_line2": null,
							"_organizationid_value": "11afca7f-025d-ee11-a382-000d3a25be4d",
							"versionnumber": 1919045,
							"address2_upszone": null,
							"address2_longitude": null,
							"address1_fax": null,
							"customizationprefix": "msdyn",
							"_createdonbehalfby_value": null,
							"_modifiedby_value": "a9a41605-b57b-4283-9122-984cd61a83f0",
							"createdon": "2023-09-23T23:50:48Z",
							"address2_country": null,
							"description": "Dynamics 365",
							"address2_addressid": null,
							"address1_shippingmethodcode": null,
							"address1_postofficebox": null,
							"address1_upszone": null,
							"address1_addresstypecode": null,
							"address1_country": null,
							"entityimageid": null,
							"entityimage": null,
							"_createdby_value": "a9a41605-b57b-4283-9122-984cd61a83f0",
							"address1_telephone3": null,
							"address1_telephone2": null,
							"address1_city": null,
							"address1_telephone1": null
						}
					}
				]
			}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				data "powerplatform_solutions" "all" {
					environment_name = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", "ProductivityToolsAnchor"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", "ProductivityTools"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", "2023-10-10T08:09:56Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", "2023-10-10T08:09:58Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", "2023-10-10T08:09:56Z"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", "9.2.1.1020"),
					resource.TestCheckResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", "70edca66-e4c2-4384-92e0-4300465c1894"),
				),
			},
		},
	})
}

func TestAccSolutionsDataSource_Validate_Read(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name     = "testaccsolutionsdatasource"
					location         = "europe"
					language_code    = "1033"
					currency_code    = "USD"
					environment_type = "Sandbox"
					domain           = "testaccsolutionsdatasource"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}

				data "powerplatform_solutions" "all" {
					environment_name = powerplatform_environment.development.environment_name
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.created_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.modified_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.install_time", regexp.MustCompile(powerplatform_helpers.TimeRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.is_managed", regexp.MustCompile(`^(true|false)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_solutions.all", "solutions.0.id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
				),
			},
		},
	})
}
