package powerplatform

import (
	"context"
	"fmt"
	"strings"

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

	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
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
	EnvironmentClient EnvironmentClient
	ProviderTypeName  string
	TypeName          string
}

type EnvironmentResourceModel struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Url             types.String `tfsdk:"url"`
	Domain          types.String `tfsdk:"domain"`
	Location        types.String `tfsdk:"location"`
	EnvironmentType types.String `tfsdk:"environment_type"`
	//CommonDataServiceDatabaseType types.String `tfsdk:"common_data_service_database_type"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	LanguageName    types.Int64  `tfsdk:"language_code"`
	CurrencyCode    types.String `tfsdk:"currency_code"`
	//IsCustomControlInCanvasAppsEnabled types.Bool   `tfsdk:"is_custom_control_in_canvas_apps_enabled"`
	Version          types.String `tfsdk:"version"`
	Templates        []string     `tfsdk:"templates"`
	TemplateMetadata types.String `tfsdk:"template_metadata"`
	LinkedAppType    types.String `tfsdk:"linked_app_type"`
	LinkedAppId      types.String `tfsdk:"linked_app_id"`
	LinkedAppUrl     types.String `tfsdk:"linked_app_url"`
}

func (r *EnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		MarkdownDescription: "PowerPlatform environment",
		Description:         "PowerPlatform environment",

		Attributes: map[string]schema.Attribute{
			//"id": schema.StringAttribute{
			//	Computed: true,
			//},
			"currency_code": schema.StringAttribute{
				Description:         "Unique currency code",
				MarkdownDescription: "Unique currency name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(EnvironmentCurrencyCodes...),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id 	(guid)",
				Description:         "Unique environment id (guid)",
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
					stringvalidator.OneOf(EnvironmentLocations...),
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
					stringvalidator.OneOf(EnvironmentTypes...),
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
					int64validator.OneOf(EnvironmentLanguages...),
				},
			},
			"version": schema.StringAttribute{
				Description:         "Version of the environment",
				MarkdownDescription: "Version of the environment",
				Computed:            true,
			},
			"templates": schema.ListAttribute{
				Description:         "The selected instance provisioning template (if any)",
				MarkdownDescription: "The selected instance provisioning template (if any)",
				Optional:            true,
				ElementType:         types.StringType,
				// Validators: []validator.String{
				// 	stringvalidator.OneOf(EnvironmentCurrencyCodes...),
				// },
			},
			"template_metadata": schema.StringAttribute{
				Description:         "JSON representation of the environment deployment metadata",
				MarkdownDescription: "JSON representation of the environment deployment metadata",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"linked_app_type": schema.StringAttribute{
				Description:         "The type of the linked D365 application",
				MarkdownDescription: "The type of the linked D365 application",
				Computed:            true,
			},
			"linked_app_id": schema.StringAttribute{
				Description:         "The GUID of the linked D365 application",
				MarkdownDescription: "The GUID of the linked D365 application",
				Computed:            true,
			},
			"linked_app_url": schema.StringAttribute{
				Description:         "The URL of the linked D365 application",
				MarkdownDescription: "The URL of the linked D365 application",
				Computed:            true,
			},
		},
	}
}

func (r *EnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBapi := req.ProviderData.(*clients.ProviderClient).BapiApi.Client
	clientDv := req.ProviderData.(*clients.ProviderClient).DataverseApi.Client

	if clientBapi == nil || clientDv == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.EnvironmentClient = NewEnvironmentClient(clientBapi, clientDv)
}

