package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
	powerplatform_modifiers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
)

var _ resource.Resource = &EnvironmentResource{}
var _ resource.ResourceWithImportState = &EnvironmentResource{}

func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment",
	}
}

type EnvironmentResource struct {
	BapiApiClient    bapi.BapiClientInterface
	ProviderTypeName string
	TypeName         string
}

type EnvironmentResourceModel struct {
	Id              types.String `tfsdk:"id"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	DisplayName     types.String `tfsdk:"display_name"`
	Url             types.String `tfsdk:"url"`
	Domain          types.String `tfsdk:"domain"`
	Location        types.String `tfsdk:"location"`
	EnvironmentType types.String `tfsdk:"environment_type"`
	//CommonDataServiceDatabaseType types.String `tfsdk:"common_data_service_database_type"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	LanguageName    types.Int64  `tfsdk:"language_code"`
	CurrencyName    types.String `tfsdk:"currency_code"`
	//IsCustomControlInCanvasAppsEnabled types.Bool   `tfsdk:"is_custom_control_in_canvas_apps_enabled"`
	Version types.String `tfsdk:"version"`
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "PowerPlatform environment",
		Description:         "PowerPlatform environment",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"currency_code": schema.StringAttribute{
				Description:         "Unique currency code",
				MarkdownDescription: "Unique currency name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(models.EnvironmentCurrencyCodes...),
				},
			},
			"environment_name": schema.StringAttribute{
				MarkdownDescription: "Unique environment name 	(guid)",
				Description:         "Unique environment name (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Description:         "Display name",
				Required:            true,
			},
			"url": schema.StringAttribute{
				Description:         "Url of the environment",
				MarkdownDescription: "Url of the environment",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				Description:         "Domain name of the environment",
				MarkdownDescription: "Domain name of the environment",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Description:         "Location of the environment (europe, unitedstates etc.)",
				MarkdownDescription: "Location of the environment (europe, unitedstates etc.)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(models.EnvironmentLocations...),
				},
			},
			"environment_type": schema.StringAttribute{
				Description:         "Type of the environment (Sandbox, Production etc.)",
				MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(models.EnvironmentTypes...),
				},
			},
			"organization_id": schema.StringAttribute{
				Description:         "Unique organization id (guid)",
				MarkdownDescription: "Unique organization id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_group_id": schema.StringAttribute{
				Description:         "Unique security group id (guid).  For an empty security group, set this property to 0000000-0000-0000-0000-000000000000",
				MarkdownDescription: "Unique security group id (guid).  For an empty security group, set this property to `0000000-0000-0000-0000-000000000000`",
				Required:            true,
			},
			"language_code": schema.Int64Attribute{
				Description:         "Unique language LCID (integer)",
				MarkdownDescription: "Unique language LCID (integer)",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					powerplatform_modifiers.RequireReplaceIntAttributePlanModifier(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(models.EnvironmentLanguages...),
				},
			},
			"version": schema.StringAttribute{
				Description:         "Version of the environment",
				MarkdownDescription: "Version of the environment",
				Computed:            true,
			},
		},
	}
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PowerPlatformProvider).BapiApi.Client.(bapi.BapiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.BapiApiClient = client
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envToCreate := models.EnvironmentCreateDto{
		Location: plan.Location.ValueString(),
		Properties: models.EnvironmentCreatePropertiesDto{
			DisplayName:    plan.DisplayName.ValueString(),
			DataBaseType:   "CommonDataService",
			BillingPolicy:  "",
			EnvironmentSku: plan.EnvironmentType.ValueString(),
			LinkedEnvironmentMetadata: models.EnvironmentCreateLinkEnvironmentMetadataDto{
				BaseLanguage:    int(plan.LanguageName.ValueInt64()),
				DomainName:      plan.Domain.ValueString(),
				SecurityGroupId: plan.SecurityGroupId.ValueString(),
				Currency: models.EnvironmentCreateCurrency{
					Code: plan.CurrencyName.ValueString(),
				},
			},
		},
	}

	envDto, err := r.BapiApiClient.CreateEnvironment(ctx, envToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env := ConvertFromEnvironmentDto(*envDto)

	plan.Id = env.EnvironmentName
	plan.EnvironmentName = env.EnvironmentName
	plan.DisplayName = env.DisplayName
	plan.OrganizationId = env.OrganizationId
	plan.SecurityGroupId = env.SecurityGroupId
	plan.LanguageName = env.LanguageName
	plan.CurrencyName = types.StringValue(envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)
	plan.Domain = env.Domain
	plan.Url = env.Url
	plan.EnvironmentType = env.EnvironmentType
	plan.Version = env.Version

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.EnvironmentName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envDto, err := r.BapiApiClient.GetEnvironment(ctx, state.EnvironmentName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	env := ConvertFromEnvironmentDto(*envDto)

	state.Id = env.EnvironmentName
	state.DisplayName = env.DisplayName
	state.OrganizationId = env.OrganizationId
	state.SecurityGroupId = env.SecurityGroupId
	state.LanguageName = env.LanguageName
	state.Domain = env.Domain
	state.Url = env.Url
	state.EnvironmentType = env.EnvironmentType
	state.Version = env.Version
	state.LanguageName = env.LanguageName
	state.Location = env.Location

	//TODO move to separate function
	ctx = tflog.SetField(ctx, "environment_name", state.EnvironmentName.ValueString())
	ctx = tflog.SetField(ctx, "display_name", state.DisplayName.ValueString())
	ctx = tflog.SetField(ctx, "url", state.Url.ValueString())
	ctx = tflog.SetField(ctx, "domain", state.Domain.ValueString())
	ctx = tflog.SetField(ctx, "location", state.Location.ValueString())
	ctx = tflog.SetField(ctx, "environment_type", state.EnvironmentType.ValueString())
	ctx = tflog.SetField(ctx, "organization_id", state.OrganizationId.ValueString())
	ctx = tflog.SetField(ctx, "security_group_id", state.SecurityGroupId.ValueString())
	ctx = tflog.SetField(ctx, "language_code", state.LanguageName.ValueInt64())
	ctx = tflog.SetField(ctx, "currency_name", state.CurrencyName.ValueString())
	ctx = tflog.SetField(ctx, "version", state.Version.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with environment_name %s", r.ProviderTypeName, state.EnvironmentName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *EnvironmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DisplayName.ValueString() != state.DisplayName.ValueString() ||
		plan.SecurityGroupId.ValueString() != state.SecurityGroupId.ValueString() ||
		plan.Domain.ValueString() != state.Domain.ValueString() {

		envToUpdate := models.EnvironmentDto{
			Id:       plan.Id.ValueString(),
			Name:     plan.EnvironmentName.ValueString(),
			Type:     plan.EnvironmentType.ValueString(),
			Location: plan.Location.ValueString(),
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    plan.DisplayName.ValueString(),
				EnvironmentSku: plan.EnvironmentType.ValueString(),
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					SecurityGroupId: plan.SecurityGroupId.ValueString(),
					DomainName:      plan.Domain.ValueString(),
				},
			},
		}

		envDto, err := r.BapiApiClient.UpdateEnvironment(ctx, plan.EnvironmentName.ValueString(), envToUpdate)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
			return
		}

		env := ConvertFromEnvironmentDto(*envDto)

		plan.Id = env.EnvironmentName
		plan.DisplayName = env.DisplayName
		plan.OrganizationId = env.OrganizationId
		plan.SecurityGroupId = env.SecurityGroupId
		plan.LanguageName = env.LanguageName
		plan.Domain = env.Domain
		plan.Url = env.Url
		plan.EnvironmentType = env.EnvironmentType
		plan.Version = env.Version
		plan.LanguageName = env.LanguageName
		plan.Location = env.Location
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *EnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.BapiApiClient.DeleteEnvironment(ctx, state.EnvironmentName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("environment_name"), req, resp)
}
