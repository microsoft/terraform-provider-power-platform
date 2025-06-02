# Title

Lack of validation for struct fields, leading to possible invalid data.

##

`/workspaces/terraform-provider-power-platform/internal/services/languages/dto.go`

## Problem

Structs defined in the file, such as `languageDto` and `languagePropertiesDto`, do not contain any validation or annotation rules for required fields nor data formats. This allows potentially invalid data to enter the system without constraints.

## Impact

The lack of validation creates a risk for the following issues:

- Data corruption due to unvalidated inputs.
- Inability to enforce key business rules at the data layer.
- Increased likelihood of runtime errors caused by invalid data.

Severity: **Critical**

## Location

Every struct field declared in the DTOs:
1. `languagesArrayDto.Value`
2. `languageDto.Name`
3. `languageDto.ID`
4. `languageDto.Type`
5. `languagePropertiesDto.LocaleID`
6. `languagePropertiesDto.LocalizedName`
7. `languagePropertiesDto.DisplayName`
8. `languagePropertiesDto.IsTenantDefault`

## Code Issue

```go
type languageDto struct {
	Name       string                `json:"name"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Properties languagePropertiesDto `json:"properties"`
}

type languagePropertiesDto struct {
	LocaleID        int64  `json:"localeId"`
	LocalizedName   string `json:"localizedName"`
	DisplayName     string `json:"displayName"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
```

## Fix

Add validation logic to the struct definition using suitable libraries like `go-playground/validator` or similar approaches that enforce required constraints and validate data formats.

```go
import "github.com/go-playground/validator/v10"

type languageDto struct {
	Name       string                `json:"name" validate:"required"`
	ID         string                `json:"id" validate:"required,uuid"`
	Type       string                `json:"type" validate:"required,oneof=text audio video"`
	Properties languagePropertiesDto `json:"properties" validate:"required"`
}

type languagePropertiesDto struct {
	LocaleID        int64  `json:"localeId" validate:"required,min=1"`
	LocalizedName   string `json:"localizedName" validate:"required"`
	DisplayName     string `json:"displayName" validate:"required"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
```

Install the validation library and use it in the code:

- Initialize the `validator` object.
- Execute validation checks before further processing structs.

```go
validate := validator.New()
err := validate.Struct(languageDtoInstance) // Perform validation on the instance
if err != nil {
    // Handle validation error
}
```

This ensures validation on all fields before data is processed further.
