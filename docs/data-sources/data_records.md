---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_data_records Data Source - Power Platform"
subcategory: ""
description: |-
  Resource for retrieving data records from Dataverse using (OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results].
---

# powerplatform_data_records (Data Source)

Resource for retrieving data records from Dataverse using (OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results].

## Example Usage

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "env" {
  display_name     = "powerplatform_data_record_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_data_records" "data_query" {
  environment_id    = powerplatform_environment.env.id
  entity_collection = "systemusers"
  filter            = "isdisabled eq false"
  select            = ["firstname", "lastname", "domainname"]
  top               = 2
  order_by          = "lastname asc"

  expand = [
    {
      navigation_property = "systemuserroles_association",
      select              = ["name"],
    },
    {
      navigation_property = "teammembership_association",
      select              = ["name"],
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entity_collection` (String) Value of the enitiy (collection of the query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#entity-collections]. Example:

 * $metadata#systemusers 

*systemusers 

*systemusers(<GUID>) 

*systemusers(<GUID>)/systemuserroles_association 

*contacts(firstname='Joe',emailaddress1='joe@contoso.com') when using (alternate key(s))[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/use-alternate-key-reference-record?tabs=webapi] for single record retrieval
- `environment_id` (String) Id of the Power Platform environment

### Optional

- `apply` (String) Apply the aggregation function to the data records. 

More information on (OData Apply)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#aggregate-data]
- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `return_total_rows_count` (Boolean) Should total records count be also retrived. 

More information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]
- `saved_query` (String) predefined saved query to be used for filtering the data records. 

More information on (Saved Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/retrieve-and-execute-predefined-queries]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]
- `user_query` (String) Predefined user query to be used for filtering the data records. 

More information on (Saved Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/retrieve-and-execute-predefined-queries]

### Read-Only

- `rows` (Dynamic) Columns of the data record table
- `total_rows_count` (Number) Total number of records if attribute `return_total_rows_count` is set to `true`
- `total_rows_count_limit_exceeded` (Boolean) Is total records count limit exceeded. 

More information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]

<a id="nestedatt--expand"></a>
### Nested Schema for `expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand"></a>
### Nested Schema for `expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `expand` (Attributes List) Expand the navigation property of the entity collection. 

More information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables] (see [below for nested schema](#nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand))
- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]

<a id="nestedatt--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand--expand"></a>
### Nested Schema for `expand.expand.expand.expand.expand.expand.expand.expand.expand.expand.expand`

Required:

- `navigation_property` (String) Navigation property of the entity collection. 

More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]

Optional:

- `filter` (String) Filter the data records. 

More information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]
- `order_by` (String) Order the data records. 

More information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]
- `select` (List of String) List of columns to be selected from record(s) defined in entity collection. 

More information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]
- `top` (Number) Number of records to be retrieved. 

More information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]












<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
