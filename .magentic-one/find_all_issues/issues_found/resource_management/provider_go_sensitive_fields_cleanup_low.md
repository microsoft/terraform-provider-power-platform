# Title

Resource/Memory Leak Risk: No Cleanup for Sensitive Fields in ProviderConfig

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

Sensitive data such as `ClientSecret`, `ClientCertificateRaw`, `ClientCertificatePassword` are assigned to the global provider configuration and are not explicitly zeroed out or cleaned up after usage. In long-lived processes, this increases the risk that sensitive values could leak or be read from memory if the process is compromised.

## Impact

Potential increased attack surface for secrets in memory, especially when running in a single process for long periods (such as with the Terraform provider plugin server). Severity: **low**.

## Location

Throughout the assignment in `Configure` and configure* functions:

```go
p.Config.ClientSecret = clientSecret
p.Config.ClientCertificateRaw = cert
p.Config.ClientCertificatePassword = clientCertificatePassword
// etc.
```

## Fix

Where possible, clear these fields from memory when they are no longer needed. For example, implement a cleanup function and call it after configuration or after authentication is complete:

```go
func (c *ProviderConfig) CleanupSensitive() {
    if c.ClientSecret != nil {
        for i := range c.ClientSecret {
            c.ClientSecret[i] = 0
        }
    }
    if c.ClientCertificateRaw != nil {
        for i := range c.ClientCertificateRaw {
            c.ClientCertificateRaw[i] = 0
        }
    }
    if c.ClientCertificatePassword != nil {
        for i := range c.ClientCertificatePassword {
            c.ClientCertificatePassword[i] = 0
        }
    }
}
```

Call this method (or similar) when sensitive data is no longer needed.
