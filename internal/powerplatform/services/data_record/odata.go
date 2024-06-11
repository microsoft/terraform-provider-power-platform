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

	if model.Top.ValueInt64Pointer() != nil {
		topString := strconv.Itoa(int(*model.Top.ValueInt64Pointer()))
		appendQuery(&resultQuery, buildODataTopPart(&topString))
	}

	if model.ReturnTotalRowsCount.ValueBoolPointer() != nil && *model.ReturnTotalRowsCount.ValueBoolPointer() {
		headers["Prefer"] = "odata.include-annotations=\"Microsoft.Dynamics.CRM.totalrecordcount,Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded\""
		countTrueString := "$count=true"
		appendQuery(&resultQuery, &countTrueString)
	}

	appendQuery(&resultQuery, buildExpandODataQueryPart(model.Expand))

	if len(resultQuery) > 0 {
		return fmt.Sprintf("%s?%s", model.EntityCollection.ValueString(), resultQuery), headers, nil
	} else {
		return model.EntityCollection.ValueString(), headers, nil
	}
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
		result := ""
		for i := 0; i < len(expandQueryStrings); i++ {
			result += fmt.Sprintf("%s,", expandQueryStrings[i])
		}
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

	if model.Top.ValueInt64Pointer() != nil {
		top := strconv.Itoa(int(*model.Top.ValueInt64Pointer()))
		topString := buildODataTopPart(&top)
		if topString != nil {
			if len(resultQuery) > 0 {
				resultQuery += ";"
			}
			resultQuery += *topString
		}
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

func buildODataTopPart(top *string) *string {
	resultQuery := ""
	if top != nil {
		resultQuery = fmt.Sprintf("$top=%s", *top)
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
