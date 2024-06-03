// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

type ExpandModel struct {
	NavigationProperty types.String  `tfsdk:"navigation_property"`
	Select             []string      `tfsdk:"select"`
	Filter             types.String  `tfsdk:"filter"`
	OrderBy            types.String  `tfsdk:"order_by"`
	Top                types.Int64   `tfsdk:"top"`
	Expand             []ExpandModel `tfsdk:"expand"`
}

type DataRecordListDataSourceModel struct {
	EnvironmentId               types.String  `tfsdk:"environment_id"`
	EntityCollection            types.String  `tfsdk:"entity_collection"`
	Select                      []string      `tfsdk:"select"`
	Filter                      types.String  `tfsdk:"filter"`
	Apply                       types.String  `tfsdk:"apply"`
	OrderBy                     types.String  `tfsdk:"order_by"`
	Top                         types.Int64   `tfsdk:"top"`
	ReturnTotalRowsCount        types.Bool    `tfsdk:"return_total_rows_count"`
	TotalRowsCount              types.Int64   `tfsdk:"total_rows_count"`
	TotalRowsCountLimitExceeded types.Bool    `tfsdk:"total_rows_count_limit_exceeded"`
	SavedQuery                  types.String  `tfsdk:"saved_query"`
	UserQuery                   types.String  `tfsdk:"user_query"`
	Expand                      []ExpandModel `tfsdk:"expand"`
	Rows                        types.Dynamic `tfsdk:"rows"`
}

func (d *DataRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

var navigationPropertySchema = schema.StringAttribute{
	MarkdownDescription: "Navigation property of the entity collection. \n\nMore information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]",
	Required:            true,
}

var selectListAttributeSchema = schema.ListAttribute{
	MarkdownDescription: "List of columns to be selected from record(s) defined in entity collection. \n\nMore information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]",
	Required:            false,
	Optional:            true,
	ElementType:         types.StringType,
}

var filterSchema = schema.StringAttribute{
	MarkdownDescription: "Filter the data records. \n\nMore information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]",
	Required:            false,
	Optional:            true,
}

var orderbySchema = schema.StringAttribute{
	MarkdownDescription: "Order the data records. \n\nMore information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]",
	Required:            false,
	Optional:            true,
}

var topSchema = schema.Int64Attribute{
	MarkdownDescription: "Number of records to be retrieved. \n\nMore information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]",
	Required:            false,
	Optional:            true,
}

func returnExpandSchema(depth int) *schema.ListNestedAttribute {
	description := "Expand the navigation property of the entity collection. \n\nMore information on (OData Expand)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#join-tables]"
	if depth == 0 {
		return &schema.ListNestedAttribute{
			MarkdownDescription: description,
			Optional:            true,
			Required:            false,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"navigation_property": navigationPropertySchema,
					"select":              selectListAttributeSchema,
					"filter":              filterSchema,
					"order_by":            orderbySchema,
					"top":                 topSchema,
				},
			},
		}
	} else {
		return &schema.ListNestedAttribute{
			MarkdownDescription: description,
			Optional:            true,
			Required:            false,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"navigation_property": navigationPropertySchema,
					"select":              selectListAttributeSchema,
					"filter":              filterSchema,
					"order_by":            orderbySchema,
					"top":                 topSchema,
					"expand":              returnExpandSchema(depth - 1),
				},
			},
		}
	}
}

