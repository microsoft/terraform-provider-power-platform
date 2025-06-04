# Title

Unwrapped Errors on I/O and Certificate Decoding

##

/workspaces/terraform-provider-power-platform/internal/helpers/cert.go

## Problem

Throughout the file, errors are returned directly from lower-level library/system functions (like `os.ReadFile`, `pkcs12.DecodeChain`, base64 decode). This exposes raw error messages, which can be less useful for consumers of this package, making debugging less contextual and error handling less robust. Wrapping these errors with higher-level context provides a clearer indication of where and why the failure occurred.

## Impact

Returning unwrapped errors in exported functions leads to less maintainable and debuggable code, especially when this package is integrated with larger projects. Severity: **medium**.

## Location

Multiple locations, specifically in:

- `GetCertificateRawFromCertOrFilePath`
- `convertBase64ToByte`
- `convertByteToCert`

## Code Issue

```go
pfx, err := os.ReadFile(certificateFilePath)
if err != nil {
    return "", err
}
```
and

```go
return nil, nil, err
```
and

```go
return pfx, fmt.Errorf("could not decode base64 certificate data: %w", err)
```
## Fix

Add context to error messages using `fmt.Errorf("context: %w", err)` especially in exported functions.

```go
pfx, err := os.ReadFile(certificateFilePath)
if err != nil {
    return "", fmt.Errorf("failed to read certificate file '%s': %w", certificateFilePath, err)
}
```
and in `convertByteToCert`:

```go
key, cert, _, err := pkcs12.DecodeChain(certData, password)
if err != nil {
    return nil, nil, fmt.Errorf("failed to decode PKCS12 certificate chain: %w", err)
}
```
and in functions where returning library error:

```go
return pfx, fmt.Errorf("could not decode base64 certificate data: %w", err)
```
