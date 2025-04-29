# Title

Magic string usage in HTTP mock response.

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record_test.go

## Problem

Hardcoded strings are used multiple times to mock HTTP responses in the unit tests (e.g., URLs, `@odata.context`). These strings are likely referenced in multiple places, and if their values need to change, they will require manual updates in all occurrences.

## Impact

Using hardcoded strings or "magic strings" can lead to maintainability issues. Changes to these values might require a significant refactor, increasing the risk of introducing bugs. Severity: **medium**.

## Location

Found in several places within the test functions, such as `TestUnitDataRecordDatasource_Validate_Expand_Query` and `TestUnitDataRecordDatasource_Validate_SavedQuery`.

## Code Issue

```go
httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$select=fullname%2Cfirstname%2Clastname&$filter=firstname+eq+%27contact1%27&$expand=contact_customer_contacts($select=fullname;$expand=contact_customer_contacts($select=fullname;$expand=account_primary_contact($select=name;$expand=contact_customer_accounts($select=fullname),primarycontactid($select=fullname))))",
```

## Fix

Introduce constants or a configuration object to define these values. Replace hardcoded strings with references to these constants or configurations. This approach minimizes refactoring and reduces error prevalence when the values change.

```go
const (
    testBaseURL          = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2"
    testEnvironmentID    = "00000000-0000-0000-0000-000000000001"
    testODataContext     = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#contacts"
)

httpmock.RegisterResponder("GET", testBaseURL+"/contacts?$select=fullname%2Cfirstname%2Clastname&$filter=firstname+eq+%27contact1%27&$expand=contact_customer_contacts($select=fullname;$expand=contact_customer_contacts($select=fullname;$expand=account_primary_contact($select=name;$expand=contact_customer_accounts($select=fullname),primarycontactid($select=fullname))))",
```

This improves maintainability by centralizing the configuration variables.