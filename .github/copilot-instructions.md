# GitHub Copilot Custom Instructions – Terraform Provider Power Platform

These instructions guide GitHub Copilot to follow our project's conventions and best practices when suggesting code. They cover how to format code, name resources and attributes, structure implementations, and write tests in this repository. By following these guidelines, Copilot's suggestions should align with the project's style and help contributors produce high-quality, consistent code. Always consider existing patterns in the repository—when in doubt, review similar resources or tests for reference and keep the new code idiomatic to the project's practices.

## Development Setup & Workflow

- Use the provided **Makefile** commands for all build and test tasks:
  - `make install` to compile the provider code.
  - `make lint` to run linters and ensure code style compliance.
  - `make unittest` to run all unit tests (optionally use `TEST=<prefix>` to run tests matching a name prefix, e.g. `make unittest TEST=Environment` to run tests named with that prefix). This filters tests by regex `^(TestAcc|TestUnit)<prefix>`.
  - `make acctest TEST=<prefix>` to run acceptance tests (integration tests) matching a prefix. Always provide a specific test prefix to limit scope, and run these tests **only with user consent** (they run against real cloud resources). Note that `make acctest` automatically sets `TF_ACC=1` (no need to set it manually).
  - `make userdocs` to regenerate documentation
  - `make precommit` to run all checks once code is ready to commit
  - `make coverage` to run all unit tests and output a code coverage report. It also shows the files that have changed on this branch to help target coverage suggestions to files in the current PR.
