// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewDataRecordResource() resource.Resource {
	return &DataRecordResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "data_record",
		},
	}
}

func (r *DataRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataRecordResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		DynamicColumns(
			path.Root("columns").Expression(),
		),
	}
}

func (r *DataRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"disable_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "If true, the resource will either set isdisabled to true or statecode to 1 with a PATCH request, before attempting to delete the record.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Dynamics 365 environment",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"table_logical_name": schema.StringAttribute{
				MarkdownDescription: "Logical name of the data record table",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"columns": schema.DynamicAttribute{
				MarkdownDescription: "Columns of the data record table",
				Required:            true,
			},
		},
	}
}

func (r *DataRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.DataRecordClient = newDataRecordClient(providerClient.Api)
}

func (r *DataRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan DataRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Why are we setting plan to itself? Remove?
	plan.Id = types.StringValue(plan.Id.ValueString())
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
	plan.Columns = types.DynamicValue(plan.Columns)

	stateColumns := plan.Columns.String()
	mapColumns, err := convertResourceModelToMap(&stateColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
		return
	}

	dr, err := r.DataRecordClient.ApplyDataRecord(ctx, plan.Id.ValueString(), plan.EnvironmentId.ValueString(), plan.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.Id = types.StringValue(dr.Id)

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state DataRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newColumns, err := r.DataRecordClient.GetDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
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

	tflog.Debug(ctx, fmt.Sprintf("READ: %s with table_name %s", r.FullTypeName(), state.TableLogicalName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DataRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan DataRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state DataRecordResourceModel
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
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
		return
	}

	plan.Id = types.StringValue(dr.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *DataRecordResourceModel
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

	if state.DisableOnDestroy.ValueBool() {
		entityAttr, err := r.DataRecordClient.GetEntityAttributesDefinition(ctx, state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when getting entity attributes definition %s", r.FullTypeName()), err.Error())
			return
		}
		var containsIsDisableAttr, containsStateCode bool
		attributes := map[string]any{}
		for _, value := range entityAttr {
			if strings.Compare(value.LogicalName, "isdisabled") == 0 {
				containsIsDisableAttr = true
			}
			if strings.Compare(value.LogicalName, "statecode") == 0 {
				containsStateCode = true
			}
		}
		// in some cases both attributes are present, in that case we will set statecode to 1 and isdisabled will be ignored.
		if (containsStateCode && !containsIsDisableAttr) || (containsStateCode && containsIsDisableAttr) {
			attributes["statecode"] = 1
			// attributes["statuscode"] = 2  we can't set statuscode (status reason) because it may be customized by the user.
		} else if containsIsDisableAttr {
			attributes["isdisabled"] = true
		} else {
			tflog.Debug(ctx, fmt.Sprintf("No statecode or isdisabled attribute found for %s", r.FullTypeName()))
		}

		if len(attributes) > 0 {
			_, err = r.DataRecordClient.ApplyDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), attributes)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when disabling %s", r.FullTypeName()), err.Error())
				return
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("No statecode or isdisabled attribute found for %s", r.FullTypeName()))
		}
	}

	err = r.DataRecordClient.DeleteDataRecord(ctx, state.Id.ValueString(), state.EnvironmentId.ValueString(), state.TableLogicalName.ValueString(), mapColumns)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
		return
	}
}

func (r *DataRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertResourceModelToMap(columnsAsString *string) (mapColumns map[string]any, err error) {
	if columnsAsString == nil {
		return nil, nil
	}

	replacedColumns := strings.ReplaceAll(*columnsAsString, `<null>`, `""`)
	columnsAsString = &replacedColumns

	jsonColumns, err := json.Marshal(columnsAsString)
	if err != nil {
		return nil, err
	}
	unquotedJsonColumns, err := strconv.Unquote(string(jsonColumns))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(unquotedJsonColumns), &mapColumns)
	if err != nil {
		return nil, err
	}
	return mapColumns, nil
}

func caseBool(ctx context.Context, columnValue any, attrValue map[string]attr.Value, attrType map[string]attr.Type, key string) {
	value, ok := columnValue.(bool)
	if !ok {
		tflog.Debug(ctx, "caseBool: failed to cast value to bool", map[string]any{"key": key, "value_type": fmt.Sprintf("%T", columnValue)})
		return
	}
	attrValue[key] = types.BoolValue(value)
	attrType[key] = types.BoolType
}

func caseInt64(ctx context.Context, columnValue any, attrValue map[string]attr.Value, attrType map[string]attr.Type, key string) {
	value, ok := columnValue.(int64)
	if !ok {
		tflog.Debug(ctx, "caseInt64: failed to cast value to int64", map[string]any{"key": key, "value_type": fmt.Sprintf("%T", columnValue)})
		return
	}
	attrValue[key] = types.Int64Value(value)
	attrType[key] = types.Int64Type
}

