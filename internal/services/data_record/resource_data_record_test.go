// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccDataRecordResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
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

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.test_env]
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

					depends_on = [time_sleep.wait_for_dataverse]
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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

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
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `
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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				
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
				Config: `
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
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `
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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `

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
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
					  contactid = "00000000-0000-0000-0000-000000000010"
					  firstname          = "contact1"
					}

				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
				      contactid = "00000000-0000-0000-0000-000000000020"
					  firstname          = "contact2"
					}

					depends_on = [powerplatform_data_record.data_record_sample_contact1]
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000030"
					  firstname          = "contact3"
					}

					depends_on = [powerplatform_data_record.data_record_sample_contact2]
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
				Config: `
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

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000020"
					  firstname          = "contact2"
					}

				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000030"
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
				Config: `
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

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000020"
					  firstname          = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact3" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000030"
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
					mocks.TestStateValueMatch(primarycontactidStep1, primarycontactidStep2, func(a, b *mocks.StateValue) error {
						if a.Value == b.Value {
							return errors.New("expected primarycontactid from before and after change are equal, but a change was expected")
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

func TestAccDataRecordResource_Validate_Remove_Relationship_Item_Disassociate_Keeps_Related_Row(t *testing.T) {
	environmentConfig := `
				resource "powerplatform_environment" "test_env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
					timeouts = {
					  create = "60m"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000010"
					  firstname = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "contact"
					columns = {
					  contactid = "00000000-0000-0000-0000-000000000020"
					  firstname = "contact2"
					}

					depends_on = [powerplatform_data_record.data_record_sample_contact1]
				}
	`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: environmentConfig + `
				resource "powerplatform_data_record" "data_record_account" {
					environment_id                = powerplatform_environment.test_env.id
					table_logical_name            = "account"
					relationship_removal_behavior = "disassociate"
					columns = {
						name = "Sample Account"

						contact_customer_accounts = toset([
							{
								table_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
								data_record_id     = powerplatform_data_record.data_record_sample_contact1.id
							},
							{
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id     = powerplatform_data_record.data_record_sample_contact2.id
							}
						])
					}
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"contact_customer_accounts": knownvalue.SetSizeExact(2),
						})),
				},
			},
			{
				// Remove contact1 from the relationship only. contact1 stays in the configuration,
				// so the related row itself must still exist in Dataverse after the apply.
				Config: environmentConfig + `
				resource "powerplatform_data_record" "data_record_account" {
					environment_id                = powerplatform_environment.test_env.id
					table_logical_name            = "account"
					relationship_removal_behavior = "disassociate"
					columns = {
						name = "Sample Account"

						contact_customer_accounts = toset([
							{
								table_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
								data_record_id     = powerplatform_data_record.data_record_sample_contact2.id
							}
						])
					}
				}

				data "powerplatform_data_records" "removed_contact" {
					environment_id    = powerplatform_environment.test_env.id
					entity_collection = "contacts"
					filter            = "contactid eq 00000000-0000-0000-0000-000000000010"
					select            = ["contactid", "firstname"]

					depends_on = [powerplatform_data_record.data_record_account]
				}

				data "powerplatform_data_records" "kept_contact" {
					environment_id    = powerplatform_environment.test_env.id
					entity_collection = "contacts"
					filter            = "contactid eq 00000000-0000-0000-0000-000000000020"
					select            = ["contactid", "firstname"]

					depends_on = [powerplatform_data_record.data_record_account]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
						knownvalue.MapPartial(map[string]knownvalue.Check{
							"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.MapExact(map[string]knownvalue.Check{
									"table_logical_name": knownvalue.StringExact("contact"),
									"data_record_id":     knownvalue.StringExact("00000000-0000-0000-0000-000000000020"),
								}),
							}),
						})),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.kept_contact", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.removed_contact", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.removed_contact", "rows.0.contactid", "00000000-0000-0000-0000-000000000010"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.removed_contact", "rows.0.firstname", "contact1"),
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
			default:
				contactId = "11111111-0000-0000-0000-000000000000"
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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
				Config: `
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
				Config: `
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

// removeRelationshipItemMock is an in-memory Dataverse for an account with a contact_customer_accounts relationship.
// It records which contacts are linked to the account, which contact rows were deleted and which links were removed
// with a $ref DELETE while the second test step is applied.
type removeRelationshipItemMock struct {
	mu                    sync.Mutex
	linkedContacts        map[string]bool
	deletedContacts       map[string]bool
	disassociatedContacts map[string]bool
	applyingSecondStep    bool
	// deletedRowsStayReadable keeps serving a deleted contact on GET, so a test can assert on the requests the
	// provider sent without the refresh of the still-configured contact resource failing afterwards.
	deletedRowsStayReadable bool
}

