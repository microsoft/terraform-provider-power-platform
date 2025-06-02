# Title

Composed Terraform config lacks explicit resource cleanup step

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

The acceptance test configures resources and runs assertions, but lacks a dedicated test step ensuring the destruction of created resources (i.e., a `Destroy: true` step or explicit verification of destroy actions). If a failure occurs before resource teardown, or if Terraform lifecycle bugs arise, resources may persist after the test.

## Impact

- Risk of orphaned resources (cloud resources, users, groups) after test runs, leading to wasted quota, additional costs, and pollution of the test environment.
- Makes resource usage unpredictable in CI pipelines.
- May make test suite non-repeatable due to run-once side effects.
- Severity: medium.

## Location

Within `TestAccConnectionsShareResource_Validate_Create`, the test steps:

```go
Steps: []resource.TestStep{
    {
        ResourceName: "powerplatform_connection_share.share_with_user1",
        Config: `
            ... all resource configs ...
        `,
        Check: resource.ComposeAggregateTestCheckFunc(
            ...
        ),
    },
},
```

## Code Issue

```go
Steps: []resource.TestStep{
    {
        ResourceName: "powerplatform_connection_share.share_with_user1",
        Config: `
            ...resources...
        `,
        Check: resource.ComposeAggregateTestCheckFunc(
            ...
        ),
    },
},
```

## Fix

Add a destroy verification step to the test, for example:

```go
Steps: []resource.TestStep{
    {
        ResourceName: "powerplatform_connection_share.share_with_user1",
        Config: `
            ...resources...
        `,
        Check: resource.ComposeAggregateTestCheckFunc(
            ...
        ),
    },
    {
        ResourceName: "powerplatform_connection_share.share_with_user1",
        Destroy:      true,
        Check: func(s *terraform.State) error {
            // Optionally verify nothing remains, or leave empty to verify deletion
            return nil
        },
    },
},
```

This guarantees that Terraform attempts resource destroy and errors are surfaced if anything fails to clean up.

---

This issue will be saved to `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/resource_connection_share_test_missing_destroy_step_medium.md`.
