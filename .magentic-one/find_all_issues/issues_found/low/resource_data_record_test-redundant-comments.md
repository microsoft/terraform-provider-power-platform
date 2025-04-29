# Title

Redundant Comments

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The file contains redundant or unnecessary comments, such as large commented-out configurations that add clutter without any explanatory benefit.

## Impact

Excessive comments hinder code readability and can confuse developers about what parts of the code are significant. Comments need to provide value or context. Severity: Low.

## Location

Example: Step 2 in `TestUnitDataRecordResource_Validate_Disable_On_Delete`:
```go
// resource "powerplatform_data_record" "mailbox" {
//     disable_on_destroy = true
//     environment_id     = "00000000-0000-0000-0000-000000000001"
//     table_logical_name = "mailbox"
//     columns = {
//         name         = "my mailbox"
//         emailaddress = "contoso@contoso.com"
//     }
// }
```

## Code Issue

```go
// resource "powerplatform_data_record" "mailbox" {
//     disable_on_destroy = true
// }
```

## Fix

Remove commented-out sections if they are not providing meaningful information or context related to the code's logic. Example:

```go
// Removed unnecessary comments or sections to improve code readability.
```