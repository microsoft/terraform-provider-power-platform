// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &DataRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &DataRecordDataSource{}
)

func NewDataRecordDataSource() datasource.DataSource {
	return &DataRecordDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_data_records",
	}
}

type DataRecordDataSource struct {
	DataRecordClient DataRecordClient
	ProviderTypeName string
	TypeName         string
}

func (d *DataRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *DataRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for retrieving data records from Dataverse using (OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results].",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Power Platform environment",
				Required:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "(OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results] to filter the data records",
				Required:            true,
			},
			"items": schema.DynamicAttribute{
				Description: "Columns of the data record table",
				Computed:    true,
			},
			// "records": schema.ListNestedAttribute{
			// 	MarkdownDescription: "List of data records",
			// 	Computed:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"columns": schema.DynamicAttribute{
			// 				Description: "Columns of the data record table",
			// 				Computed:    true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func (d *DataRecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*api.ProviderClient).Api

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.DataRecordClient = NewDataRecordClient(client)
}

type DataRecordListDataSourceModel struct {
	EnvironmentId types.String  `tfsdk:"environment_id"`
	Query         types.String  `tfsdk:"query"`
	Items         types.Dynamic `tfsdk:"items"`
	//Records       []DataRecordDataSourceModel `tfsdk:"records"`
}

type DataRecordDataSourceModel struct {
	Columns types.Dynamic `tfsdk:"columns"`
}

func (d *DataRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataRecordListDataSourceModel
	var config DataRecordListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	attributeTypes1 := make(map[string]attr.Type)
	attributes1 := make(map[string]attr.Value)

	attributeTypes1["col1"] = types.StringType
	attributes1["col1"] = types.StringValue("test")

	v1, _ := types.ObjectValue(attributeTypes1, attributes1)

	attributeTypes2 := make(map[string]attr.Type)
	attributes2 := make(map[string]attr.Value)

	attributeTypes2["col1"] = types.StringType
	attributes2["col1"] = types.StringValue("test")

	v2, _ := types.ObjectValue(attributeTypes2, attributes2)

	var elements = []attr.Value{
		types.DynamicValue(v1),
		types.DynamicValue(v2),
	}
	aaa, _ := types.ListValue(types.DynamicType, elements)

	state.Items = types.DynamicValue(aaa)

	// records, err := d.DataRecordClient.GetDataRecordsByODataQuery(ctx, config.EnvironmentId.String(), config.Query.String())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Failed to get data records", err.Error())
	// 	return
	// }

	// for _, record := range records {
	// 	columns, err := convertColumnsToState2(ctx, &d.DataRecordClient, config.EnvironmentId.String(), "systemuser", record)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Failed to convert columns to state", err.Error())
	// 		return
	// 	}

	// 	v, _ := types.ObjectValue(objectTypeItemList, map[string]attr.Value{
	// 		"columns": types.DynamicValue(columns),
	// 	})

	// 	listValues = append(listValues, v)
	// 	listTypes = append(listTypes, tupleElementType)

	// 	nestedObjectType := types.TupleType{
	// 		ElemTypes: listTypes,
	// 	}
	// 	nestedObjectValue, _ := types.TupleValue(listTypes, listValues)

	// 	// state.Records = append(state.Records, DataRecordDataSourceModel{
	// 	// 	Columns: types.DynamicValue(columns),
	// 	// })
	// }

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func convertColumnsToState2(ctx context.Context, apiClient *DataRecordClient, environmentId string, tableLogicalName string, columns map[string]interface{}) (*basetypes.ObjectValue, error) {
	var objectType = map[string]attr.Type{
		"entity_logical_name": types.StringType,
		"data_record_id":      types.StringType,
	}

	var old_columns map[string]interface{}
	jsonColumns, _ := json.Marshal(columns)
	unquotedJsonColumns, _ := strconv.Unquote(string(jsonColumns))
	json.Unmarshal([]byte(unquotedJsonColumns), &old_columns)

	attributeTypes := make(map[string]attr.Type)
	attributes := make(map[string]attr.Value)

	for key, value := range old_columns {
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
				entityLogicalName := "" //apiClient.GetEntityRelationTableName(ctx, environmentId, tableLogicalName, key)
				dataRecordId := v

				nestedObjectType := types.ObjectType{
					AttrTypes: objectType,
				}
				nestedObjectValue, _ := types.ObjectValue(
					objectType,
					map[string]attr.Value{
						"entity_logical_name": types.StringValue(entityLogicalName),
						"data_record_id":      types.StringValue(dataRecordId),
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
			for _, value := range value.([]interface{}) {
				item := value.(map[string]interface{})

				entityLogicalName := "" /// apiClient.GetEntityRelationTableName(ctx, environmentId, tableLogicalName, key)
				dataRecordId := item["data_record_id"].(string)

				v, _ := types.ObjectValue(objectType, map[string]attr.Value{
					"entity_logical_name": types.StringValue(entityLogicalName),
					"data_record_id":      types.StringValue(dataRecordId),
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
	return &columnField, nil

	//currentState.EnvironmentId = types.StringValue(environmentId)
	//currentState.TableLogicalName = types.StringValue(tableLogicalName)
	//currentState.Columns = types.DynamicValue(column_field)
}
