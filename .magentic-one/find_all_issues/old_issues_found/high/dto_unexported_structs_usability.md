# Title  
Unexported Structs Limit Usability Outside the Package  

##  
/workspaces/terraform-provider-power-platform/internal/services/currencies/dto.go  

## Problem  
The structs (`currenciesDto`, `currenciesArrayDto`, `currenciesPropertiesDto`) defined in this file are unexported because their names start with a lowercase letter. This restricts their usage to the `currencies` package only, preventing other packages from accessing or utilizing these types.  

## Impact  
This limitation reduces the flexibility and usability of the code. If another package needs to use these DTOs, developers will either need to duplicate the types or modify this code, which increases technical debt. Severity: **high**  

## Location  

File: `/workspaces/terraform-provider-power-platform/internal/services/currencies/dto.go`  

## Code Issue  

```go  
type currenciesDto struct {
  Value []currenciesArrayDto `json:"value"`
}

type currenciesArrayDto struct {
  Name       string                  `json:"name"`
  ID         string                  `json:"id"`
  Type       string                  `json:"type"`
  Properties currenciesPropertiesDto `json:"properties"`
}

type currenciesPropertiesDto struct {
  Code            string `json:"code"`
  Symbol          string `json:"symbol"`
  IsTenantDefault bool   `json:"isTenantDefault"`
}
```  

## Fix  

Make the names of the structs exported so that they can be used by other packages. To do this, capitalize the first letter of the struct names.  

```go  
type CurrenciesDto struct {
  Value []CurrenciesArrayDto `json:"value"`
}

type CurrenciesArrayDto struct {
  Name       string                  `json:"name"`
  ID         string                  `json:"id"`
  Type       string                  `json:"type"`
  Properties CurrenciesPropertiesDto `json:"properties"`
}

type CurrenciesPropertiesDto struct {
  Code            string `json:"code"`
  Symbol          string `json:"symbol"`
  IsTenantDefault bool   `json:"isTenantDefault"`
}
```