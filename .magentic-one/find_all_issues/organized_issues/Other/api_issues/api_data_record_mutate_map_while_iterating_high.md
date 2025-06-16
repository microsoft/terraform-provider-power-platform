# Title

Resource leak: `columns` map mutation in loop can cause runtime nondeterminism

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go

## Problem

In `ApplyDataRecord`, there is a loop over the `columns` map where keys are deleted during iteration (`delete(columns, key)`), and new items are being added to the same map inline. In Go, deleting from a map while iterating over it directly leads to undefined and unpredictable behavior, which may cause subtle bugs, such as skipping keys or panicking at runtime.

## Impact

**Severity: High**

- Can cause undefined runtime behavior, hard-to-diagnose bugs, and possible panics.
- Data sent to API may be incomplete or corrupted.
- Terraform runs could behave nondeterministically due to missing or extra map entries.

## Location

```go
for key, value := range columns {
	if nestedMap, ok := value.(map[string]any); ok {
		delete(columns, key)
		if len(nestedMap) > 0 {
			tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
			if err != nil {
				return nil, err
			}
			// ...
			columns[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("/%s(%s)", entityDefinition.LogicalCollectionName, dataRecordId)
		}
	} else if nestedMapList, ok := value.([]any); ok {
		delete(columns, key)
		relations[key] = nestedMapList
	}
}
```

## Code Issue

```go
for key, value := range columns {
	if nestedMap, ok := value.(map[string]any); ok {
		delete(columns, key)
		if len(nestedMap) > 0 {
			tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
			if err != nil {
				return nil, err
			}

			entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
			if err != nil {
				return nil, err
			}

			columns[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("/%s(%s)", entityDefinition.LogicalCollectionName, dataRecordId)
		}
	} else if nestedMapList, ok := value.([]any); ok {
		delete(columns, key)
		relations[key] = nestedMapList
	}
}
```

## Fix

Build a list of mutations, and apply them after the loop. For example:

```go
// Collect keys to delete and values to add
var keysToDelete []string
var bindingsToAdd = make(map[string]any)
var relations = make(map[string]any)

for key, value := range columns {
    if nestedMap, ok := value.(map[string]any); ok {
        keysToDelete = append(keysToDelete, key)
        if len(nestedMap) > 0 {
            tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
            if err != nil {
                return nil, err
            }
            entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
            if err != nil {
                return nil, err
            }
            bindingsToAdd[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("/%s(%s)", entityDefinition.LogicalCollectionName, dataRecordId)
        }
    } else if nestedMapList, ok := value.([]any); ok {
        keysToDelete = append(keysToDelete, key)
        relations[key] = nestedMapList
    }
}
for _, k := range keysToDelete {
    delete(columns, k)
}
for k, v := range bindingsToAdd {
    columns[k] = v
}
```

This approach guarantees consistency and avoids mutating the map during iteration.

---

Save as:
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/resource_management/api_data_record_mutate_map_while_iterating_high.md`
