# Terraform Data Source Implementation Template for Power Platform Provider

## OVERVIEW

Your task is to implement a new Terraform data source for the Power Platform Terraform Provider. This resource should manage a specific Power Platform object. It must adhere to the established coding patterns of the repository and use the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).  The attached userstory contains the requirements that need to be built.

## PRINCIPLES & RULES

- Include the appropriate license header
- Use the correct package name (matching the folder except for tests) at the beginning of source code files
- Follow the naming conventions (e.g. `powerplatform_${resource_name}`).
- Reuse helper functions and client methods already available in the repository.
- Write both unit tests (with HTTP mocks) and acceptance tests.  Mocked JSON must be in separate files and not inline with the test.
- Use `resp.Diagnostics.AddError` to report errors when available

---

# User Story
