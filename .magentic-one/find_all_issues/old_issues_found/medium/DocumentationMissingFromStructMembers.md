# Title

Missing Documentation for Struct Members

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go`

## Problem

Most struct members are not accompanied by comments documenting their purpose and usage. Clear documentation aids readability and ensures that developers using or modifying the code understand its intended behavior.

## Impact

The lack of comments on struct members degrades code maintainability and could lead to incorrect assumptions or misuse by developers unfamiliar with the context. Severity: **medium**.

## Location

Occurrences throughout the file, including but not limited to the below example:

### Example 1

Struct `TenantSettingsResourceModel`.

## Code Issue

```go
type TenantSettingsResourceModel struct {
	Timeouts                                       timeouts.Value `tfsdk:"timeouts"`
	Id                                             types.String   `tfsdk:"id"`
	WalkMeOptOut                                   types.Bool     `tfsdk:"walk_me_opt_out"`
	DisableNPSCommentsReachout                     types.Bool     `tfsdk:"disable_nps_comments_reachout"`
	DisableNewsletterSendout                       types.Bool     `tfsdk:"disable_newsletter_sendout"`
	DisableEnvironmentCreationByNonAdminUsers      types.Bool     `tfsdk:"disable_environment_creation_by_non_admin_users"`
	DisablePortalsCreationByNonAdminUsers          types.Bool     `tfsdk:"disable_portals_creation_by_non_admin_users"`
	DisableSurveyFeedback                          types.Bool     `tfsdk:"disable_survey_feedback"`
	DisableTrialEnvironmentCreationByNonAdminUsers types.Bool     `tfsdk:"disable_trial_environment_creation_by_non_admin_users"`
	DisableCapacityAllocationByEnvironmentAdmins   types.Bool     `tfsdk:"disable_capacity_allocation_by_environment_admins"`
	DisableSupportTicketsVisibleByAllUsers         types.Bool     `tfsdk:"disable_support_tickets_visible_by_all_users"`
	PowerPlatform                                  types.Object   `tfsdk:"power_platform"`
}
```

## Fix

Add inline documentation comments to each member of the struct. Example:

```go
type TenantSettingsResourceModel struct {
	// Timeouts defines the customizable timeout settings for operations.
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	
	// Id is a unique identifier for a tenant setting instance.
	Id types.String   `tfsdk:"id"`

	// WalkMeOptOut indicates whether WalkMe features are opted out.
	WalkMeOptOut types.Bool     `tfsdk:"walk_me_opt_out"`

	// DisableNPSCommentsReachout disables the reachout for NPS comments.
	DisableNPSCommentsReachout types.Bool     `tfsdk:"disable_nps_comments_reachout"`
	// Additional fields should be similarly documented...
}
```