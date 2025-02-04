// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

var _ resource.Resource = &ManagedEnvironmentResource{}
var _ resource.ResourceWithImportState = &ManagedEnvironmentResource{}

func NewManagedEnvironmentResource() resource.Resource {
	return &ManagedEnvironmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "managed_environment",
		},
	}
}

func (r *ManagedEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *ManagedEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Manages a \"Managed Environment\" and associated settings",
		MarkdownDescription: "Manages a [Managed Environment](https://learn.microsoft.com/power-platform/admin/managed-environment-overview) and associated settings. A Power Platform Managed Environment is a suite of premium capabilities that allows administrators to manage Power Platform at scale with more control, less effort, and more insights. Once an environment is managed, it unlocks additional features across the Power Platform",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique managed environment settings id (guid)",
				Description:         "Unique managed environment settings id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid), of the environment that is managed by these settings",
				Description:         "Unique environment id (guid), of the environment that is managed by these settings",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protection_level": schema.StringAttribute{
				MarkdownDescription: "Protection level",
				Description:         "Protection level",
				Computed:            true,
			},
			"is_usage_insights_disabled": schema.BoolAttribute{
				MarkdownDescription: "[Weekly insights digest for the environment](https://learn.microsoft.com/power-platform/admin/managed-environment-usage-insights)",
				Description:         "Weekly insights digest for the environment",
				Required:            true,
			},
			"is_group_sharing_disabled": schema.BoolAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared. See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details.",
				Description:         "Limits how widely canvas apps can be shared",
				Required:            true,
			},
			"limit_sharing_mode": schema.StringAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared.  See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details",
				Description:         "Limits how widely canvas apps can be shared.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ExcludeSharingToSecurityGroups", "NoLimit"),
				},
			},
			"max_limit_user_sharing": schema.Int64Attribute{
				MarkdownDescription: "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'",
				Description:         "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'. See [Managed Environment sharing limits](https://learn.microsoft.com/power-platform/admin/managed-environment-sharing-limits) for more details",
				Required:            true,
			},
			"solution_checker_mode": schema.StringAttribute{
				MarkdownDescription: "Automatically verify solution checker results for security and reliability issues before solution import.  See [Solution Checker enforcement](https://learn.microsoft.com/power-platform/admin/managed-environment-solution-checker) for more details.",
				Description:         "Automatically verify solution checker results for security and reliability issues before solution import.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("None", "Warn", "Block"),
				},
			},
			"suppress_validation_emails": schema.BoolAttribute{
				MarkdownDescription: "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Description:         "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Required:            true,
			},
			"solution_checker_rule_overrides": schema.SetAttribute{
				MarkdownDescription: "List of rules to exclude from solution checker.  See [Solution Checker enforcement](https://learn.microsoft.com/power-platform/admin/managed-environment-solution-checker) for more details.",
				Description:         "List of rules to exclude from solution checker",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("meta-remove-dup-reg", "meta-avoid-reg-no-attribute", "meta-avoid-reg-retrieve", "meta-remove-inactive", "web-avoid-unpub-api", "web-avoid-modals", "web-avoid-crm2011-service-odata", "web-avoid-crm2011-service-soap", "web-avoid-browser-specific-api", "web-avoid-2011-api", "web-use-relative-uri", "web-use-async", "web-avoid-window-top", "web-use-client-context", "web-use-navigation-api", "web-use-offline", "web-use-grid-api", "web-avoid-isactivitytype", "meta-avoid-silverlight", "meta-avoid-retrievemultiple-annotation", "web-remove-debug-script", "web-use-strict-mode", "web-use-strict-equality-operators", "web-avoid-eval", "app-formula-issues-high", "app-formula-issues-medium", "app-formula-issues-low", "app-use-delayoutput-text-input", "app-reduce-screen-controls", "app-include-accessible-label", "app-include-alternative-input", "app-avoid-autostart", "app-include-captions", "app-make-focusborder-visible", "app-include-helpful-control-setting", "app-avoid-interactive-html", "app-include-readable-screen-name", "app-include-state-indication-text", "app-include-tab-order", "app-include-tab-index", "flow-avoid-recursive-loop", "flow-avoid-invalid-reference", "flow-outlook-attachment-missing-info", "meta-include-missingunmanageddependencies", "web-remove-alert", "web-remove-console", "web-use-global-context", "web-use-org-setting", "app-testformula-issues-high", "app-testformula-issues-medium", "app-testformula-issues-low", "flow-avoid-connection-mode", "web-avoid-with", "web-avoid-loadtheme", "web-use-getsecurityroleprivilegesinfo", "web-sdl-no-cookies", "web-sdl-no-document-domain", "web-sdl-no-document-write", "web-sdl-no-html-method", "web-sdl-no-inner-html", "web-sdl-no-insecure-url", "web-sdl-no-msapp-exec-unsafe", "web-sdl-no-postmessage-star-origin", "web-sdl-no-winjs-html-unsafe", "connector-validate-brandcolor", "connector-validate-iconimage", "connector-validate-swagger-isproperjson", "connector-validate-swagger", "connector-validate-swagger-extended", "connector-validate-title", "connector-validate-connectionparam-isproperjson", "connector-validate-connectionparameters", "connector-validate-connectionparam-oauth2idp", "meta-license-sales-sdkmessages", "meta-license-sales-entity-operations", "meta-license-sales-customcontrols", "web-use-appsidepane-api", "meta-license-fieldservice-sdkmessages", "meta-license-fieldservice-entity-operations", "meta-license-fieldservice-customcontrols", "meta-avoid-managed-entity-assets", "meta-include-unmanaged-entity-assets", "connector-validate-hexadecimalbrandcolor", "connector-validate-pngiconimage", "connector-validate-iconsize", "connector-validate-backgroundwithbrandiconcolor", "web-unsupported-syntax")),
				},
			},
			"maker_onboarding_markdown": schema.StringAttribute{
				MarkdownDescription: "First-time Power Apps makers will see this content in the Studio.  See [Maker welcome content](https://learn.microsoft.com/power-platform/admin/welcome-content) for more details.",
				Description:         "First-time Power Apps makers will see this content in the Studio",
				Required:            true,
			},
			"maker_onboarding_url": schema.StringAttribute{
				MarkdownDescription: "Maker onboarding 'Learn more' URL. See [Maker welcome content](https://learn.microsoft.com/power-platform/admin/welcome-content) for more details.",
				Description:         "Maker onboarding 'Learn more' URL",
				Required:            true,
			},
		},
	}
}

