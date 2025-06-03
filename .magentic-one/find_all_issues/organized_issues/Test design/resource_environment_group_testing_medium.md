# Title

Missing Unit Tests for Custom CRUD Logic

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

There is no indication within the file or adjacent code (such as build tags or test imports) that individual unit tests exist for the custom logic in the CRUD (`Create`, `Read`, `Update`, `Delete`) handlers in this file.

While the plugin framework and Terraform acceptance tests provide some coverage, direct unit tests for custom error handling, API response handling, and field population logic can prevent regressions and ensure logic correctness.

## Impact

- Increases risk of regression when logic or dependency contracts change.
- Makes it harder to refactor or optimize logic safely.
- Medium severity as acceptance tests likely cover the functional flows, but nuanced edge cases and internal logic might be missed.

**Severity:** medium

## Location

Whole file, especially CRUD function implementations.

## Code Issue

_No code block: This is an absence of direct testing._

## Fix

Add unit tests that directly test the resource's `Create`, `Read`, `Update`, and `Delete` logic using mocks/stubs for the client. Ensure cases like:
- Partial failures in `Delete`
- State/model conversion from API to state and vice versa
- Error propagation from API to diagnostics

can be tested at a function/unit level.
