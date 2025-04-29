# Title

Hardcoded Values

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The code contains numerous hardcoded values, such as GUIDs and other strings, within test configurations for Terraform resources.

## Impact

Hardcoded values make the code less maintainable and harder to adapt, as changes to constants require manual updates in multiple locations. Using hardcoded values also diminishes flexibility during testing. Severity: Medium.

## Location

Throughout the file, in test configurations for resources such as `powerplatform_data_record_sample_contact1` and `powerplatform_data_record_account`.
Example:

```go
    environment_id     = "00000000-0000-0000-0000-000000000001"
    table_logical_name = "contact"
    columns = {
      firstname = "John"
      lastname = "Doe"
      telephone1 = "555-555-5555"
```

## Code Issue

```go
    resource "powerplatform_data_record" "data_record_sample_contact1" {
        environment_id     = "00000000-0000-0000-0000-000000000001"
        table_logical_name = "contact"
        columns = {
          firstname = "John"
          lastname = "Doe"
          telephone1 = "555-555-5555"
        }
    }
```

## Fix

Refactor to use constants or variables for repeated values. This will make code testing and modification easier and more consistent. Here is an example fix:

```go
const (
     environmentID = "00000000-0000-0000-0000-000000000001"
     tableLogicalNameContact = "contact"
     firstNameJohn = "John"
)

resource "powerplatform_data_record" "data_record_sample_contact1" {
    environment_id     = environmentID
    table_logical_name = tableLogicalNameContact
    columns = {
      firstname = firstNameJohn
      lastname = "Doe"
      telephone1 = "555-555-5555"
    }
}
```