# Issue 4

## Title

Hardcoded Error Messages Limiting Localization and Adaptability

##

`/workspaces/terraform-provider-power-platform/internal/helpers/cert.go`

## Problem

Error messages are hardcoded directly into multiple functions, such as `GetCertificateRawFromCertOrFilePath` and `ConvertBase64ToCert`. This approach makes it challenging to adapt the codebase for different locales and reduces maintainability in case of error message updates.

## Impact

- Lack of flexibility to support localization for international users.
- Increased difficulty in maintaining consistent message styles across the codebase.
- Severity: **Medium**

## Location

- `GetCertificateRawFromCertOrFilePath`
- `ConvertBase64ToCert`
- Helper functions like `convertByteToCert`

## Code Issue

```go
return "", errors.New("either client_certificate base64 or certificate_file_path must be provided")
```

## Fix

Use a centralized error message constant map or struct to manage error messages. This allows easier updates and helps in localization setups.

Example:

```go
var ErrorMessages = map[string]string{
    "CertInputMissing": "Either client_certificate base64 or certificate_file_path must be provided.",
    // Add other error messages here...
}

// Usage
return "", errors.New(ErrorMessages["CertInputMissing"])
```

Benefits:
1. Error messages are reusable and easier to maintain.
2. Facilitates translation/localization for international users.
3. Promotes better code organization and consistency.