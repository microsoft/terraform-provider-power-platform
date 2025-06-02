# Title

Unmocked HTTP call dependencies in test configuration

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest_test.go

## Problem

The acceptance test (`TestAccTestRest_Validate_Create`) seems to rely on access to certain external or live resources by using Terraform resource configuration with URLs derived from the dynamic test environment. If any HTTP endpoints or services are unavailable or have state drift, the test may fail unpredictably. Good practice is either to isolate all network calls with extensive mocking or clearly comment that this test is non-deterministic due to environment dependencies.

## Impact

Reliance on live HTTP endpoints for acceptance (and especially for what might be unit) tests can make tests flaky, slow, and difficult to debug, as well as potentially polluting real infrastructure. This is a medium severity concern for maintainability and CI reliability; it can slow development and create confusing failures.

## Location

Within the `Config` block in acceptance test steps:

```go
resource "powerplatform_rest" "query" {
    create = {
        scope   = "${powerplatform_environment.env.dataverse.url}/.default"
        url     = "${powerplatform_environment.env.dataverse.url}/api/data/v9.2/accounts?$select=name,accountid"
        method  = "POST"
        body    = local.body
        headers = local.headers
    }
    ...
}
```

## Fix

For robust and reliable tests:
- Ensure all HTTP(s) network dependencies are either isolated by a test double/mock (as is done in the unit test), or 
- Clearly indicate in code comments that this test is designed for running in a known, isolated test environment (not as an automated CI run for every PR).
- For better determinism in acceptance testing, use local test servers or frameworks to mock external service responses.
  
For example, comment explicitly in the test function:

```go
// NOTE: This acceptance test relies on live Power Platform infrastructure. Failures may occur due to external state.
```

Or, refactor to utilize mocks (as shown in the unit test) where feasible. This approach improves reliability and makes the scope and intended use of the test clearer to maintainers and contributors.
