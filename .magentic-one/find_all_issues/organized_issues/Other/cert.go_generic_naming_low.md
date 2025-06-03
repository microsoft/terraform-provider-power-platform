# Title

Generic Naming in `convertBase64ToByte` and `convertByteToCert`

##

/workspaces/terraform-provider-power-platform/internal/helpers/cert.go

## Problem

`convertBase64ToByte` and `convertByteToCert` use generic and slightly misleading names. For example, "Byte" in Go is conventionally pluralized as "Bytes" in method names when referring to a slice. Similarly, file-private helper names should be more descriptive about their domain, e.g., "PKCS12" or "Cert" instead of simply "Byte".

## Impact

Unclear naming can hinder readability and maintainability, particularly for new contributors. Severity: **low**.

## Location

```go
func convertBase64ToByte(b64 string) ([]byte, error)
func convertByteToCert(certData []byte, password string) ([]*x509.Certificate, crypto.PrivateKey, error)
```

## Code Issue

```go
func convertBase64ToByte(b64 string) ([]byte, error)

func convertByteToCert(certData []byte, password string) ([]*x509.Certificate, crypto.PrivateKey, error)
```

## Fix

Rename:
- `convertBase64ToByte` → `decodeBase64ToBytes`
- `convertByteToCert` → `decodePKCS12CertAndKey`

```go
func decodeBase64ToBytes(b64 string) ([]byte, error)

func decodePKCS12CertAndKey(certData []byte, password string) ([]*x509.Certificate, crypto.PrivateKey, error)
```
