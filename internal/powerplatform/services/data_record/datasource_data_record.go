// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

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

type DataRecordListDataSourceModel struct {
	EnvironmentId                  types.String  `tfsdk:"environment_id"`
	EntityCollection               types.String  `tfsdk:"entity_collection"`
	Select                         []string      `tfsdk:"select"`
	Top                            types.Int64   `tfsdk:"top"`
	ReturnTotalRecordsCount        types.Bool    `tfsdk:"return_total_records_count"`
	TotalRecordsCount              types.Int64   `tfsdk:"total_records_count"`
	TotalRecordsCountLimitExceeded types.Bool    `tfsdk:"total_records_count_limit_exceeded"`
	Query                          types.String  `tfsdk:"query"`
	Items                          types.Dynamic `tfsdk:"items"`
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
			"entity_collection": schema.StringAttribute{
				MarkdownDescription: "Value of the enitiy (collection of the query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#entity-collections]. " +
					"Example:\n\n * $metadata#systemusers \n\n*systemusers \n\n*systemusers(<GUID>) \n\n*systemusers(<GUID>)/systemuserroles_association " +
					"\n\n*contacts(firstname='Joe',emailaddress1='joe@contoso.com') when using (alternate key(s))[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/use-alternate-key-reference-record?tabs=webapi] for single record retrieval",
				Required: true,
			},
			"select": schema.ListAttribute{
				MarkdownDescription: "List of columns to be selected from record(s) defined in entity collection. \n\nMore information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]",
				Required:            false,
				Optional:            true,
				ElementType:         types.StringType,
			},
			"top": schema.Int64Attribute{
				MarkdownDescription: "Number of records to be retrieved. \n\nMore information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]",
				Required:            false,
				Optional:            true,
			},
			"return_total_records_count": schema.BoolAttribute{
				MarkdownDescription: "Should total records count be also retrived. \n\nMore information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]",
				Required:            false,
				Optional:            true,
			},
			"total_records_count": schema.Int64Attribute{
				MarkdownDescription: "Total number of records if attribute `return_total_records_count` is set to `true`",
				Computed:            true,
			},
			"total_records_count_limit_exceeded": schema.BoolAttribute{
				MarkdownDescription: "Is total records count limit exceeded. \n\nMore information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]",
				Computed:            true,
			},
			"items": schema.DynamicAttribute{
				Description: "Columns of the data record table",
				Computed:    true,
			},

			"query": schema.StringAttribute{
				MarkdownDescription: "(OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results] to filter the data records",
				Required:            false,
				Optional:            true,
			},
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

func (d *DataRecordDataSource) buildODataQueryFromModel(model *DataRecordListDataSourceModel) (string, map[string]string, error) {
	var resultQuery = ""
	var headers = make(map[string]string)

	if len(model.Select) > 0 {
		resultQuery = fmt.Sprintf("$select=%s", model.Select[0])
		for i := 1; i < len(model.Select); i++ {
			resultQuery = fmt.Sprintf("%s,%s", resultQuery, model.Select[i])
		}
	}

	if model.Top.ValueInt64Pointer() != nil {
		if len(resultQuery) > 0 {
			resultQuery += "&"
		}
		resultQuery += fmt.Sprintf("$top=%d", *model.Top.ValueInt64Pointer())
	}

	if model.ReturnTotalRecordsCount.ValueBoolPointer() != nil && *model.ReturnTotalRecordsCount.ValueBoolPointer() {
		headers["Prefer"] = "odata.include-annotations=\"Microsoft.Dynamics.CRM.totalrecordcount,Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded\""
		if len(resultQuery) > 0 {
			resultQuery += "&"
		}
		resultQuery += "$count=true"
	}

	if len(resultQuery) > 0 {
		return fmt.Sprintf("%s?%s", model.EntityCollection.ValueString(), resultQuery), headers, nil
	} else {
		return model.EntityCollection.ValueString(), headers, nil
	}
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

	query, headers, err := d.buildODataQueryFromModel(&config)
	tflog.Warn(ctx, fmt.Sprintf("Query: %s", query))
	tflog.Warn(ctx, fmt.Sprintf("Headers: %v", headers))
	if err != nil {
		resp.Diagnostics.AddError("Failed to build OData query", err.Error())
	}

	records, totalrecords, totalRecordsCountLimitExceeded, err := d.DataRecordClient.GetDataRecordsByODataQuery(ctx, config.EnvironmentId.ValueString(), query, headers)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get data records", err.Error())
		return
	}

	if totalrecords != nil {
		state.TotalRecordsCount = types.Int64Value(*totalrecords)
	}
	if totalRecordsCountLimitExceeded != nil {
		state.TotalRecordsCountLimitExceeded = types.BoolValue(*totalRecordsCountLimitExceeded)
	}

	var elements = []attr.Value{}
	for _, record := range records {
		columns, err := convertColumnsToState2(ctx, &d.DataRecordClient, config.EnvironmentId.ValueString(), "systemuser", "systemuserid", record)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert columns to state", err.Error())
			return
		}
		elements = append(elements, types.DynamicValue(columns))

	}

	elementTypes := []attr.Type{}
	for range elements {
		elementTypes = append(elementTypes, types.DynamicType)
	}
	items, _ := types.TupleValue(elementTypes, elements)
	state.Items = types.DynamicValue(items)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func convertColumnsToState2(ctx context.Context, apiClient *DataRecordClient, environmentId, tableLogicalName, primaryFieldName string, columns map[string]interface{}) (*basetypes.ObjectValue, error) {
	var objectType = map[string]attr.Type{
		"table_logical_name": types.StringType,
		"data_record_id":     types.StringType,
	}

	// var old_columns map[string]interface{}
	// jsonColumns, _ := json.Marshal(columns)
	// unquotedJsonColumns, err := strconv.Unquote(string(jsonColumns))
	// if err != nil {
	// 	return nil, err
	// }
	// json.Unmarshal([]byte(unquotedJsonColumns), &old_columns)

	attributeTypes := make(map[string]attr.Type)
	attributes := make(map[string]attr.Value)

	for key, value := range columns {
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
				entityLogicalName := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
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
			setObjectValues := []attr.Value{}
			var setObjectType = types.ObjectType{
				AttrTypes: objectType,
			}
			recordId := columns[primaryFieldName]
			relationMap, err := apiClient.GetRelationData(ctx, recordId.(string), environmentId, tableLogicalName, key)
			if err != nil {
				return nil, err
			}

			for _, rawItem := range relationMap {
				item := rawItem.(map[string]interface{})

				relationTableLogicalName := apiClient.GetEntityRelationDefinitionInfo(ctx, environmentId, tableLogicalName, key)
				dataRecordId := ""

				for itemKey, itemValue := range item {
					if itemKey != "@odata.etag" && itemKey != "createdon" {
						dataRecordId = itemValue.(string)
					}
				}

				setObjectValues = append(setObjectValues, types.ObjectValueMust(objectType,
					map[string]attr.Value{
						"table_logical_name": types.StringValue(relationTableLogicalName),
						"data_record_id":     types.StringValue(dataRecordId),
					}))
			}

			setValue, _ := types.SetValue(setObjectType, setObjectValues)
			attributes[key] = setValue
			attributeTypes[key] = types.SetType{ElemType: setObjectType}
		}
	}

	columnField, _ := types.ObjectValue(attributeTypes, attributes)
	return &columnField, nil
}