func (d *DataRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// selectListAttributeSchema := schema.ListAttribute{
	// 	MarkdownDescription: "List of columns to be selected from record(s) defined in entity collection. \n\nMore information on (OData Select)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#select-columns]",
	// 	Required:            false,
	// 	Optional:            true,
	// 	ElementType:         types.StringType,
	// }

	// filterSchema := schema.StringAttribute{
	// 	MarkdownDescription: "Filter the data records. \n\nMore information on (OData Filter)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#filter-rows]",
	// 	Required:            false,
	// 	Optional:            true,
	// }

	// orderbySchema := schema.StringAttribute{
	// 	MarkdownDescription: "Order the data records. \n\nMore information on (OData Order By)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#order-rows]",
	// 	Required:            false,
	// 	Optional:            true,
	// }

	// topSchema := schema.Int64Attribute{
	// 	MarkdownDescription: "Number of records to be retrieved. \n\nMore information on (OData Top)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options]",
	// 	Required:            false,
	// 	Optional:            true,
	// }

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
			"select":   selectListAttributeSchema,
			"expand":   returnExpandSchema(5),
			"filter":   filterSchema,
			"order_by": orderbySchema,
			"top":      topSchema,
			"apply": schema.StringAttribute{
				MarkdownDescription: "Apply the aggregation function to the data records. \n\nMore information on (OData Apply)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#aggregate-data]",
				Required:            false,
				Optional:            true,
			},

			"saved_query": schema.StringAttribute{
				MarkdownDescription: "predefined saved query to be used for filtering the data records. \n\nMore information on (Saved Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/retrieve-and-execute-predefined-queries]",
				Required:            false,
				Optional:            true,
			},
			"user_query": schema.StringAttribute{
				MarkdownDescription: "Predefined user query to be used for filtering the data records. \n\nMore information on (Saved Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/retrieve-and-execute-predefined-queries]",
				Required:            false,
				Optional:            true,
			},

			"return_total_rows_count": schema.BoolAttribute{
				MarkdownDescription: "Should total records count be also retrived. \n\nMore information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]",
				Required:            false,
				Optional:            true,
			},
			"total_rows_count": schema.Int64Attribute{
				MarkdownDescription: "Total number of records if attribute `return_total_rows_count` is set to `true`",
				Computed:            true,
			},
			"total_rows_count_limit_exceeded": schema.BoolAttribute{
				MarkdownDescription: "Is total records count limit exceeded. \n\nMore information on (OData Count)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#count-number-of-rows]",
				Computed:            true,
			},
			"rows": schema.DynamicAttribute{
				Description: "Columns of the data record table",
				Computed:    true,
			},
		},
	}
}

func (d *DataRecordDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("user_query"),
			path.MatchRoot("saved_query"),
		),
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

