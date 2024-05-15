// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccDataRecordResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mock_helpers.TestName() + `"
					location         = "europe"
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
							address1_latitude   = 47.639583
							description         = "This is the description of the sample account"
							revenue             = 5000000
							accountcategorycode = 1
						
							primarycontactid = {
								entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
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
							"address1_latitude":   knownvalue.Float64Exact(47.639583),
							"description":         knownvalue.StringExact("This is the description of the sample account"),
							"revenue":             knownvalue.Float64Exact(5000000),
							"accountcategorycode": knownvalue.Int64Exact(1),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"entity_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":      knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.ListExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"entity_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":      knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
								}),
							}),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_account", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_account", "table_logical_name", "account"),
				),
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Delete_Relationships(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mock_helpers.TestName() + `"
					location         = "europe"
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
								entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = [
								{
									entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
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
								"entity_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":      knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.ListExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"entity_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":      knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
								}),
							}),
						})),
				},
			},
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mock_helpers.TestName() + `"
					location         = "europe"
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
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{

					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
							//"primarycontactid":          knownvalue.MapSizeExact(0),
							//"contact_customer_accounts": knownvalue.ListSizeExact(0),
						})),
				},
			},
		},
	})
}
