// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	powerplatform_modifiers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
)

var _ resource.Resource = &SolutionResource{}
var _ resource.ResourceWithImportState = &SolutionResource{}

func NewSolutionResource() resource.Resource {
	return &SolutionResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_solution",
	}
}

type SolutionResource struct {
	SolutionClient   SolutionClient
	ProviderTypeName string
	TypeName         string
}

type SolutionResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	SolutionFileChecksum types.String `tfsdk:"solution_file_checksum"`
	SettingsFileChecksum types.String `tfsdk:"settings_file_checksum"`
	EnvironmentId        types.String `tfsdk:"environment_id"`
	SolutionName         types.String `tfsdk:"solution_name"`
	SolutionVersion      types.String `tfsdk:"solution_version"`
	SolutionFile         types.String `tfsdk:"solution_file"`
	SettingsFile         types.String `tfsdk:"settings_file"`
	IsManaged            types.Bool   `tfsdk:"is_managed"`
	DisplayName          types.String `tfsdk:"display_name"`
}

func (r *SolutionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *SolutionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for importing solutions in Power Platform environments",
		MarkdownDescription: "Resource for importing exporting solutions in Power Platform environments.  This is the equivalent of the [`pac solution import`](https://learn.microsoft.com/power-platform/developer/cli/reference/solution#pac-solution-import) command in the Power Platform CLI.",
		Attributes: map[string]schema.Attribute{
			"solution_file_checksum": schema.StringAttribute{
				MarkdownDescription: "Checksum of the solution file",
				Description:         "Checksum of the solution file",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					powerplatform_modifiers.SyncAttributePlanModifier("solution_file"),
				},
			},
			"solution_file": schema.StringAttribute{
				MarkdownDescription: "Path to the solution file",
				Description:         "Path to the solution file",
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
				MarkdownDescription: "Path to the settings file. The settings file uses the same format as pac cli. See https://learn.microsoft.com/power-platform/alm/conn-ref-env-variables-build-tools#deployment-settings-file for more details",
				Description:         "Path to the settings file. The settings file uses the same format as pac cli. See https://learn.microsoft.com/power-platform/alm/conn-ref-env-variables-build-tools#deployment-settings-file for more details",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the solution",
				Description:         "Unique identifier of the solution",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"solution_name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the solution",
				Description:         "Unique name of the solution",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the environment where the solution is imported",
				Description:         "Id of the environment where the solution is imported",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the solution",
				Description:         "Display name of the solution",
				Computed:            true,
			},
			"is_managed": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the solution is managed or not",
				Description:         "Indicates whether the solution is managed or not",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"solution_version": schema.StringAttribute{
				MarkdownDescription: "Version of the solution",
				Description:         "Version of the solution",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SolutionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.SolutionClient = NewSolutionClient(clientApi)
}

func (r *SolutionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *SolutionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	solution := r.importSolution(ctx, plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.SolutionName = types.StringValue(solution.Name)
	plan.SolutionVersion = types.StringValue(solution.Version)
	plan.IsManaged = types.BoolValue(solution.IsManaged)
	plan.DisplayName = types.StringValue(solution.DisplayName)
	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solution.Name))

	plan.SettingsFileChecksum = types.StringUnknown()
	if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.SettingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
		} else {
			plan.SettingsFileChecksum = types.StringValue(value)
			tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
		}
	} else {
		plan.SettingsFileChecksum = types.StringNull()
	}

	plan.SolutionFileChecksum = types.StringUnknown()
	if !plan.SolutionFile.IsNull() && !plan.SolutionFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.SolutionFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for solution file", err.Error())
		} else {
			plan.SolutionFileChecksum = types.StringValue(value)
			tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of solution file: %s", value))
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *SolutionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *SolutionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	solutions, err := r.SolutionClient.GetSolutions(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	solutionFound := false
	for _, solution := range solutions {
		if solution.Name == state.SolutionName.ValueString() {
			state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), solution.Name))
			state.SolutionName = types.StringValue(solution.Name)
			//TODO test a case when solution version changes
			state.SolutionVersion = types.StringValue(solution.Version)
			state.IsManaged = types.BoolValue(solution.IsManaged)
			state.DisplayName = types.StringValue(solution.DisplayName)
			solutionFound = true
			break
		}
	}

	if !solutionFound {

		state.Id = types.StringNull()
		state.SolutionName = types.StringNull()
		state.SolutionVersion = types.StringNull()
		state.IsManaged = types.BoolNull()
		state.DisplayName = types.StringNull()
		state.SettingsFileChecksum = types.StringNull()
		state.SolutionFileChecksum = types.StringNull()

		tflog.Debug(ctx, fmt.Sprintf("Solution %s not found", state.SolutionName.ValueString()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *SolutionResource) importSolution(ctx context.Context, plan *SolutionResourceModel, diagnostics *diag.Diagnostics) *SolutionDto {

	s := ImportSolutionDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: true,
		ComponentParameters:              make([]interface{}, 0),
	}

	solutionContent, err := os.ReadFile(plan.SolutionFile.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", plan.SolutionFile.ValueString()), err.Error())
	}

	settingsContent := make([]byte, 0)
	//todo check if settings file is not empty in .tf
	if plan.SettingsFile.ValueString() != "" {
		settingsContent, err = os.ReadFile(plan.SettingsFile.ValueString())
		if err != nil {
			diagnostics.AddError(fmt.Sprintf("Client error when reading settings file %s", plan.SettingsFile.ValueString()), err.Error())
		}
	}

	dvExits, err := r.SolutionClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return nil
	}

	solution, err := r.SolutionClient.CreateSolution(ctx, plan.EnvironmentId.ValueString(), s, solutionContent, settingsContent)
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when importing solution %s", plan.SolutionFile), err.Error())
	}
	return solution
}

func (r *SolutionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	var plan *SolutionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *SolutionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	solution := r.importSolution(ctx, plan, &resp.Diagnostics)

	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solution.Name))

	plan.SolutionName = types.StringValue(solution.Name)
	plan.SolutionVersion = types.StringValue(solution.Version)
	plan.IsManaged = types.BoolValue(solution.IsManaged)
	plan.DisplayName = types.StringValue(solution.DisplayName)

	plan.SettingsFileChecksum = types.StringUnknown()
	if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.SettingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
		} else {
			plan.SettingsFileChecksum = types.StringValue(value)
		}
	}

	plan.SolutionFileChecksum = types.StringUnknown()
	if !plan.SolutionFile.IsNull() && !plan.SolutionFile.IsUnknown() {
		value, err := powerplatform_helpers.CalculateMd5(plan.SolutionFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for solution file", err.Error())
		} else {
			plan.SolutionFileChecksum = types.StringValue(value)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *SolutionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *SolutionResourceModel
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !state.EnvironmentId.IsNull() && !state.SolutionName.IsNull() {
		err := r.SolutionClient.DeleteSolution(ctx, state.EnvironmentId.ValueString(), state.SolutionName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *SolutionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
