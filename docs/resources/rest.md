---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_rest Resource - powerplatform"
subcategory: ""
description: |-
  Resource to execute web api requests. There are four distinct operations, that you can define idepenetly. The HTTP response' body of the operation, that was called as last, will be returned in 'output.body' \n\n:
  * Create: will be called once during the lifecycle of the resource (first 'terraform apply')
  * Read: terraform will call this operation every time during 'plan' and 'apply' to get the current state of the resource
  * Update: will be called every time during 'terraform apply' if the resource has changed (change done by the user or different values returned by the 'read' operation than those in the current state)
  * Destroy: will be called once during the lifecycle of the resource (last 'terraform destroy')
  \n\nYOu don't have to define all the operations but there are some things to consider:
  * lack of 'create' operation will result in no reasource being created. If you only need to read values consider using datasource 'powerplatform_rest_query' instead
  * lack of 'read' operation will result in no resource changes being tracked. That means that the 'update' operation will never be called
  * lack of destroy will couse, that the resource will not be deleted during 'terraform destroy'
---

# powerplatform_rest (Resource)

Resource to execute web api requests. There are four distinct operations, that you can define idepenetly. The HTTP response' body of the operation, that was called as last, will be returned in 'output.body' \n\n:
		* Create: will be called once during the lifecycle of the resource (first 'terraform apply')
		* Read: terraform will call this operation every time during 'plan' and 'apply' to get the current state of the resource
		* Update: will be called every time during 'terraform apply' if the resource has changed (change done by the user or different values returned by the 'read' operation than those in the current state)
		* Destroy: will be called once during the lifecycle of the resource (last 'terraform destroy')
		\n\nYOu don't have to define all the operations but there are some things to consider:
		* lack of 'create' operation will result in no reasource being created. If you only need to read values consider using datasource 'powerplatform_rest_query' instead
		* lack of 'read' operation will result in no resource changes being tracked. That means that the 'update' operation will never be called
		* lack of destroy will couse, that the resource will not be deleted during 'terraform destroy'

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
  display_name     = "sample_data_environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_rest" "install_sample_data" {
  create = {
    scope                = "${powerplatform_environment.env.dataverse.url}/.default"
    url                  = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/InstallSampleData"
    method               = "POST"
    expected_http_status = [204]
  }
  destroy = {
    scope                = "${powerplatform_environment.env.dataverse.url}/.default"
    url                  = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/UninstallSampleData"
    method               = "POST"
    expected_http_status = [204]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `create` (Attributes) Create operation (see [below for nested schema](#nestedatt--create))
- `destroy` (Attributes) Destroy operation (see [below for nested schema](#nestedatt--destroy))
- `read` (Attributes) Read operation (see [below for nested schema](#nestedatt--read))
- `update` (Attributes) Update operation (see [below for nested schema](#nestedatt--update))

### Read-Only

- `id` (String) Unique id (guid)
- `output` (Attributes) Response body after executing the web api request (see [below for nested schema](#nestedatt--output))

<a id="nestedatt--create"></a>
### Nested Schema for `create`

Required:

- `method` (String) HTTP method
- `scope` (String) Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)
- `url` (String) Absolute url of the api call

Optional:

- `body` (String) Body of the request
- `expected_http_status` (List of Number) Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.
- `headers` (Attributes List) Headers of the request (see [below for nested schema](#nestedatt--create--headers))

<a id="nestedatt--create--headers"></a>
### Nested Schema for `create.headers`

Required:

- `name` (String) Header name
- `value` (String) Header value



<a id="nestedatt--destroy"></a>
### Nested Schema for `destroy`

Required:

- `method` (String) HTTP method
- `scope` (String) Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)
- `url` (String) Absolute url of the api call

Optional:

- `body` (String) Body of the request
- `expected_http_status` (List of Number) Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.
- `headers` (Attributes List) Headers of the request (see [below for nested schema](#nestedatt--destroy--headers))

<a id="nestedatt--destroy--headers"></a>
### Nested Schema for `destroy.headers`

Required:

- `name` (String) Header name
- `value` (String) Header value



<a id="nestedatt--read"></a>
### Nested Schema for `read`

Required:

- `method` (String) HTTP method
- `scope` (String) Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)
- `url` (String) Absolute url of the api call

Optional:

- `body` (String) Body of the request
- `expected_http_status` (List of Number) Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.
- `headers` (Attributes List) Headers of the request (see [below for nested schema](#nestedatt--read--headers))

<a id="nestedatt--read--headers"></a>
### Nested Schema for `read.headers`

Required:

- `name` (String) Header name
- `value` (String) Header value



<a id="nestedatt--update"></a>
### Nested Schema for `update`

Required:

- `method` (String) HTTP method
- `scope` (String) Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)
- `url` (String) Absolute url of the api call

Optional:

- `body` (String) Body of the request
- `expected_http_status` (List of Number) Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.
- `headers` (Attributes List) Headers of the request (see [below for nested schema](#nestedatt--update--headers))

<a id="nestedatt--update--headers"></a>
### Nested Schema for `update.headers`

Required:

- `name` (String) Header name
- `value` (String) Header value



<a id="nestedatt--output"></a>
### Nested Schema for `output`

Read-Only:

- `body` (String) Response body after executing the web api request