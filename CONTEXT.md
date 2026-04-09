# Project Context

## Overview

- Repository: `terraform-provider-power-platform`
- Purpose: Terraform provider for Microsoft Power Platform, implemented with the Terraform Plugin Framework.
- Entry point: `main.go`
- Provider registration: `internal/provider/provider.go`
- Main implementation area: `internal/services/<service_name>`

## Current Structure

- `internal/provider`: provider schema, configuration, client wiring, resource and data source registration.
- `internal/services`: one folder per resource or data source area. Each service typically contains:
  - `resource_*.go` and/or `datasource_*.go`
  - `api_*.go` for API client operations
  - `models.go` and `dto.go`
  - `*_test.go`
  - `tests/...` fixtures for mocked HTTP responses
- `docs`: generated provider docs.
- `examples`: sample Terraform configurations.
- `devdocs`: contributor guidance for schema, testing, security, logging, and release flow.

## Development Conventions Observed

- Resources use Terraform Plugin Framework schema types and plan modifiers.
- Tests are split between:
  - unit tests using `httpmock`
  - acceptance tests using real provider flows
- Docs are generated with `tfplugindocs`.
- The repository prefers strongly typed schemas where the API surface is stable enough to justify it.

## Solution Resource Notes

- Current resource: `powerplatform_solution`
- Main files:
  - `internal/services/solution/resource_solution.go`
  - `internal/services/solution/api_solution.go`
  - `internal/services/solution/dto.go`
  - `internal/services/solution/models.go`
  - `docs/resources/solution.md`
- Current behavior:
  - Reads `solution_file` from a local filesystem path only.
  - Optionally reads a PAC-style `settings_file` JSON from a local filesystem path.
  - Stages the solution zip with `StageSolution`.
  - Imports it with `ImportSolutionAsync`.
  - Maps settings JSON into Dataverse import `ComponentParameters`.
  - Stores `solution_file_checksum` and `settings_file_checksum` in state.
  - Uses file checksum changes to trigger update planning behavior.
  - Reads back `display_name`, `is_managed`, and `solution_version` from Dataverse after import.
- Important limitation:
  - The current interface is file-driven, not intent-driven.
  - Drift/update detection is tied to local file content hash, not declared solution version.
  - `is_managed` is only observed after import, not enforced as an invariant up front.

## Current Investigation

- Current feature branch focus: implement a new `powerplatform_managed_solution` resource derived from `powerplatform_solution`.
- Final shape chosen for this branch:
  - constrain the resource to managed solutions only
  - support a local filesystem path or remote URL as the package source
  - use explicit `unique_name` and `version` inputs for identity and update planning
  - expose `connection_references` as a `map(string)` from logical name to connection id
  - do not manage environment variables through the resource
  - inspect the solution package at apply time to validate metadata and deployment prerequisites before import

## Managed Solution Design Notes

- Plan-time drift/update detection is identity-driven, not file-hash-driven.
- `source` is treated as an apply-time delivery input rather than desired-state identity.
- Source-only changes are intentionally ignored when `environment_id`, `unique_name`, `version`, and `connection_references` are unchanged.
- Package inspection is done directly from the zip contents:
  - `solution.xml` provides unique name, version, managed flag, and dependency metadata
  - `customizations.xml` provides connection reference declarations
  - `environmentvariabledefinitions/*/environmentvariabledefinition.xml` provides environment variable references and defaults
  - `environmentvariablevalues.json` is rejected if it carries packaged current values
- Environment variables are treated as environment-owned configuration:
  - if the package defines a default, the import may rely on it
  - if the package does not define a default, the variable must already exist in the target environment
  - the managed solution resource does not attempt to set or update environment variable values
- Solution dependencies are parsed from `ImportExportXml.SolutionManifest.MissingDependencies.*.MissingDependency.Required[type="66"]/@solution`.
- Duplicate dependency entries are collapsed to the highest required version before validation against installed environment solutions.

## Current State

- Added a new `powerplatform_managed_solution` resource under `internal/services/managed_solution`.
- The resource is now registered in `internal/provider/provider.go` and `internal/provider/provider_test.go`.
- Current managed solution behavior implemented:
  - desired identity is `environment_id + unique_name + version`
  - `source` supports either local `path` or remote `url`
  - `source` is treated as an apply-time delivery input; source-only changes are ignored in `ModifyPlan`
  - package metadata is parsed directly from `solution.xml`
  - package connection references are parsed from `customizations.xml`
  - package environment variable definitions are parsed from `environmentvariabledefinitions/*/environmentvariabledefinition.xml`
  - packaged environment variable current values are detected from `environmentvariablevalues.json` and rejected
  - dependency requirements are parsed from `MissingDependencies -> Required[@type="66"] -> @solution`
  - duplicate dependency requirements are collapsed to the highest required version
  - target environment solutions are checked to ensure each dependency is installed at an equal or higher version
  - package connection references must be fully mapped via `connection_references`
  - environment variables are not managed by the resource; they must already be satisfiable via default or existing environment value
- Added:
  - unit tests for package inspection, dependency validation, environment variable validation, and resource create happy path
  - example usage under `examples/resources/powerplatform_managed_solution`
  - draft docs under `docs/resources/managed_solution.md`
  - changelog fragment under `.changes/unreleased`
- Verification completed on the scoped change set:
  - `go test ./internal/services/managed_solution ./internal/provider`
  - `golangci-lint run ./internal/services/managed_solution/... ./internal/provider/...`
- Repo-wide `make precommit` was not completed in this branch because it fails on pre-existing lint debt outside the managed solution change set, and the user explicitly asked not to broaden the PR scope.

## Design Constraints Worth Preserving

- `environment_id` should remain required and should still force replacement if changed.
- The resource should fail loudly when the environment has no Dataverse.
- The managed-only invariant should be explicit and enforced, not treated as best effort.
- The new schema should align with existing Plugin Framework patterns already used in the repository.
