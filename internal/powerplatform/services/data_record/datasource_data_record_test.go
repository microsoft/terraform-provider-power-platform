// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jarcoal/httpmock"

	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func BootstrapDataRecordTest(name string) string {
	return `
resource "powerplatform_environment" "data_env" {
	display_name     = "` + name + `"
	location         = "unitedstates"
	environment_type = "Sandbox"
	dataverse = {
	  language_code     = "1033"
	  currency_code     = "USD"
	  security_group_id = "00000000-0000-0000-0000-000000000000"
	}
}

resource "null_resource" "wait_60_seconds" {
	provisioner "local-exec" {
		command = "sleep 120"
	}
	depends_on = [powerplatform_environment.data_env]
}

resource "powerplatform_data_record" "contact1" {
  environment_id     = powerplatform_environment.data_env.id 
  table_logical_name = "contact"

  columns = {
    contactid = "00000000-0000-0000-0000-000000000001"
    firstname = "contact1"
    lastname  = "contact1"

    contact_customer_contacts = [
      {
        table_logical_name = powerplatform_data_record.contact2.table_logical_name
        data_record_id     = powerplatform_data_record.contact2.columns.contactid
      },
      {
        table_logical_name = powerplatform_data_record.contact3.table_logical_name
        data_record_id     = powerplatform_data_record.contact3.columns.contactid
      }
    ]
  }

  depends_on = [
    null_resource.wait_60_seconds
]
}

resource "powerplatform_data_record" "contact2" {
  environment_id     = powerplatform_environment.data_env.id 
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000002"
    firstname = "contact2"
    lastname  = "contact2"

    contact_customer_contacts = [
      {
        table_logical_name = powerplatform_data_record.contact4.table_logical_name
        data_record_id     = powerplatform_data_record.contact4.columns.contactid
      }
    ]
  }
}

resource "powerplatform_data_record" "contact3" {
  environment_id     = powerplatform_environment.data_env.id 
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000003"
    firstname = "contact3"
    lastname  = "contact3"
  }
  depends_on = [
    null_resource.wait_60_seconds
  ]
}

resource "powerplatform_data_record" "contact4" {
  environment_id     = powerplatform_environment.data_env.id 
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000004"
    firstname = "contact4"
    lastname  = "contact4"
    account_primary_contact = [
      {
        table_logical_name = powerplatform_data_record.account1.table_logical_name
        data_record_id     = powerplatform_data_record.account1.columns.accountid
      }
    ]
  }
  depends_on = [
    null_resource.wait_60_seconds
  ]
}

resource "powerplatform_data_record" "account1" {
  environment_id     = powerplatform_environment.data_env.id 
  table_logical_name = "account"
  columns = {
    accountid = "00000000-0000-0000-0000-000000000010"
    name      = "account1"
    contact_customer_accounts = [
      {
        table_logical_name = powerplatform_data_record.contact5.table_logical_name
        data_record_id     = powerplatform_data_record.contact5.columns.contactid
      }
    ]
  }

  depends_on = [
    null_resource.wait_60_seconds
  ]
}

resource "powerplatform_data_record" "contact5" {
  environment_id     = powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000005"
    firstname = "contact5"
    lastname  = "contact5"
  }
  
  depends_on = [
    null_resource.wait_60_seconds
  ]
}`
}

