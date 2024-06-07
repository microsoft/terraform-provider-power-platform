// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func BootstrapDataRecordTest(name string) string {
	return `
// resource "powerplatform_environment" "data_env" {
// 	display_name     = "` + name + `"
// 	location         = "europe"
// 	environment_type = "Sandbox"
// 	dataverse = {
// 	  language_code     = "1033"
// 	  currency_code     = "USD"
// 	  security_group_id = "00000000-0000-0000-0000-000000000000"
// 	}
//   }

resource "powerplatform_data_record" "contact1" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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
}

resource "powerplatform_data_record" "contact2" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000003"
    firstname = "contact3"
    lastname  = "contact3"
  }
}

resource "powerplatform_data_record" "contact4" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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
}



resource "powerplatform_data_record" "account1" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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
}

resource "powerplatform_data_record" "contact5" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000005"
    firstname = "contact5"
    lastname  = "contact5"
  }
}`
}

func TestAccDataRecordDatasource_Validate_Expand_Query(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + BootstrapDataRecordTest(mock_helpers.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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

func TestAccDataRecordDatasource_Validate_Single_Record_Expand_Query(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + BootstrapDataRecordTest(mock_helpers.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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

func TestAccDataRecordDatasource_Validate_Top(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + BootstrapDataRecordTest(mock_helpers.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
						entity_collection = "contacts"
						select            = ["fullname","firstname","lastname"]
						//Top = 2
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
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "rows.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "total_rows_count", "5"),
					resource.TestCheckResourceAttr("data.powerplatform_data_records.data_query", "total_rows_count_limit_exceeded", "false"),
				),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Apply(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + BootstrapDataRecordTest(mock_helpers.TestName()) +
					`
					data "powerplatform_data_records" "data_query" {
						environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
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

func TestAccDataRecordDatasource_Validate_Create2(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "saved_view" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "savedqueries"
					select            = ["name", "savedqueryid", "returnedtypecode"]
					filter            = "returnedtypecode eq 'systemuser' and name eq 'Enabled Users'"
					top               = 1
				  }
				  
				  data "powerplatform_data_records" "example_data_records" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					saved_query       = data.powerplatform_data_records.saved_view.rows[0].savedqueryid
					select            = ["fullname", "internalemailaddress", "domainname"]
					top               = 3
				  }
				  
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create3(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_rows" "example_data_rows" {
					environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					select            = ["firstname", "lastname", "createdon"]
					//top               = 2
					return_total_rows_count = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create4(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "saved_view" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "savedqueries"
					select            = ["name", "savedqueryid", "returnedtypecode"]
					filter            = "returnedtypecode eq 'systemuser' and name eq 'Enabled Users'"
					top               = 1
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create5(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					apply             = "groupby((isdisabled),aggregate($count as count))"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create6(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					select            = ["internalemailaddress", "systemuserid"]
					filter            = "internalemailaddress ne null"
					order_by          = "internalemailaddress"
					top               = 5
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create7(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					select            = ["fullname", "systemuserid"]
					expand = [
					  {
						navigation_property = "systemuserroles_association"
						select              = ["roleid", "name"]
						// filter, select, orderby, top, expand
					  },
					  {
						navigation_property = "teammembership_association"
						select              = ["teamid", "name"]
					  }
					]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccDataRecordDatasource_Validate_Create8(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_records" "example_data_records" {
					environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers"
					select            = ["fullname", "systemuserid"]
					filter = "main eq 1"
					expand = [
					  {
						navigation_property = "level1"
						select              = ["l1", "l1"]
						filter = "l1 eq 1"
						expand = [
						  {
							navigation_property = "level2"
							select              = ["l2", "l2"]
							filter = "l2 eq 1"
							expand = [
							  {
								navigation_property = "level3"
								select              = ["l3", "l3"]
								filter = "l3 eq 1"
							  },
							  {
								navigation_property = "level3a"
								select              = ["l3a", "l3a"]
								filter = "l3a eq 1"
							  },
							]
						  },
						]
					  },
					  {
						navigation_property = "teammembership_association"
						select              = ["teamid", "name"]
						fitler = "teamid eq 1"
						expand = [
							{
								navigation_property = "teamroles_association"
								select              = ["roleid", "name"]
								filer = "roleid eq 1"
								order_by = "name desc, roleid asc"
								top = 2
							}
						]
					  }
					]
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
