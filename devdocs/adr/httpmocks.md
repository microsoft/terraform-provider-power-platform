# ADR: Use of httpmock for HTTP Request Mocking in Unit Tests

## Status

Accepted

## Context

Our Terraform provider interacts extensively with external HTTP APIs. Reliable and isolated unit tests require mocking these HTTP interactions to avoid dependencies on external services, ensure test determinism, and improve test execution speed.

Previously, unit tests loaded data directly into data transfer objects (DTOs), bypassing actual HTTP request-response cycles. This approach had several limitations:

- Wire serialization and deserialization logic was not exercised, potentially hiding serialization-related bugs.
- Difficulty in simulating specific HTTP responses, error conditions, and network failures.
- Inability to test HTTP client behavior comprehensively, including handling of HTTP status codes, headers, and timeouts.

Unit tests never invoked external APIs over the network, but the approach still lacked realism and robustness.

## Decision

We will adopt the [httpmock](https://pkg.go.dev/github.com/jarcoal/httpmock) library for mocking HTTP requests in our unit tests. Httpmock provides a straightforward and robust mechanism to intercept HTTP requests made by Go's standard `net/http` client and return predefined responses.

Key reasons for selecting httpmock:

- Simple integration with existing Go HTTP clients.
- Ability to define precise request-response mappings.
- Support for simulating various HTTP scenarios, including errors and timeouts.
- Active community support and clear documentation.

## Consequences

### Positive

- Improved test reliability and determinism.
- Faster test execution by eliminating external HTTP calls.
- Easier simulation of edge cases, error conditions, and network failures.
- Comprehensive testing of serialization and deserialization logic, HTTP status handling, and client behavior.
- Reduced complexity and maintenance overhead compared to custom mocks.

### Negative

- Additional dependency introduced into the project.
- Developers must familiarize themselves with the httpmock API.

## Implementation Guidelines

- Include httpmock as a test dependency in the project.
- Enable httpmock at the beginning of each test requiring HTTP mocking and disable it after the test completes.
- Clearly document mocked HTTP interactions within each test case.
- Create Acceptance Tests that mirror unit tests to detect changes in external APIs.
