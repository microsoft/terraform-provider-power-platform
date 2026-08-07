// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithModifyPlan = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

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
		MarkdownDescription: "Resource for deploying managed Power Platform solutions using solution identity (`unique_name` + `version`) instead of package checksums or delivery locations. Changing only a local path or signed URL does not replay an unchanged version. Initial creation installs the managed package or adopts an exact already-installed managed version only when no target-specific component parameters must be applied. Exact or same-version target bindings use an ordinary managed re-import; a higher version uses Dataverse stage-and-upgrade so omitted components are removed, and lower versions are rejected before import. Import state uses `{environment_id}_{solution_id}` and the first configured apply adopts the package source without replaying only when no target bindings require reconciliation. The resource verifies managed package identity, connection bindings, environment-variable satisfiability, and solution dependencies before import. Import-start requests are never retried after an ambiguous response.",
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
						MarkdownDescription: "Remote URL to the managed solution zip package. Marked sensitive because artifact URLs commonly embed credentials such as SAS tokens.",
						Optional:            true,
						Sensitive:           true,
					},
				},
			},
			"connection_references": schema.MapAttribute{
				MarkdownDescription: "Map of connection reference logical name to environment connection id. Every connection reference declared by the package must be bound here.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"environment_variables": schema.MapAttribute{
				MarkdownDescription: "Map of environment variable schema name to target-specific value supplied as a managed import component parameter. Configured values satisfy package definitions that have no default or existing target value.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"skip_product_update_dependencies": schema.BoolAttribute{
				MarkdownDescription: "Skip Dataverse product-update dependency processing during managed import. The package graph remains responsible for satisfying declared solution dependencies.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"publish_all_customizations": schema.BoolAttribute{
				MarkdownDescription: "Publish all Dataverse customizations after the managed import completes. This is opt-in because managed solution import already publishes its own solution components.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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

	if plan.EnvironmentId.ValueString() == state.EnvironmentId.ValueString() &&
		plan.UniqueName.ValueString() == state.UniqueName.ValueString() &&
		plan.Version.ValueString() == state.Version.ValueString() &&
		plan.ConnectionReferences.Equal(state.ConnectionReferences) &&
		plan.EnvironmentVariables.Equal(state.EnvironmentVariables) &&
		plan.SkipProductUpdateDependencies.Equal(state.SkipProductUpdateDependencies) &&
		plan.PublishAllCustomizations.Equal(state.PublishAllCustomizations) &&
		sourcesAreEquivalent(plan.Source, state.Source) {
		plan.Source = cloneSourceModel(state.Source)
		plan.Id = state.Id
		plan.SolutionId = state.SolutionId
		plan.DisplayName = state.DisplayName
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

func sourcesAreEquivalent(plan *SourceModel, state *SourceModel) bool {
	if plan == nil || state == nil {
		return plan == nil && state == nil
	}

	// The managed solution's unique name + version is the deployment identity.
	// Source is apply-time delivery plumbing: local graph workspaces and signed
	// URLs legitimately change between runs without changing package identity.
	// Package metadata is validated against the configured identity before any
	// import, so a locator-only change must not replay the managed deployment.
	return hasDeliverySource(plan) && hasDeliverySource(state)
}

func hasDeliverySource(source *SourceModel) bool {
	return normalizedSourceString(source.Path) != "" || normalizedSourceURL(source.URL) != ""
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

func reconcileSolutionVersion(declared types.String, remote string) types.String {
	if !declared.IsNull() && !declared.IsUnknown() {
		declaredVersion, declaredErr := normalizeSolutionVersion(declared.ValueString())
		remoteVersion, remoteErr := normalizeSolutionVersion(remote)
		if declaredErr == nil && remoteErr == nil && declaredVersion == remoteVersion {
			return declared
		}
	}

	return types.StringValue(normalizeSolutionVersionOrOriginal(remote))
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	solutionState := r.applyManagedSolution(ctx, &plan, config.Source, true, &resp.Diagnostics)
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

	var solutionState *solution.SolutionDto
	var err error
	if !state.UniqueName.IsNull() && !state.UniqueName.IsUnknown() && strings.TrimSpace(state.UniqueName.ValueString()) != "" {
		solutionState, err = r.Client.GetSolutionUniqueName(ctx, state.EnvironmentId.ValueString(), state.UniqueName.ValueString())
	} else if !state.SolutionId.IsNull() && !state.SolutionId.IsUnknown() && strings.TrimSpace(state.SolutionId.ValueString()) != "" {
		solutionState, err = r.Client.GetSolutionById(ctx, state.EnvironmentId.ValueString(), state.SolutionId.ValueString())
	} else {
		resp.Diagnostics.AddError("Managed solution state is incomplete", "Neither unique_name nor solution_id is available to refresh the managed solution.")
		return
	}
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}
	if !solutionState.IsManaged {
		resp.Diagnostics.AddError("Managed solution required", fmt.Sprintf("Solution %q in environment %q is unmanaged and cannot be adopted by %s.", solutionState.Name, state.EnvironmentId.ValueString(), r.FullTypeName()))
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), solutionState.Id))
	state.SolutionId = types.StringValue(solutionState.Id)
	state.DisplayName = types.StringValue(solutionState.DisplayName)
	state.UniqueName = types.StringValue(solutionState.Name)
	state.Version = reconcileSolutionVersion(state.Version, solutionState.Version)

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

	var config ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Imported state cannot recover the package delivery source from Dataverse. The first
	// configured apply after import therefore adopts an exact installed managed version and
	// records the configured source without replaying the import. A normal update does not adopt:
	// a higher package version uses stage-and-upgrade, while a same-version target-binding change
	// uses ImportSolutionAsync because Dataverse upgrades require a higher version.
	allowExactAdoption := state.Source == nil ||
		(normalizedSourceString(state.Source.Path) == "" && normalizedSourceString(state.Source.URL) == "")
	solutionState := r.applyManagedSolution(ctx, &plan, config.Source, allowExactAdoption, &resp.Diagnostics)
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

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	parts := strings.SplitN(req.ID, "_", 2)
	if len(parts) != 2 || !isGuid(parts[0]) || !isGuid(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in format 'environment_id_solution_id', got %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("solution_id"), parts[1])...)
}

