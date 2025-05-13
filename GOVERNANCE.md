# Governance

## Purpose & Scope
The purpose of this document is to show how the Terraform Provider for Power Platform makes technical and process decisions.

## Roles
- **Product owner** – manages the roadmap and release cadence, gathers stakeholder input.
- **Maintainers** – review/merge PRs, run CI/CD, approve breaking changes.
- **Contributors** – anyone submitting issues or PRs.

## Decision Process
Complex proposals start as GitHub Discussions to seek community consensus. When consensus emerges, a PR implements the change.

## Pull-request Workflow
- Each request needs at least one maintainer review + all CI checks green.
- Breaking changes require explicit approval from a majority of maintainers.
- If a dispute remains unresolved, a simple majority vote of maintainers breaks the tie.

## Conflict Resolution
1. Discuss in the PR/issue or Discussion thread.
2. If still unresolved, hold a maintainers’ vote; majority decision is final.

## Adding or Removing Maintainers
Nomination via Discussion → majority maintainer vote → update CODEOWNERS.

## Release Process
The maintainer creates a signed tag via GoReleaser; product owner coordinates release notes and roadmap updates.

## Amending Governance
Create a PR with proposed changes → majority maintainer approval + acknowledgment from the product owner.
