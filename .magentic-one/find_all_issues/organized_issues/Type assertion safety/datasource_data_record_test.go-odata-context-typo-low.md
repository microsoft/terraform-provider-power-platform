# Title

Potential Typo in @odata.context Property in Mock HTTP Response

##

internal/services/data_record/datasource_data_record_test.go

## Problem

In the unit test mock for the `accounts` expand test, the property `@odata.context` in the mock HTTP response has an unusual endpoint path:
```
"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/a/v9.2/$metadata#accounts"
```
Note the segment `/api/a/v9.2/` (extra `a/`), which appears to be a typo and likely should be `/api/data/v9.2/`. If this string is checked or parsed by the client code being tested (e.g., for OData version correctness or entity type assertions), it could cause the test to pass or fail incorrectly.

## Impact

A typo here can make the test less realistic (the actual API would never return this endpoint), potentially hiding real bugs or creating false failures/pass scenarios if code under test tries to validate the OData metadata structure. 

Severity: Low, unless the code under test does strict validation of `@odata.context`, in which case it is Medium/High.

## Location

Function: `TestUnitDataRecordDatasource_Validate_Expand_Lookup`
In the responder for `"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/accounts?...`:

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/a/v9.2/$metadata#accounts","value":[]}`), nil
```

## Fix

Correct this segment to match the true Dynamics endpoint (`api/data/v9.2`):
```go
return httpmock.NewStringResponse(http.StatusOK, `{"@odata.context":"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/$metadata#accounts","value":[]}`), nil
```
This ensures the test mock is as realistic and robust as possible.

Save as a testing/data consistency issue.