func (r *Resource) applyManagedSolution(ctx context.Context, plan *ResourceModel, deliverySource *SourceModel, allowExactAdoption bool, diagnostics *diag.Diagnostics) *ResourceModel {
	dataverseExists, err := r.Client.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
		return nil
	}
	if !dataverseExists {
		diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return nil
	}

	sourcePath, cleanup, err := resolveSourceToPath(ctx, deliverySource)
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
	configuredVariables, variableDiags := expandStringMap(plan.EnvironmentVariables)
	diagnostics.Append(variableDiags...)
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
	if err := validateEnvironmentVariables(pkg.EnvironmentVariables, configuredVariables, existingValues); err != nil {
		diagnostics.AddError("Managed solution environment variable validation failed", err.Error())
		return nil
	}

	installedSolutions, err := r.Client.GetInstalledSolutions(ctx, plan.EnvironmentId.ValueString())
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
	componentParameters = append(componentParameters, buildEnvironmentVariableParameters(configuredVariables)...)
	if len(componentParameters) == 0 {
		componentParameters = nil
	}

	content, err := readSourceContent(sourcePath)
	if err != nil {
		diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", sourcePath), err.Error())
		return nil
	}

	installed, installedErr := r.Client.GetSolutionUniqueName(ctx, plan.EnvironmentId.ValueString(), plan.UniqueName.ValueString())
	if installedErr != nil && !errors.Is(installedErr, customerrors.ErrObjectNotFound) {
		diagnostics.AddError("Unable to inspect existing managed solution", installedErr.Error())
		return nil
	}
	if installedErr == nil && !installed.IsManaged {
		diagnostics.AddError("Managed solution required", fmt.Sprintf("Solution %q is already installed as unmanaged and cannot be adopted or upgraded as managed.", installed.Name))
		return nil
	}

	operation := importOperationInstall
	if installedErr == nil {
		adoptExact, decidedOperation, decisionErr := managedSolutionImportDecision(
			installed.Version,
			normalizedConfiguredVersion,
			allowExactAdoption,
			len(componentParameters) > 0)
		if decisionErr != nil {
			diagnostics.AddError("Managed solution version regression", decisionErr.Error())
			return nil
		}
		if adoptExact {
			result := stateFromSolution(plan, installed)
			if plan.PublishAllCustomizations.ValueBool() {
				if err := r.Client.PublishAllCustomizations(ctx, plan.EnvironmentId.ValueString()); err != nil {
					diagnostics.AddError("Unable to publish Dataverse customizations after managed solution adoption", err.Error())
					return nil
				}
			}
			return result
		}
		operation = decidedOperation
	}

	solutionState, err := r.Client.ApplyManagedSolution(
		ctx,
		plan.EnvironmentId.ValueString(),
		content,
		componentParameters,
		operation,
		plan.SkipProductUpdateDependencies.ValueBool())
	if err != nil {
		diagnostics.AddError("Unable to apply managed solution", err.Error())
		return nil
	}
	if plan.PublishAllCustomizations.ValueBool() {
		if err := r.Client.PublishAllCustomizations(ctx, plan.EnvironmentId.ValueString()); err != nil {
			diagnostics.AddError("Unable to publish Dataverse customizations after managed solution import", err.Error())
			return nil
		}
	}

	return stateFromSolution(plan, solutionState)
}

