// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package rest

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/modifiers"
)

func NewDataverseWebApiResource() resource.Resource {
	return &DataverseWebApiResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "rest",
		},
	}
}

type DataverseWebApiResource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataverseWebApiResourceModel struct {
	Timeouts timeouts.Value            `tfsdk:"timeouts"`
	Id       types.String              `tfsdk:"id"`
	Create   *DataverseWebApiOperation `tfsdk:"create"`
	Update   *DataverseWebApiOperation `tfsdk:"update"`
	Destroy  *DataverseWebApiOperation `tfsdk:"destroy"`
	Read     *DataverseWebApiOperation `tfsdk:"read"`
	Output   types.Object              `tfsdk:"output"`
}

type DataverseWebApiOperation struct {
	Scope              types.String                             `tfsdk:"scope"`
	Method             types.String                             `tfsdk:"method"`
	Url                types.String                             `tfsdk:"url"`
	Body               types.String                             `tfsdk:"body"`
	Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
	ExpectedHttpStatus []int                                    `tfsdk:"expected_http_status"`
}

type DataverseWebApiOperationHeaderResource struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *DataverseWebApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *DataverseWebApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		DeprecationMessage: "This resource is deprecated. Please consider using specific resources for your use case, or opening a GitHub issue requesting a new resource for your use case.",
		MarkdownDescription: `Resource to execute web api requests. There are four distinct operations, that you can define independently. The HTTP response' body of the operation, that was called as last, will be returned in 'output.body' \n\n:
		* Create: will be called once during the lifecycle of the resource (first 'terraform apply')
		* Read: terraform will call this operation every time during 'plan' and 'apply' to get the current state of the resource
		* Update: will be called every time during 'terraform apply' if the resource has changed (change done by the user or different values returned by the 'read' operation than those in the current state)
		* Destroy: will be called once during the lifecycle of the resource (last 'terraform destroy')
		\n\nYou don't have to define all the operations but there are some things to consider:
		* lack of 'create' operation will result in no reasource being created. If you only need to read values consider using datasource 'powerplatform_rest_query' instead
		* lack of 'read' operation will result in no resource changes being tracked. That means that the 'update' operation will never be called
		* lack of destroy will cause that the resource to not be deleted during 'terraform destroy'`,
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id (guid)",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create":  r.buildOperationSchema("Create operation"),
			"update":  r.buildOperationSchema("Update operation"),
			"destroy": r.buildOperationSchema("Destroy operation"),
			"read":    r.buildOperationSchema("Read operation"),
			"output": schema.SingleNestedAttribute{
				MarkdownDescription: "Response after executing the web api request",
				Computed:            true,
				Optional:            true,
				Required:            false,
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						MarkdownDescription: "Response body after executing the web api request",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.ForceStringValueUnknownModifier(),
						},
					},
				},
			},
		},
	}
}

func (r *DataverseWebApiResource) buildOperationSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: description,
		Required:            false,
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{
				MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
				Required:            true,
				Optional:            false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expected_http_status": schema.SetAttribute{
				ElementType:         types.Int32Type,
				MarkdownDescription: "Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.",
				Required:            false,
				Optional:            true,
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "HTTP method",
				Required:            true,
				Optional:            false,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Absolute url of the api call",
				Required:            true,
				Optional:            false,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Body of the request",
				Required:            false,
				Optional:            true,
			},
			"headers": schema.ListNestedAttribute{
				MarkdownDescription: "Headers of the request",
				Required:            false,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Header name",
							Required:            true,
							Optional:            false,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Header value",
							Required:            true,
							Optional:            false,
						},
					},
				},
			},
		},
	}
}

func (r *DataverseWebApiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.DataRecordClient = newWebApiClient(client.Api)
}

func (r *DataverseWebApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan DataverseWebApiResourceModel
	resp.State.Get(ctx, &plan)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := DataverseWebApiResourceModel{}
	state.Timeouts = plan.Timeouts
	state.Create = plan.Create
	state.Update = plan.Update
	state.Destroy = plan.Destroy
	state.Read = plan.Read
	state.Output = plan.Output

	state.Id = types.StringValue(strconv.Itoa(int(time.Now().UnixMilli())))
	if plan.Create != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Create)
		if err != nil {
			resp.Diagnostics.AddError("Error executing create operation", err.Error())
			return
		}
		state.Output = bodyWrapped

		if plan.Read != nil {
			bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Read)
			if err != nil {
				resp.Diagnostics.AddError("Error executing read operation", err.Error())
				return
			}
			state.Output = bodyWrapped
		}
	} else {
		state.Output = r.NullOutputValue()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DataverseWebApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState := DataverseWebApiResourceModel{}
	newState.Timeouts = state.Timeouts
	newState.Id = state.Id
	newState.Create = state.Create
	newState.Update = state.Update
	newState.Destroy = state.Destroy
	newState.Read = state.Read
	newState.Output = state.Output

	if state.Read != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Read)
		if err != nil {
			resp.Diagnostics.AddError("Error executing read operation", err.Error())
			return
		}

		if state.Output.String() != bodyWrapped.String() {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("true"))
		} else {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("false"))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *DataverseWebApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Update != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Update)
		if err != nil {
			resp.Diagnostics.AddError("Error executing update operation", err.Error())
			return
		}
		plan.Output = bodyWrapped

		if plan.Read != nil {
			bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Read)
			if err != nil {
				resp.Diagnostics.AddError("Error executing read operation", err.Error())
				return
			}
			plan.Output = bodyWrapped
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataverseWebApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Destroy != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Destroy)
		if err != nil {
			resp.Diagnostics.AddError("Error executing destroy operation", err.Error())
			return
		}
		state.Output = bodyWrapped
	}
}

func (r *DataverseWebApiResource) NullOutputValue() basetypes.ObjectValue {
	return types.ObjectValueMust(map[string]attr.Type{
		"body": types.StringType,
	}, map[string]attr.Value{
		"body": types.StringValue("{}"),
	})
}
