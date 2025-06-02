# Title: Inefficient Application Filtering Logic in `Read` Method

## Path to file
`/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages.go`

## Problem
The filtering logic for applications uses `continue` within a loop, effectively iterating through all applications even if the filter criteria are not met. This can waste processing resources when there are many applications, especially when the filtering criteria are strict.

## Impact
- **Performance Degradation:** Unnecessary processing consumed iterating over items that can be skipped earlier.
- **Readability Issue:** Using `continue` here makes the code less efficient and slightly harder to follow.

Severity: **Medium**

## Location
Function `Read`, Code snippet:
```go
for _, application := range applications {
    if (state.Name.ValueString() != "" && state.Name.ValueString() != application.Name) ||
        (state.PublisherName.ValueString() != "" && state.PublisherName.ValueString() != application.PublisherName) {
        continue
    }
    state.Applications = append(state.Applications, EnvironmentApplicationPackageDataSourceModel{
        ApplicationId:         types.StringValue(application.ApplicationId),
        Name:                  types.StringValue(application.Name),
        UniqueName:            types.StringValue(application.UniqueName),
        Version:               types.StringValue(application.Version),
        Description:           types.StringValue(application.Description),
        PublisherId:           types.StringValue(application.PublisherId),
        PublisherName:         types.StringValue(application.PublisherName),
        LearnMoreUrl:          types.StringValue(application.LearnMoreUrl),
        State:                 types.StringValue(application.State),
        ApplicationVisibility: types.StringValue(application.ApplicationVisibility),
    })
}
```

## Code Issue
```go
for _, application := range applications {
    if (state.Name.ValueString() != "" && state.Name.ValueString() != application.Name) ||
        (state.PublisherName.ValueString() != "" && state.PublisherName.ValueString() != application.PublisherName) {
        continue
    }
    state.Applications = append(state.Applications, EnvironmentApplicationPackageDataSourceModel{
        ApplicationId:         types.StringValue(application.ApplicationId),
        Name:                  types.StringValue(application.Name),
        UniqueName:            types.StringValue(application.UniqueName),
        Version:               types.StringValue(application.Version),
        Description:           types.StringValue(application.Description),
        PublisherId:           types.StringValue(application.PublisherId),
        PublisherName:         types.StringValue(application.PublisherName),
        LearnMoreUrl:          types.StringValue(application.LearnMoreUrl),
        State:                 types.StringValue(application.State),
        ApplicationVisibility: types.StringValue(application.ApplicationVisibility),
    })
}
```

## Fix
Use `filter` logic as a preprocess step on `applications`, drastically reducing unnecessary iterations and avoiding the need for `continue`.

```go
filteredApplications := make([]api.Application, 0)
for _, application := range applications {
    if (state.Name.ValueString() == "" || state.Name.ValueString() == application.Name) &&
        (state.PublisherName.ValueString() == "" || state.PublisherName.ValueString() == application.PublisherName) {
        filteredApplications = append(filteredApplications, application)
    }
}

for _, application := range filteredApplications {
    state.Applications = append(state.Applications, EnvironmentApplicationPackageDataSourceModel{
        ApplicationId:         types.StringValue(application.ApplicationId),
        Name:                  types.StringValue(application.Name),
        UniqueName:            types.StringValue(application.UniqueName),
        Version:               types.StringValue(application.Version),
        Description:           types.StringValue(application.Description),
        PublisherId:           types.StringValue(application.PublisherId),
        PublisherName:         types.StringValue(application.PublisherName),
        LearnMoreUrl:          types.StringValue(application.LearnMoreUrl),
        State:                 types.StringValue(application.State),
        ApplicationVisibility: types.StringValue(application.ApplicationVisibility),
    })
}
```

Explanation:
- Filters applications matching `Name` and `PublisherName`.
- Reduces the number of iterations required for the `applications` loop by preprocessing.
- Improves clarity and performance of code.
