# Title

Missing validation for struct fields

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/dto.go`

## Problem

The struct fields in the DTO (Data Transfer Objects) such as `Id`, `Type`, `DisplayName`, and `Description`, lack any form of validation. For example, the `Id` field in `environmentGroupPrincipalDto` could be empty or malformed, and there is no mechanism in place to ensure that the `DisplayName` and `Description` fields of `environmentGroupDto` are provided or conform to specific rules.

## Impact

This can lead to situations where invalid data is passed to or from the system, potentially causing runtime errors, reduced system reliability, and additional debugging overhead. **Severity: High**

## Location

Struct definitions across the file.

## Code Issue

```go
type environmentGroupPrincipalDto struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type environmentGroupDto struct {
	DisplayName string                       `json:"displayName"`
	Description string                       `json:"description"`
	Id          string                       `json:"id,omitempty"`
	CreatedTime string                       `json:"createdTime,omitempty"`
	CreatedBy   environmentGroupPrincipalDto `json:"createdBy,omitempty"`
}
```

## Fix

Introduce validation mechanisms for required fields. One option is to use libraries or manual validation to ensure data integrity for these DTOs. Example:

```go
import (
	"errors"
	"strings"
)

func (dto *environmentGroupDto) Validate() error {
	if strings.TrimSpace(dto.DisplayName) == "" {
		return errors.New("DisplayName is required")
	}
	if strings.TrimSpace(dto.Description) == "" {
		return errors.New("Description is required")
	}
	if strings.TrimSpace(dto.Id) == "" {
		return errors.New("Id is required")
	}
	return nil
}
```

Implement similar functions for other DTOs, and ensure proper usage wherever an object of these types is created or used.
