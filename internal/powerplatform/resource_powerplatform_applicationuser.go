package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	powerplatform_modifiers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
)

//var _ resource.Resource = &ApplicationUserResource{}
//var _ resource.ResourceWithImportState = &ApplicationUserResource{}

/*
func NewApplicationUserResource() resource.Resource {
	return &ApplicationUserResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_ApplicationUser",
	}

}
*/

type ApplicationUserResource struct {
	BapiApiClient    powerplatform_bapi.ApiClientInterface
	ProviderTypeName string
	TypeName         string
}

type ApplicationUserResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	ApplicationUserFileChecksum types.String `tfsdk:"ApplicationUser_file_checksum"`
	SettingsFileChecksum        types.String `tfsdk:"settings_file_checksum"`
	EnvironmentName             types.String `tfsdk:"environment_id"`
	ApplicationUserName         types.String `tfsdk:"ApplicationUser_name"`
	ApplicationUserVersion      types.String `tfsdk:"ApplicationUser_version"`
	ApplicationUserFile         types.String `tfsdk:"ApplicationUser_file"`
	SettingsFile                types.String `tfsdk:"settings_file"`
	IsManaged                   types.Bool   `tfsdk:"is_managed"`
	DisplayName                 types.String `tfsdk:"display_name"`
}

func (r *ApplicationUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *ApplicationUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for importing ApplicationUsers in Power Platform environments",
		Description:         "Resource for importing exporting ApplicationUsers in Power Platform environments",
		Attributes: map[string]schema.Attribute{
			"applicationUser_file_checksum": schema.StringAttribute{
				MarkdownDescription: "Checksum of the applicationUser file",
				Description:         "Checksum of the applicationUser file",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					powerplatform_modifiers.SyncAttributePlanModifier("applicationUser_file"),
				},
			},
			"applicationUser_file": schema.StringAttribute{
				MarkdownDescription: "Path to the applicationUser file",
				Description:         "Path to the applicationUser file",
				Required:            true,
			},
			"settings_file_checksum": schema.StringAttribute{
				MarkdownDescription: "Checksum of the settings file",
				Description:         "Checksum of the settings file",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					powerplatform_modifiers.SyncAttributePlanModifier("settings_file"),
				},
			},
			"settings_file": schema.StringAttribute{
				MarkdownDescription: "Path to the settings file. The settings file uses the same format as pac cli. See https://learn.microsoft.com/en-us/power-platform/alm/conn-ref-env-variables-build-tools#deployment-settings-file for more details",
				Description:         "Path to the settings file. The settings file uses the same format as pac cli. See https://learn.microsoft.com/en-us/power-platform/alm/conn-ref-env-variables-build-tools#deployment-settings-file for more details",
				Optional:            true,
			},

			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the applicationUser",
				Description:         "Unique identifier of the applicationUser",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"applicationUser_name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the applicationUser",
				Description:         "Unique name of the applicationUser",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Name of the environment where the applicationUser is imported",
				Description:         "Name of the environment where the applicationUser is imported",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the applicationUser",
				Description:         "Display name of the applicationUser",
				Computed:            true,
			},
			"is_managed": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the applicationUser is managed or not",
				Description:         "Indicates whether the applicationUser is managed or not",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"applicationUser_version": schema.StringAttribute{
				MarkdownDescription: "Version of the applicationUser",
				Description:         "Version of the applicationUser",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApplicationUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PowerPlatformProvider).bapiClient.(powerplatform_bapi.ApiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.BapiApiClient = client
}

