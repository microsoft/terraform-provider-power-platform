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
	"github.com/microsoft/terraform-provider-power-platform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	modifier "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/modifiers"
)

func NewDataverseWebApiResource() resource.Resource {
	return &DataverseWebApiResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_rest",
	}
}

type DataverseWebApiResource struct {
	DataRecordClient WebApiClient
	ProviderTypeName string
	TypeName         string
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
	ExpectedHttpStatus []int64                                  `tfsdk:"expected_http_status"`
}

type DataverseWebApiOperationHeaderResource struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *DataverseWebApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataverseWebApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
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
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id (guid)",
				Description:         "Unique id (guid)",
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
							modifier.ForceStringValueUnknownModifier(),
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
			"expected_http_status": schema.ListAttribute{
				ElementType:         types.Int64Type,
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
	r.DataRecordClient = NewWebApiClient(clientApi)

}

func (r *DataverseWebApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DataverseWebApiResourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state := DataverseWebApiResourceModel{}
	state.Create = plan.Create
	state.Update = plan.Update
	state.Destroy = plan.Destroy
	state.Read = plan.Read
	state.Output = plan.Output

	timeout, diags := plan.Timeouts.Create(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DataverseWebApiResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newState := DataverseWebApiResourceModel{}
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
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *DataverseWebApiResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.TypeName))
}

func (r *DataverseWebApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *DataverseWebApiResourceModel
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if state.Destroy != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Destroy)
		if err != nil {
			resp.Diagnostics.AddError("Error executing destroy operation", err.Error())
			return
		}
		state.Output = bodyWrapped
	}
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.TypeName))

}

func (r *DataverseWebApiResource) NullOutputValue() basetypes.ObjectValue {
	return types.ObjectValueMust(map[string]attr.Type{
		"body": types.StringType,
	}, map[string]attr.Value{
		"body": types.StringValue("{}"),
	})
}
