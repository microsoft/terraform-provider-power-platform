# Project Context

## Overview

- Repository: `github.com/microsoft/terraform-provider-power-platform`
- Purpose: Terraform provider for administering Microsoft Power Platform resources and selected Dataverse-backed capabilities.
- Language/runtime: Go, using `hashicorp/terraform-plugin-framework`.
- Current branch in this worktree: `feature/application_user`

## Repository Structure

- `internal/provider`: provider schema, auth/config wiring, resource/data source registration.
- `internal/services/<service>`: service-scoped resources, data sources, API clients, DTOs, tests.
- `internal/api`: shared HTTP/auth client plumbing.
- `docs` and `examples`: generated docs and user-facing examples.

## Current Workstream

- We are adding a new first-class Dataverse application-user resource modeled on the environment-scoped Dataverse `systemusers` API, not the BAP `addAppUser` shortcut.
- The target resource shape is:
  - `powerplatform_environment_application_user`
  - required: `environment_id`, `application_id`
  - optional: `business_unit_id`, `security_roles`, `timeouts`
  - computed: `id`, `system_user_id`, resolved role assignments
- Current design decisions:
  - create application users by `POST /api/data/v9.2/systemusers` with:
    - `applicationid`
    - `accessmode = "4"` (hard-coded for now)
    - `businessunitid@odata.bind`
    - `isdisabled = false`
  - resolve requested role names to environment-specific `roleid` values and store those resolved IDs in state for nested assignment identity
  - assign/remove roles via `systemuserroles_association/$ref`
  - disable before delete, then delete the `systemuser`
- We are not using `addAppUser` for the new resource because that path is a narrower bootstrap/admin shortcut and does not match the generic application-user flow observed in the admin portal.

## Current Implementation Status

- Implemented `powerplatform_environment_application_user` under `internal/services/application`.
- Registered the resource in the provider and provider tests.
- Added user-facing docs and an example configuration.
- The resource currently:
  - creates Dataverse application users through `POST /api/data/v9.2/systemusers`
  - hard-codes `accessmode = 4`
  - defaults `business_unit_id` to the root business unit when omitted
  - resolves configured `security_roles` by exact role name within the selected business unit
  - stores resolved environment-specific role IDs in computed state via `resolved_security_roles`
  - disables then deletes the backing `systemuser` on destroy
  - fails create with an import-oriented diagnostic if the application user already exists remotely
  - rolls back a newly created remote application user if later create steps fail before Terraform state is written
- Verification status:
  - `go build ./...` passes
  - `go test ./internal/services/application ./internal/provider` passes
  - a broader `go test ./...` sweep was running at the time of this update

## Existing Related Resources

- `powerplatform_environment_application_admin` exists today and uses BAP `addAppUser`.
- `powerplatform_user` uses Dataverse `systemusers` plus `systemuserroles_association` for normal users.
- The `powerplatform_data_record` examples already demonstrate the desired application-user CRUD shape using raw Dataverse records.

## Open Design Constraints

- Role names are the right user-facing input, but the provider must fail loudly on ambiguous or missing name resolution.
- Resolved environment-specific `roleid` values should be persisted in state for reconciliation.
- `business_unit_id` should default to the root business unit when omitted.
