# Title

Redundant Structs: Duplicate Resource and Datasource Policy Models

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/models.go

## Problem

The file defines two nearly identical structs: `dataLossPreventionPolicyDatasourceModel` and `dataLossPreventionPolicyResourceModel`. Both structs contain the same set of fields (except for possibly minor differences, for example, the order or type of timeouts). Maintaining these two almost identical structs can cause unnecessary redundancy and future maintenance issues if a field is updated in one struct but not the other.

## Impact

Severity: Medium

Duplicated code increases maintenance overhead and the risk of subtle bugs if fields or tags are changed in one struct but not the other. It also makes code less readable and harder to refactor.

## Location

```go
type dataLossPreventionPolicyDatasourceModel struct { ... }

type dataLossPreventionPolicyResourceModel struct { ... }
```

## Code Issue

```go
type dataLossPreventionPolicyDatasourceModel struct {
	Timeouts timeouts.Value                            `tfsdk:"timeouts"`
	Policies []dataLossPreventionPolicyDatasourceModel `tfsdk:"policies"`
    ...
}

type dataLossPreventionPolicyResourceModel struct {
	Timeouts timeouts.Value                            `tfsdk:"timeouts"`
    ...
}
```

## Fix

Consider consolidating these structs into a single model where possible, parameterizing differences (like `Timeouts` type) if needed. For example:

```go
type dataLossPreventionPolicyModel struct {
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
	Id                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	DefaultConnectorsClassification   types.String   `tfsdk:"default_connectors_classification"`
	EnvironmentType                   types.String   `tfsdk:"environment_type"`
	CreatedBy                         types.String   `tfsdk:"created_by"`
	CreatedTime                       types.String   `tfsdk:"created_time"`
	LastModifiedBy                    types.String   `tfsdk:"last_modified_by"`
	LastModifiedTime                  types.String   `tfsdk:"last_modified_time"`
	Environments                      types.Set      `tfsdk:"environments"`
	NonBusinessConnectors             types.Set      `tfsdk:"non_business_connectors"`
	BusinessGeneralConnectors         types.Set      `tfsdk:"business_connectors"`
	BlockedConnectors                 types.Set      `tfsdk:"blocked_connectors"`
	CustomConnectorsPatterns          types.Set      `tfsdk:"custom_connectors_patterns"`
}
```
Use type aliases or wrappers if a slight distinction is needed between datasource and resource.
