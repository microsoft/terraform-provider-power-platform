# Title
Missing HTTP Status Validation in Test Cases

## 
/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem
In the test cases, the HTTP status returned from mock responses is not explicitly checked to match expected values. For example, the validation of HTTP status `200` is expected but not enforced in the test case implementations by default. This can cause scenarios where a failed HTTP response still passes the test inappropriately.

## Impact
Failure to validate the HTTP status code can lead to incorrect assumptions in the test results. This impacts the reliability of the tests and creates potential for false positives, where a test passes despite an incorrect HTTP response. 

Severity: **Medium**

## Location
- Lines 42-52 in the `resource.TestStep` block for the `TestAccDatasourceRestQuery_WhoAmI_Using_Scope` function.

## Code Issue
The test case provided does not include explicit HTTP status validation:
```go
data "powerplatform_rest_query" "webapi_query" {
    scope                = "${powerplatform_environment.env.dataverse.url}/.default"
    url                  = "${powerplatform_environment.env.dataverse.url}api/data/v9.2/WhoAmI"
    method               = "GET"
    expected_http_status = [200]
}
```

Note the lack of explicit checks for the HTTP status returned by the responder in the mock setup.

## Fix
Add explicit validation in the `resource.ComposeAggregateTestCheckFunc` or another appropriate segment of the test to ensure the returned status matches `200`. Here’s an updated version of the test case:

```go
resource.Test(t, resource.TestCase{
    ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
        {
            Config: `
            resource "powerplatform_environment" "env" {
                display_name     = "` + mocks.TestName() + `"
                location         = "unitedstates"
                environment_type = "Sandbox"
                dataverse = {
                    language_code     = "1033"
                    currency_code     = "USD"
                    security_group_id = "00000000-0000-0000-0000-000000000000"
                }
            }

            data "powerplatform_rest_query" "webapi_query" {
                scope                = "${powerplatform_environment.env.dataverse.url}/.default"
                url                  = "${powerplatform_environment.env.dataverse.url}api/data/v9.2/WhoAmI"
                method               = "GET"
            }
            `,
            Check: resource.ComposeAggregateTestCheckFunc(
                resource.TestMatchResourceAttr("data.powerplatform_rest_query.webapi_query", "output.status_code", 200),
                resource.TestMatchResourceAttr("data.powerplatform_rest_query.webapi_query", "output.body", regexp.MustCompile(whoAmIResponseRegex)),
            ),
        },
    },
})
```

This ensures the HTTP status code `200` is confirmed in addition to body validation, improving the reliability and accuracy of test results.
