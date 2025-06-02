# Title

Variable Naming Inconsistent

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The variable naming conventions in the file are inconsistent, with some variables using camelCase (e.g., `contactAtIndex1Step1`) and others using snake_case (`disable_on_destroy`).

## Impact

Inconsistent naming conventions make the code harder to read and maintain, especially in larger codebases. This issue is of low severity but contributes to reduced code clarity.

## Location

Throughout the file, such as variables `contactAtIndex1Step1`, `disable_on_destroy`. Example:

```go
contactAtIndex1Step1
contactAtIndex1Step2
primarycontactidStep1
```  

In configuration code:
```go
resource "powerplatform_data_record" "mailbox" {
    disable_on_destroy = true
    environment_id     = "00000000-0000-0000-0000-000000000001"
```

## Code Issue

```go
contactAtIndex1Step1
contactAtIndex1Step2
primarycontactidStep1
```

## Fix

Adopt a single consistent naming convention across the codebase. Example:
- Use camelCase for variable names consistently:

```go
contactAtIndex1Step1
contactAtIndex1Step2
primaryContactIdStep1
```

OR

- Use snake_case consistently:

```go
contact_at_index1_step1
contact_at_index1_step2
primary_contact_id_step1
```