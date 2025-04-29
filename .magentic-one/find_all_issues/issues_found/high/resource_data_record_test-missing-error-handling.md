# Title

Missing Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

In certain places, such as mock responders, errors are not properly handled which might lead to unexpected behaviors. One specific instance is reading the body of incoming requests in `readBodyBuff` without error checks at each point.

## Impact

A lack of proper error handling can lead to runtime panics or untraceable failures. This diminishes the reliability of the tests and could lead to time-consuming debugging sessions. Severity: High.

## Location

Specific instance in the function `readBodyBuff`.

## Code Issue

```go
func readBodyBuff(req *http.Request) (string, error) {
    r, err := req.GetBody()
    if err != nil {
        return "", err
    }
    defer r.Close()

    buf := new(bytes.Buffer)
    if _, err := buf.ReadFrom(r); err != nil {
        return "", err
    }
    return buf.String(), nil
}
```

## Fix

Improve error handling by wrapping result of `req.GetBody()` and `buf.ReadFrom(r)` into meaningful error messages or logs along with continuing checks.

```go
func readBodyBuff(req *http.Request) (string, error) {
    r, err := req.GetBody()
    if err != nil {
        return "", fmt.Errorf("error getting body from request: %w", err)
    }
    defer func() {
        if rErr := r.Close(); rErr != nil {
            log.Printf("error closing request body: %v", rErr)
        }
    }()

    buf := new(bytes.Buffer)
    if _, err := buf.ReadFrom(r); err != nil {
        return "", fmt.Errorf("error reading from body buffer: %w", err)
    }

    return buf.String(), nil
}
```