func (d *DataRecordDataSource) buildExpandQueryFilterPart(model *ExpandModel, subExpandValueString *string) *string {
	resultQuery := ""

	s := d.buildODataSelectPart(model.Select)
	if s != nil {
		resultQuery += *s
	}
	f := d.buildODataFilterPart(model.Filter.ValueStringPointer())
	if f != nil {
		resultQuery += *f
	}
	o := d.buildODataOrderByPart(model.OrderBy.ValueStringPointer())
	if o != nil {
		resultQuery += *o
	}

	if subExpandValueString != nil {
		if len(resultQuery) > 0 {
			resultQuery += ";"
		}
		resultQuery += *subExpandValueString
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) buildExpandODataQueryPathRecursive(model []ExpandModel) *string {
	if model == nil {
		return nil
	}

	s := make([]string, 0)
	for _, m := range model {

		expandString := d.buildExpandODataQueryPathRecursive(m.Expand)
		expandQueryFilterString := d.buildExpandQueryFilterPart(&m, expandString)

		if expandQueryFilterString != nil {
			s = append(s, fmt.Sprintf("$expand=%s(%s)", m.NavigationProperty.ValueString(), *expandQueryFilterString))
		}
		// if expandQueryFilterString != nil {
		// 	s = append(s, fmt.Sprintf("$expand=%s(%s)", m.NavigationProperty.ValueString(), *expandQueryFilterString))
		// } else {
		// 	s = append(s, fmt.Sprintf("$expand=%s", m.NavigationProperty.ValueString()))
		// }
	}

	if len(s) > 0 {
		aaa := ""
		for i := 0; i < len(s); i++ {
			aaa += fmt.Sprintf("%s,", s[i])
		}
		aaa = strings.TrimSuffix(aaa, ",")
		return &aaa
	}
	return nil
}

func (d *DataRecordDataSource) buildODataSelectPart(selectPart []string) *string {
	resultQuery := ""
	if len(selectPart) > 0 {
		resultQuery = fmt.Sprintf("$select=%s", selectPart[0])
		for i := 1; i < len(selectPart); i++ {
			resultQuery = fmt.Sprintf("%s,%s", resultQuery, selectPart[i])
		}
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) buildODataFilterPart(filter *string) *string {
	resultQuery := ""
	if filter != nil {
		encoded := url.Values{}
		encoded.Add("$filter", *filter)
		resultQuery += encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) buildODataOrderByPart(orderBy *string) *string {
	resultQuery := ""
	if orderBy != nil {
		encoded := url.Values{}
		encoded.Add("$orderby", *orderBy)
		resultQuery += encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) buildODataTopPart(top *string) *string {
	resultQuery := ""
	if top != nil {
		resultQuery = fmt.Sprintf("$top=%s", *top)
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) buildOdataApplyPart(apply *string) *string {
	resultQuery := ""
	if apply != nil {
		encoded := url.Values{}
		encoded.Add("$apply", *apply)
		resultQuery += encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func (d *DataRecordDataSource) appendQuery(query, part *string) {
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query += *part
	}
}

func (d *DataRecordDataSource) buildODataQueryFromModel(model *DataRecordListDataSourceModel) (string, map[string]string, error) {
	var resultQuery = ""
	var headers = make(map[string]string)

	// if len(model.Select) > 0 {
	// 	resultQuery = fmt.Sprintf("$select=%s", model.Select[0])
	// 	for i := 1; i < len(model.Select); i++ {
	// 		resultQuery = fmt.Sprintf("%s,%s", resultQuery, model.Select[i])
	// 	}
	// }
	d.appendQuery(&resultQuery, d.buildODataSelectPart(model.Select))
	// if s != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	resultQuery += *s
	// }

	//if len(model.Expand) > 0 {
	d.appendQuery(&resultQuery, d.buildExpandODataQueryPathRecursive(model.Expand))
	// if s != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	resultQuery += *s
	// }
	//}

	d.appendQuery(&resultQuery, d.buildODataFilterPart(model.Filter.ValueStringPointer()))
	// if f != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	resultQuery += *f
	// }
	// if model.Filter.ValueStringPointer() != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	encoded := url.Values{}
	// 	encoded.Add("$filter", model.Filter.ValueString())
	// 	resultQuery += encoded.Encode()
	// }

	// if model.Apply.ValueStringPointer() != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	encoded := url.Values{}
	// 	encoded.Add("$apply", model.Apply.ValueString())
	// 	resultQuery += encoded.Encode()
	// }
	d.appendQuery(&resultQuery, d.buildOdataApplyPart(model.Apply.ValueStringPointer()))

	// if model.OrderBy.ValueStringPointer() != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	encoded := url.Values{}
	// 	encoded.Add("$orderby", model.OrderBy.ValueString())
	// 	resultQuery += encoded.Encode()
	// }
	d.appendQuery(&resultQuery, d.buildODataOrderByPart(model.OrderBy.ValueStringPointer()))
	// if ob != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	resultQuery += *ob
	// }

	// if model.Top.ValueInt64Pointer() != nil {
	// 	if len(resultQuery) > 0 {
	// 		resultQuery += "&"
	// 	}
	// 	resultQuery += fmt.Sprintf("$top=%d", *model.Top.ValueInt64Pointer())
	// }

	//TODO
	//d.appendQuery(&resultQuery, d.buildODataTopPart(model.Top.ValueInt64Pointer()))

	if model.ReturnTotalRowsCount.ValueBoolPointer() != nil && *model.ReturnTotalRowsCount.ValueBoolPointer() {
		headers["Prefer"] = "odata.include-annotations=\"Microsoft.Dynamics.CRM.totalrecordcount,Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded\""
		countTrueString := "$count=true"
		d.appendQuery(&resultQuery, &countTrueString)
		// if len(resultQuery) > 0 {
		// 	resultQuery += "&"
		// }
		// resultQuery += "$count=true"
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
		state.TotalRowsCount = types.Int64Value(*totalrecords)
	}
	if totalRecordsCountLimitExceeded != nil {
		state.TotalRowsCountLimitExceeded = types.BoolValue(*totalRecordsCountLimitExceeded)
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
	rows, _ := types.TupleValue(elementTypes, elements)
	state.Rows = types.DynamicValue(rows)

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
