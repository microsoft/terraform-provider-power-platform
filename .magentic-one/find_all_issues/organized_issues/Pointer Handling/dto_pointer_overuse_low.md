# Overuse of Pointers for Scalar Types

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go

## Problem

The struct `environmentSettingsDto` uses pointers for all its fields, even for simple types like `int64` and `bool`. If it’s not required to distinguish between “not set” and “zero value,” this pointer usage adds unnecessary complexity and memory allocation overhead.

## Impact

Unnecessary pointer usage increases complexity and the need for nil checks, which could result in runtime bugs and more verbose code. It also increases heap allocations and can degrade performance, especially when many such structs are used.  
**Severity:** low (unless API nullability exactly requires it).

## Location

```go
MaxUploadFileSize *int64 `json:"maxuploadfilesize,omitempty"`
...
AuditRetentionPeriodV2 *int32 `json:"auditretentionperiodv2,omitempty"`
...
```

## Code Issue

```go
MaxUploadFileSize *int64  `json:"maxuploadfilesize,omitempty"`
AuditRetentionPeriodV2 *int32  `json:"auditretentionperiodv2,omitempty"`
...
```

## Fix

Use value types for scalars unless the API requires distinguishing between “unset” (null) and “zero.” Example fix:

```go
MaxUploadFileSize int64  `json:"maxuploadfilesize,omitempty"`
AuditRetentionPeriodV2 int32  `json:"auditretentionperiodv2,omitempty"`
...
```
