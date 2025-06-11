# Api General Issues - Merged Issues

## ISSUE 1

# Title

Misspelling in Variable Name 'connetionsArray'

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

The variable `connetionsArray` in the `GetConnections` method is misspelled. It should be `connectionsArray` to accurately convey its purpose and maintain naming consistency.

## Impact

Misspelled variable names reduce code readability, can cause confusion among maintainers, and undermine code quality. Severity: Low.

## Location

Line containing:

```go
connetionsArray := connectionArrayDto{}
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connetionsArray)
if err != nil {
	return nil, err
}

return connetionsArray.Value, nil
```

## Code Issue

```go
connetionsArray := connectionArrayDto{}
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connetionsArray)
if err != nil {
	return nil, err
}

return connetionsArray.Value, nil
```

## Fix

Update the variable name to use the correct spelling ("connectionsArray") everywhere in the function.

```go
connectionsArray := connectionArrayDto{}
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectionsArray)
if err != nil {
	return nil, err
}

return connectionsArray.Value, nil
```


---

## ISSUE 2

# Function naming: `covertDlpPolicyToPolicyModel` and `covertDlpPolicyToPolicyModelDto`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The function `covertDlpPolicyToPolicyModel` is likely a typo of `convertDlpPolicyToPolicyModel`. Consistent and correct naming improves readability and maintainability.

## Impact

Lowers developer experience and codebase quality. (Severity: Low)

## Location

Function and usages throughout file

## Code Issue

```go
func covertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	...
}
```

## Fix

Rename function definition and all usages to `convertDlpPolicyToPolicyModel`.

```go
func convertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	...
}
```



---

## ISSUE 3

# Spelling mistake: `covertDlpPolicyToPolicyModelDto` should be `convertDlpPolicyToPolicyModelDto`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The function name `covertDlpPolicyToPolicyModelDto` is likely a typo and should be `convertDlpPolicyToPolicyModelDto` to match naming consistency (other conversions are `convert...` as well).

## Impact

Lower readability, confusion for maintainers, decreased discoverability for function. (Severity: Low)

## Location

Used on lines 40, 78, 89, 142

## Code Issue

```go
v, err := covertDlpPolicyToPolicyModelDto(policy)
```

Also other locations:
```go
return covertDlpPolicyToPolicyModel(policy)
```

## Fix

Rename all instances and the function definition itself to `convertDlpPolicyToPolicyModelDto` and `convertDlpPolicyToPolicyModel`.

```go
v, err := convertDlpPolicyToPolicyModelDto(policyDto)
// ... update the actual function names in source and usage accordingly.
```



---

## ISSUE 4

# Unexported Type Name Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The `client` type is defined as a struct with a lowercase name, making it unexported. In Go, if the intention is to use this type outside the `languages` package, it should be exported (i.e., named `Client`). If it is deliberately unexported, this is not an issue, but the naming should be reviewed for intent and clarity.

## Impact

If the `client` type is supposed to be used by other packages, keeping it unexported prevents access from outside the package. Severity: **low** (unless package boundaries require it to be exported).

## Location

```go
type client struct {
	Api *api.Client
}
```

## Code Issue

```go
type client struct {
	Api *api.Client
}
```

## Fix

If the type should be exported for reuse, capitalize its name:

```go
type Client struct {
	Api *api.Client
}
```


---

## ISSUE 5

# Incorrect Function Naming Convention

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go

## Problem

The function is named `newPowerAppssClient`, which is inconsistent with Go naming conventions and likely a typo (double 's' in "Appss").

## Impact

This can cause confusion for maintainers and may introduce subtle bugs when this function is called elsewhere. Severity: Low.

## Location

Line 14 in the file.

## Code Issue

```go
func newPowerAppssClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}
```

## Fix

Correct the function name to use the singular "PowerApps":

```go
func newPowerAppsClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}
```


---

## ISSUE 6

# Struct and receiver naming does not follow Go conventions

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The `client` struct and its receiver variable are lowercase and generic, conflicting with Go naming conventions. According to Go best practices:
- Type names should be upper camel case (exported/public) or at least capitalized internally to distinguish types from variables.
- Receiver names should be short, typically the first letter or abbreviation of the type (`c` for `Client`).

The current conventions may create readability and maintainability issues, especially as the codebase grows and `client` could clash with other variables or imports.

## Impact

Severity: Medium. Hinders code readability/maintainability and risks future naming conflicts.

## Location

Type and method declarations:

## Code Issue

```go
type client struct {
	Api *api.Client
}

func (client *client) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
   // ...
}

func (client *client) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
   // ...
}
```

## Fix

Rename the struct to `Client` and update receivers from `client` to `c`:

```go
type Client struct {
	Api *api.Client
}

func (c *Client) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
   // ...
}

func (c *Client) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
   // ...
}
```


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