// managedSolutionImportDecision keeps package lifecycle and target binding lifecycle distinct.
// Exact installed versions may be adopted only during create/import-state convergence when there
// are no target-specific component parameters to prove or apply. A normal update, or an exact-version
// create carrying connection/environment bindings, uses ImportSolutionAsync so target wiring is made
// authoritative without pretending that immutable package content changed. A higher version uses
// StageAndUpgradeAsync so omitted managed components are removed.
func managedSolutionImportDecision(installedVersion, configuredVersion string, allowExactAdoption, hasComponentParameters bool) (adoptExact bool, operation importOperation, err error) {
	installed, err := normalizedSolutionVersionParts(installedVersion)
	if err != nil {
		return false, importOperationInstall, fmt.Errorf("installed managed solution version %q is invalid: %w", installedVersion, err)
	}
	configured, err := normalizedSolutionVersionParts(configuredVersion)
	if err != nil {
		return false, importOperationInstall, fmt.Errorf("configured managed solution version %q is invalid: %w", configuredVersion, err)
	}

	comparison := 0
	for index := range installed {
		if configured[index] > installed[index] {
			comparison = 1
			break
		}
		if configured[index] < installed[index] {
			comparison = -1
			break
		}
	}

	switch {
	case comparison < 0:
		return false, importOperationInstall, fmt.Errorf(
			"configured version %q is lower than installed version %q; managed solution deployment must move forward",
			configuredVersion,
			installedVersion)
	case comparison > 0:
		return false, importOperationStageAndUpgrade, nil
	default:
		return allowExactAdoption && !hasComponentParameters, importOperationInstall, nil
	}
}

func normalizedSolutionVersionParts(raw string) ([4]int, error) {
	normalized, err := normalizeSolutionVersion(raw)
	if err != nil {
		return [4]int{}, err
	}

	parts := [4]int{}
	for index, segment := range strings.Split(normalized, ".") {
		parts[index], _ = strconv.Atoi(segment) // normalizeSolutionVersion already validated every segment.
	}
	return parts, nil
}

func stateFromSolution(plan *ResourceModel, solutionState *solution.SolutionDto) *ResourceModel {
	result := *plan
	result.Id = types.StringValue(fmt.Sprintf("%s_%s", plan.EnvironmentId.ValueString(), solutionState.Id))
	result.SolutionId = types.StringValue(solutionState.Id)
	result.DisplayName = types.StringValue(solutionState.DisplayName)
	result.UniqueName = types.StringValue(solutionState.Name)
	result.Version = reconcileSolutionVersion(plan.Version, solutionState.Version)
	return &result
}

func isGuid(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
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
