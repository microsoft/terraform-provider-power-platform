package powerplaform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	EnvironmentClient ManagedEnvironmentClient
	ProviderTypeName  string
	TypeName          string
}

type ManagedEnvironmentResourceModel struct {
	Id                       types.String `tfsdk:"id"`
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
			"protection_level": schema.StringAttribute{
				MarkdownDescription: "Protection level",
				Description:         "Protection level",
				Required:            true,
				Default:             stringdefault.StaticString("Standard"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Standard", "Basic"),
				},
			},
			"is_usage_insights_disabled": schema.BoolAttribute{
				MarkdownDescription: "Weekly inishgts digest for the environment",
				Description:         "Weekly inishgts digest for the environment",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"is_group_sharing_disabled": schema.BoolAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared",
				Description:         "Limits how widely canvas apps can be shared",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"limit_sharing_mode": schema.StringAttribute{
				MarkdownDescription: "Limits how widely canvas apps can be shared",
				Description:         "Limits how widely canvas apps can be shared",
				Default:             stringdefault.StaticString("NoLimit"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ExcludeSharingToSecurityGroups", "NoLimit"),
				},
			},
			"max_limit_user_sharing": schema.Int64Attribute{
				MarkdownDescription: "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'True', then this values should be '-1'",
				Description:         "Limits how many users can share canvas apps. if 'is_group_sharing_disabled' is 'True', then this values should be '-1'",
				Default:             int64default.StaticInt64(-1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"solution_checker_mode": schema.StringAttribute{
				MarkdownDescription: "Automatically verify solution checker results for security and reliability issues before solution import",
				Description:         "Automatically verify solution checker results for security and reliability issues before solution import",
				Default:             stringdefault.StaticString("None"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("None", "Warn", "Block"),
				},
			},
			"suppress_validation_emails": schema.BoolAttribute{
				MarkdownDescription: "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Description:         "Send emails only when a solution is blocked. If 'False', you'll also get emails when there are warnings",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"maker_onboarding_markdown": schema.StringAttribute{
				MarkdownDescription: "First-time Power Apps makers will see this content in the Studio",
				Description:         "First-time Power Apps makers will see this content in the Studio",
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"maker_onboarding_url": schema.StringAttribute{
				MarkdownDescription: "Maker onboarding 'Learn more' URL",
				Description:         "Maker onboarding 'Learn more' URL",
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ManagedEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// envToCreate := EnvironmentCreateDto{
	// 	Location: plan.Location.ValueString(),
	// 	Properties: EnvironmentCreatePropertiesDto{
	// 		DisplayName:    plan.DisplayName.ValueString(),
	// 		DataBaseType:   "CommonDataService",
	// 		BillingPolicy:  "",
	// 		EnvironmentSku: plan.EnvironmentType.ValueString(),
	// 		LinkedEnvironmentMetadata: EnvironmentCreateLinkEnvironmentMetadataDto{
	// 			BaseLanguage:    int(plan.LanguageName.ValueInt64()),
	// 			DomainName:      plan.Domain.ValueString(),
	// 			SecurityGroupId: plan.SecurityGroupId.ValueString(),
	// 			Currency: EnvironmentCreateCurrency{
	// 				Code: plan.CurrencyCode.ValueString(),
	// 			},
	// 			//Templates:        plan.Templates,
	// 			//TemplateMetadata: EnvironmentCreateTemplateMetadata{},
	// 		},
	// 	},
	// }

	// envDto, err := r.EnvironmentClient.CreateEnvironment(ctx, envToCreate)
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
	// 	return
	// }

	// env := ConvertFromEnvironmentDto(*envDto, plan.CurrencyCode.ValueString())

	// plan.Id = env.EnvironmentId
	// plan.DisplayName = env.DisplayName
	// plan.OrganizationId = env.OrganizationId
	// plan.SecurityGroupId = env.SecurityGroupId
	// plan.LanguageName = env.LanguageName
	// plan.CurrencyCode = types.StringValue(envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)
	// plan.Domain = env.Domain
	// plan.Url = env.Url
	// plan.EnvironmentType = env.EnvironmentType
	// plan.Version = env.Version
	// plan.LinkedAppType = env.LinkedAppType
	// plan.LinkedAppId = env.LinkedAppId
	// plan.LinkedAppUrl = env.LinkedAppURL

	// tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

	// resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	// tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ManagedEnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// envDto, err := r.EnvironmentClient.GetEnvironment(ctx, state.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
	// 	return
	// }

	// defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
	// if err != nil {
	// 	resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", envDto.Name), err.Error())
	// } else {
	// 	state.CurrencyCode = types.StringValue(defaultCurrency.IsoCurrencyCode)
	// }

	// env := ConvertFromEnvironmentDto(*envDto, state.CurrencyCode.ValueString())

	// state.Id = env.EnvironmentId

	// //TODO move to separate function
	// ctx = tflog.SetField(ctx, "id", state.Id.ValueString())

	// tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with id %s", r.ProviderTypeName, state.Id.ValueString()))

	// resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *ManagedEnvironmentResource

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ManagedEnvironmentResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// if plan.DisplayName.ValueString() != state.DisplayName.ValueString() ||
	// 	plan.SecurityGroupId.ValueString() != state.SecurityGroupId.ValueString() ||
	// 	plan.Domain.ValueString() != state.Domain.ValueString() {

	// 	envToUpdate := EnvironmentDto{
	// 		Id:       plan.Id.ValueString(),
	// 		Name:     plan.DisplayName.ValueString(),
	// 		Type:     plan.EnvironmentType.ValueString(),
	// 		Location: plan.Location.ValueString(),
	// 		Properties: EnvironmentPropertiesDto{
	// 			DisplayName:    plan.DisplayName.ValueString(),
	// 			EnvironmentSku: plan.EnvironmentType.ValueString(),
	// 			LinkedEnvironmentMetadata: LinkedEnvironmentMetadataDto{
	// 				SecurityGroupId: plan.SecurityGroupId.ValueString(),
	// 				DomainName:      plan.Domain.ValueString(),
	// 			},
	// 		},
	// 	}
	// 	if !plan.LinkedAppId.IsNull() && plan.LinkedAppId.ValueString() != "" {
	// 		envToUpdate.Properties.LinkedAppMetadata = &LinkedAppMetadataDto{
	// 			Type: plan.LinkedAppType.ValueString(),
	// 			Id:   plan.LinkedAppId.ValueString(),
	// 			Url:  plan.LinkedAppUrl.ValueString(),
	// 		}
	// 	} else {
	// 		envToUpdate.Properties.LinkedAppMetadata = nil
	// 	}

	// 	envDto, err := r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), envToUpdate)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
	// 		return
	// 	}

	// 	env := ConvertFromEnvironmentDto(*envDto, plan.CurrencyCode.ValueString())

	// 	plan.Id = env.EnvironmentId
	// 	plan.DisplayName = env.DisplayName
	// 	plan.OrganizationId = env.OrganizationId
	// 	plan.SecurityGroupId = env.SecurityGroupId
	// 	plan.LanguageName = env.LanguageName
	// 	plan.Domain = env.Domain
	// 	plan.Url = env.Url
	// 	plan.CurrencyCode = env.CurrencyCode
	// 	plan.EnvironmentType = env.EnvironmentType
	// 	plan.Version = env.Version
	// 	plan.LanguageName = env.LanguageName
	// 	plan.Location = env.Location
	// 	plan.LinkedAppType = env.LinkedAppType
	// 	plan.LinkedAppId = env.LinkedAppId
	// 	plan.LinkedAppUrl = env.LinkedAppURL
	// }

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

	// err := r.EnvironmentClient.DeleteEnvironment(ctx, state.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
	// 	return
	// }

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ManagedEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
