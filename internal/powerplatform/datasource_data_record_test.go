// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataRecordDatasource_Validate_Create(t *testing.T) {

	t.Setenv("TF_ACC", "1")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_data_rows" "example_data_rows" {
					environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
					entity_collection = "systemusers(1f70a364-5019-ef11-840b-002248ca35c3)"
					return_total_rows_count = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(),
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
