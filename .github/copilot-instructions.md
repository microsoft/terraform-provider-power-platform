# GitHub Copilot Custom Instructions – Terraform Provider Power Platform

These instructions guide GitHub Copilot to follow our project’s conventions and best practices when suggesting code. They cover how to format code, name resources and attributes, structure implementations, and write tests in this repository.

## Development Setup & Workflow

- Use the provided **Makefile** commands for all build and test tasks:
  - `make install` to compile the provider code.
  - `make lint` to run linters and ensure code style compliance.
  - `make unittest` to run all unit tests (optionally use `TEST=<prefix>` to run tests matching a name prefix, e.g. `make unittest TEST=Environment` to run tests named with that prefix). This filters tests by regex `^(TestAcc|TestUnit)<prefix>`.
  - `make acctest TEST=<prefix>` to run acceptance tests (integration tests) matching a prefix. Always provide a specific test prefix to limit scope, and run these tests **only with user consent** (they run against real cloud resources). Note that `make acctest` automatically sets `TF_ACC=1` (no need to set it manually).
  - `make userdocs` to regenerate documentation
  - `make precommit` to run all checks once code is ready to commit
- Always run the above `make` commands from the repository root (e.g. in the `/workspaces/terraform-provider-power-platform` directory).
- **Never run** `terraform init` inside the provider repo. Terraform is only used in examples or tests; initializing in the provider directory is not needed and may cause conflicts.
- Do not manually edit files under the `/docs` folder. These files are auto-generated from the schema `MarkdownDescription` attributes. Instead, update schema's `MarkdownDescription` in code and run `make userdocs` to regenerate documentation.
- To try out an example configuration, navigate to its directory under `/examples` and run `terraform apply -auto-approve` (ensure you’ve built the provider and set it in your Terraform plugins path beforehand).

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

## Naming Conventions

