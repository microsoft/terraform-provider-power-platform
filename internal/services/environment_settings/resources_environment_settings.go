// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &EnvironmentSettingsResource{}
var _ resource.ResourceWithImportState = &EnvironmentSettingsResource{}

func NewEnvironmentSettingsResource() resource.Resource {
	return &EnvironmentSettingsResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_settings",
		},
	}
}

const SERVICE_TAGS_NAMES = "ApiManagement,AppConfiguration,AppService,ActionGroup,AppServiceManagement,ApplicationInsightsAvailability,AutonomousDevelopmentPlatform,AzureActiveDirectory,AzureAdvancedThreatProtection,AzureArcInfrastructure,AzureAttestation,AzureBackup,AzureBotService,AzureCognitiveSearch,AzureConnectors,AzureContainerRegistry,AzureCosmosDB,AzureDataExplorerManagement,AzureDataLake,AzureDatabricks,AzureDevOps,AzureDevSpaces,AzureDeviceUpdate,AzureDigitalTwins,AzureEventGrid,AzureHealthcareAPIs,AzureInformationProtection,AzureIoTHub,AzureKeyVault,AzureLoadTestingInstanceManagement,AzureMachineLearning,AzureMachineLearningInference,AzureManagedGrafana,AzureMonitorForSAP,AzureMonitor,AzureOpenDatasets,AzurePortal,AzureRemoteRendering,AzureResourceManager,AzureSecurityCenter,AzureSentinel,AzureSignalR,AzureSiteRecovery,AzureSphere,AzureSpringCloud,AzureStack,AzureTrafficManager,AzureUpdateDelivery,AzureWebPubSub,BatchNodeManagement,ChaosStudio,CognitiveServicesFrontend,CognitiveServicesManagement,ContainerAppsManagement,DataFactory,Dynamics365BusinessCentral,Dynamics365ForMarketingEmail,Dynamics365FraudProtection,EOPExternalPublishedIPs,EventHub,GatewayManager,Grafana,GuestAndHybridManagement,HDInsight,KustoAnalytics,LogicApps,M365ManagementActivityApi,M365ManagementActivityApiWebhook,Marketplace,MicrosoftAzureFluidRelay,MicrosoftCloudAppSecurity,MicrosoftContainerRegistry,MicrosoftDefenderForEndpoint,MicrosoftPurviewPolicyDistribution,OneDsCollector,PowerBI,PowerPlatformPlex,PowerQueryOnline,SCCservice,Scuba,SecurityCopilot,SerialConsole,ServiceBus,ServiceFabric,Sql,SqlManagement,Storage,StorageMover,StorageSyncService,VideoIndexer,WindowsAdminCenter,WindowsVirtualDesktop"

