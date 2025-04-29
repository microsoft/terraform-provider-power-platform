# Title
Inefficient URL Construction in Multiple Resource Blocks

## Problem
The URL construction within `powerplatform_rest.query` (create, read, update, destroy) blocks are repetitive and lack modularity. Each operation manually constructs the URLs, leading to unnecessary duplication and increased maintenance effort.

## Impact
This approach affects maintainability and readability. Any change in the URL structure will require modifying it in multiple places, increasing the risk of errors. This is a **medium severity** issue.

## Location
Found in `resource_rest_test.go`, lines within the `powerplatform_rest.query` definition.

## Code Issue
```go
create = {
    scope   = "${powerplatform_environment.env.dataverse.url}/.default"
    url     = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/accounts?$select=name,accountid"
    method  = "POST"
    body    = local.body
    headers = local.headers
}
read = {
    scope  = "${powerplatform_environment.env.dataverse.url}/.default"
    url    = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
    method = "GET"
}
update = {
    scope   = "${powerplatform_environment.env.dataverse.url}/.default"
    url     = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
    method  = "PATCH"
    body    = local.body
    headers = local.headers
}
destroy = {
    scope  = "${powerplatform_environment.env.dataverse.url}/.default"
    url    = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)"
    method = "DELETE"
}
```

## Fix
Introduce helper functions or constants for constructing URLs and scopes.
This eliminates duplication and centralizes logic for managing URLs.

```go
locals {
    base_url = "${powerplatform_environment.env.dataverse.url}"
    base_scope = "${powerplatform_environment.env.dataverse.url}/.default"
    account_path = "/api/data/v9.2/accounts"
}

create = {
    scope   = local.base_scope
    url     = local.base_url + local.account_path + "?$select=name,accountid"
    method  = "POST"
    body    = local.body
    headers = local.headers
}
read = {
    scope  = local.base_scope
    url    = local.base_url + local.account_path + "(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
    method = "GET"
}
update = {
    scope   = local.base_scope
    url     = local.base_url + local.account_path + "(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
    method  = "PATCH"
    body    = local.body
    headers = local.headers
}
destroy = {
    scope  = local.base_scope
    url    = local.base_url + local.account_path + "(00000000-0000-0000-0000-000000000001)"
    method = "DELETE"
}
```