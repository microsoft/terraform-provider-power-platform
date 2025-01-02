// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewCopilotStudioApplicationInsightsResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connection",
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
		MarkdownDescription: "Manages the [Application Insights configuration for a Copilot](https://learn.microsoft.com/en-us/microsoft-copilot-studio/advanced-bot-framework-composer-capture-telemetry?tabs=webApp).",
		Attributes: map[string]schema.Attribute{
			"botId": schema.StringAttribute{
				MarkdownDescription: "The ID of the Copilot for which the Application Insights configuration is to be managed.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment ID for the Power Platform environment where the Copilot exists",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application_insights_connection_string": schema.StringAttribute{
				MarkdownDescription: "The connection string for the target Application Insights resource in Azure. If needed, follow [these instructions](https://learn.microsoft.com/en-us/azure/azure-monitor/app/connection-strings?tabs=net#find-your-connection-string) to find your connection string.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_sensitive_information": schema.BoolAttribute{
				MarkdownDescription: "Whether to log sensitive properties such as user ID, name, and text.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_activities": schema.BoolAttribute{
				MarkdownDescription: "Whether to log details of incoming/outgoing messages and events.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_actions": schema.BoolAttribute{
				MarkdownDescription: "Whether to log an event each time a node within a topic is executed.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"network_isolation": schema.StringAttribute{
				MarkdownDescription: "The network isolation setting for the Application Insights resource connection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"errors": schema.SetAttribute{
				MarkdownDescription: "Any errors arising from the Application Insights configuration.",
				Computed:            true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(ctx, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
	}

	// You can't really create a config, so treat a create as an update
	appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating/updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	newState, err := convertAppInsightsConfigModelFromDto(appInsightsConfigDto)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting Copilot Studio Application Insights configuration to source model", err.Error())
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Updated application insights for bot '%s' in environment '%s'", plan.BotId.ValueString(), plan.EnvironmentId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.getCopilotStudioAppInsightsConfiguration(ctx, state.EnvironmentId.String(), state.BotId.String())
	if err != nil {
		if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	newState, err := convertAppInsightsConfigModelFromDto(appInsightsConfigDto)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting Copilot Studio Application Insights configuration to source model", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ: Copilot Studio bot ID '%s' in environment ID '%s'", state.BotId.ValueString(), state.EnvironmentId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan *ResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(ctx, *plan)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
	}

	// You can't really create a config, so treat a create as an update
	appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating/updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	newState, err := convertAppInsightsConfigModelFromDto(appInsightsConfigDto)
	if err != nil {
		resp.Diagnostics.AddError("Error when converting Copilot Studio Application Insights configuration to source model", err.Error())
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("Updated application insights for bot '%s' in environment '%s'", plan.BotId.ValueString(), plan.EnvironmentId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Log that the delete operation is not supported
	tflog.Warn(ctx, "Delete operation is not supported for the Copilot Studio Application Insights configuration resource")
	resp.State.RemoveResource(ctx)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
