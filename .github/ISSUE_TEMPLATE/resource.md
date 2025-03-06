---
name: "New Resource Request"
about: "Request a new resource to be added to the Terraform provider."
labels: [enhancement,resource,triage]
assignees: ""

---

## User Story

As a **Terraform user**, I want to manage **[resource type]** through Terraform so that I can **[benefit/business value]**.

Use case: [Describe specific scenario where this resource would be valuable in infrastructure automation]

## Resource

- **Resource Name:** `powerplatform_[your resource name]`
- **Service Name:** `[service name]`
- **Documentation Link:**

### Potential Terraform Configuration

```hcl
# Sample Terraform config that describes how the new resource might look.

resource "example_resource" {
  name = "example" # required
  parameter1 = "value1"
  enabled = false
  items = toset([
    { 
      name = "item name" # required, must be 3 characters or more
    }
  ])
}

```

### Additional Validation Rules

## API documentation

| Action | Verb | URL | Status Codes | Comments |
|--------|------|-----|----------------------|----------|
| Create | POST | /api/v1/resources | 201 | |
| Read   | GET  | /api/v1/resources/{id} | 200 | |
| Update | PUT  | /api/v1/resources/{id} | 200 | |
| Delete | DELETE | /api/v1/resources/{id} | 204 | |

### JSON

```json
{}
```

## Definition of Done

- [ ] Data Transfer Objects (dtos) in `dto.go`
- [ ] Resource Model in `model.go`
- [ ] API Client functions in `api_{name}.go`
- [ ] Resource Implementation in `resource_{name}.go`
- [ ] Unit Tests in `resource_{name}_test.go` for Happy Path, Error conditions, boundry cases
- [ ] Acceptance Tests in `resource_{name}_test.go` for Happy Path
- [ ] Resource Added to `provider.go` and `provider_test.go`
- [ ] Example in the `/examples` folder
- [ ] Schema documented using `MarkdownDescription`
- [ ] Change log entry `changie new -k added`
- [ ] Run `make precommit` before PR

See the [contributing guide](/CONTRIBUTING.md?) for more information about what's expected for contributions.
