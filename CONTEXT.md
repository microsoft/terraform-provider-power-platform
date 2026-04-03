# Context

## Repository purpose

- Terraform provider for Microsoft Power Platform and Dataverse.
- Built with the Terraform Plugin Framework.
- Manages environments, solutions, publishers, users, connections, and related Power Platform resources through Power Platform admin APIs, BAPI, and Dataverse Web API.

## Current operational branches

- `codex/preview-integration`
  - working integration branch used to build forked preview binaries for Azure DevOps pipeline consumption
  - carries in-flight work around git integration resources, unmanaged solutions, and publishers
- `codex/fix-environment-security-group-update`
  - clean bugfix branch cut from `upstream/main`
  - fixes environment update behavior so non-Developer environment updates preserve the planned `dataverse.security_group_id`

## Release path used by shared Azure Pipelines

- Forked preview binaries are published from GitHub workflow:
  - `.github/workflows/fork_provider_binaries.yml`
- Preview assets are released as GitHub prereleases with tags like:
  - `fork-v4.1.1-adam-preview.6`
- Shared Azure DevOps pipeline templates download those prerelease zip assets directly.

## Bugfix currently in progress

- Existing environment resource update logic in `internal/services/environment/resource_environment.go`
  - function `updateExistingDataverse(...)`
  - was nulling out `LinkedEnvironmentMetadata.SecurityGroupId` for non-Developer environments during update
- Result:
  - Terraform planned a security-group change
  - provider update request dropped that value
  - post-apply read still returned the old group id
  - Terraform failed with `Provider produced inconsistent result after apply`
- Clean fix:
  - preserve the planned `dataverse.security_group_id` during non-Developer updates
- Verified locally on the clean bugfix branch with:
  - `go test ./internal/services/environment/...`

## Notes

- This repo uses `main` upstream and should keep that primary branch naming.
- Merge strategy should remain fast-forward or merge commit only; never rebase.
