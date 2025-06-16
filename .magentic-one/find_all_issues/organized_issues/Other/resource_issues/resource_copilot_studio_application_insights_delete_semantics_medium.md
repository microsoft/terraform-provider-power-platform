# Title

Potential Logic Issue: "Delete" operation clears configuration rather than deleting

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

In the `Delete` method, instead of performing an actual delete API operation, the method creates an update DTO with empty strings or `false` values for all fields and submits an update request, clearing out the configuration rather than deleting the resource itself. This is highlighted by the comment:  
```go
// You can't really create a config, so treat a create as an update
```
and similarly in Delete.

This may be a workaround due to API limitations, but it's not clearly documented, and users may not expect that a deleted resource is simply \"zeroed out\" in the backend.

## Impact

Severity: **Medium**

This discrepancy could lead to confusion and potential data leakage or unexpected charges, since the resource might still exist in a \"cleared\" but present state. It is not deleted from the backend, only emptied, which does not always meet user expectations or compliance needs.

If this is due to an API constraint, a more explicit comment and documentation are required in-code and in documentation.

## Location

- Method: `Delete`

## Code Issue

```go
var state *ResourceModel

resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

if resp.Diagnostics.HasError() {
	return
}

appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*state)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
}
appInsightsConfigToCreate.AppInsightsConnectionString = ""
appInsightsConfigToCreate.IncludeSensitiveInformation = false
appInsightsConfigToCreate.IncludeActivities = false
appInsightsConfigToCreate.IncludeActions = false

_, err = r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, state.BotId.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
	return
}
resp.State.RemoveResource(ctx)
```

## Fix

If the API does not support deleting the configuration, document this behavior clearly both in the method and the Terraform resource documentation.

Consider adding a warning `tflog.Warn` in the `Delete` method about the behavior. For example:

```go
// The API does not support deleting an Application Insights configuration; instead, this 'deletes' the configuration by setting all fields to empty/false.
tflog.Warn(ctx, "Delete called: Will clear configuration but not remove resource from backend, as delete operation is not supported by platform/API.")
```
And update resource documentation to note this important behavior up front for users.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_copilot_studio_application_insights_delete_semantics_medium.md`
