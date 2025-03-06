---
name: "New Data Source Request"
about: "Request a new data source to be added to the Terraform provider."
labels: [enhancement, data source, triage]
assignees: ""

---

## User Story

As a **Terraform user**, I want to query **[data source type]** through Terraform so that I can **[benefit/business value]**.

Use case: [Describe specific scenario where this data source would be valuable in infrastructure automation]

## Data Source

- **Data Source Name:** `powerplatform_[your data source name]`
- **Service Name:** `[service name]`
- **Documentation Link:**

### Potential Terraform Configuration

```hcl
# Sample Terraform config that describes how the new data source might look.

data "example_data_source" {
  # input parameters
  filter = "example" # optional

  # output schema
  groups = toset({
    name = ""
    description = ""
  })
}

```

### Output Schema

```hcl
{

}
```

### Additional Validation Rules

## API documentation

| Action | Verb | URL | Status Codes | Comments |
|--------|------|-----|----------------------|----------|
| Read   | GET  | /api/v1/resources/{id} | 200 | |

### JSON

```json
{}
```

## Definition of Done

- [ ] Data Transfer Objects (dtos) in `dto.go`
- [ ] Data Source Model in `model.go`
- [ ] API Client functions in `api_{name}.go`
- [ ] Data Source implementation in `datasource_{name}.go`
- [ ] Unit Tests in `datasource_{name}_test.go` for Happy Path, Error conditions, boundry cases
- [ ] Acceptance Tests in `datasource_{name}_test.go` for Happy Path
- [ ] Data Source added to `provider.go` and `provider_test.go`
- [ ] Working example in the `/examples` folder
- [ ] Schema documented using `MarkdownDescription` with relevant links to product documentation
- [ ] Change log entry `changie new -k added`
- [ ] Run `make precommit` before PR

See the [contributing guide](/CONTRIBUTING.md?) for more information about what's expected for contributions.
