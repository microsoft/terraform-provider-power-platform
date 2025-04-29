# Title

Missing Error Handling in Validation Logic for Provider Configuration

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

The validation logic within `configureClientCertificate`, `configureClientSecret`, and similar methods do not adequately handle validation errors. For instance, when certificate retrieval fails, it is logged but not appropriately handled.

## Impact

Severity: High

- Leads to configuration issues when invalid values are provided.
- Users may experience degraded functionality without clear feedback.

## Location

Functions dealing with provider configuration.

## Code Issue Examples

```go
cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
if err != nil {
    resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
}
p.Config.ClientCertificateRaw = cert
p.Config.ClientCertificatePassword = clientCertificatePassword
```

## Fix

Introduce explicit error handling for cases like these to avoid delayed failures:

```go
cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
if err != nil {
    resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
    return  // Ensure failure is explicitly handled.
}
p.Config.ClientCertificateRaw = cert
p.Config.ClientCertificatePassword = clientCertificatePassword
```