/*
	func (r *ApplicationUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
		var plan *ApplicationUserResourceModel

		tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

		if resp.Diagnostics.HasError() {
			return
		}

		applicationUser := r.ImportApplicationUserDto(ctx, plan, &resp.Diagnostics)

		if resp.Diagnostics.HasError() {
			return
		}

		plan.ApplicationUserName = types.StringValue(applicationUser.Name)
		plan.ApplicationUserVersion = types.StringValue(applicationUser.Version)
		plan.IsManaged = types.BoolValue(applicationUser.IsManaged)
		plan.DisplayName = types.StringValue(applicationUser.DisplayName)
		plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentName.ValueString(), applicationUser.Name))

		tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
	}
*/
func (r *ApplicationUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ApplicationUserResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	applicationUsers, err := r.BapiApiClient.GetApplicationUser(ctx, state.EnvironmentName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	applicationUserFound := false
	for _, applicationUser := range applicationUsers {
		if applicationUser.Name == state.ApplicationUserName.ValueString() {
			state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentName.ValueString(), applicationUser.Name))
			state.ApplicationUserName = types.StringValue(applicationUser.Name)
			//TODO test a case when applicationUser version changes
			state.ApplicationUserVersion = types.StringValue(applicationUser.Version)
			state.IsManaged = types.BoolValue(applicationUser.IsManaged)
			state.DisplayName = types.StringValue(applicationUser.DisplayName)
			applicationUserFound = true
			break
		}
	}

	if !applicationUserFound {

		state.Id = types.StringNull()
		state.ApplicationUserName = types.StringNull()
		state.ApplicationUserVersion = types.StringNull()
		state.IsManaged = types.BoolNull()
		state.DisplayName = types.StringNull()
		state.SettingsFileChecksum = types.StringNull()
		state.ApplicationUserFileChecksum = types.StringNull()

		tflog.Debug(ctx, fmt.Sprintf("ApplicationUser %s not found", state.ApplicationUserName.ValueString()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

/*
func (r *ApplicationUserResource) importApplicationUser(ctx context.Context, plan *ApplicationUserResourceModel, diagnostics *diag.Diagnostics) *models.ApplicationUserDto {

	s := models.ImportApplicationUserDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: true,
		ComponentParameters:              make([]interface{}, 0),
	}

	applicationUserContent, err := os.ReadFile(plan.ApplicationUserFile.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when reading applicationUser file %s", plan.ApplicationUserFile.ValueString()), err.Error())
	}

	settingsContent := make([]byte, 0)
	//todo check if settings file is not empty in .tf
	if plan.SettingsFile.ValueString() != "" {
		settingsContent, err = os.ReadFile(plan.SettingsFile.ValueString())
		if err != nil {
			diagnostics.AddError(fmt.Sprintf("Client error when reading settings file %s", plan.SettingsFile.ValueString()), err.Error())
		}
	}

	applicationUser, err := r.BapiApiClient.CreateApplicationUser(ctx, plan.EnvironmentName.ValueString(), s, applicationUserContent, settingsContent)
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when importing applicationUser %s", plan.ApplicationUserFile), err.Error())
	}
	return applicationUser
}
*/
/*
func (r *ApplicationUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	var plan *ApplicationUserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ApplicationUserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	applicationUser := r.importApplicationUser(ctx, plan, &resp.Diagnostics)

	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentName.ValueString(), applicationUser.Name))

	plan.ApplicationUserName = types.StringValue(applicationUser.Name)
	plan.ApplicationUserVersion = types.StringValue(applicationUser.Version)
	plan.IsManaged = types.BoolValue(applicationUser.IsManaged)
	plan.DisplayName = types.StringValue(applicationUser.DisplayName)

	plan.SettingsFileChecksum = types.StringUnknown()
	if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.SettingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
		} else {
			plan.SettingsFileChecksum = types.StringValue(value)
		}
	}

	plan.ApplicationUserFileChecksum = types.StringUnknown()
	if !plan.ApplicationUserFile.IsNull() && !plan.ApplicationUserFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.ApplicationUserFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for applicationUser file", err.Error())
		} else {
			plan.ApplicationUserFileChecksum = types.StringValue(value)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}
*/
/*
func (r *ApplicationUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ApplicationUserResourceModel
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !state.EnvironmentName.IsNull() && !state.ApplicationUserName.IsNull() {
		err := r.BapiApiClient.DeleteApplicationUser(ctx, state.EnvironmentName.ValueString(), state.ApplicationUserName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.ProviderTypeName), err.Error())
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}
*/
/*
func (r *ApplicationUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
*/
