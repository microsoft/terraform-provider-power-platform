// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithModifyPlan = &Resource{}

func NewManagedSolutionResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "managed_solution",
		},
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for deploying managed Power Platform solutions using solution identity (`unique_name` + `version`) instead of package checksums. The resource verifies that the package is managed, validates package identity at apply time, requires explicit connection reference bindings, checks that referenced environment variables are already satisfiable in the target environment, and validates solution dependencies before import.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the managed solution deployment in format `{environment_id}_{solution_id}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the environment where the managed solution is imported.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"unique_name": schema.StringAttribute{
				MarkdownDescription: "Unique name of the managed solution that must match the package metadata exactly.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Desired solution version. Terraform uses this as the deployment trigger instead of package checksums.",
				Required:            true,
			},
			"source": schema.SingleNestedAttribute{
				MarkdownDescription: "Artifact source for the managed solution package. This is an apply-time delivery input. Set exactly one of `path` or `url`.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						MarkdownDescription: "Local filesystem path to the managed solution zip package.",
						Optional:            true,
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "Remote URL to the managed solution zip package.",
						Optional:            true,
					},
				},
			},
			"connection_references": schema.MapAttribute{
				MarkdownDescription: "Map of connection reference logical name to environment connection id. Every connection reference declared by the package must be bound here.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the installed solution.",
				Computed:            true,
			},
			"solution_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse solution id of the installed solution.",
				Computed:            true,
			},
		},
	}
}

func (r *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("source").AtName("path"),
			path.MatchRoot("source").AtName("url"),
		),
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
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

	r.Client = NewManagedSolutionClient(client.Api)
}

func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planRefs, diags := expandStringMap(plan.ConnectionReferences)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefs, diags := expandStringMap(state.ConnectionReferences)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.EnvironmentId.ValueString() == state.EnvironmentId.ValueString() &&
		plan.UniqueName.ValueString() == state.UniqueName.ValueString() &&
		plan.Version.ValueString() == state.Version.ValueString() &&
		mapsEqual(planRefs, stateRefs) &&
		sourcesAreEquivalent(plan.Source, state.Source) {
		plan.Source = cloneSourceModel(state.Source)
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

func sourcesAreEquivalent(plan *SourceModel, state *SourceModel) bool {
	if plan == nil || state == nil {
		return plan == nil && state == nil
	}

	planPath := normalizedSourceString(plan.Path)
	statePath := normalizedSourceString(state.Path)
	if planPath != "" || statePath != "" {
		return planPath != "" && planPath == statePath
	}

	planURL := normalizedSourceURL(plan.URL)
	stateURL := normalizedSourceURL(state.URL)
	return planURL != "" && planURL == stateURL
}

func normalizedSourceString(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}

	return value.ValueString()
}

func normalizedSourceURL(value types.String) string {
	raw := normalizedSourceString(value)
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

func normalizeSolutionVersion(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("version must be a non-empty string")
	}

	segments := strings.Split(trimmed, ".")
	if len(segments) < 1 || len(segments) > 4 {
		return "", fmt.Errorf("version %q must contain between one and four numeric segments", raw)
	}

	normalized := [4]string{"0", "0", "0", "0"}
	for index, segment := range segments {
		value := strings.TrimSpace(segment)
		if value == "" {
			return "", fmt.Errorf("version %q contains an empty segment", raw)
		}

		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return "", fmt.Errorf("version %q contains non-numeric segment %q", raw, segment)
		}

		normalized[index] = strconv.Itoa(parsed)
	}

	return strings.Join(normalized[:], "."), nil
}

func normalizeSolutionVersionOrOriginal(raw string) string {
	normalized, err := normalizeSolutionVersion(raw)
	if err != nil {
		return raw
	}

	return normalized
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionState := r.applyManagedSolution(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, solutionState)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionState, err := r.Client.GetSolutionUniqueName(ctx, state.EnvironmentId.ValueString(), state.UniqueName.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), solutionState.Id))
	state.SolutionId = types.StringValue(solutionState.Id)
	state.DisplayName = types.StringValue(solutionState.DisplayName)
	state.UniqueName = types.StringValue(solutionState.Name)
	state.Version = types.StringValue(normalizeSolutionVersionOrOriginal(solutionState.Version))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionState := r.applyManagedSolution(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, solutionState)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.EnvironmentId.IsNull() || state.SolutionId.IsNull() {
		return
	}

	if err := r.Client.DeleteSolution(ctx, state.EnvironmentId.ValueString(), state.SolutionId.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *Resource) applyManagedSolution(ctx context.Context, plan *ResourceModel, diagnostics *diag.Diagnostics) *ResourceModel {
	dataverseExists, err := r.Client.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
		return nil
	}
	if !dataverseExists {
		diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return nil
	}

	sourcePath, cleanup, err := resolveSourceToPath(ctx, plan.Source)
	if err != nil {
		diagnostics.AddError("Unable to resolve managed solution source", err.Error())
		return nil
	}
	defer cleanup()

	pkg, err := inspectSolutionPackage(sourcePath)
	if err != nil {
		diagnostics.AddError("Unable to inspect managed solution package", err.Error())
		return nil
	}

	if !pkg.IsManaged {
		diagnostics.AddError("Managed solution package required", "The supplied solution package is not managed.")
		return nil
	}
	if pkg.UniqueName != plan.UniqueName.ValueString() {
		diagnostics.AddError("Managed solution unique name mismatch", fmt.Sprintf("Configured unique_name %q does not match package unique name %q.", plan.UniqueName.ValueString(), pkg.UniqueName))
		return nil
	}
	normalizedConfiguredVersion, err := normalizeSolutionVersion(plan.Version.ValueString())
	if err != nil {
		diagnostics.AddError("Managed solution version mismatch", fmt.Sprintf("Configured version %q is invalid: %s", plan.Version.ValueString(), err.Error()))
		return nil
	}
	normalizedPackageVersion, err := normalizeSolutionVersion(pkg.Version)
	if err != nil {
		diagnostics.AddError("Managed solution version mismatch", fmt.Sprintf("Package version %q is invalid: %s", pkg.Version, err.Error()))
		return nil
	}
	if normalizedPackageVersion != normalizedConfiguredVersion {
		diagnostics.AddError("Managed solution version mismatch", fmt.Sprintf("Configured version %q does not match package version %q.", plan.Version.ValueString(), pkg.Version))
		return nil
	}

	configuredRefs, diags := expandStringMap(plan.ConnectionReferences)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return nil
	}

	if err := validateConnectionReferences(pkg.ConnectionReferences, configuredRefs); err != nil {
		diagnostics.AddError("Managed solution connection reference validation failed", err.Error())
		return nil
	}

	existingValues, err := r.Client.GetExistingEnvironmentVariableValues(ctx, plan.EnvironmentId.ValueString(), sortedEnvironmentVariableNames(pkg.EnvironmentVariables))
	if err != nil {
		diagnostics.AddError("Unable to verify environment variable satisfiability", err.Error())
		return nil
	}
	if err := validateEnvironmentVariables(pkg.EnvironmentVariables, existingValues); err != nil {
		diagnostics.AddError("Managed solution environment variable validation failed", err.Error())
		return nil
	}

	installedSolutions, err := r.Client.GetSolutions(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		diagnostics.AddError("Unable to verify installed solution dependencies", err.Error())
		return nil
	}
	if err := validateDependencies(pkg.Dependencies, installedSolutions); err != nil {
		diagnostics.AddError("Managed solution dependency validation failed", err.Error())
		return nil
	}

	connections, err := r.Client.getConnections(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		diagnostics.AddError("Unable to resolve connection references", err.Error())
		return nil
	}

	componentParameters, err := buildConnectionParameters(pkg.ConnectionReferences, configuredRefs, connections)
	if err != nil {
		diagnostics.AddError("Unable to build managed solution connection reference bindings", err.Error())
		return nil
	}

	content, err := readSourceContent(sourcePath)
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", sourcePath), err.Error())
		return nil
	}

	solutionState, err := r.Client.CreateManagedSolution(ctx, plan.EnvironmentId.ValueString(), content, componentParameters)
	if err != nil {
		diagnostics.AddError("Unable to import managed solution", err.Error())
		return nil
	}

	result := *plan
	result.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solutionState.Id))
	result.SolutionId = types.StringValue(solutionState.Id)
	result.DisplayName = types.StringValue(solutionState.DisplayName)
	result.UniqueName = types.StringValue(solutionState.Name)
	result.Version = types.StringValue(normalizeSolutionVersionOrOriginal(solutionState.Version))

	return &result
}

func expandStringMap(value types.Map) (map[string]string, diag.Diagnostics) {
	result := make(map[string]string)
	if value.IsNull() || value.IsUnknown() {
		return result, nil
	}

	var converted map[string]string
	diags := value.ElementsAs(context.Background(), &converted, false)
	if diags.HasError() {
		return nil, diags
	}

	if converted == nil {
		return result, nil
	}

	return converted, nil
}

func sortedEnvironmentVariableNames(values map[string]packageEnvironmentVariable) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
