# Title

Duplicated Logic for Building API URLs With Query Strings

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

The process of building API URLs with query strings (scheme, host, path, setting RawQuery) is repeated throughout the file. This code duplication can lead to subtle errors or inconsistencies if logic changes in only some places. Duplicated snippets add cognitive load and violate DRY (Don't Repeat Yourself) principles.

## Impact

Severity: Medium

Duplication makes refactoring harder, inflates the codebase, and increases testing surface area. Changes in query encoding, endpoint base paths, or global API versioning may become inconsistent across the implementation.

## Location

See similar code in almost every public method, such as `GetDataverseUserBySystemUserId`, `GetDataverseUsers`, `GetDataverseUserByAadObjectId`, `GetEnvironmentUserByAadObjectId`, and others. Example:

## Code Issue

```go
apiUrl := &url.URL{
	Scheme: constants.HTTPS,
	Host:   environmentHost,
	Path:   "/api/data/v9.2/systemusers",
}
values := url.Values{}
values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
apiUrl.RawQuery = values.Encode()
```

## Fix

Extract to helper(s):

```go
func buildApiUrl(scheme, host, path string, query url.Values) string {
	apiUrl := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	apiUrl.RawQuery = query.Encode()
	return apiUrl.String()
}
```

This allows future global changes and makes tests/maintenance much easier.
