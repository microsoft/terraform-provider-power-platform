// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
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
							address1_latitude   = 47.63
							description         = "This is the description of the sample account"
							revenue             = 5000000
							accountcategorycode = 1
						
							primarycontactid = {
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = toset([
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							])
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
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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

func TestAccDataRecordResource_Validate_Update(t *testing.T) {
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
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
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
					resource.TestMatchResourceAttr("powerplatform_data_record.data_record_sample_contact1", "environment_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
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
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
							}
						
							contact_customer_accounts = toset([
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								}
							])
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{

					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("Sample Account"),
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringRegexp(regexp.MustCompile(powerplatform_helpers.GuidRegex)),
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
						})),
				},
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Update_Relationships(t *testing.T) {

	var primarycontactidStep1 = &mock_helpers.StateValue{}
	var primarycontactidStep2 = &mock_helpers.StateValue{}

	var contact1Id = &mock_helpers.StateValue{}
	var contact2Id = &mock_helpers.StateValue{}
	var contact3Id = &mock_helpers.StateValue{}

	var contactAtIndex1Step1 = &mock_helpers.StateValue{}
	var contactAtIndex2Step1 = &mock_helpers.StateValue{}

	var contactAtIndex1Step2 = &mock_helpers.StateValue{}
	var contactAtIndex2Step2 = &mock_helpers.StateValue{}

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

							contact_customer_accounts = toset([
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								}
							])
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact1", tfjsonpath.New("id"), mock_helpers.GetStateValue(contact1Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact2", tfjsonpath.New("id"), mock_helpers.GetStateValue(contact2Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_sample_contact3", tfjsonpath.New("id"), mock_helpers.GetStateValue(contact3Id)),
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     mock_helpers.GetStateValue(primarycontactidStep1),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mock_helpers.GetStateValue(contactAtIndex1Step1),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mock_helpers.GetStateValue(contactAtIndex2Step1),
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

							contact_customer_accounts = toset([
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							])
						}
					}
				`,

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"primarycontactid": knownvalue.MapExact(map[string]knownvalue.Check{
								"table_logical_name": knownvalue.StringExact("contact"),
								"data_record_id":     mock_helpers.GetStateValue(primarycontactidStep2),
							}),
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								0: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mock_helpers.GetStateValue(contactAtIndex1Step2),
								}),
								1: knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     mock_helpers.GetStateValue(contactAtIndex2Step2),
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

							contact_customer_accounts = toset([
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
								},
								{
									table_logical_name = powerplatform_data_record.data_record_sample_contact3.table_logical_name
									data_record_id      = powerplatform_data_record.data_record_sample_contact3.id
								}
							])
						}
					}
				`,

				Check: resource.ComposeAggregateTestCheckFunc(
					mock_helpers.TestStateValueMatch(primarycontactidStep1, primarycontactidStep2, func(a, b *mock_helpers.StateValue) error {
						if a.Value == b.Value {
							return fmt.Errorf("expected primarycontactid from before and after change are equal, but a change was expected. '%s' == '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mock_helpers.TestStateValueMatch(contactAtIndex1Step1, contact1Id, func(a, b *mock_helpers.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step1 expected that the first item in contact_customer_accounts will be contact1. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mock_helpers.TestStateValueMatch(contactAtIndex2Step1, contact2Id, func(a, b *mock_helpers.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step1 expected that the second item in contact_customer_accounts will be contact2. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mock_helpers.TestStateValueMatch(contactAtIndex1Step2, contact2Id, func(a, b *mock_helpers.StateValue) error {
						if a.Value != b.Value {
							return fmt.Errorf("step2 expected that the first item in contact_customer_accounts will be contact2. '%s' != '%s'", a.Value, b.Value)
						}
						return nil
					}),
					mock_helpers.TestStateValueMatch(contactAtIndex2Step2, contact3Id, func(a, b *mock_helpers.StateValue) error {
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