func caseFloat64(ctx context.Context, columnValue any, attrValue map[string]attr.Value, attrType map[string]attr.Type, key string) {
	value, ok := columnValue.(float64)
	if !ok {
		tflog.Debug(ctx, "caseFloat64: failed to cast value to float64", map[string]any{"key": key, "value_type": fmt.Sprintf("%T", columnValue)})
		return
	}
	attrValue[key] = types.Float64Value(value)
	attrType[key] = types.Float64Type
}

func caseString(ctx context.Context, columnValue any, attrValue map[string]attr.Value, attrType map[string]attr.Type, key string) {
	value, ok := columnValue.(string)
	if !ok {
		tflog.Debug(ctx, "caseString: failed to cast value to string", map[string]any{"key": key, "value_type": fmt.Sprintf("%T", columnValue)})
		return
	}
	attrValue[key] = types.StringValue(value)
	attrType[key] = types.StringType
}

func caseMapStringOfAny(ctx context.Context, columnValue any, attrValue map[string]attr.Value, attrType map[string]attr.Type, key, entityLogicalName string, objectType map[string]attr.Type) {
	value, ok := columnValue.(string)
	if !ok {
		tflog.Debug(ctx, "caseMapStringOfAny: failed to cast value to string", map[string]any{"key": key, "value_type": fmt.Sprintf("%T", columnValue)})
		return
	}
	dataRecordId := value
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
	attrType[key] = nestedObjectType
	attrValue[key] = nestedObjectValue
}

func caseArrayOfAny(ctx context.Context, attrValue map[string]attr.Value, attrType map[string]attr.Type,
	apiClient *client, objectType map[string]attr.Type, key, environmentId, tableLogicalName, recordid string) error {
	var listValues []attr.Value
	tupleElementType := types.ObjectType{
		AttrTypes: objectType,
	}

	relationMap, err := apiClient.GetRelationData(ctx, environmentId, tableLogicalName, recordid, key)
	if err != nil {
		return errors.New("error getting relation data: " + err.Error())
	}

	for _, rawItem := range relationMap {
		item, ok := rawItem.(map[string]any)
		if !ok {
			return errors.New("error asserting rawItem to map[string]any")
		}

		relationTableLogicalName, err := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
		if err != nil {
			return errors.New("error getting entity relation definition info: " + err.Error())
		}
		entDefinition, err := getEntityDefinition(ctx, apiClient, environmentId, relationTableLogicalName)
		if err != nil {
			return errors.New("error getting entity definition: " + err.Error())
		}

		dataRecordId, ok := item[entDefinition.PrimaryIDAttribute].(string)
		if !ok {
			return errors.New("error asserting dataRecordId to string")
		}

		v, _ := types.ObjectValue(objectType, map[string]attr.Value{
			"table_logical_name": types.StringValue(relationTableLogicalName),
			"data_record_id":     types.StringValue(dataRecordId),
		})
		listValues = append(listValues, v)
	}

	nestedObjectValue, _ := types.SetValue(tupleElementType, listValues)
	attrValue[key] = nestedObjectValue
	attrType[key] = types.SetType{
		ElemType: tupleElementType,
	}

	return nil
}

func (r *DataRecordResource) convertColumnsToState(ctx context.Context, apiClient *client, environmentId, tableLogicalName string, recordid, recordColumns *string, columns map[string]any) (*basetypes.DynamicValue, error) {
	var objectType = map[string]attr.Type{
		"table_logical_name": types.StringType,
		"data_record_id":     types.StringType,
	}

	mapColumns, err := convertResourceModelToMap(recordColumns)
	if err != nil {
		return nil, errors.New("error converting columns to map: " + err.Error())
	}

	attributeTypes := make(map[string]attr.Type)
	attributes := make(map[string]attr.Value)

	for key, value := range mapColumns {
		switch value.(type) {
		case bool:
			caseBool(ctx, columns[key], attributes, attributeTypes, key)
		case int64:
			caseInt64(ctx, columns[key], attributes, attributeTypes, key)
		case float64:
			caseFloat64(ctx, columns[key], attributes, attributeTypes, key)
		case string:
			caseString(ctx, columns[key], attributes, attributeTypes, key)
		case map[string]any:
			entityLogicalName, err := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
			if err != nil {
				return nil, errors.New("error getting entity relation definition info: " + err.Error())
			}
			caseMapStringOfAny(ctx, columns[fmt.Sprintf("_%s_value", key)], attributes, attributeTypes, key, entityLogicalName, objectType)
		case []any:
			err := caseArrayOfAny(ctx, attributes, attributeTypes, apiClient, objectType, key, environmentId, tableLogicalName, *recordid)
			if err != nil {
				return nil, err
			}
		}
	}
	columnField, diags := types.ObjectValue(attributeTypes, attributes)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to create object value: %v", diags)
	}
	result := types.DynamicValue(columnField)
	return &result, nil
}