- Always run the above `make` commands from the repository root (e.g. in the `/workspaces/terraform-provider-power-platform` directory).
- **Never run** `terraform init` inside the provider repo. Terraform is only used in examples or tests; initializing in the provider directory is not needed and may cause conflicts.
- Do not manually edit files under the `/docs` folder. These files are auto-generated from the schema `MarkdownDescription` attributes. Instead, update schema's `MarkdownDescription` in code and run `make userdocs` to regenerate documentation.
- To try out an example configuration, navigate to its directory under `/examples` and run `terraform apply -auto-approve` (ensure you've built the provider and set it in your Terraform plugins path beforehand).

## File and Folder Structure

### Service Organization

- Organize all service implementations within the `internal/services` directory, with each service in its own subdirectory.
- Name service directories using lowercase words with underscores (e.g., `tenant_settings`, `environment_templates`).
- Choose service names that reflect the Power Platform domain they represent.

### Service Files

Each service directory MUST contain:

- **Models File**: Create a single `models.go` file containing all data models and DTO conversion functions.
- **API Client File**: Name as `api_<service_name>.go` (e.g., `api_licensing.go`).
- **Data Transfer Object File**: Create a single `dto.go` file containing all DTO objects used by the client to represent JSON sent to or received from the API
- **Test Files**: Name as `resource_<resource_name>_test.go` or `datasource_<data_source_name>_test.go` and place in the same directory. Must have a test file for every resource and data source file in a service.
- **Mock Data Files**: Place test JSON fixtures in `tests/resource/<test_scenario>/` or `tests/datasource/<test_scenario>/` subdirectories. Name JSON test files according to the pattern `<method>_<object>.json` (e.g., `get_environment.json`, `post_lifecycle.json`). JSON test files may have a numerical index if multiple API calls are made in a test scenario.

Each service SHALL contain one or more resources or data sources:

- **Resource Files**: Name as `resource_<resource_name>.go` (e.g., `resource_environment.go`).
- **Data Source Files**: Name as `datasource_<data_source_name>.go` (e.g., `datasource_tenant_settings.go`).

### Example Files

The `/examples` directory provides usage examples for the provider, resources, and data sources:

- Organize examples in three top-level directories:
  - `data-sources/` - Contains examples for all data source types
  - `resources/` - Contains examples for all resource types
  - `provider/` - Contains provider configuration examples

- Create a subdirectory for each resource or data source type under its respective category:
  - Name subdirectories exactly matching the resource/data source name (e.g., `powerplatform_environment`)
  - Include the resource prefix (`powerplatform_`) in the directory name

- Include the following files in each resource example directory:
  - `resource.tf` - Required, Basic usage example for the resource
  - `variables.tf` - Optional, for examples with variable inputs
  - `outputs.tf` - Optional, for examples that return outputs
  - `import.sh` - Optional, For resources that support import, include example import command.
  - Optional, Additional `.tf` files for more complex or specific use cases. Ask for human approval before creating non-standard `.tf` files.

- Include the following files in each data source example directory:
  - `data-source.tf` - Required, Basic usage example for the data source
  - `variables.tf` - Optional, for examples with variable inputs
  - `outputs.tf` - Optional, for examples that return outputs

- For import script examples:
  - Use comments to explain the import syntax
  - Include placeholder values (e.g., `00000000-0000-0000-0000-000000000000`) for IDs
  - Show the proper resource name format (e.g., `powerplatform_environment.example`)

- Keep example code simple and focused on demonstrating one clear use case per file

## Naming Conventions

- **Resource and Data Source Names:** Follow the existing naming pattern of prefixing with `powerplatform_`. For example, an environment resource is named `powerplatform_environment`. Use lowercase with underscores for Terraform resource/data names.
- **Attribute Naming:** Name resource attributes to match Power Platform terminology. Prefer the modern, user-friendly terms used in the current Power Platform API/UX/[Official Documentation](https://learn.microsoft.com/en-us/power-platform/admin/) over deprecated names. Keep names concise but descriptive of their function in the resource.
- **Test Function Naming:** Name test functions with a prefix indicating their type. **Acceptance test** functions should start with `TestAcc` and **unit test** functions with `TestUnit` (this allows filtering tests by type). Also, name test files' package with a `_test` suffix (e.g. `package environment_test`) to ensure tests access the provider only via its public interface.
- **Data Transfer Objects:** Define all DTO structures in `dto.go` with a `Dto` suffix (e.g., `EnvironmentDto`).
- **Conversion Functions:** Implement conversion functions named exactly as `convertDtoToModel` and `convertModelToDto` in `models.go`.
- **Client Factory:** When implementing a client factory, name it `New<Service>Client` (e.g., `NewSolutionClient`).
- **Resource/Data Source Factory:** For each resource and data source, create a new function named `New<ResourceName>Resource` or `New<DataSourceName>DataSource` that returns the appropriate type.

## Comments and Documentation

- Include the appropriate license header at the top of all new code files
- Write Go comments only on exported functions, types, and methods to explain their purpose, parameters, and return values when it adds clarity.
- Focus comments on **why** something is done if it's not obvious from the code.
- Avoid redundant comments that just restate the code or don't provide additional insight.
- When defining resource or data source schema, **always use** the `MarkdownDescription` field for documentation. Do **not** use the deprecated `Description` field. Markdown descriptions will be used to auto-generate docs, so make them clear and user-friendly, and include links to topics in the [official Power Platform docs](https://learn.microsoft.com/en-us/power-platform/admin/) when helpful.

## Code Organization and Implementation Guidelines

### Frameworks

- **Terraform Plugin Framework:** Use [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) for implementing resources and data sources. Avoid using legacy Terraform SDK constructs.
- **Azure Identity Client:** Use [Azure Identity Client Module for Go](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity) for implementing authentication methods

### Common Utilities

- **API Layer:** Use the common API client functionality in `internal/api` for making Power Platform API calls:
  - Use service-specific clients that build upon the common API layer
  - Handle API errors gracefully and return detailed error diagnostics
  - Use the retry mechanisms provided by the API layer for transient failures

- **Constants:** Reference centralized constants from `internal/constants/constants.go` instead of hardcoding values for:
  - API endpoints and URL paths
  - Common string literals and configuration keys
  - Status codes and enum values used across the provider

- **Error Handling:** Leverage the error types defined in `internal/customerrors` for consistent error handling:
  - Use appropriate error types for different failure scenarios (authentication, validation, etc.)
  - Include contextual information in error messages to aid troubleshooting
  - Wrap API errors with additional context using the provided helper functions

- **Custom Types:** Utilize the custom types defined in `internal/customtypes` for specialized data handling:
  - Use custom Terraform schema types where appropriate
  - Leverage provided plan modifiers and validators for custom types
  - Follow the patterns established for marshaling/unmarshaling custom types

- **Validators:** Apply common validators from `github.com/hashicorp/terraform-plugin-framework-validators` package or `internal/validators` to ensure consistent validation logic:
  - Use provided validators for common validation requirements (UUID format, string length, etc.)
  - Chain validators for attributes that need multiple validation rules
  - Add resource-specific validation only when generic validators are insufficient

- **Helper Functions:** Make use of utility functions in `internal/helpers` to reduce duplication:
  - Use the helper functions for common tasks like state management and data conversion
  - Leverage the provided resource base types and embedded functionality
  - Follow established patterns for logging, attribute access, and diagnostics handling

### Best Practices

- **Method Scope:** Methods that are not used outside the namespace scope should be kept private (unexported).

- **API Interaction:**
  - Use the service-specific clients from the provider for all API calls.
  - Handle asynchronous operations with proper polling and timeouts.
  - Validate input values before sending API requests when needed.
  - Ensure that you always pass context `ctx` into long-running or asynchronous operations like API calls

- **Error Handling:**
  - Add context to API errors using the provider's error types from `internal/customerrors`.
  - Return detailed diagnostics with `resp.Diagnostics.AddError()` for user-friendly messages.
  - Distinguish between different error types (authentication, validation, not found, etc.).
  - Log API responses and errors at debug level using `tflog.Debug` for troubleshooting.

- **Request Context:**
  - Resources and Data Sources should call `ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)` and `defer exitContext()` near the beginning of any method from resource or datasource interfaces.

### Guidelines for Resources

#### Resource Structure and Interfaces

- Implement `resource.Resource` interface for all resources.
- Implement `resource.ResourceWithImportState` for resources supporting import.
- Implement `resource.ResourceWithValidateConfig` when custom validation is needed.
- Structure resources in a consistent pattern by ordering methods: `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`, `ImportState`.
- Embed `helpers.TypeInfo` in your resource struct to inherit standard functionality.
- Add required client fields to your resource struct to access APIs.

#### Resource Schema Definition

- Define complete schemas with proper attribute types (String, Int64, Bool, etc.).
- Mark attributes explicitly as `Required`, `Optional`, or `Computed`.
- Use `Computed: true` for server-generated fields like IDs.
- Use `Optional: true` with `Computed: true` for fields that can be specified or defaulted by the service.
- Apply `RequiresReplace` plan modifier to immutable attributes that necessitate resource recreation when changed.
- Apply `UseStateForUnknown` modifier to computed fields to prevent unnecessary diffs during planning.
- Include standard timeouts using `github.com/hashicorp/terraform-plugin-framework-timeouts`.
- Write clear `MarkdownDescription` for each attribute (do not use the deprecated `Description` field).

#### Resource State Management

- In `Create`, populate state with all resource attributes after successful creation.
- In `Read`, refresh the full state based on the current resource values from the API.
- Check for deleted resources in `Read` - when API returns 404, call `resp.State.RemoveResource(ctx)`.
- In `Update`, apply only the changed attributes and refresh state afterwards.
- Return early with appropriate diagnostics when operations cannot complete successfully.

#### Resource Validation

- Apply built-in validators from `github.com/hashicorp/terraform-plugin-framework-validators` for attribute constraints.
- Implement resource-level validation in the `ValidateConfig` method when validation involves multiple attributes.
- Add custom validators only when built-in validators are insufficient.
- Provide clear validation error messages that explain the specific constraint and how to fix it.

### Guidelines for Data Sources

#### Data Source Structure and Interfaces

- Implement the `datasource.DataSource` and `datasource.DataSourceWithConfigure` interfaces for all data sources.
- Order data source methods consistently: `Metadata`, `Schema`, `Configure`, `Read`.
- Embed `helpers.TypeInfo` in your data source struct to inherit standard functionality.
- Add required client fields to your data source struct to access APIs.
- Name factory functions as `New<DataSourceName>DataSource()` (e.g., `NewSolutionsDataSource`).

#### Data Source Schema Definition

- Mark all attributes as `Computed: true` since data sources are read-only by design.
- For optional filter parameters, use `Required: false` and `Optional: true`.
- Define nested schemas for complex return types using appropriate collection types:
  - Use `schema.ListNestedAttribute` for collections of objects like "environments", "applications", etc.
  - Use `schema.SingleNestedAttribute` for single complex objects.
- Only include Read timeouts in timeouts schema (omit Create, Update, Delete).
- Use `map[string]schema.Attribute` for schema attributes that allow extensible field sets.
- Include output-only fields that will assist users in identifying or using the data in further resources.
- For primary list attributes (e.g., "applications", "environments"), use the plural form as the attribute name.

#### Data Source Query Parameters

- For data sources that filter results, define explicit filter attributes:
  - Common patterns include `name`, `publisher_name`, `environment_id`, etc.
  - Keep filter attributes simple and intuitive based on how the target API implements filtering.
- Support sensible combinations of filter parameters that match Power Platform API capabilities.
- Document filter parameters with clear examples in the `MarkdownDescription`.

#### Data Source Read Implementation

- Parse all input filter parameters from state at the beginning of the Read method.
- Include context propagation: `ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)` with a matching `defer exitContext()`.
- Validate any required filter parameters and return appropriate diagnostic errors.
- Use the appropriate client method to retrieve data based on filter criteria.
- For empty results, set an empty list rather than returning an error.
- Transform API responses to data models using the appropriate conversion functions.
- Set all fields in the state model, even those that might be nil or empty.
- Log API calls using `tflog.Debug` statements to assist troubleshooting.
- For list-type data sources, return a consistent response structure even when results vary in size.

#### Testing Data Sources

- Test all supported filter combinations in unit tests.
- Verify that filtered results return the expected subset of data.
- Test edge cases like empty results, single results, and large result sets.
- For collection data sources, test accessing nested attributes and verify attribute counts.
- Ensure acceptance tests use non-destructive read-only operations.
- For data sources that return lists, test accessing list items with collection syntax.

#### Data Source Documentation and Examples

- Include a representative example in the `/examples/data-sources/{data_source_name}/` directory.
- For data sources with filter parameters, include examples showing different filtering options.
- Showcase practical use-cases with supporting resources when applicable.
- Describe the purpose of the data source clearly in the schema's `MarkdownDescription`.
- Link to relevant Power Platform documentation that explains the underlying API or service.

## Logging

Use the Terraform plugin logger (`tflog`) for logging within resource implementations.

- **Debug Level:** Add `tflog.Debug` statements to log major operations or key values for troubleshooting. Use debug logs liberally to aid in diagnosing issues during development or when users enable verbose logging.
- **Info Level:** Use `tflog.Info` sparingly, only for important events or high-level information that could be useful in normal logs. Too much info-level logging can clutter output, so prefer Debug for most internal details.
- **No Print/Printf:** Do not use `fmt.Println`/`log.Printf` or similar for logging. The `tflog` system ensures logs are structured and can be filtered by Terraform log level.
- **Sensitive Data:** Never log sensitive information (credentials, PII, etc.). Ensure that debug logs do not expose secrets or user data.
- **Request Context:** Do not trace the entry/exit of interface methods in resources or data source.  Instead use `EnterRequestContext` and `exitContext`

## Testing Best Practices

- **Unit Tests:** For each new resource or data source, write unit tests covering all operations and edge cases. Use the `jarcoal/httpmock` library (already in the project) to simulate HTTP API responses.
  - Register **mock responders** for every HTTP call that the Create, Read, Update, or Delete functions will make. Each test step should set up the expected API responses (e.g. mock the POST response for Create, GET for Read, etc.).
  - **Test Steps Lifecycle:** Structure unit tests in sequential steps to simulate resource lifecycle transitions:
    - **Step 1 (Create):** Call the resource's Create, then Read. Verify that after Create, the state read back includes all the created fields/attributes.
    - **Step 2 (Update):** Call Read (to get current state), then Update, then Read again. Ensure the first Read in this step matches the final state from the previous step, and the final Read reflects the updates applied.
    - **Step 3 (Delete):** Call Delete, then Read. After deletion, the final Read should return a "not found" error (e.g. 404) indicating the resource is gone.
    - If the resource supports import, write a dedicated test (single step) that calls the Read (or Import) with a given `ImportStateId` and verifies Terraform state import logic.
  - Include negative test cases: simulate API errors (like 403 Forbidden or 500 Internal Server Error) and ensure the provider surfaces appropriate errors. Also test validation logic (e.g., providing an invalid parameter returns an error).
  - Place JSON fixtures for mock responses in the appropriate test data directory (e.g. `internal/services/<service>/test/<resource>/<scenario>/response.json`). **Do not use real customer data** in tests – anonymize any IDs or personal info in your dummy data.
  - Name unit test functions with the `TestUnit` prefix as mentioned, and keep them in a `_test.go` file using the `<package>_test` package name.
  - All the JSON response for unit tests should be stored in .json files:
    - Files should be placed in a folder with a name corresponding to the Unit Test that is being used. Folder name should omit `UnitTest` in its name.
    - Each Unit Test folder with .json files should be stored at `services\{service_name}\test\resource` or `services\{service_name}\test\datasource` with all other resource and/or datasource .go files.
    - The .json file name should consist of the mock request method (`get`, `post`, `delete`) followed by `_` and name of the returned mock object name or action.
    - The file names have to be sensible without empty spaces and special characters.

- **Acceptance Tests:** Add acceptance tests for any new resource covering the same scenarios as unit tests, but against real Power Platform resources. These tests live in files with the `TestAcc...` prefix and require real credentials.
  - Wrap any acceptance test with appropriate pre-check functions and environment variable checks so it skips if not configured.
  - Ensure each acceptance test cleans up after itself. Use `CheckDestroy` functions to verify that resources are actually deleted in Azure/Power Platform after the test run.
  - Keep acceptance tests focused and isolated (use separate environment or resource names to avoid conflicts).

- **Test Coverage:** Aim for **at least 80%** code coverage for unit tests on new code. `make unittest` will return a coverage score by service and overall. Focus on the service that is currently being worked on when adding tests to improve coverage.

- **Examples and Documentation:** Whenever a new resource or data source is added, provide an example configuration under the `/examples` directory to demonstrate usage. This helps both in documentation and in manually verifying the resource behavior. After implementing and testing, run `make userdocs` to update the documentation in `/docs` from your schema comments.
