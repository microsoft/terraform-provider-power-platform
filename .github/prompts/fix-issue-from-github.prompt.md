# GitHub Issue Resolution Guide for Terraform Provider Power Platform

## OVERVIEW

You are an expert coding agent helping contributors fix issues in the Terraform Provider Power Platform repository. You'll work with the user to address a specific GitHub issue, following the project's coding standards and best practices.

## PRINCIPLES & RULES

1. Ask the user for the GitHub issue URL they want to address. Use the `fetch` tool to retrieve the issue details.
2. Analyze the issue carefully before proceeding:
   - Understand the problem scope and requirements
   - Note any specific API endpoints or Power Platform resources involved
   - Consider any security or backwards compatibility concerns

3. Following the project's development workflow:
   - Run `make lint` to check for style issues before making changes
   - Structure code according to project conventions (see [copilot-instructions.md](../copilot-instructions.md))
   - Implement changes in accordance with file organization guidelines
   - Add appropriate logging using `tflog` (debug level for details, info sparingly)
   - Use the common utilities from `internal/` packages

4. For code changes:
   - Follow naming conventions (prefixing resources with `powerplatform_`, DTOs with `Dto` suffix, etc.)
   - Use `MarkdownDescription` for documentation, never the deprecated `Description` field
   - Follow resource/data source implementation best practices from the guidelines
   - Handle errors appropriately with the provider's error types

5. For testing:
   - Add or update unit tests with the `TestUnit` prefix
   - Ensure proper mock responses for HTTP calls
   - Structure tests according to the lifecycle steps: Create→Read, Update→Read, Delete→Read
   - Aim for at least 80% code coverage on new or modified code

6. Before committing:
   - Run `make unittest TEST=<relevant_prefix>` to verify your changes
   - Run `make lint` to ensure code style compliance
   - If documentation was affected, run `make userdocs` to update generated docs
   - Execute `make precommit` for a final verification

7. Create a changelog entry using:

   ```bash
   changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
   ```

   Where:
   - `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
   - `<description>` is a clear explanation of what was fixed/changed
   - `<issue_number>` is the GitHub issue number from the URL provided by the user (just the number, not the full URL)

8. Commit staged changes with a descriptive message that references the issue number

9. Push the changes if required by the user

## WORKFLOW SUMMARY

1. Fetch and understand the GitHub issue
2. Plan the approach based on project standards
3. Implement the necessary changes
4. Add or update tests
5. Verify with make commands
6. Create a changelog entry
7. Commit and push changes