- **Resource and Data Source Names:** Follow the existing naming pattern of prefixing with `powerplatform_`. For example, an environment resource is named `powerplatform_environment`. Use lowercase with underscores for Terraform resource/data names.
- **Attribute Naming:** Name resource attributes to match Power Platform terminology. Prefer the modern, user-friendly terms used in the current Power Platform API/UX/[Official Documentation](https://learn.microsoft.com/en-us/power-platform/admin/) over deprecated names. Keep names concise but descriptive of their function in the resource.
- **Test Function Naming:** Name test functions with a prefix indicating their type. **Acceptance test** functions should start with `TestAcc` and **unit test** functions with `TestUnit` (this allows filtering tests by type). Also, name test files’ package with a `_test` suffix (e.g. `package environment_test`) to ensure tests access the provider only via its public interface.
- Define all DTO (Data Transfer Object) structures in `dto.go` with a `Dto` suffix (e.g., `EnvironmentDto`).
- Implement conversion functions named exactly as `convertDtoToModel` and `convertModelToDto` in `models.go`.
- When implementing a client factory, name it `New<Service>Client` (e.g., `NewSolutionClient`).

## Comments and Documentation

- Include the appropriate license header at the top of all new code files
- Write Go comments only on exported functions, types, and methods to explain their purpose, parameters, and return values when it adds clarity.
- Focus comments on **why** something is done if it’s not obvious from the code.
- Avoid redundant comments that just restate the code or don’t provide additional insight.
- Remember that schema's field `MarkdownDescription` will become user docs, so keep those clear and up-to-date.

## Code Organization and Implementation Guidelines

Follow these guidelines to maintain consistency with the established project patterns and ensure your code integrates seamlessly with the existing codebase.

- When defining resource or data source schema, **always use** the `MarkdownDescription` field for documentation. Do **not** use the deprecated `Description` field. Markdown descriptions will be used to auto-generate docs, so make them clear and user-friendly, and include links to topics in the [official Power Platform docs](https://learn.microsoft.com/en-us/power-platform/admin/) when helpful.

- **Terraform Plugin Framework:** Use [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) patterns for implementing resources and data sources.
  - Avoid using legacy Terraform SDK constructs.
  - Work with `schema.Schema` fields
  - Return `diag.Diagnostics` in resource functions as per Plugin Framework conventions.
- **Azure Identity Client:** Use [Azure Identity Client Module for Go](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity) for implementing authentication methods

### Data Source Specific Guidelines

### Resource Specific Guidelines

- **Resource Structure:** Implement all required methods: `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`, and when applicable, `ImportState`.
  - If the resource supports import, implement an `ImportState` function that properly parses the import ID.
  - Use `RequiresReplace` in the schema for any attribute changes that require recreating the resource (e.g. immutable keys like an environment’s name or ID, if changing them necessitates a new resource).
  - Use plan modifiers like `UseStateForUnknown` for computed fields (e.g., IDs) to avoid unnecessary diffs.

- **Schema Definitions:**
  - **Required/Optional Settings:** Clearly mark attributes as `Required`, `Optional`, or `Computed` as appropriate. Use `Computed: true` for fields like IDs that are set by the server, and use `Optional + Computed` for fields that can be set or defaulted by the service (so Terraform doesn’t show spurious diffs).

- **API Interactions:** Use the appropriate Power Platform SDK/client provided by the project to make API calls. Handle errors gracefully:
  - Return detailed error diagnostics for API call failures.
  - If an API 404 (Not Found) is encountered during a Read, translate it to a "resource not found" state so Terraform can handle resource removal.
  - Avoid hard-coding values; use constants or enums defined in the SDK for things like statuses or types when available.
- **Helper Functions & Structure:** Keep resource code organized:
  - Place any data model definitions or complex object translation logic in a **`models.go`** file within the resource’s package or service module. Define structs for request/response payloads or any mapping logic there.
  - Encapsulate repetitive tasks (e.g. common API error handling, or converting between Terraform schema values and API DTO objects) in helper functions. If a helper is only used internally in one package, make it unexported (private).
  - Functions that convert between API data transfer objects (DTOs) and internal models should be named consistently as `convertDtoToModel` and `convertModelToDto` and kept in the `models.go` file. This ensures a predictable pattern across resources for data transformations.
  - Any function not used outside its package should be unexported (start with a lowercase letter) to limit scope.

1. Define a resource struct that embeds `helpers.TypeInfo` and includes any necessary client fields.
2. 
3. Keep all unexported helper functions with the resource implementation.
4. For each resource and data source, create a new function named `New<ResourceName>Resource` or `New<DataSourceName>DataSource` that returns the appropriate type.

## Logging

 Use the Terraform plugin logger (`tflog`) for logging within resource implementations.

- **Debug Level:** Add `tflog.Debug` statements to log entry/exit of major operations or key values for troubleshooting. Use debug logs liberally to aid in diagnosing issues during development or when users enable verbose logging.
- **Info Level:** Use `tflog.Info` sparingly, only for important events or high-level information that could be useful in normal logs. Too much info-level logging can clutter output, so prefer Debug for most internal details.
- **No Print/Printf:** Do not use `fmt.Println`/`log.Printf` or similar for logging. The `tflog` system ensures logs are structured and can be filtered by Terraform log level.
- **Sensitive Data:** Never log sensitive information (credentials, PII, etc.). Ensure that debug logs do not expose secrets or user data.

## Testing Best Practices

- **Unit Tests:** For each new resource or data source, write unit tests covering all operations and edge cases. Use the `jarcoal/httpmock` library (already in the project) to simulate HTTP API responses.
  - Register **mock responders** for every HTTP call that the Create, Read, Update, or Delete functions will make. Each test step should set up the expected API responses (e.g. mock the POST response for Create, GET for Read, etc.).
  - **Test Steps Lifecycle:** Structure unit tests in sequential steps to simulate resource lifecycle transitions:
    - **Step 1 (Create):** Call the resource’s Create, then Read. Verify that after Create, the state read back includes all the created fields/attributes.
    - **Step 2 (Update):** Call Read (to get current state), then Update, then Read again. Ensure the first Read in this step matches the final state from the previous step, and the final Read reflects the updates applied.
    - **Step 3 (Delete):** Call Delete, then Read. After deletion, the final Read should return a “not found” error (e.g. 404) indicating the resource is gone.
    - If the resource supports import, write a dedicated test (single step) that calls the Read (or Import) with a given `ImportStateId` and verifies Terraform state import logic.
  - Include negative test cases: simulate API errors (like 403 Forbidden or 500 Internal Server Error) and ensure the provider surfaces appropriate errors. Also test validation logic (e.g., providing an invalid parameter returns an error).
  - Place JSON fixtures for mock responses in the appropriate test data directory (e.g. `internal/services/<service>/test/<resource>/<scenario>/response.json`). **Do not use real customer data** in tests – anonymize any IDs or personal info in your dummy data.
  - Name unit test functions with the `TestUnit` prefix as mentioned, and keep them in a `_test.go` file using the `<package>_test` package name.
- **Acceptance Tests:** Add acceptance tests for any new resource covering the same scenarios as unit tests, but against real Power Platform resources. These tests live in files with the `TestAcc...` prefix and require real credentials.
  - Wrap any acceptance test with appropriate pre-check functions and environment variable checks so it skips if not configured.
  - Ensure each acceptance test cleans up after itself. Use `CheckDestroy` functions to verify that resources are actually deleted in Azure/Power Platform after the test run.
  - Keep acceptance tests focused and isolated (use separate environment or resource names to avoid conflicts).
- **Test Coverage:** Aim for **at least 80%** code coverage for unit tests on new code. `make unittest` will return a coverage score by service and overall.  Focus on the service that is currently being worked on when adding tests to improve coverage.
- **Examples and Documentation:** Whenever a new resource or data source is added, provide an example configuration under the `/examples` directory to demonstrate usage. This helps both in documentation and in manually verifying the resource behavior. After implementing and testing, run `make userdocs` to update the documentation in `/docs` from your schema comments.

By following these guidelines, Copilot’s suggestions should align with the project’s style and help contributors produce high-quality, consistent code. Always consider existing patterns in the repository—when in doubt, review similar resources or tests for reference and keep the new code idiomatic to the project’s practices.

## Examples

- For import script examples:
  - Use comments to explain the import syntax
  - Include placeholder values (e.g., `00000000-0000-0000-0000-000000000000`) for IDs
  - Show the proper resource name format (e.g., `powerplatform_environment.example`)

- Keep example code simple and focused on demonstrating one clear use case per file

# Terraform Provider Power Platform Instructions


Do not direcly edit the files under `/docs` because they are auto-generated from `MarkdownDescription` on schemas using `make userdocs`
test files should have `_test` appended to their package name
To run an example, `cd` to its working directory and run `terraform apply -auto-approve`

When creating Schema for the resource or datasource, use MarkdownDescription and never use Description attribute
Methods that are not used outside the namespace scope, should be kept private.
Helper methods that covert DTO to model and model to DTO should be in models.go file
The DTO structures should always have `Dto` suffix and be in models.go file
Use tflog.Debug for logging unless there is something really important (tflog.Info) or an error/warning
Comments on methods should provide information about how to use it, its parameters, and expected results. Omit comments that don't substantially improve the readability of the code.
The functions to convert Dto and model objects should be always named `convertDtoToModel` and `convertModelToDto` and should be placed in models.go file

## Testing

When writing unit tests for resources you must register mock responders for every step of the process:

- Test steps will call the `Create`, `Read`, `Update`, and `Delete` methods in the resource.  All the API calls made in those functions need to be mocked for each time the operation is called.
- First test step will call `Create` then `Read` methods. Test expects that JSON from read include the changes applied in create
- Subsequent test steps will call `Read`, then `Update`, then `Read`. First read should match what was read at the end of previous step.  The final read should include the changes applied in the update step.
- Steps that delete a resource (or omit a previously configured resource from config) will call `Delete` then `Read`.  Unless otherwise specified the final read will return a 404
- Tests that import state should be a single test step and only call `Read` using the `ImportStateId`
- All the JSON response for unit test should be stored in .json files.
  - Files should be placed in a folder with a name coresponding to the Unit Test that is being used. Folder name should ommit `UnitTest` in its name.
  - Each Unit Test folder with .json files be stored at `services\{service_name}\test\resource` or `services\{service_name}\datasource` with all other resource and/or datasource .go files.
  - The .json file name should consist of the mock request method (`get`, 'post', `delete`) followed by `_` and name of the returned mock object name or action.
  - The file names have to be sensible without empty spaces and special characters.
