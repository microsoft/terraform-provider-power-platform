# Title

Improper Error Handling in configureClientCertificate Function

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

In the `configureClientCertificate` function, if an error occurs when calling `helpers.GetCertificateRawFromCertOrFilePath`, an error is added to the response diagnostics, but the function continues to execute and modifies the configuration state by setting `ClientCertificateRaw`, `ClientCertificatePassword`, `TenantId`, and `ClientId`.

## Impact

This can result in the provider being in a partially configured, invalid state. The code should return immediately after encountering a critical error to prevent inconsistent or faulty configuration. Severity: **high**.

## Location

Lines around the following fragment in `configureClientCertificate`:

```go
cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
if err != nil {
    resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
}
p.Config.ClientCertificateRaw = cert
p.Config.ClientCertificatePassword = clientCertificatePassword
p.Config.TenantId = tenantId
p.Config.ClientId = clientId
```

## Fix

Return early if an error is detected, preventing further assignment to the configuration:

```go
cert, err := helpers.GetCertificateRawFromCertOrFilePath(clientCertificate, clientCertificateFilePath)
if err != nil {
    resp.Diagnostics.AddAttributeError(path.Root("client_certificate"), "Error getting certificate", err.Error())
    return  // Prevent further action if error occurs
}
p.Config.ClientCertificateRaw = cert
p.Config.ClientCertificatePassword = clientCertificatePassword
p.Config.TenantId = tenantId
p.Config.ClientId = clientId
```
