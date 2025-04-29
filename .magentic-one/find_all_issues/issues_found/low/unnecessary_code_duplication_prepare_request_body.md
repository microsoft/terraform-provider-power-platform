# Title

Unnecessary Code Duplication in PrepareRequestBody Method

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The `PrepareRequestBody` method contains redundant checks and logic for handling `body` as a string pointer or as a generic object that needs marshalling. This duplication increases code verbosity and maintenance costs.

## Impact

- Impacts maintainability by making it harder to understand or modify the logic.
- Keeps the function longer than necessary, potentially increasing cognitive load for new developers.

Severity: **Low**

## Location

Within the `PrepareRequestBody` function:

```go
if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
    if strp, ok := body.(*string); ok {
        bodyBuffer = strings.NewReader(*strp)
    } else {
        bodyBytes, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        bodyBuffer = bytes.NewBuffer(bodyBytes)
    }
}
```

## Code Issue

```go
if body != nil && (reflect.ValueOf(body).Kind() != reflect.Ptr || !reflect.ValueOf(body).IsNil()) {
    if strp, ok := body.(*string); ok {
        bodyBuffer = strings.NewReader(*strp)
    } else {
        bodyBytes, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        bodyBuffer = bytes.NewBuffer(bodyBytes)
    }
}
```

## Fix

Refactor the code to reduce duplication by consolidating the conditions and logic.

```go
if body == nil || (reflect.ValueOf(body).Kind() == reflect.Ptr && reflect.ValueOf(body).IsNil()) {
    return nil, nil
}

switch v := body.(type) {
case *string:
    bodyBuffer = strings.NewReader(*v)
default:
    bodyBytes, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }
    bodyBuffer = bytes.NewBuffer(bodyBytes)
}
```

This refactored version uses a `switch` statement to handle type-based cases, reducing the code indentation and repetition.