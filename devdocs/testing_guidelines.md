## Testing Guidelines

These guidelines detail best practices for testing the Terraform Provider for Microsoft Power Platform, based on official [Terraform Plugin Testing guidance](https://developer.hashicorp.com/terraform/plugin/testing). They cover comprehensive testing strategies, including unit tests, acceptance tests, and various established testing patterns. Additionally, the document outlines specific approaches to ensure consistent test coverage, robust validation of provider behaviors, and effective handling of different authentication and cloud contexts.

### Testing Patterns

Following established testing patterns ensures consistency and reliability in Terraform provider tests. These patterns help in designing robust test cases that improve maintainability and test quality.

- **Happy Path Testing**: Verify basic create, read, update, and delete operations.
- **Error Handling Validation**: Include test cases confirming proper error messages for incorrect inputs, invalid configurations, or API errors.
- **Property-Based Testing**: Validate multiple variations of inputs to ensure broad coverage. These tests are especially important in scenarios where properties may be required together, are mutually exclusive, or involve optional values, maximum-length strings, and boundary conditions, ensuring the provider correctly handles complex validation logic and edge cases.
- **Idempotency Testing**: Ensure running `terraform apply` twice results in no changes.
- **Import Verification**: Include `ImportState: true` and `ImportStateVerify: true` in test cases to fully validate import functionality.
- **Parallel Test Execution**: Design tests to run concurrently without resource conflicts. Avoid reusing resource names across tests to prevent naming collisions. Using randomized names can be helpful to ensure uniqueness and avoid conflicts.
- **Authentication Context**: Validate provider behavior across different authentication contexts, such as service principal versus user-based authentication. APIs may sometimes behave differently depending on the user context, making it important to test these variations explicitly.
- **Cloud Environment Testing**: Test across various cloud environments (e.g., public clouds vs. government clouds) to ensure compatibility and robustness. Refer to the [Terraform Provider documentation on non-public clouds](https://registry.terraform.io/providers/microsoft/power-platform/latest/docs/guides/nonpublic_clouds) for guidance on configuring tests in different environments.

For more information, refer to [Terraform Plugin Testing Patterns](https://developer.hashicorp.com/terraform/plugin/testing/testing-patterns).

### Unit Tests

Unit tests are critical for fast feedback. They should cover the logic of CRUD operations by mocking out HTTP calls to the Power Platform APIs. All unit tests for a given resource or data source are located in the `/internal/<resource_or_datasource>_test.go` file. JSON for mocked API responses can often be captured using network traces from browser debugging tools. Utilizing features in the Power Platform Admin Center can provide examples of how the Power Platform UI interacts with these APIs. While a test-first approach is not mandatory, it can be beneficial, particularly when using mocked data.&#x20;

- This provider uses the `httpmock` library to simulate API responses. This approach was chosen because it allows us to also test serialization and deserialization of API requests and responses.
- Common mocks (e.g., creating a Power Platform Environment) are encapsulated in helper functions for reuse, such as the `ActivateEnvironmentHttpMocks` function.
- When adding new unit tests, place static JSON responses in `internal/services/<service>/test/<resource_or_datasource>/<test_name>/...`.
- **Do not include real personal data** – anonymize IDs, emails, tenant IDs, phone numbers, etc.
- Cover **create, read, update, and delete** behaviors, along with conditional logic (e.g., forced recreation).
- Include **negative tests** for API errors (e.g., 403 Forbidden), invalid inputs, and boundary conditions.
- When creating mocked JSON responses, you can reuse existing ones by duplicating them into your `<test_name>` folder.

#### Running Unit Tests

To run all unit tests:

```bash
make unittest
```

To run a single unit test:

```bash
TF_ACC=0 go test -v ./... -run TestUnit<test_name>
```

### Acceptance Tests (Integration Tests)

Every unit test covering a new feature or fix should have a corresponding acceptance test validating the same use case against real infrastructure. Acceptance tests ensure provider correctness in actual Power Platform environments.

- Acceptance tests (files with `TestAcc...`) call real APIs and require valid credentials and a test tenant.
- **Tests create real resources** – ensure proper cleanup after execution.
  - Use **CheckDestroy** to verify resource deletion post `terraform destroy`.

#### Test Pre-checks

Acceptance tests require authentication and appropriate test environment setup. This section is particularly useful for testing different authentication styles, ensuring the provider supports various authentication mechanisms reliably.

- The test suite includes **pre-checks** (`testAccPreCheck(t)`) to validate required environment variables before running tests.
- Ensure necessary credentials are set before executing acceptance tests.
- Missing variables will result in skipped tests to prevent erroneous failures.

#### Running Acceptance Tests

To run all acceptance tests:

```bash
make acctest
```

To run a single acceptance test:

```bash
TF_ACC=1 go test -v ./... -run TestAcc<test_name>
```

### Test Coverage

The project expects **high test coverage (80% or more)** for new contributions.

- Measure coverage using `go test -cover` or reviewing CI results.
- The project uses **Codecov.io** to monitor test coverage in CI. View the current coverage report at the [Codecov Dashboard](https://app.codecov.io/gh/microsoft/terraform-provider-power-platform).
- Code coverage on Codecov.io is updated by the CI workflow.
- Focus on core logic rather than Terraform framework internals.

> [!NOTE] The tests require permissions on specific folders. These permissions are assigned when creating your container. If you encounter permission issues, rebuild your development container or run the following commands to assign necessary permissions:

```bash
sudo chown -R vscode /workspaces/terraform-provider-power-platform/
sudo chown -R vscode /go/pkg
```

