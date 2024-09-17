// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
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
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
	EnvironmentId               types.String   `tfsdk:"environment_id"`
	EntityCollection            types.String   `tfsdk:"entity_collection"`
	Select                      []string       `tfsdk:"select"`
	Filter                      types.String   `tfsdk:"filter"`
	Apply                       types.String   `tfsdk:"apply"`
	OrderBy                     types.String   `tfsdk:"order_by"`
	Top                         types.Int64    `tfsdk:"top"`
	ReturnTotalRowsCount        types.Bool     `tfsdk:"return_total_rows_count"`
	TotalRowsCount              types.Int64    `tfsdk:"total_rows_count"`
	TotalRowsCountLimitExceeded types.Bool     `tfsdk:"total_rows_count_limit_exceeded"`
	SavedQuery                  types.String   `tfsdk:"saved_query"`
	UserQuery                   types.String   `tfsdk:"user_query"`
	Expand                      []ExpandModel  `tfsdk:"expand"`
	Rows                        types.Dynamic  `tfsdk:"rows"`
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

//nolint:unused-receiver
func (d *DataRecordDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for retrieving data records from Dataverse using (OData Query)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#page-results].",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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
			"expand":   returnExpandSchema(10),
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

//nolint:unused-receiver
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

func (d *DataRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataRecordListDataSourceModel
	var config DataRecordListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query, headers, err := BuildODataQueryFromModel(&config)
	tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
	tflog.Debug(ctx, fmt.Sprintf("Headers: %v", headers))
	if err != nil {
		resp.Diagnostics.AddError("Failed to build OData query", err.Error())
	}
	tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))

	queryRespnse, err := d.DataRecordClient.GetDataRecordsByODataQuery(ctx, config.EnvironmentId.ValueString(), query, headers)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get data records", err.Error())
		return
	}

	if queryRespnse.TotalRecord != nil {
		state.TotalRowsCount = types.Int64Value(*queryRespnse.TotalRecord)
	}
	if queryRespnse.TotalRecordLimitExceeded != nil {
		state.TotalRowsCountLimitExceeded = types.BoolValue(*queryRespnse.TotalRecordLimitExceeded)
	}

	var elements = []attr.Value{}
	for _, record := range queryRespnse.Records {

		columns, err := d.convertColumnsToState(record)
		if err != nil {
			resp.Diagnostics.AddError("Failed to convert columns to state", err.Error())
			return
		}
		if columns != nil {
			elements = append(elements, types.DynamicValue(columns))
		}
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

func (d *DataRecordDataSource) convertColumnsToState(columns map[string]interface{}) (*basetypes.DynamicValue, error) {
	if columns == nil {
		return nil, nil
	}
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
			typ, val, _ := d.buildObjectValueFromX(columns[key].(map[string]interface{}))
			tupleElementType := types.ObjectType{
				AttrTypes: typ,
			}
			v, _ := types.ObjectValue(typ, val)
			attributes[key] = v
			attributeTypes[key] = tupleElementType
		case []interface{}:
			typeObj, valObj := d.buildExpandObject(columns[key].([]interface{}))
			attributeTypes[key] = typeObj
			attributes[key] = valObj
		}
	}

	columnField, _ := types.ObjectValue(attributeTypes, attributes)
	result := types.DynamicValue(columnField)
	return &result, nil
}

func (d *DataRecordDataSource) buildObjectValueFromX(columns map[string]interface{}) (map[string]attr.Type, map[string]attr.Value, error) {

	knownObjectType := map[string]attr.Type{}
	knownObjectValue := map[string]attr.Value{}

	for key, value := range columns {
		switch value.(type) {
		case bool:
			v, ok := columns[key].(bool)
			if ok {
				knownObjectType[key] = types.BoolType
				knownObjectValue[key] = types.BoolValue(v)
			}
		case int64:
			v, ok := columns[key].(int64)
			if ok {
				knownObjectType[key] = types.Int64Type
				knownObjectValue[key] = types.Int64Value(v)
			}
		case float64:
			v, ok := columns[key].(float64)
			if ok {
				knownObjectType[key] = types.Float64Type
				knownObjectValue[key] = types.Float64Value(v)
			}
		case string:
			v, ok := columns[key].(string)
			if ok {
				knownObjectType[key] = types.StringType
				knownObjectValue[key] = types.StringValue(v)
			}
		case map[string]interface{}:
			typ, val, _ := d.buildObjectValueFromX(columns[key].(map[string]interface{}))
			tupleElementType := types.ObjectType{
				AttrTypes: typ,
			}
			v, _ := types.ObjectValue(typ, val)
			knownObjectValue[key] = v
			knownObjectType[key] = tupleElementType
		case []interface{}:
			typeObj, valObj := d.buildExpandObject(columns[key].([]interface{}))
			knownObjectValue[key] = valObj
			knownObjectType[key] = typeObj
		}
	}
	return knownObjectType, knownObjectValue, nil
}

func (d *DataRecordDataSource) buildExpandObject(items []interface{}) (basetypes.TupleType, basetypes.TupleValue) {
	var listTypes []attr.Type
	var listValues []attr.Value
	for _, item := range items {

		typ, val, _ := d.buildObjectValueFromX(item.(map[string]interface{}))
		tupleElementType := types.ObjectType{
			AttrTypes: typ,
		}
		v, _ := types.ObjectValue(typ, val)
		listValues = append(listValues, v)
		listTypes = append(listTypes, tupleElementType)

	}
	nestedObjectType := types.TupleType{
		ElemTypes: listTypes,
	}
	nestedObjectValue, _ := types.TupleValue(listTypes, listValues)
	return nestedObjectType, nestedObjectValue
}