func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Power Platform Settings for a given environment. They control various aspects of Power Platform features and behaviors, See [Environment Settings Overview](https://learn.microsoft.com/power-platform/admin/admin-settings) for more details.  While this resource provides a limited set of settings, many of the settings in an environment are stored as Dataverse records and can be managed using `powerplatform_data_record` resource.  See the [data record resource documentation](./data_record) for examples of how to manage more environment settings like business units, roles, and more.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment Id",
				Required:            true,
			},
			"audit_and_logs": schema.SingleNestedAttribute{
				MarkdownDescription: "Audit and Logs",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"plugin_trace_log_setting": schema.StringAttribute{
						MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All. See [Plugin Trace Log Settings Overview](https://learn.microsoft.com/power-apps/developer/data-platform/logging-tracing) for more details.",
						Optional:            true, Computed: true,
						Default: stringdefault.StaticString("Off"),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("Off", "Exception", "All"),
						},
					},
					"audit_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Audit Settings. See [Audit Settings Overview](https://learn.microsoft.com/power-platform/admin/system-settings-dialog-box-auditing-tab) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"is_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is audit enabled",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"is_user_access_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is user access audit enabled",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"is_read_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is read audit enabled",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"log_retention_period_in_days": schema.Int32Attribute{
								MarkdownDescription: "Retain these logs for a value between 31 days and 24855 days, value of '-1' means logs will be retained forever. See [Start/stop auditing for an environment and set retention policy](https://learn.microsoft.com/power-platform/admin/manage-dataverse-auditing#startstop-auditing-for-an-environment-and-set-retention-policy) You can set a retention period for how long audit logs are kept in an environment. Under Retain these logs for, choose the period of time you wish to retain the logs.",
								Optional:            true, Computed: true,
								Default: int32default.StaticInt32(-1),
								Validators: []validator.Int32{
									int32validator.Any(int32validator.Between(31, 24855), int32validator.OneOf(-1)),
								},
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"email": schema.SingleNestedAttribute{
				MarkdownDescription: "Email",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"email_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Email Settings. See [Email Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-email) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"max_upload_file_size_in_bytes": schema.Int64Attribute{
								MarkdownDescription: "Maximum file size that can be uploaded to the environment",
								Optional:            true, Computed: true,
								Default: int64default.StaticInt64(5120000),
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"product": schema.SingleNestedAttribute{
				MarkdownDescription: "Product",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"behavior_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Behavior Settings.See [Behavior Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-behavior) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
								MarkdownDescription: "Show dashboard cards in expanded state",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"features": schema.SingleNestedAttribute{
						MarkdownDescription: "Features. See [Features Overview](https://learn.microsoft.com/power-platform/admin/settings-features) for more details.",
						Optional:            true, Computed: true,
						Attributes: map[string]schema.Attribute{
							"power_apps_component_framework_for_canvas_apps": schema.BoolAttribute{
								MarkdownDescription: "Power Apps component framework for canvas apps",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_powerapps_maker_bot": schema.BoolAttribute{
								MarkdownDescription: "Enable new AI-powered Copilot features for people who make apps. [Learn more](https://go.microsoft.com/fwlink/?linkid=2223555)",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_access_to_session_transcripts_for_copilot_studio": schema.BoolAttribute{
								MarkdownDescription: "Allow agent access owners and editors to see session transcripts from conversations interactions in their agents.",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_transcript_recording_for_copilot_studio": schema.BoolAttribute{
								MarkdownDescription: "Allow conversation transcripts and their associated metadata to be saved in Dataverse (required for enhanced reporting).",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_copilot_studio_share_data_with_viva_insights": schema.BoolAttribute{
								MarkdownDescription: "Allow Copilot Studio to share data with Viva Insights.",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_copilot_studio_cross_geo_share_data_with_viva_insights": schema.BoolAttribute{
								MarkdownDescription: "Allow cross-geo sharing of aggregated analytics data if your tenant preferred data location for Viva Insights is different from the location of your Copilot Studio environment.",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_preview_and_experimental_ai_models": schema.BoolAttribute{
								MarkdownDescription: "Let people use AI models that are experimental or in preview to make agents, prompts, apps, flow and more in Copilot Studio. [Learn more](https://go.microsoft.com/fwlink/?linkid=2331400)",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ai_prompts": schema.BoolAttribute{
								MarkdownDescription: "Enable the AI prompts feature in Power Platform and Copilot Studio. [Learn more](https://go.microsoft.com/fwlink/?linkid=2283738)",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"security": schema.SingleNestedAttribute{
						MarkdownDescription: "Security. See [Security Overview](https://learn.microsoft.com/en-us/power-platform/admin/settings-privacy-security) for more details.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enable_ip_based_cookie_binding": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based cookie binding",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ip_based_firewall_rule": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based firewall rule",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"allowed_ip_range_for_firewall": schema.SetAttribute{
								MarkdownDescription: "Allowed IP range for firewall",
								Optional:            true, Computed: true,
								Default:     setdefault.StaticValue(types.SetNull(types.StringType)),
								ElementType: types.StringType,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
							"allowed_service_tags_for_firewall": schema.SetAttribute{
								MarkdownDescription: "Allowed service tags for firewall",
								Optional:            true, Computed: true,
								Default:     setdefault.StaticValue(types.SetNull(types.StringType)),
								ElementType: types.StringType,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(stringvalidator.OneOf(append([]string{""}, strings.Split(SERVICE_TAGS_NAMES, ",")...)...)),
								},
							},
							"allow_application_user_access": schema.BoolAttribute{
								MarkdownDescription: "Allow application user access",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"allow_microsoft_trusted_service_tags": schema.BoolAttribute{
								MarkdownDescription: "Allow Microsoft trusted service tags",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"enable_ip_based_firewall_rule_in_audit_mode": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based firewall rule in audit mode",
								Optional:            true, Computed: true,
								Default: booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"reverse_proxy_ip_addresses": schema.SetAttribute{
								MarkdownDescription: "Reverse proxy IP addresses",
								Optional:            true, Computed: true,
								Default:     setdefault.StaticValue(types.SetNull(types.StringType)),
								ElementType: types.StringType,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EnvironmentSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	client := providerClient.Api
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Provider client Api is nil. Please report this issue to the provider developers.",
		)
		return
	}

	r.EnvironmentSettingClient = newEnvironmentSettingsClient(providerClient.Api)
}

func (r *EnvironmentSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settingsToUpdate, err := convertFromEnvironmentSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings model", err.Error())
		return
	}

	dvExits, err := r.EnvironmentSettingClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
		return
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return
	}

	envSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), *settingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment settings", fmt.Sprintf("Error creating environment settings: %s", err.Error()),
		)
		return
	}

	state, err := convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings", err.Error())
		return
	}
	state.Id = plan.EnvironmentId
	state.EnvironmentId = plan.EnvironmentId

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", state.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	envSettings, err := r.EnvironmentSettingClient.GetEnvironmentSettings(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
		return
	}

	newState, err := convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, state.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings", err.Error())
		return
	}
	newState.Id = state.EnvironmentId
	newState.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *EnvironmentSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentSettingsResourceModel
	var state EnvironmentSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	envSettingsToUpdate, err := convertFromEnvironmentSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings model", err.Error())
		return
	}

	environmentSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), *envSettingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment settings", fmt.Sprintf("Error updating environment settings: %s", err.Error()),
		)
		return
	}

	plan, err = convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](environmentSettings, plan.Timeouts)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings", err.Error())
		return
	}
	plan.Id = state.EnvironmentId
	plan.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	// Do nothing on purpose
}

func (r *EnvironmentSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_id"), req, resp)
}
