# Title

Misnamed and Unexported Struct `linkEnterprosePolicyDto`

##

`internal/services/enterprise_policy/dto.go`

## Problem

The struct `linkEnterprosePolicyDto` is poorly named and uses an inconsistent naming convention with a typo ("Enterprose" instead of "Enterprise"), which could lead to confusion or incorrect implementation. Additionally, it is unexported, meaning it cannot be used outside the package, which might limit its utility depending on its intended purpose.

## Impact

The typo creates ambiguity and lower readability. The unexported nature of the struct may create issues for other parts of the application if the struct is supposed to be used more broadly. Severity: Medium.

## Location

Line defining `linkEnterprosePolicyDto` (starting likely from `type`).

## Code Issue

Current structure of the code:

```go
type linkEnterprosePolicyDto struct {
	SystemId string `json:"systemId"`
}
```

## Fix

Rename the struct to `LinkEnterprisePolicyDTO` for proper naming conventions, capitalization, and typo correction. Additionally, decide based on your requirements whether it should be exported.

```go
type LinkEnterprisePolicyDTO struct {
	SystemId string `json:"systemId"`
}
```

Explanation: 
- The updated name uses camel case and corrects the typo from "Enterprose" to "Enterprise". 
- By capitalizing the first letter, the struct becomes exported, ensuring compatibility with packages outside the current one if applicable.
