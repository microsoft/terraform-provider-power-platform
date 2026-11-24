# Agent instructions

These instructions apply to every coding agent (GitHub Copilot, Claude Code, Codex and others) that changes this repository.
Detailed conventions live in [.github/copilot-instructions.md](.github/copilot-instructions.md), [DEVELOPER.md](DEVELOPER.md) and [devdocs/testing_guidelines.md](devdocs/testing_guidelines.md).

## How every change is tested

The agent that makes a change runs all three layers itself before it opens or updates a pull request:

1. **Unit tests**: `make unittest TEST=<prefix>` for the tests the change touches, then `make unittest` for the whole suite.
   A fix starts with a unit test that fails before the fix and passes after it.
2. **Acceptance tests**: `make acctest TEST=<prefix>` for the affected resources and data sources.
   Always pass a specific `TEST=` prefix; never run the whole acceptance suite.
3. **Live infrastructure check**: build the provider locally (`make install`), point Terraform at it with `dev_overrides`, and apply a real configuration that exercises the change against the test tenant.
   For a bug, reproduce it this way before the fix and confirm it is gone after the fix.

Live runs use the test tenant credentials already configured for acceptance tests, never a production tenant.
Name every resource you create with the test name prefix, never modify resources you did not create, and confirm that everything you created was deleted after each run.
The one exception is a singleton resource such as `powerplatform_tenant_settings`, whose create updates existing settings and whose delete restores them: capture its current state as a baseline before the live check, and after the run verify that the state matches the baseline again.
If the baseline cannot be captured or the restoration cannot be verified safely, stop and escalate to a human instead of testing.
Record the commands and results of all three layers in the pull request description.

Also run `make lint`, and `make userdocs` when a schema description changes.
