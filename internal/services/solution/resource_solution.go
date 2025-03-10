// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewSolutionResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solution",
		},
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for importing exporting solutions in Power Platform environments.  This is the equivalent of the [`pac solution import`](https://learn.microsoft.com/power-platform/developer/cli/reference/solution#pac-solution-import) command in the Power Platform CLI.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"solution_file_checksum": schema.StringAttribute{
				MarkdownDescription: "Checksum of the solution file",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.SyncAttributePlanModifier("solution_file"),
					modifiers.SyncAttributePlanModifier("solution_file"),
				},
			},
			"solution_file": schema.StringAttribute{
				MarkdownDescription: "Path to the solution file",
				Required:            true,
			},
			"settings_file_checksum": schema.StringAttribute{
				MarkdownDescription: "Checksum of the settings file",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.SyncAttributePlanModifier("settings_file"),
					modifiers.SyncAttributePlanModifier("settings_file"),
				},
			},
			"settings_file": schema.StringAttribute{
				MarkdownDescription: "Path to the settings file. The settings file uses the same format as pac cli. See https://learn.microsoft.com/power-platform/alm/conn-ref-env-variables-build-tools#deployment-settings-file for more details",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the solution",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the environment where the solution is imported",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the solution",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.SetStringValueToUnknownIfChecksumsChangeModifier([]string{"solution_file", "solution_file_checksum"}, []string{"settings_file", "settings_file_checksum"}),
				},
			},
			"is_managed": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the solution is managed or not",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.SetBoolValueToUnknownIfChecksumsChangeModifier([]string{"solution_file", "solution_file_checksum"}, []string{"settings_file", "settings_file_checksum"}),
				},
			},
			"solution_version": schema.StringAttribute{
				MarkdownDescription: "Version of the solution",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.SetStringValueToUnknownIfChecksumsChangeModifier([]string{"solution_file", "solution_file_checksum"}, []string{"settings_file", "settings_file_checksum"}),
				},
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.SolutionClient = NewSolutionClient(client.Api)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	solution := r.importSolution(ctx, plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.SolutionVersion = types.StringValue(solution.Version)
	plan.IsManaged = types.BoolValue(solution.IsManaged)
	plan.DisplayName = types.StringValue(solution.DisplayName)

	plan.SettingsFileChecksum = types.StringNull()
	if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
		value, err := helpers.CalculateSHA256(plan.SettingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
		} else {
			plan.SettingsFileChecksum = types.StringValue(value)
			tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
		}
	}

	plan.SolutionFileChecksum = types.StringUnknown()
	if !plan.SolutionFile.IsNull() && !plan.SolutionFile.IsUnknown() {
		value, err := helpers.CalculateSHA256(plan.SolutionFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for solution file", err.Error())
		} else {
			plan.SolutionFileChecksum = types.StringValue(value)
			tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of solution file: %s", value))
		}
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solution.Id))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *ResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	solutionId := getSolutionId(state.Id.ValueString())
	solution, err := r.SolutionClient.GetSolutionById(ctx, state.EnvironmentId.ValueString(), solutionId)
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	if solution == nil {
		state.Id = types.StringNull()
		state.SolutionVersion = types.StringNull()
		state.IsManaged = types.BoolNull()
		state.DisplayName = types.StringNull()
		state.SettingsFileChecksum = types.StringNull()
		state.SolutionFileChecksum = types.StringNull()

		tflog.Debug(ctx, fmt.Sprintf("Solution %s not found", solutionId))
		return
	}
	state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), solution.Id))
	state.SolutionVersion = types.StringValue(solution.Version)
	state.IsManaged = types.BoolValue(solution.IsManaged)
	state.DisplayName = types.StringValue(solution.DisplayName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) importSolution(ctx context.Context, plan *ResourceModel, diagnostics *diag.Diagnostics) *SolutionDto {
	solutionContent, err := os.ReadFile(plan.SolutionFile.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", plan.SolutionFile.ValueString()), err.Error())
	}

	cwd, _ := os.Getwd()
	tflog.Debug(ctx, fmt.Sprintf("Current working directory: %s", cwd))

	settingsContent := make([]byte, 0)
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

	solution, err := r.SolutionClient.CreateSolution(ctx, plan.EnvironmentId.ValueString(), solutionContent, settingsContent)
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when importing solution %s", plan.SolutionFile), err.Error())
	}
	return solution
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solution := r.importSolution(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solution.Id))

	plan.SolutionVersion = types.StringValue(solution.Version)
	plan.IsManaged = types.BoolValue(solution.IsManaged)
	plan.DisplayName = types.StringValue(solution.DisplayName)

	plan.SettingsFileChecksum = types.StringNull()
	if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
		value, err := helpers.CalculateSHA256(plan.SettingsFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
		} else {
			plan.SettingsFileChecksum = types.StringValue(value)
		}
	}

	plan.SolutionFileChecksum = types.StringUnknown()
	if !plan.SolutionFile.IsNull() && !plan.SolutionFile.IsUnknown() {
		value, err := helpers.CalculateSHA256(plan.SolutionFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning("Issue when calculating checksum for solution file", err.Error())
		} else {
			plan.SolutionFileChecksum = types.StringValue(value)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.EnvironmentId.IsNull() && !state.Id.IsNull() {
		solutionId := getSolutionId(state.Id.ValueString())
		err := r.SolutionClient.DeleteSolution(ctx, state.EnvironmentId.ValueString(), solutionId)

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getSolutionId(id string) string {
	split := strings.Split(id, "_")
	return split[len(split)-1]
}