const (
	removeRelationshipItemContact1 = "00000000-0000-0000-0000-000000000010"
	removeRelationshipItemContact2 = "00000000-0000-0000-0000-000000000011"
)

func registerRemoveRelationshipItemMocks() *removeRelationshipItemMock {
	const envHost = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com"
	const fixtures = "tests/resource/Validate_Update_Relationships/"
	const accountId = "00000000-0000-0000-0000-000000000020"

	m := &removeRelationshipItemMock{
		linkedContacts:        map[string]bool{},
		deletedContacts:       map[string]bool{},
		disassociatedContacts: map[string]bool{},
	}

	sortedLinkedContacts := func() []string {
		ids := make([]string, 0, len(m.linkedContacts))
		for id := range m.linkedContacts {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		return ids
	}

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File(fixtures+"get_environment_00000000-0000-0000-0000-000000000001.json").String()))

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File(fixtures+"get_entitydefinition_contact.json").String()))

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/EntityDefinitions%28LogicalName=%27account%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File(fixtures+"get_entitydefinition_account.json").String()))

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/EntityDefinitions(LogicalName='account')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File(fixtures+"get_entitydefinition_account.json").String()))

	httpmock.RegisterResponder("POST", envHost+`/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			contactId := ""
			switch string(body) {
			case `{"firstname":"contact1"}`:
				contactId = removeRelationshipItemContact1
			case `{"firstname":"contact2"}`:
				contactId = removeRelationshipItemContact2
			default:
				return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
			}
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("%s/api/data/v9.2/contacts(%s)", envHost, contactId))
			return resp, nil
		})

	for _, contactId := range []string{removeRelationshipItemContact1, removeRelationshipItemContact2} {
		httpmock.RegisterResponder("GET", fmt.Sprintf("%s/api/data/v9.2/contacts%%28%s%%29", envHost, contactId),
			func(req *http.Request) (*http.Response, error) {
				m.mu.Lock()
				defer m.mu.Unlock()
				if m.deletedContacts[contactId] && !m.deletedRowsStayReadable {
					resp := httpmock.NewStringResponse(http.StatusNotFound, "")
					resp.Request = req
					return resp, nil
				}
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("%sget_contact_%s.json", fixtures, contactId)).String()), nil
			})
	}

	httpmock.RegisterResponder("POST", envHost+`/api/data/v9.2/accounts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("%s/api/data/v9.2/accounts(%s)", envHost, accountId))
			return resp, nil
		})

	httpmock.RegisterResponder("PATCH", envHost+`/api/data/v9.2/accounts%28`+accountId+`%29`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", fmt.Sprintf("%s/api/data/v9.2/accounts(%s)", envHost, accountId))
			return resp, nil
		})

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/accounts%28`+accountId+`%29`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File(fixtures+"get_account_00000000-0000-0000-0000-000000000020_1.json").String()))

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/accounts%28`+accountId+`%29/contact_customer_accounts?$select=createdon`,
		func(req *http.Request) (*http.Response, error) {
			m.mu.Lock()
			defer m.mu.Unlock()
			rows := []map[string]any{}
			for _, id := range sortedLinkedContacts() {
				rows = append(rows, map[string]any{"createdon": "2024-06-12T14:13:17Z", "contactid": id})
			}
			return httpmock.NewJsonResponse(http.StatusOK, map[string]any{"value": rows})
		})

	httpmock.RegisterResponder("GET", envHost+`/api/data/v9.2/accounts%28`+accountId+`%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			m.mu.Lock()
			defer m.mu.Unlock()
			refs := []map[string]any{}
			for _, id := range sortedLinkedContacts() {
				refs = append(refs, map[string]any{"@odata.id": fmt.Sprintf("%s/api/data/v9.2/contacts(%s)", envHost, id)})
			}
			return httpmock.NewJsonResponse(http.StatusOK, map[string]any{"value": refs})
		})

	httpmock.RegisterResponder("POST", envHost+`/api/data/v9.2/accounts%28`+accountId+`%29/contact_customer_accounts/$ref`,
		func(req *http.Request) (*http.Response, error) {
			body := map[string]string{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return nil, err
			}
			contactId := regexp.MustCompile(`contacts\(([0-9a-f-]+)\)$`).FindStringSubmatch(body["@odata.id"])
			if contactId == nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
			}
			m.mu.Lock()
			defer m.mu.Unlock()
			m.linkedContacts[contactId[1]] = true
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	// Disassociate: DELETE <account>/contact_customer_accounts(<contact>)/$ref only removes the link.
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`/api/data/v9\.2/accounts(?:\(|%28)`+accountId+`(?:\)|%29)/contact_customer_accounts(?:\(|%28)([0-9a-f-]+)(?:\)|%29)/\$ref$`),
		func(req *http.Request) (*http.Response, error) {
			contactId := httpmock.MustGetSubmatch(req, 1)
			m.mu.Lock()
			defer m.mu.Unlock()
			if m.applyingSecondStep {
				m.disassociatedContacts[contactId] = true
			}
			delete(m.linkedContacts, contactId)
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	// Deleting a contact row removes the row itself, and with it the link to the account.
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`/api/data/v9\.2/contacts(?:\(|%28)([0-9a-f-]+)(?:\)|%29)$`),
		func(req *http.Request) (*http.Response, error) {
			contactId := httpmock.MustGetSubmatch(req, 1)
			m.mu.Lock()
			defer m.mu.Unlock()
			if m.applyingSecondStep {
				m.deletedContacts[contactId] = true
			}
			delete(m.linkedContacts, contactId)
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`/api/data/v9\.2/accounts(?:\(|%28)`+accountId+`(?:\)|%29)$`),
		httpmock.NewStringResponder(http.StatusNoContent, ""))

	return m
}

// removeRelationshipItemConfig returns a configuration with two contacts and an account that links the given contacts
// through contact_customer_accounts. behaviorConfig is inserted into the account resource as is.
func removeRelationshipItemConfig(behaviorConfig string, linkedContactResources ...string) string {
	links := ""
	for _, contactResource := range linkedContactResources {
		links += `
							{
								table_logical_name = powerplatform_data_record.` + contactResource + `.table_logical_name
								data_record_id     = powerplatform_data_record.` + contactResource + `.id
							},`
	}
	return `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname = "contact1"
					}
				}

				resource "powerplatform_data_record" "data_record_sample_contact2" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
					  firstname = "contact2"
					}
				}

				resource "powerplatform_data_record" "data_record_account" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "account"
					` + behaviorConfig + `
					columns = {
						name = "Sample Account"

						contact_customer_accounts = toset([` + links + `
						])
					}
				}
	`
}

func removeRelationshipItemFirstStep(m *removeRelationshipItemMock, behaviorConfig string) resource.TestStep {
	return resource.TestStep{
		Config: removeRelationshipItemConfig(behaviorConfig, "data_record_sample_contact1", "data_record_sample_contact2"),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
				knownvalue.MapPartial(map[string]knownvalue.Check{
					"contact_customer_accounts": knownvalue.SetSizeExact(2),
				})),
		},
		Check: func(_ *terraform.State) error {
			m.mu.Lock()
			defer m.mu.Unlock()
			m.applyingSecondStep = true
			return nil
		},
	}
}

var removeRelationshipItemOnlyContact2Linked = statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("columns"),
	knownvalue.MapPartial(map[string]knownvalue.Check{
		"contact_customer_accounts": knownvalue.SetExact([]knownvalue.Check{
			knownvalue.MapExact(map[string]knownvalue.Check{
				"table_logical_name": knownvalue.StringExact("contact"),
				"data_record_id":     knownvalue.StringExact(removeRelationshipItemContact2),
			}),
		}),
	}))

func TestUnitDataRecordResource_Validate_Remove_Relationship_Item_Disassociate_Keeps_Related_Row(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	m := registerRemoveRelationshipItemMocks()
	const behaviorConfig = `relationship_removal_behavior = "disassociate"`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			removeRelationshipItemFirstStep(m, behaviorConfig),
			{
				// contact1 is removed from the relationship but stays in the configuration.
				Config: removeRelationshipItemConfig(behaviorConfig, "data_record_sample_contact2"),
				ConfigStateChecks: []statecheck.StateCheck{
					removeRelationshipItemOnlyContact2Linked,
					statecheck.ExpectKnownValue("powerplatform_data_record.data_record_account", tfjsonpath.New("relationship_removal_behavior"), knownvalue.StringExact("disassociate")),
				},
				Check: func(_ *terraform.State) error {
					m.mu.Lock()
					defer m.mu.Unlock()
					m.applyingSecondStep = false
					if len(m.deletedContacts) != 0 {
						return fmt.Errorf("expected removing a contact from contact_customer_accounts to keep the contact row, but these contacts were deleted: %v", m.deletedContacts)
					}
					if !m.disassociatedContacts[removeRelationshipItemContact1] || len(m.disassociatedContacts) != 1 {
						return fmt.Errorf("expected only contact1 to be disassociated from the account with a $ref DELETE, got: %v", m.disassociatedContacts)
					}
					return nil
				},
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Remove_Relationship_Item_Default_Deletes_Related_Row(t *testing.T) {
	for name, behaviorConfig := range map[string]string{
		"not set":                        "",
		"explicit delete_related_record": `relationship_removal_behavior = "delete_related_record"`,
	} {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// The default behavior deletes the contact1 row while the contact1 resource is still configured. The
			// mock keeps serving that row on GET so the step can assert on the requests that were sent.
			m := registerRemoveRelationshipItemMocks()
			m.deletedRowsStayReadable = true

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					removeRelationshipItemFirstStep(m, behaviorConfig),
					{
						// contact1 is removed from the relationship but stays in the configuration.
						Config: removeRelationshipItemConfig(behaviorConfig, "data_record_sample_contact2"),
						ConfigStateChecks: []statecheck.StateCheck{
							removeRelationshipItemOnlyContact2Linked,
						},
						Check: func(_ *terraform.State) error {
							m.mu.Lock()
							defer m.mu.Unlock()
							m.applyingSecondStep = false
							if !m.deletedContacts[removeRelationshipItemContact1] || len(m.deletedContacts) != 1 {
								return fmt.Errorf("expected the default behavior to delete only the contact1 row, got deleted contacts: %v", m.deletedContacts)
							}
							if len(m.disassociatedContacts) != 0 {
								return fmt.Errorf("expected the default behavior to delete the related row and not send a $ref DELETE, but these contacts were disassociated: %v", m.disassociatedContacts)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestUnitDataRecordResource_Validate_Relationship_Removal_Behavior_Invalid_Value(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_data_record" "data_record_account" {
					environment_id                = "00000000-0000-0000-0000-000000000001"
					table_logical_name            = "account"
					relationship_removal_behavior = "delete"
					columns = {
						name = "Sample Account"
					}
				}
				`,
				ExpectError: regexp.MustCompile(`(?s)relationship_removal_behavior.*value must be one of`),
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Disable_On_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	readBodyBuff := func(req *http.Request) (string, error) {
		r, err := req.GetBody()
		if err != nil {
			return "", err
		}
		defer r.Close()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	var contactDisabledExecuted bool
	var businessUnitDisabledExecuted bool

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27mailbox%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_entitydefinition_mailbox.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/mailboxes%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_mailbox_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/mailboxes`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/mailboxes(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='mailbox')/Attributes?$select=LogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_entitydefinition_mailbox_attributes.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/mailboxes%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			body, err := readBodyBuff(req)
			if err != nil {
				return nil, err
			}

			if body == "{\"statecode\":1}" {
				contactDisabledExecuted = true
			}

			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/mailboxes(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27businessunit%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_entitydefinition_businessunit.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits(00000000-0000-0000-0000-000000000001)")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='businessunit')/Attributes?$select=LogicalName",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_entitydefinition_businessunit_attributes.json").String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits%2800000000-0000-0000-0000-000000000001%29`,
		func(req *http.Request) (*http.Response, error) {
			body, err := readBodyBuff(req)
			if err != nil {
				return nil, err
			}

			if body == "{\"isdisabled\":true}" {
				businessUnitDisabledExecuted = true
			}
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits(00000000-0000-0000-0000-000000000001)")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits%2800000000-0000-0000-0000-000000000001%29`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Disable_On_Delete/get_businessunit_00000000-0000-0000-0000-000000000001.json").String())
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/businessunits(00000000-0000-0000-0000-000000000001)")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_data_record" "mailbox" {
					disable_on_destroy = true
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "mailbox"
					columns = {
						name         = "contoso"
						emailaddress = "contoso@contoso.com"
					}
				}
				
				resource "powerplatform_data_record" "business_unit" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "businessunit"
					disable_on_destroy = true
					columns = {
						name       = "test"
					}
				}`,
			},
			{
				// commenting out the resource to trigger the delete
				Config: ` 
				// resource "powerplatform_data_record" "mailbox" {
				// 	disable_on_destroy = true
				// 	environment_id     = "00000000-0000-0000-0000-000000000001"
				// 	table_logical_name = "mailbox"
				// 	columns = {
				// 		name         = "my mailbox"
				// 		emailaddress = "contoso@contoso.com"
				// 	}
				// }
				
				// resource "powerplatform_data_record" "business_unit" {
				// 	environment_id     = "00000000-0000-0000-0000-000000000001"
				// 	table_logical_name = "businessunit"
				// 	disable_on_destroy = true
				// 	columns = {
				// 		name       = "test"
				// 	}
				// }`,
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						if !contactDisabledExecuted {
							return errors.New("expected mailbox to be disabled on delete")
						}
						return nil
					},
					func(_ *terraform.State) error {
						if !businessUnitDisabledExecuted {
							return errors.New("expected business unit to be disabled on delete")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccDataRecordResource_Validate_Disable_On_Delete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               false,
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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

				data "powerplatform_data_records" "root_business_unit" {
					environment_id    = powerplatform_environment.test_env.id
					entity_collection = "businessunits"
					filter            = "parentbusinessunitid eq null"
					select            = ["name"]
				}

				resource "powerplatform_data_record" "mailbox" {
					disable_on_destroy = true
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "mailbox"
					columns = {
						name         = "my mailbox"
						emailaddress = "contoso@contoso.com"
					}
				}
				
				resource "powerplatform_data_record" "business_unit" {
					environment_id     = powerplatform_environment.test_env.id
					table_logical_name = "businessunit"
					disable_on_destroy = true
					columns = {
						name       = "test"
				        parentbusinessunitid = {
							table_logical_name = "businessunit"
							data_record_id     =  one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
						}
					}
				}`,
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var getEnvironmentCallCount = 0
	// The data_record client's internal getEnvironment() (exact URL with api-version only).
	// Step 1 (Create + Reads) makes ~5 calls that succeed; step 2 refresh returns 404
	// to simulate the environment being deleted out-of-band.
	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			getEnvironmentCallCount++
			if getEnvironmentCallCount <= 5 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/Validate_Read_ParentDeleted/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	// The EnvironmentClient's GetEnvironment() returns 404 (URL includes $expand query params).
	// This is only called during tier-2 parent check when the data record read fails.
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001\?.+expand`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/Validate_Read_ParentDeleted/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create/get_contact_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
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
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the data record successfully.
				Config: `
				resource "powerplatform_data_record" "data_record_sample_contact1" {
					environment_id     = "00000000-0000-0000-0000-000000000001"
					table_logical_name = "contact"
					columns = {
						firstname     = "John"
						lastname      = "Doe"
						telephone1    = "555-555-5555"
						emailaddress1 = "johndoe@contoso.com"
						anniversary   = "2024-04-10"
						annualincome  = 1234.56
						birthdate     = "2024-04-10"
						description   = "This is the description of the the terraform \n\nsample contact"
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_data_record.data_record_sample_contact1", "table_logical_name", "contact"),
				),
			},
			{
				// Step 2: refresh when the parent environment has been deleted out-of-band.
				// The GetDataRecord call fails because getEnvironment returns 404; the provider detects the parent is gone and removes the resource from state.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestUnitDataRecordResource_Validate_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Import/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions%28LogicalName=%27contact%27%29#$select=PrimaryIdAttribute,LogicalCollectionName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Import/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts%2800000000-0000-0000-0000-000000000010%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Import/get_contact_00000000-0000-0000-0000-000000000010.json").String()), nil
		})

	httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("OData-EntityId", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000010)")
			return resp, nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://00000000-0000-0000-0000-000000000001\.crm4\.dynamics\.com/api/data/v9\.2/([a-zA-Z]+)`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	config := `
	resource "powerplatform_data_record" "data_record_sample_contact1" {
		environment_id     = "00000000-0000-0000-0000-000000000001"
		table_logical_name = "contact"
		columns = {
			firstname = "John"
			lastname  = "Doe"
		}
	}`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config:        config,
				ResourceName:  "powerplatform_data_record.data_record_sample_contact1",
				ImportState:   true,
				ImportStateId: "00000000-0000-0000-0000-000000000001/contact(00000000-0000-0000-0000-000000000010)",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					state := states[0]
					if state.ID != "00000000-0000-0000-0000-000000000010" {
						return fmt.Errorf("expected id '00000000-0000-0000-0000-000000000010', got '%s'", state.ID)
					}
					if got := state.Attributes["environment_id"]; got != "00000000-0000-0000-0000-000000000001" {
						return fmt.Errorf("expected environment_id '00000000-0000-0000-0000-000000000001', got '%s'", got)
					}
					if got := state.Attributes["table_logical_name"]; got != "contact" {
						return fmt.Errorf("expected table_logical_name 'contact', got '%s'", got)
					}
					return nil
				},
			},
			{
				Config:        config,
				ResourceName:  "powerplatform_data_record.data_record_sample_contact1",
				ImportState:   true,
				ImportStateId: "00000000-0000-0000-0000-000000000010",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}
