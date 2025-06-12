# Inefficient JSON Marshalling in convertResourceModelToMap

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

The function `convertResourceModelToMap` first replaces `<null>` with `""` in the string, then marshals the string to JSON, then unquotes, then unmarshals to map. This is unnecessary and inefficient; you can simply unmarshal the original string into a map.

## Impact

- **Severity:** Medium
- Wastes CPU resources and developer time understanding the intent.
- Can introduce subtle bugs if the input string is not exactly as expected.
- Lower code clarity and unnecessarily complex implementation.

## Location

Function: `convertResourceModelToMap`

## Code Issue

```go
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
```

## Fix

Assuming `*columnsAsString` is the JSON string of the columns, just unmarshal it directly:

```go
err = json.Unmarshal([]byte(*columnsAsString), &mapColumns)
if err != nil {
    return nil, err
}
```

Full corrected function:

```go
func convertResourceModelToMap(columnsAsString *string) (mapColumns map[string]any, err error) {
    if columnsAsString == nil {
        return nil, nil
    }
    replacedColumns := strings.ReplaceAll(*columnsAsString, `<null>`, `""`)
    err = json.Unmarshal([]byte(replacedColumns), &mapColumns)
    if err != nil {
        return nil, err
    }
    return mapColumns, nil
}
```