func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EnvironmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	envToCreate := EnvironmentCreateDto{
		Location: plan.Location.ValueString(),
		Properties: EnvironmentCreatePropertiesDto{
			DisplayName:    plan.DisplayName.ValueString(),
			DataBaseType:   "CommonDataService",
			BillingPolicy:  "",
			EnvironmentSku: plan.EnvironmentType.ValueString(),
			LinkedEnvironmentMetadata: EnvironmentCreateLinkEnvironmentMetadataDto{
				BaseLanguage:    int(plan.LanguageName.ValueInt64()),
				DomainName:      plan.Domain.ValueString(),
				SecurityGroupId: plan.SecurityGroupId.ValueString(),
				Currency: EnvironmentCreateCurrency{
					Code: plan.CurrencyCode.ValueString(),
				},
				//Templates:        plan.Templates,
				//TemplateMetadata: EnvironmentCreateTemplateMetadata{},
			},
		},
	}

	envDto, err := r.EnvironmentClient.CreateEnvironment(ctx, envToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	env := ConvertFromEnvironmentDto(*envDto, plan.CurrencyCode.ValueString())

	plan.Id = env.EnvironmentId
	plan.DisplayName = env.DisplayName
	plan.OrganizationId = env.OrganizationId
	plan.SecurityGroupId = env.SecurityGroupId
	plan.LanguageName = env.LanguageName
	plan.CurrencyCode = types.StringValue(envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)
	plan.Domain = env.Domain
	plan.Url = env.Url
	plan.EnvironmentType = env.EnvironmentType
	plan.Version = env.Version
	plan.LinkedAppType = env.LinkedAppType
	plan.LinkedAppId = env.LinkedAppId
	plan.LinkedAppUrl = env.LinkedAppURL

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))

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

	envDto, err := r.EnvironmentClient.GetEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
	if err != nil {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", envDto.Name), err.Error())
	} else {
		state.CurrencyCode = types.StringValue(defaultCurrency.IsoCurrencyCode)
	}

	env := ConvertFromEnvironmentDto(*envDto, state.CurrencyCode.ValueString())

	state.Id = env.EnvironmentId
	state.DisplayName = env.DisplayName
	state.OrganizationId = env.OrganizationId
	state.SecurityGroupId = env.SecurityGroupId
	state.LanguageName = env.LanguageName
	state.Domain = env.Domain
	state.Url = env.Url
	state.CurrencyCode = env.CurrencyCode
	state.EnvironmentType = env.EnvironmentType
	state.Version = env.Version
	state.LanguageName = env.LanguageName
	state.Location = env.Location
	state.LinkedAppId = env.LinkedAppId
	state.LinkedAppType = env.LinkedAppType
	state.LinkedAppUrl = env.LinkedAppURL

	//TODO move to separate function
	ctx = tflog.SetField(ctx, "id", state.Id.ValueString())
	ctx = tflog.SetField(ctx, "display_name", state.DisplayName.ValueString())
	ctx = tflog.SetField(ctx, "url", state.Url.ValueString())
	ctx = tflog.SetField(ctx, "domain", state.Domain.ValueString())
	ctx = tflog.SetField(ctx, "location", state.Location.ValueString())
	ctx = tflog.SetField(ctx, "environment_type", state.EnvironmentType.ValueString())
	ctx = tflog.SetField(ctx, "organization_id", state.OrganizationId.ValueString())
	ctx = tflog.SetField(ctx, "security_group_id", state.SecurityGroupId.ValueString())
	ctx = tflog.SetField(ctx, "language_code", state.LanguageName.ValueInt64())
	ctx = tflog.SetField(ctx, "currency_code", state.CurrencyCode.ValueString())
	ctx = tflog.SetField(ctx, "version", state.Version.ValueString())
	ctx = tflog.SetField(ctx, "template", strings.Join(state.Templates, " "))

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with id %s", r.ProviderTypeName, state.Id.ValueString()))

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

		envToUpdate := EnvironmentDto{
			Id:       plan.Id.ValueString(),
			Name:     plan.DisplayName.ValueString(),
			Type:     plan.EnvironmentType.ValueString(),
			Location: plan.Location.ValueString(),
			Properties: EnvironmentPropertiesDto{
				DisplayName:    plan.DisplayName.ValueString(),
				EnvironmentSku: plan.EnvironmentType.ValueString(),
				LinkedEnvironmentMetadata: LinkedEnvironmentMetadataDto{
					SecurityGroupId: plan.SecurityGroupId.ValueString(),
					DomainName:      plan.Domain.ValueString(),
				},
			},
		}
		if !plan.LinkedAppId.IsNull() && plan.LinkedAppId.ValueString() != "" {
			envToUpdate.Properties.LinkedAppMetadata = &LinkedAppMetadataDto{
				Type: plan.LinkedAppType.ValueString(),
				Id:   plan.LinkedAppId.ValueString(),
				Url:  plan.LinkedAppUrl.ValueString(),
			}
		} else {
			envToUpdate.Properties.LinkedAppMetadata = nil
		}

		envDto, err := r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), envToUpdate)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.ProviderTypeName), err.Error())
			return
		}

		env := ConvertFromEnvironmentDto(*envDto, plan.CurrencyCode.ValueString())

		plan.Id = env.EnvironmentId
		plan.DisplayName = env.DisplayName
		plan.OrganizationId = env.OrganizationId
		plan.SecurityGroupId = env.SecurityGroupId
		plan.LanguageName = env.LanguageName
		plan.Domain = env.Domain
		plan.Url = env.Url
		plan.CurrencyCode = env.CurrencyCode
		plan.EnvironmentType = env.EnvironmentType
		plan.Version = env.Version
		plan.LanguageName = env.LanguageName
		plan.Location = env.Location
		plan.LinkedAppType = env.LinkedAppType
		plan.LinkedAppId = env.LinkedAppId
		plan.LinkedAppUrl = env.LinkedAppURL
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

	err := r.EnvironmentClient.DeleteEnvironment(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
