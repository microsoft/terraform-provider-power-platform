### 4. Title
Misnamed Constant for `BooleanRegex`

##  
`/workspaces/terraform-provider-power-platform/internal/helpers/regex.go`  

## Problem  
The name `BooleanRegex` suggests a match for true/false boolean values only, but the regex does not enforce case sensitivity. This could lead to mismatches if case sensitivity is required in certain contexts.  

## Impact  
This can cause unexpected matches for inputs like "TRUE" or "False" depending on where this constant is used. Severity: **medium**.  

## Location  
```go
const (  
	BooleanRegex = "^(true|false)$"  
)  
```  

## Fix  
Rename the constant or update its regex to match specific requirements, such as enforcing case sensitivity:  

```go
const (  
	// Matches lowercase "true" or "false" only.
	BooleanRegex = "^(true|false)$"  

	// OR to make it case-insensitive:
	// BooleanRegex = "(?i)^(true|false)$"
)
```