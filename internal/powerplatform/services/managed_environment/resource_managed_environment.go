package powerplatform

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
)

var _ resource.Resource = &ManagedEnvironmentResource{}
var _ resource.ResourceWithImportState = &ManagedEnvironmentResource{}

func NewManagedEnvironmentResource() resource.Resource {
	return &ManagedEnvironmentResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_managed_environment",
	}
}

type ManagedEnvironmentResource struct {
	ManagedEnvironmentClient ManagedEnvironmentClient
	ProviderTypeName         string
	TypeName                 string
}

type ManagedEnvironmentResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	EnvironmentId            types.String `tfsdk:"environment_id"`
	ProtectionLevel          types.String `tfsdk:"protection_level"`
	IsUsageInsightsDisabled  types.Bool   `tfsdk:"is_usage_insights_disabled"`
	IsGroupSharingDisabled   types.Bool   `tfsdk:"is_group_sharing_disabled"`
	MaxLimitUserSharing      types.Int64  `tfsdk:"max_limit_user_sharing"`
	LimitSharingMode         types.String `tfsdk:"limit_sharing_mode"`
	SolutionCheckerMode      types.String `tfsdk:"solution_checker_mode"`
	SuppressValidationEmails types.Bool   `tfsdk:"suppress_validation_emails"`
	//SolutionCheckerRuleOverrides  types.String `tfsdk:"solution_checker_rule_overrides"`
	MakerOnboardingUrl      types.String `tfsdk:"maker_onboarding_url"`
	MakerOnboardingMarkdown types.String `tfsdk:"maker_onboarding_markdown"`
}

func (r *ManagedEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *ManagedEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "Managed environment settings",
		Description:         "Managed environment settings",

		Attributes: map[string]schema.Attribute{
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
				MarkdownDescription: "Weekly inishgts digest for the environment",
				Description:         "Weekly inishgts digest for the environment",
				Required:            true,
			},
			"is_group_sharing_disabled": schema.BoolAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared",
				Description:         "Limits how widely canvas apps can be shared",
				Required:            true,
			},
			"limit_sharing_mode": schema.StringAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared",
				Description:         "Limits how widely canvas apps can be shared",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ExcludeSharingToSecurityGroups", "NoLimit"),
				},
			},
			"max_limit_user_sharing": schema.Int64Attribute{
				MarkdownDescription: "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'",
				Description:         "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'False', then this values should be '-1'",
				Required:            true,
			},
			"solution_checker_mode": schema.StringAttribute{
				MarkdownDescription: "Automatically verify solution checker results for security and reliability issues before solution import",
				Description:         "Automatically verify solution checker results for security and reliability issues before solution import",
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
			"maker_onboarding_markdown": schema.StringAttribute{
				MarkdownDescription: "First-time Power Apps makers will see this content in the Studio",
				Description:         "First-time Power Apps makers will see this content in the Studio",
				Required:            true,
			},
			"maker_onboarding_url": schema.StringAttribute{
				MarkdownDescription: "Maker onboarding 'Learn more' URL",
				Description:         "Maker onboarding 'Learn more' URL",
				Required:            true,
			},
		},
	}
}

func (r *ManagedEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.ManagedEnvironmentClient = NewManagedEnvironmentClient(clientApi)
}

func (r *ManagedEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard", //plan.ProtectionLevel.ValueString(),
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
				//SolutionCheckerRuleOverrides:   "",
				MakerOnboardingUrl: plan.MakerOnboardingUrl.ValueString(),
				//MakerOnboardingTimestamp:       nil
				MakerOnboardingMarkdown: plan.MakerOnboardingMarkdown.ValueString(),
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
	//plan.SolutionCheckerRuleOverrides = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides)
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.ManagedEnvironmentClient.environmentClient.GetEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
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
		//state.SolutionCheckerRuleOverrides = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides)
		state.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
		state.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)
	} else {
		state.IsGroupSharingDisabled = types.BoolUnknown()
		state.IsUsageInsightsDisabled = types.BoolUnknown()
		state.MaxLimitUserSharing = types.Int64Unknown()
		state.LimitSharingMode = types.StringUnknown()
		state.SolutionCheckerMode = types.StringUnknown()
		state.SuppressValidationEmails = types.BoolUnknown()
		//state.SolutionCheckerRuleOverrides = types.StringUnknown()
		state.MakerOnboardingUrl = types.StringUnknown()
		state.MakerOnboardingMarkdown = types.StringUnknown()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	// var state *ManagedEnvironmentResource
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	managedEnvironmentDto := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Standard", //plan.ProtectionLevel.ValueString(),
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
				//SolutionCheckerRuleOverrides:   "",
				MakerOnboardingUrl: plan.MakerOnboardingUrl.ValueString(),
				//MakerOnboardingTimestamp:       nil
				MakerOnboardingMarkdown: plan.MakerOnboardingMarkdown.ValueString(),
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
	//plan.SolutionCheckerRuleOverrides = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.SolutionCheckerRuleOverrides)
	plan.MakerOnboardingUrl = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingUrl)
	plan.MakerOnboardingMarkdown = types.StringValue(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MakerOnboardingMarkdown)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ManagedEnvironmentClient.DisableManagedEnvironment(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling managed environment %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
