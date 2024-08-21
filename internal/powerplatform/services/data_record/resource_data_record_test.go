// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func TestAccDataRecordResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  telephone1         = "555-555-5555"
					  emailaddress1      = "johndoe@contoso.com"
					  anniversary        = "2024-04-10"
					  annualincome       = 1234.56
					  birthdate          = "2024-04-10"
					  description        = "This is the description of the the terraform \n\nsample contact"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
							creditonhold        = false
							address1_latitude   = 47.63
							description         = "This is the description of the sample account"
							revenue             = 5000000
							accountcategorycode = 1
						
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John"),
							"lastname":      knownvalue.StringExact("Doe"),
							"telephone1":    knownvalue.StringExact("555-555-5555"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com"),
							"anniversary":   knownvalue.StringExact("2024-04-10"),
							"annualincome":  knownvalue.Float64Exact(1234.56),
							"birthdate":     knownvalue.StringExact("2024-04-10"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact"),
						})),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name":                knownvalue.StringExact("Sample Account"),
							"creditonhold":        knownvalue.Bool(false),
							"address1_latitude":   knownvalue.Float64Exact(47.63),
							"description":         knownvalue.StringExact("This is the description of the sample account"),
							"revenue":             knownvalue.Float64Exact(5000000),
							"accountcategorycode": knownvalue.Int64Exact(1),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
								}),
							}),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_account", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_account", "table_logical_name", "account"),
				),
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27account%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='account')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts?$select=createdon`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_account_00000000-0000-0000-0000-000000000020_contact_customer_accounts.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_account_00000000-0000-0000-0000-000000000020_ref.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_contact_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_account_00000000-0000-0000-0000-000000000020.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  telephone1         = "555-555-5555"
					  emailaddress1      = "johndoe@contoso.com"
					  anniversary        = "2024-04-10"
					  annualincome       = 1234.56
					  birthdate          = "2024-04-10"
					  description        = "This is the description of the the terraform \n\nsample contact"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
							creditonhold        = false
							address1_latitude   = 47.63
							description         = "This is the description of the sample account"
							revenue             = 5000000
							accountcategorycode = 1
						
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John"),
							"lastname":      knownvalue.StringExact("Doe"),
							"telephone1":    knownvalue.StringExact("555-555-5555"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com"),
							"anniversary":   knownvalue.StringExact("2024-04-10"),
							"annualincome":  knownvalue.Float64Exact(1234.56),
							"birthdate":     knownvalue.StringExact("2024-04-10"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact"),
						})),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name":                knownvalue.StringExact("Sample Account"),
							"creditonhold":        knownvalue.Bool(false),
							"address1_latitude":   knownvalue.Float64Exact(47.63),
							"description":         knownvalue.StringExact("This is the description of the sample account"),
							"revenue":             knownvalue.Float64Exact(5000000),
							"accountcategorycode": knownvalue.Int64Exact(1),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
								}),
							}),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_account", "environment_id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_account", "table_logical_name", "account"),
				),
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  telephone1         = "555-555-5555"
					  emailaddress1      = "johndoe@contoso.com"
					  anniversary        = "2024-04-10"
					  annualincome       = 1234.56
					  birthdate          = "2024-04-10"
					  description        = "This is the description of the the terraform \n\nsample contact"
					}
				}`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John"),
							"lastname":      knownvalue.StringExact("Doe"),
							"telephone1":    knownvalue.StringExact("555-555-5555"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com"),
							"anniversary":   knownvalue.StringExact("2024-04-10"),
							"annualincome":  knownvalue.Float64Exact(1234.56),
							"birthdate":     knownvalue.StringExact("2024-04-10"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact"),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
			},
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "John1"
					  lastname           = "Doe1"
					  telephone1         = "555-555-55551"
					  emailaddress1      = "johndoe@contoso.com1"
					  anniversary        = "2024-04-11"
					  annualincome       = 1234.51
					  birthdate          = "2024-04-11"
					  description        = "This is the description of the the terraform \n\nsample contact1"
					}
				}`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John1"),
							"lastname":      knownvalue.StringExact("Doe1"),
							"telephone1":    knownvalue.StringExact("555-555-55551"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com1"),
							"anniversary":   knownvalue.StringExact("2024-04-11"),
							"annualincome":  knownvalue.Float64Exact(1234.51),
							"birthdate":     knownvalue.StringExact("2024-04-11"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact1"),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/get_entitydefinition_contact.json").String()), nil
		})

	var inx = 0
	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update/get_contact_00000000-0000-0000-0000-000000000010_%d.json", inx)).String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  telephone1         = "555-555-5555"
					  emailaddress1      = "johndoe@contoso.com"
					  anniversary        = "2024-04-10"
					  annualincome       = 1234.56
					  birthdate          = "2024-04-10"
					  description        = "This is the description of the the terraform \n\nsample contact"
					}
				}`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John"),
							"lastname":      knownvalue.StringExact("Doe"),
							"telephone1":    knownvalue.StringExact("555-555-5555"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com"),
							"anniversary":   knownvalue.StringExact("2024-04-10"),
							"annualincome":  knownvalue.Float64Exact(1234.56),
							"birthdate":     knownvalue.StringExact("2024-04-10"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact"),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "John1"
					  lastname           = "Doe1"
					  telephone1         = "555-555-55551"
					  emailaddress1      = "johndoe@contoso.com1"
					  anniversary        = "2024-04-11"
					  annualincome       = 1234.51
					  birthdate          = "2024-04-11"
					  description        = "This is the description of the the terraform \n\nsample contact1"
					}
				}`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"firstname":     knownvalue.StringExact("John1"),
							"lastname":      knownvalue.StringExact("Doe1"),
							"telephone1":    knownvalue.StringExact("555-555-55551"),
							"emailaddress1": knownvalue.StringExact("johndoe@contoso.com1"),
							"anniversary":   knownvalue.StringExact("2024-04-11"),
							"annualincome":  knownvalue.Float64Exact(1234.51),
							"birthdate":     knownvalue.StringExact("2024-04-11"),
							"description":   knownvalue.StringExact("This is the description of the the terraform \n\nsample contact1"),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Delete_Relationships(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  emailaddress1      = "johndoe@contoso.com"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
							
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{

					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{

					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
						})),
				},
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Delete_Relationships(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27account%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='account')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts?$select=createdon`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_account_00000000-0000-0000-0000-000000000020_contact_customer_accounts.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_account_00000000-0000-0000-0000-000000000020_ref.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_contact_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Delete/get_account_00000000-0000-0000-0000-000000000020.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "John"
					  lastname           = "Doe"
					  emailaddress1      = "johndoe@contoso.com"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
							
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{

					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(helpers.GuidRegex)),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsUnitProviderConfig + `

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
						})),
				},
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Update_Relationships(t *testing.T) {

	var primarycontactidStep1 = &mocks.StateValue{}
	var primarycontactidStep2 = &mocks.StateValue{}

	var contact1Id = &mocks.StateValue{}
	var contact2Id = &mocks.StateValue{}
	var contact3Id = &mocks.StateValue{}

	var contactAtIndex1Step1 = &mocks.StateValue{}
	var contactAtIndex2Step1 = &mocks.StateValue{}

	var contactAtIndex1Step2 = &mocks.StateValue{}
	var contactAtIndex2Step2 = &mocks.StateValue{}

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
						
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("id"), mocks.GetStateValue(contact1Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact2", tfjsonpath.New("id"), mocks.GetStateValue(contact2Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact3", tfjsonpath.New("id"), mocks.GetStateValue(contact3Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     mocks.GetStateValue(primarycontactidStep1),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mocks.GetStateValue(contactAtIndex1Step1),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mocks.GetStateValue(contactAtIndex2Step1),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"

							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     mocks.GetStateValue(primarycontactidStep2),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mocks.GetStateValue(contactAtIndex1Step2),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mocks.GetStateValue(contactAtIndex2Step2),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = powerplatform_environment.test_env.id
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"

							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							]
						}
					}
				`,

				Check: resource.ComposeAggregateTestCheckFunc(
					mocks.TestStateValueMatch(primarycontactidStep1, primarycontactidStep2, func(a, b *mocks.StateValue) error {
						if a.Value == b.Value {
							return fmt.Errorf("expected primarycontactid from before and after change are equal, but a change was expected. '%s' == '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mocks.TestStateValueMatch(contactAtIndex1Step1, contact1Id, func(a, b *mocks.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step1 expected that the first item in contact_customer_accounts will be contact1. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mocks.TestStateValueMatch(contactAtIndex2Step1, contact2Id, func(a, b *mocks.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step1 expected that the second item in contact_customer_accounts will be contact2. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mocks.TestStateValueMatch(contactAtIndex1Step2, contact2Id, func(a, b *mocks.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step2 expected that the first item in contact_customer_accounts will be contact2. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mocks.TestStateValueMatch(contactAtIndex2Step2, contact3Id, func(a, b *mocks.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step2 expected that the second item in contact_customer_accounts will be contact3. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Update_Relationships(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27account%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='account')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_entitydefinition_account.json").String()), nil
		})

	contactCustomerAccountsInx := 0
	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts?$select=createdon`,
		func(req *http.Request) (*http.Response, error) {
			contactCustomerAccountsInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_Relationships/get_account_00000000-0000-0000-0000-000000000020_contact_customer_accounts_%d.json", contactCustomerAccountsInx)).String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	contactCustomerAccountsRefInx := 0
	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			contactCustomerAccountsRefInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_Relationships/get_account_00000000-0000-0000-0000-000000000020_ref_%d.json", contactCustomerAccountsRefInx)).String()), nil
		})

	accountInx := 0
	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29`,
		func(req *http.Request) (*http.Response, error) {
			accountInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Update_Relationships/get_account_00000000-0000-0000-0000-000000000020_%d.json", accountInx)).String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_contact_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000011%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_contact_00000000-0000-0000-0000-000000000011.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000012%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update_Relationships/get_contact_00000000-0000-0000-0000-000000000012.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts%2800000000-0000-0000-0000-000000000020%29`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {

			bodyAsBytes := make([]byte, req.ContentLength)
			_, err := req.Body.Read(bodyAsBytes)
			if err != nil {
				panic(err)
			}
			bodyAsString := string(bodyAsBytes)

			contactId := ""
			switch bodyAsString {
			case `{"firstname":"contact1"}`:
				contactId = "00000000-0000-0000-0000-000000000010"
			case `{"firstname":"contact2"}`:
				contactId = "00000000-0000-0000-0000-000000000011"
			case `{"firstname":"contact3"}`:
				contactId = "00000000-0000-0000-0000-000000000012"
			}

			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(%s)", contactId))
			return resp, nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000020)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"
						
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								}
							]
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("id"), knownvalue.StringExact("00000000-0000-0000-0000-000000000010")),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact2", tfjsonpath.New("id"), knownvalue.StringExact("00000000-0000-0000-0000-000000000011")),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact3", tfjsonpath.New("id"), knownvalue.StringExact("00000000-0000-0000-0000-000000000012")),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000010"),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000010"),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000011"),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"

							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							]
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000011"),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000011"),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000012"),
								}),
							}),
						})),
				},
			},
			{
				Config: provider.TestsUnitProviderConfig + `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname          = "contact3"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
						environment_id     = "00000000-0000-0000-0000-000000000001"
						table_logical_name = "account"
						columns = {
							name                = "Sample Account"

							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
							}

							contact_customer_accounts = [
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							]
						}
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000011"),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000011"),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000012"),
								}),
							}),
						})),
				},
			},
		},
	})
}