func (r *ManagedEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	clientApi := client.Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.ManagedEnvironmentClient = newManagedEnvironmentClient(clientApi)
}

func (r *ManagedEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard",
		Settings: &environment.SettingsDto{
			ExtendedSettings: environment.ExtendedSettingsDto{
				ExcludeEnvironmentFromAnalysis: strconv.FormatBool(plan.IsUsageInsightsDisabled.ValueBool()),
				IsGroupSharingDisabled:         strconv.FormatBool(plan.IsGroupSharingDisabled.ValueBool()),
				MaxLimitUserSharing:            strconv.FormatInt(plan.MaxLimitUserSharing.ValueInt64(), 10),
				DisableAiGeneratedDescriptions: "false",
				IncludeOnHomepageInsights:      "false",
				LimitSharingMode:               strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:],
				SolutionCheckerMode:            strings.ToLower(plan.SolutionCheckerMode.ValueString()),
				SuppressValidationEmails:       strconv.FormatBool(plan.SuppressValidationEmails.ValueBool()),
				MakerOnboardingUrl:             plan.MakerOnboardingUrl.ValueString(),
				MakerOnboardingMarkdown:        plan.MakerOnboardingMarkdown.ValueString(),
				SolutionCheckerRuleOverrides:   strings.Join(helpers.SetToStringSlice(plan.SolutionCheckerRuleOverrides), ","),
			},
		},
	}

	err := r.ManagedEnvironmentClient.EnableManagedEnvironment(ctx, managedEnvironmentDto, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
	plan.Id = plan.EnvironmentId
	plan.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)
	plan.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
	plan.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
	plan.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
	plan.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
	plan.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
	plan.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)
	plan.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	state.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)

	if env.Properties.GovernanceConfiguration.Settings != nil {
		maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)

		state.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
		state.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
		state.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
		state.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
		state.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
		state.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
		state.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
		state.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)
		state.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))
	} else {
		state.IsGroupSharingDisabled = types.BoolUnknown()
		state.IsUsageInsightsDisabled = types.BoolUnknown()
		state.MaxLimitUserSharing = types.Int64Unknown()
		state.LimitSharingMode = types.StringUnknown()
		state.SolutionCheckerMode = types.StringUnknown()
		state.SuppressValidationEmails = types.BoolUnknown()
		state.MakerOnboardingUrl = types.StringUnknown()
		state.MakerOnboardingMarkdown = types.StringUnknown()
		state.SolutionCheckerRuleOverrides = types.SetUnknown(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ManagedEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *ManagedEnvironmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard",
		Settings: &environment.SettingsDto{
			ExtendedSettings: environment.ExtendedSettingsDto{
				ExcludeEnvironmentFromAnalysis: strconv.FormatBool(plan.IsUsageInsightsDisabled.ValueBool()),
				IsGroupSharingDisabled:         strconv.FormatBool(plan.IsGroupSharingDisabled.ValueBool()),
				MaxLimitUserSharing:            strconv.FormatInt(plan.MaxLimitUserSharing.ValueInt64(), 10),
				DisableAiGeneratedDescriptions: "false",
				IncludeOnHomepageInsights:      "false",
				LimitSharingMode:               strings.ToLower(plan.LimitSharingMode.ValueString()[:1]) + plan.LimitSharingMode.ValueString()[1:],
				SolutionCheckerMode:            strings.ToLower(plan.SolutionCheckerMode.ValueString()),
				SuppressValidationEmails:       strconv.FormatBool(plan.SuppressValidationEmails.ValueBool()),
				MakerOnboardingUrl:             plan.MakerOnboardingUrl.ValueString(),
				MakerOnboardingMarkdown:        plan.MakerOnboardingMarkdown.ValueString(),
				SolutionCheckerRuleOverrides:   strings.Join(helpers.SetToStringSlice(plan.SolutionCheckerRuleOverrides), ","),
			},
		},
	}

	err := r.ManagedEnvironmentClient.EnableManagedEnvironment(ctx, managedEnvironmentDto, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when enabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
	plan.Id = plan.EnvironmentId
	plan.ProtectionLevel = types.StringValue(env.Properties.GovernanceConfiguration.ProtectionLevel)
	plan.IsUsageInsightsDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.ExcludeEnvironmentFromAnalysis == "true")
	plan.IsGroupSharingDisabled = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.IsGroupSharingDisabled == "true")
	plan.MaxLimitUserSharing = types.Int64Value(maxLimitUserSharing)
	plan.LimitSharingMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.LimitSharingMode[1:])
	plan.SolutionCheckerMode = types.StringValue(strings.ToUpper(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[:1]) + env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerMode[1:])
	plan.SuppressValidationEmails = types.BoolValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SuppressValidationEmails == "true")
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)
	plan.SolutionCheckerRuleOverrides = helpers.StringSliceToSet(strings.Split(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides, ","))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ManagedEnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ManagedEnvironmentClient.DisableManagedEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
}

func (r *ManagedEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
