# Title

Misuse of `fmt.Sprintf` with unescaped characters in URL construction

##

Path to the file `/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go`

## Problem

The code uses `fmt.Sprintf` for constructing URLs without escaping potentially unsafe characters. For example, when creating the `entityDefinitionApiUrl` or similar API URLs in other places, user inputs directly go into the URL paths, which could lead to vulnerabilities if special characters are included.

## Impact

This issue can lead to security vulnerabilities such as open redirect flaws, server-side request forgery (SSRF), or encoding issues while handling user or external input. Severity: **High**

## Location

Refer to the following block in the function `getEntityDefinition`.

## Code Issue

```go
entityDefinitionApiUrl := &url.URL{
	Scheme:   constants.HTTPS,
	Host:     environmentHost,
	Path:     fmt.Sprintf("/api/data/%s/EntityDefinitions(LogicalName='%s')", constants.DATAVERSE_API_VERSION, entityLogicalName),
	Fragment: "$select=PrimaryIdAttribute,LogicalCollectionName",
}
```

## Fix

Use `url.PathEscape` to safely encode user-controlled input before constructing the URL path.

```go
entityDefinitionApiUrl := &url.URL{
	Scheme: constants.HTTPS,
	Host:   environmentHost,
	Path:   fmt.Sprintf("/api/data/%s/EntityDefinitions(LogicalName='%s')", constants.DATAVERSE_API_VERSION, url.PathEscape(entityLogicalName)),
	Fragment: "$select=PrimaryIdAttribute,LogicalCollectionName",
}
```

This ensures that special characters in `entityLogicalName` are correctly escaped, mitigating security issues.