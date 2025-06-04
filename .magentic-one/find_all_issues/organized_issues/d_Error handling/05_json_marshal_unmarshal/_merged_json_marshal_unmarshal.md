# JSON Marshal/Unmarshal Issues

This document consolidates all issues related to JSON marshaling/unmarshaling and type assertion errors in the Terraform Provider for Power Platform.

## ISSUE 1

### Type Safety: Potentially Unhandled Error for JSON Unmarshal

**File:** `/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go`

**Problem:** The code unmarshals `response.BodyAsBytes` to `languages` but does not verify if the payload is indeed valid JSON or is non-empty before trying to unmarshal. While an error is returned if unmarshalling fails, a more defensive check would help with debugging.

**Impact:** Poor error diagnosis in cases of empty or malformed responses. Severity: **medium**.

**Location:**

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)

if err != nil {
 return languages, err
}
```

**Code Issue:**

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)

if err != nil {
 return languages, err
}
```

**Fix:** Optionally, check that `response.BodyAsBytes` is not empty before unmarshalling:

```go
if len(response.BodyAsBytes) == 0 {
    return languages, fmt.Errorf("empty response body")
}
err = json.Unmarshal(response.BodyAsBytes, &languages)
if err != nil {
    return languages, err
}
```

## ISSUE 2

### Unhandled Errors When Marshaling JSON in ConvertFromConnectionDto

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go`

**Problem:** In `ConvertFromConnectionDto`, calls to `json.Marshal` ignore returned errors, using the blank identifier `_`. If marshaling fails, this will silently set a possibly empty or wrong string value, introducing silent data loss or incorrect data.

**Impact:** Severity: **Medium** - Silently ignoring marshaling errors can result in invalid or empty JSON strings in the state model, making debugging difficult and potentially leading to inaccurate resource states in Terraform.

**Location:**

```go
if connection.Properties.ConnectionParametersSet != nil {
 p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
 conn.ConnectionParametersSet = types.StringValue(string(p))
}

if connection.Properties.ConnectionParameters != nil {
 p, _ := json.Marshal(connection.Properties.ConnectionParameters)
 conn.ConnectionParameters = types.StringValue(string(p))
}
```

**Code Issue:**

```go
p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
```

**Fix:** Check the error returned by `json.Marshal` and handle it accordingly:

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParametersSet)
    if err == nil {
        conn.ConnectionParametersSet = types.StringValue(string(p))
    } else {
        // Optionally log or handle the error here
    }
}

if connection.Properties.ConnectionParameters != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParameters)
    if err == nil {
        conn.ConnectionParameters = types.StringValue(string(p))
    } else {
        // Optionally log or handle the error here
    }
}
```

## ISSUE 3

### Error-prone Type Assertion for ProviderData Without Validation

**File:** `/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go`

**Problem:** In the `Configure` method, the code directly type asserts `req.ProviderData.(*api.ProviderClient)` without validating that `ProviderData` is in fact of this type. If the assertion fails, this will cause a panic and interrupt control flow, instead of gracefully handling configuration errors.

**Impact:** A failed type assertion leads to a provider panic instead of a diagnostic error. This impacts user experience and debugging and is considered a **high severity** issue for Terraform providers.

**Location:**

```go
clientApi := req.ProviderData.(*api.ProviderClient).Api
if clientApi == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )

    return
}
d.LanguagesClient = newLanguagesClient(clientApi)
```

**Code Issue:**

```go
clientApi := req.ProviderData.(*api.ProviderClient).Api
```

**Fix:** First, check that `ProviderData` is of the correct type via an assertion with the `ok` idiom:

```go
providerClient, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
clientApi := providerClient.Api
if clientApi == nil {
    resp.Diagnostics.AddError(
        "Unexpected nil Api in ProviderClient",
        "The 'Api' field on ProviderClient was nil.",
    )
    return
}
d.LanguagesClient = newLanguagesClient(clientApi)
```

## ISSUE 4

### Incorrect `MarshallTo` Usage of `json.NewDecoder().Decode()`

**File:** `/workspaces/terraform-provider-power-platform/internal/api/request.go`

**Problem:** The code passes `&obj` to `json.NewDecoder().Decode()`, but `obj` is already of type `any` (interface{}), typically a pointer to struct should be passed, not its address. Taking `&obj` gives a pointer to interface, which is almost always incorrect and leads to runtime errors.

**Impact:** Severity: High - This can cause decoding to fail, with JSON not being correctly deserialized into the provided struct.

**Location:**

```go
func (apiResponse *Response) MarshallTo(obj any) error {
 err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
 if err != nil {
  return err
 }
 return nil
}
```

**Code Issue:**

```go
 err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
```

**Fix:** Pass `obj` directly to Decode (it should already be a pointer):

```go
func (apiResponse *Response) MarshallTo(obj any) error {
 return json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(obj)
}
```

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
