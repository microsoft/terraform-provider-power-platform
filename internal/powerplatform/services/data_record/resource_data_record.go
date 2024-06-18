// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func NewDataRecordResource() resource.Resource {
	return &DataRecordResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_data_record",
	}
}

type DataRecordResource struct {
	DataRecordClient DataRecordClient
	ProviderTypeName string
	TypeName         string
}

type DataRecordResourceModel struct {
	Id               types.String  `tfsdk:"id"`
	EnvironmentId    types.String  `tfsdk:"environment_id"`
	TableLogicalName types.String  `tfsdk:"table_logical_name"`
	Columns          types.Dynamic `tfsdk:"columns"`
}

func (r *DataRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *DataRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.",
		MarkdownDescription: "The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id (guid)",
				Description:         "Unique id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Description: "Id of the Dynamics 365 environment",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"table_logical_name": schema.StringAttribute{
				Description: "Logical name of the data record table",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"columns": schema.DynamicAttribute{
				Description: "Columns of the data record table",
				Required:    true,
			},
		},
	}
}

func (r *DataRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.DataRecordClient = NewDataRecordClient(clientApi)
}

func (r *DataRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state DataRecordResourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(state.Id.ValueString())
	state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
	state.TableLogicalName = types.StringValue(state.TableLogicalName.ValueString())
	state.Columns = types.DynamicValue(state.Columns)

	stateColumns := state.Columns.String()
	mapColumns, err := convertResourceModelToMap(&stateColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
		return
	}

	dr, err := r.DataRecordClient.ApplyDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	state.Id = types.StringValue(dr.Id)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", state.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newColumns, err := r.DataRecordClient.GetDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString())
	if err != nil {
		if powerplatform_helpers.Code(err) == powerplatform_helpers.ERROR_OBJECT_NOT_FOUND {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	stateColumns := state.Columns.String()
	stateRecordId := state.Id.ValueString()
	columns, err := r.convertColumnsToState(ctx, &r.DataRecordClient, state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), &stateRecordId, &stateColumns, newColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to state: %s", err.Error()), err.Error())
		return
	}
	state.Columns = *columns

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_data_record with table_name %s", r.ProviderTypeName, state.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *DataRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(plan.Id.ValueString())
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
	plan.Columns = types.DynamicValue(plan.Columns)

	planColumns := plan.Columns.String()
	mapColumns, err := convertResourceModelToMap(&planColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
		return
	}

	dr, err := r.DataRecordClient.ApplyDataRecord(ctx, state.Id.ValueString(), plan.EnvironmentId.ValueString(), plan.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.StringValue(dr.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *DataRecordResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateColumns := state.Columns.String()
	mapColumns, err := convertResourceModelToMap(&stateColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
		return
	}

	err = r.DataRecordClient.DeleteDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *DataRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertResourceModelToMap(columnsAsString *string) (mapColumns map[string]interface{}, err error) {
	if columnsAsString == nil {
		return nil, nil
	}

	jsonColumns, err := json.Marshal(columnsAsString)
	if err != nil {
		return nil, err
	}
	unquotedJsonColumns, err := strconv.Unquote(string(jsonColumns))
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(unquotedJsonColumns), &mapColumns)
	return mapColumns, nil
}

func (r *DataRecordResource) convertColumnsToState(ctx context.Context, apiClient *DataRecordClient, environmentId, tableLogicalName string, recordid, recordColumns *string, columns map[string]interface{}) (*basetypes.DynamicValue, error) {
	var objectType = map[string]attr.Type{
		"table_logical_name": types.StringType,
		"data_record_id":     types.StringType,
	}

	mapColumns, err := convertResourceModelToMap(recordColumns)
	if err != nil {
		return nil, fmt.Errorf("error converting columns to map: %s", err.Error())
	}

	attributeTypes := make(map[string]attr.Type)
	attributes := make(map[string]attr.Value)

	for key, value := range mapColumns {
		switch value.(type) {
		case bool:
			v, ok := columns[key].(bool)
			if ok {
				attributeTypes[key] = types.BoolType
				attributes[key] = types.BoolValue(v)
			}
		case int64:
			v, ok := columns[key].(int64)
			if ok {
				attributeTypes[key] = types.Int64Type
				attributes[key] = types.Int64Value(v)
			}
		case float64:
			v, ok := columns[key].(float64)
			if ok {
				attributeTypes[key] = types.Float64Type
				attributes[key] = types.Float64Value(v)
			}
		case string:
			v, ok := columns[key].(string)
			if ok {
				attributeTypes[key] = types.StringType
				attributes[key] = types.StringValue(v)
			}
		case map[string]interface{}:
			v, ok := columns[fmt.Sprintf("_%s_value", key)].(string)
			if ok {
				entityLogicalName, err := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
				if err != nil {
					return nil, fmt.Errorf("error getting entity relation definition info: %s", err.Error())
				}
				dataRecordId := v

				nestedObjectType := types.ObjectType{
					AttrTypes: objectType,
				}
				nestedObjectValue, _ := types.ObjectValue(
					objectType,
					map[string]attr.Value{
						"table_logical_name": types.StringValue(entityLogicalName),
						"data_record_id":     types.StringValue(dataRecordId),
					},
				)

				attributeTypes[key] = nestedObjectType
				attributes[key] = nestedObjectValue
			}
		case []interface{}:
			var listTypes []attr.Type
			var listValues []attr.Value
			tupleElementType := types.ObjectType{
				AttrTypes: objectType,
			}

			relationMap, err := apiClient.GetRelationData(ctx, environmentId, tableLogicalName, *recordid, key)
			if err != nil {
				return nil, fmt.Errorf("error getting relation data: %s", err.Error())
			}

			for _, rawItem := range relationMap {
				item := rawItem.(map[string]interface{})

				relationTableLogicalName, err := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
				if err != nil {
					return nil, fmt.Errorf("error getting entity relation definition info: %s", err.Error())
				}
				entDefinition, err := GetEntityDefinition(ctx, apiClient, environmentId, relationTableLogicalName)
				if err != nil {
					return nil, fmt.Errorf("error getting entity definition: %s", err.Error())
				}

				dataRecordId := item[entDefinition.PrimaryIDAttribute].(string)

				v, _ := types.ObjectValue(objectType, map[string]attr.Value{
					"table_logical_name": types.StringValue(relationTableLogicalName),
					"data_record_id":     types.StringValue(dataRecordId),
				})
				listValues = append(listValues, v)
				listTypes = append(listTypes, tupleElementType)
			}

			nestedObjectType := types.TupleType{
				ElemTypes: listTypes,
			}
			nestedObjectValue, _ := types.TupleValue(listTypes, listValues)

			attributes[key] = nestedObjectValue
			attributeTypes[key] = nestedObjectType
		}
	}

	columnField, _ := types.ObjectValue(attributeTypes, attributes)
	result := types.DynamicValue(columnField)
	return &result, nil
}
