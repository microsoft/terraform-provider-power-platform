// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package powerplatform

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func BuildODataQueryFromModel(model *DataRecordListDataSourceModel) (string, map[string]string, error) {
	var resultQuery = ""
	var headers = make(map[string]string)

	appendQuery(&resultQuery, buildODataSelectPart(model.Select))
	appendQuery(&resultQuery, buildODataFilterPart(model.Filter.ValueStringPointer()))
	appendQuery(&resultQuery, buildOdataApplyPart(model.Apply.ValueStringPointer()))
	appendQuery(&resultQuery, buildODataOrderByPart(model.OrderBy.ValueStringPointer()))
	appendQuery(&resultQuery, buildODataTopPart(model.Top.ValueInt64Pointer()))
	appendQuery(&resultQuery, buildTotalRowsCountPart(headers, model.ReturnTotalRowsCount.ValueBoolPointer()))
	appendQuery(&resultQuery, buildExpandODataQueryPart(model.Expand))

	if len(resultQuery) > 0 {
		return fmt.Sprintf("%s?%s", model.EntityCollection.ValueString(), resultQuery), headers, nil
	} else {
		return model.EntityCollection.ValueString(), headers, nil
	}
}

func buildTotalRowsCountPart(headers map[string]string, returnTotalRowsCount *bool) *string {
	if returnTotalRowsCount != nil && *returnTotalRowsCount {
		headers["Prefer"] = "odata.include-annotations=\"Microsoft.Dynamics.CRM.totalrecordcount,Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded\""
		countTrueString := "$count=true"
		return &countTrueString
	}
	return nil
}

func buildExpandODataQueryPart(model []ExpandModel) *string {
	if model == nil {
		return nil
	}

	expandQueryStrings := make([]string, 0)
	for _, m := range model {
		expandString := buildExpandODataQueryPart(m.Expand)
		expandQueryFilterString := buildExpandFilterQueryPart(&m, expandString)

		if expandQueryFilterString != nil {
			expandQueryStrings = append(expandQueryStrings, fmt.Sprintf("%s(%s)", m.NavigationProperty.ValueString(), *expandQueryFilterString))
		} else {
			expandQueryStrings = append(expandQueryStrings, m.NavigationProperty.ValueString())

		}
	}

	if len(expandQueryStrings) > 0 {
		result := strings.Join(expandQueryStrings, ",")
		result = "$expand=" + strings.TrimSuffix(result, ",")
		return &result
	}
	return nil
}

func buildExpandFilterQueryPart(model *ExpandModel, subExpandValueString *string) *string {
	resultQuery := ""

	selectString := buildODataSelectPart(model.Select)
	if selectString != nil {
		resultQuery += *selectString
	}
	filterString := buildODataFilterPart(model.Filter.ValueStringPointer())
	if filterString != nil {
		if len(resultQuery) > 0 {
			resultQuery += ";"
		}
		resultQuery += *filterString
	}
	orderByString := buildODataOrderByPart(model.OrderBy.ValueStringPointer())
	if orderByString != nil {
		if len(resultQuery) > 0 {
			resultQuery += ";"
		}
		resultQuery += *orderByString
	}

	topString := buildODataTopPart(model.Top.ValueInt64Pointer())
	if topString != nil {
		if len(resultQuery) > 0 {
			resultQuery += ";"
		}
		resultQuery += *topString
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

func appendQuery(query, part *string) {
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query += *part
	}
}

func buildODataSelectPart(selectPart []string) *string {
	resultQuery := ""
	if len(selectPart) > 0 {
		resultQuery = "$select=" + strings.Join(selectPart, ",")
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func buildODataFilterPart(filter *string) *string {
	resultQuery := ""
	if filter != nil {
		encoded := url.Values{}
		encoded.Add("filter", *filter)
		resultQuery += "$" + encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func buildODataOrderByPart(orderBy *string) *string {
	resultQuery := ""
	if orderBy != nil {
		encoded := url.Values{}
		encoded.Add("orderby", *orderBy)
		resultQuery += "$" + encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func buildODataTopPart(top *int64) *string {
	resultQuery := ""
	if top != nil {
		topString := strconv.Itoa(int(*top))
		resultQuery = fmt.Sprintf("$top=%s", topString)
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}

func buildOdataApplyPart(apply *string) *string {
	resultQuery := ""
	if apply != nil {
		encoded := url.Values{}
		encoded.Add("apply", *apply)
		resultQuery += "$" + encoded.Encode()
	}
	if resultQuery == "" {
		return nil
	}
	return &resultQuery
}