func TestAccDataRecordDatasource_Validate_Expand_Query(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = powerplatform_environment.data_env.id
						entity_collection = "contacts"
						filter            = "firstname eq 'contact1'"
						select            = ["fullname","firstname","lastname"]
						expand = [
						  {
							navigation_property = "contact_customer_contacts"
							select              = ["fullname"]
							expand = [
							  {
								navigation_property = "contact_customer_contacts"
								select              = ["fullname"]
								expand = [
								  {
									navigation_property = "account_primary_contact"
									select              = ["name"]
									expand = [
									  {
										navigation_property = "contact_customer_accounts"
										select              = ["fullname"]
									  },
									  {
										navigation_property = "primarycontactid"
										select              = ["fullname"]
									  }
									]
								  }
								]
							  }
							]
						  },
						]
					  
						depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
						  powerplatform_data_record.account1,
						]
					  }
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.powerplatform_data_records.data_query", tfjsonpath.New("rows"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapPartial(map[string]knownvalue.Check{
								"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000001"),
								"firstname": knownvalue.StringExact("contact1"),
								"lastname":  knownvalue.StringExact("contact1"),
								"contact_customer_contacts": knownvalue.ListExact([]knownvalue.Check{
									0: knownvalue.MapExact(map[string]knownvalue.Check{
										"fullname":  knownvalue.StringExact("contact2 contact2"),
										"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000002"),
										"contact_customer_contacts": knownvalue.ListExact([]knownvalue.Check{
											0: knownvalue.MapExact(map[string]knownvalue.Check{
												"fullname":  knownvalue.StringExact("contact4 contact4"),
												"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000004"),
												"account_primary_contact": knownvalue.ListExact([]knownvalue.Check{
													0: knownvalue.MapExact(map[string]knownvalue.Check{
														"accountid": knownvalue.StringExact("00000000-0000-0000-0000-000000000010"),
														"name":      knownvalue.StringExact("account1"),
														"contact_customer_accounts": knownvalue.ListExact([]knownvalue.Check{
															0: knownvalue.MapExact(map[string]knownvalue.Check{
																"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000005"),
																"fullname":  knownvalue.StringExact("contact5 contact5")}),
														}),
														"primarycontactid": knownvalue.ObjectExact(map[string]knownvalue.Check{
															"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000004"),
															"fullname":  knownvalue.StringExact("contact4 contact4"),
														}),
													}),
												}),
											}),
										}),
									}),
									1: knownvalue.MapExact(map[string]knownvalue.Check{
										"fullname":                  knownvalue.StringExact("contact3 contact3"),
										"contactid":                 knownvalue.StringExact("00000000-0000-0000-0000-000000000003"),
										"contact_customer_contacts": knownvalue.ListSizeExact(0),
									}),
								}),
							}),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.contactid", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.firstname", "contact1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.lastname", "contact1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.fullname", "contact1 contact1"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_Expand_Query(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$select=fullname%2Cfirstname%2Clastname&$filter=firstname+eq+%27contact1%27&$expand=contact_customer_contacts($select=fullname;$expand=contact_customer_contacts($select=fullname;$expand=account_primary_contact($select=name;$expand=contact_customer_accounts($select=fullname),primarycontactid($select=fullname))))",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						filter            = "firstname eq 'contact1'"
						select            = ["fullname","firstname","lastname"]
						expand = [
							{
							navigation_property = "contact_customer_contacts"
							select              = ["fullname"]
							expand = [
								{
								navigation_property = "contact_customer_contacts"
								select              = ["fullname"]
								expand = [
									{
									navigation_property = "account_primary_contact"
									select              = ["name"]
									expand = [
										{
										navigation_property = "contact_customer_accounts"
										select              = ["fullname"]
										},
										{
										navigation_property = "primarycontactid"
										select              = ["fullname"]
										}
									]
									}
								]
								}
							]
							},
						]
					}`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_Single_Record_Expand_Query(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = powerplatform_environment.data_env.id
						entity_collection = "contacts(00000000-0000-0000-0000-000000000001)"
						select            = ["fullname","firstname","lastname"]
						expand = [
						  {
							navigation_property = "contact_customer_contacts"
							select              = ["fullname"]
							expand = [
							  {
								navigation_property = "contact_customer_contacts"
								select              = ["fullname"]
							  }
							]
						  },
						]
					  
						depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
						  powerplatform_data_record.account1,
						]
					  }
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.powerplatform_data_records.data_query", tfjsonpath.New("rows"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapPartial(map[string]knownvalue.Check{
								"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000001"),
								"firstname": knownvalue.StringExact("contact1"),
								"lastname":  knownvalue.StringExact("contact1"),
								"contact_customer_contacts": knownvalue.ListExact([]knownvalue.Check{
									0: knownvalue.MapPartial(map[string]knownvalue.Check{
										"fullname":                  knownvalue.StringExact("contact2 contact2"),
										"contactid":                 knownvalue.StringExact("00000000-0000-0000-0000-000000000002"),
										"contact_customer_contacts": knownvalue.ListSizeExact(0),
									}),
									1: knownvalue.MapPartial(map[string]knownvalue.Check{
										"fullname":                  knownvalue.StringExact("contact3 contact3"),
										"contactid":                 knownvalue.StringExact("00000000-0000-0000-0000-000000000003"),
										"contact_customer_contacts": knownvalue.ListSizeExact(0),
									}),
								}),
							}),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.contactid", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.firstname", "contact1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.lastname", "contact1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.fullname", "contact1 contact1"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_Single_Record_Expand_Query(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts(00000000-0000-0000-0000-000000000001)?$select=fullname%2Cfirstname%2Clastname&$expand=contact_customer_contacts($select=fullname;$expand=contact_customer_contacts($select=fullname))",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts(00000000-0000-0000-0000-000000000001)"
						select            = ["fullname","firstname","lastname"]
						expand = [
						  {
							navigation_property = "contact_customer_contacts"
							select              = ["fullname"]
							expand = [
							  {
								navigation_property = "contact_customer_contacts"
								select              = ["fullname"]
							  }
							]
						  },
						]
					  }`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_Top(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = powerplatform_environment.data_env.id
						entity_collection = "contacts"
						select            = ["fullname","firstname","lastname"]
						return_total_rows_count = true
					  
						depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
						  powerplatform_data_record.account1,
						]
					  }
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "5"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "total_rows_count", "5"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "total_rows_count_limit_exceeded", "false"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_Top(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$select=fullname%2Cfirstname%2Clastname&$count=true",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						select            = ["fullname","firstname","lastname"]
						return_total_rows_count = true
					  }`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_Apply(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = powerplatform_environment.data_env.id
						entity_collection = "contacts"
						apply             = "groupby((statuscode),aggregate($count as count))"
					  
						depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
						  powerplatform_data_record.account1,
						]
					  }
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.statuscode", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.count", "5"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_Apply(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$apply=groupby%28%28statuscode%29%2Caggregate%28%24count+as+count%29%29",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						apply             = "groupby((statuscode),aggregate($count as count))"
					  
					  }`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_OrderBy(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = powerplatform_environment.data_env.id
						entity_collection = "contacts"
						order_by             = "firstname desc, lastname desc"
					  
						depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
						]
					  }
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "5"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.firstname", "contact5"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.1.firstname", "contact4"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.2.firstname", "contact3"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.3.firstname", "contact2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.4.firstname", "contact1"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_OrderBy(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$orderby=firstname+desc%2C+lastname+desc",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						order_by             = "firstname desc, lastname desc"
					  }`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_SavedQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) + `
				data "powerplatform_data_records" "saved_view" {
					environment_id    = powerplatform_environment.data_env.id
					entity_collection = "savedqueries"
					select            = ["name", "savedqueryid", "returnedtypecode"]
					filter            = "returnedtypecode eq 'contact' and name eq 'All Contacts'"
					top               = 1

					depends_on = [
						  powerplatform_data_record.contact1,
						  powerplatform_data_record.contact2,
						  powerplatform_data_record.contact3,
						  powerplatform_data_record.contact4,
						  powerplatform_data_record.contact5,
					]
				  }
				  
				  data "powerplatform_data_records" "data_query" {
					environment_id    = powerplatform_environment.data_env.id
					entity_collection = "contacts"
					saved_query       = data.powerplatform_data_records.saved_view.rows[0].savedqueryid
					select            = ["fullname"]
				  }
				  
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "5"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_SavedQuery(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?savedQuery=00000000-0000-0000-0000-000000000002&$select=fullname",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						saved_query       = "00000000-0000-0000-0000-000000000002"
						select            = ["fullname"]
					}`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_UserQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) + `
				resource "powerplatform_data_record" "userquery" {
					environment_id     = powerplatform_environment.data_env.id
					table_logical_name = "userquery"
					columns = {
						name             = "user_query_acceptance_test"
						userqueryid      = "00000000-0000-0000-0000-000000000021"
						fetchxml         = "<fetch version=\"1.0\" output-format=\"xml-platform\" mapping=\"logical\" no-lock=\"false\" userqueryid=\"00000000-0000-0000-0000-000000000021\"><entity name=\"contact\"><attribute name=\"statecode\" /><attribute name=\"entityimage_url\" /><attribute name=\"fullname\" /><attribute name=\"contactid\" /></entity></fetch>"
						querytype        = 0
						returnedtypecode = "contact"
						layoutxml        = "<grid name=\"resultset\" object=\"2\" jump=\"fullname\" select=\"1\" icon=\"false\" preview=\"1\"><row name=\"result\" id=\"contactid\"><cell name=\"fullname\" width=\"300\"/><cell name=\"emailaddress1\" width=\"150\"/></row></grid>"
					}
				}
				
				data "powerplatform_data_records" "data_query" {
				environment_id    = powerplatform_environment.data_env.id
				entity_collection = "contacts"
				user_query        = powerplatform_data_record.userquery.columns.userqueryid
				select            = ["fullname", "firstname", "lastname"]

				depends_on = [
						powerplatform_data_record.contact1,
						powerplatform_data_record.contact2,
						powerplatform_data_record.contact3,
						powerplatform_data_record.contact4,
						powerplatform_data_record.contact5,
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "5"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_UserQuery(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?userQuery=00000000-0000-0000-0000-000000000002&$select=fullname%2Cfirstname%2Clastname",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27contacts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_contact.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "contacts"
						user_query        = "00000000-0000-0000-0000-000000000002"
						select            = ["fullname", "firstname", "lastname"]

						
					}`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}

func TestAccDataRecordDatasource_Validate_Expand_Lookup(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + BootstrapDataRecordTest(mocks.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
					environment_id    = powerplatform_environment.data_env.id
					entity_collection = "accounts"
					filter            = "accountid eq 00000000-0000-0000-0000-000000000010"
					select            = ["name", "accountid", "owninguser"]
					expand = [
						{
						navigation_property = "primarycontactid"
						select              = ["contactid", "firstname", "lastname"],
						},
						{
						navigation_property = "owningbusinessunit"
						select              = ["createdon", "name"],
						}
					]
					depends_on = [
						powerplatform_data_record.contact1,
						powerplatform_data_record.contact2,
						powerplatform_data_record.contact3,
						powerplatform_data_record.contact4,
						powerplatform_data_record.contact5,
						powerplatform_data_record.account1,
						]
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.powerplatform_data_records.data_query", tfjsonpath.New("rows"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"@odata.etag": knownvalue.NotNull(),
								"accountid":   knownvalue.StringExact("00000000-0000-0000-0000-000000000010"),
								"name":        knownvalue.StringExact("account1"),
								"owningbusinessunit": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"businessunitid": knownvalue.NotNull(),
									"createdon":      knownvalue.NotNull(),
									"name":           knownvalue.NotNull(),
								}),
								"primarycontactid": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"contactid": knownvalue.StringExact("00000000-0000-0000-0000-000000000004"),
									"firstname": knownvalue.StringExact("contact4"),
									"lastname":  knownvalue.StringExact("contact4"),
								}),
							}),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.0.name", "account1"),
				),
			},
		},
	})
}

func TestUnitDataRecordDatasource_Validate_Expand_Lookup(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	isOdataQueryRun := false

	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts?$select=name%2Caccountid%2Cowninguser&$filter=accountid+eq+00000000-0000-0000-0000-000000000010&$expand=primarycontactid($select=contactid%2Cfirstname%2Clastname),owningbusinessunit($select=createdon%2Cname)",
		func(req *http.Request) (*http.Response, error) {
			isOdataQueryRun = true
			return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/a/v9.2/$metadata#accounts","value":[]}`), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/EntityDefinitions?%24filter=LogicalCollectionName+eq+%27accounts%27&%24select=PrimaryIdAttribute%2CLogicalCollectionName%2CLogicalName`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_entitydefinition_account.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig +
					`data "powerplatform_data_records" "data_query" {
						environment_id    = "00000000-0000-0000-0000-000000000001"
						entity_collection = "accounts"
						filter            = "accountid eq 00000000-0000-0000-0000-000000000010"
						select            = ["name", "accountid", "owninguser"]
						expand = [
							{
								navigation_property = "primarycontactid"
								select              = ["contactid", "firstname", "lastname"],
							},
							{
								navigation_property = "owningbusinessunit"
								select              = ["createdon", "name"],
							}
						]
					}		
					`,
			},
		},
	})

	if !isOdataQueryRun {
		t.Errorf("Odata query should have been run in '%s' unit test", mocks.TestName())
	}
}
