# Developer Guide

The Terraform Provider for Power Platform extends Terraform's capabilities to allow Terraform to manage Power Platform infrastructure and services.  The provider is built on the modern [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and NOT on the the older Terraform SDK.  Ensure that you are referencing the correct [Plugin Framework documentation](https://developer.hashicorp.com/terraform/plugin/framework) when developing for this provider.

If you want to contribute to the provider, refer to the [Contributing Guide](/CONTRIBUTING.md) which can help you learn about the different types of contributions you can make to the repo.  The following documentation will help developers get setup and prepared to make code contributions to the repo.

## Summary of Developer Resources

To help you effectively contribute to the Terraform Provider for Power Platform, we've documented several key guidelines and decisions that shape our development practices:

- **[Contributing Guide](./CONTRIBUTING.md)**: Provides clear instructions on how to contribute new resources, data sources, or improvements. It includes a detailed pull request checklist to ensure your contributions meet our quality standards.

- **[Testing Guidelines](./devdocs/testing_guidelines.md)**: Outlines best practices for writing robust unit and acceptance tests. It covers testing patterns, mocking HTTP interactions using `httpmock`, and ensuring comprehensive test coverage across various scenarios, including error handling, idempotency, and multi-cloud compatibility.

- **[Security Guidelines](./devdocs/security_guidelines.md)**: Highlights essential security practices, including secure coding, handling sensitive data, dependency management, and responsible vulnerability disclosure.

- **[Observability and Logging](./devdocs/observability_guidelines.md)**: Describes how to implement structured logging using Terraform's `tflog` package, ensuring effective debugging and observability without exposing sensitive information.

- **[Schema Design Guidelines](./devdocs/schema_guidelines.md)**: Outlines best practices for designing resource and data source schemas, including choosing attribute types, structuring schemas, and deciding between loosely typed and strongly typed resources.

- **[Feature Request Guidelines](./devdocs/feature_request_guidelines.md)**: Provides criteria for evaluating new feature requests, clarifying the provider's scope, and distinguishing between infrastructure-as-code and application lifecycle management scenarios.

- **[Community and Support](./devdocs/community_support.md)**: Explains how to engage with maintainers and the community effectively, including channels for support, reporting issues, and contributing to discussions.

We encourage you to explore these resources to familiarize yourself with our development philosophies and best practices. Following these guidelines helps maintain a consistent, secure, and robust provider that benefits the entire community.

## Developer Workflow

Once you decide to contribute back to this repository by fixing a bug or adding a feature, your workflow will be as follows:

1. Fork this repository and open it in codespaces or local devcontainer
1. Start working in the devcontainer on your changes (commands: `make install`, `terraform plan`, `terraform apply`)
    - Completely new features should be located in a new `/internal/services/<new_service_name>` folder.
1. Add and/or update unit and acceptance tests. Tests for new features should be created in a new resource/datasource_test.go file (commands: `make unittest`, `make acctest`)
    - When working on a bug, remember to add a new unit and acceptance test(s) covering your use case if that test does not exist yet.
    - When working on a new feature, add unit and acceptance tests covering the [happy path](https://en.wikipedia.org/wiki/Happy_path) for your feature, ideally also some edge cases. If your feature enhances an existing resource/datasource, add/change validation of your new properties in all tests that use that resource/datasource.
1. Create/Update examples in `/examples/...` folder(s)
    - When working on enhancement, remember to add new enhancement properties to all existing examples using that resource/datasource, especially if it is a required property.
    - When creating a new resource/datasource, create new examples showcasing how to use it.
1. Regenerate the docs, run liners, run unit tests locally (commands: `make precommit`)
1. Create a changelog record for your contribution (commands: `changie new`)
1. Raise a pull request from your fork back to this repository

> [!NOTE]
> Core maintainers with write access to the `microsoft/terraform-provider-power-platform` repository do not need to fork.  They may create branches in the repository.

## Development Environment Options

You have two recommended options for setting up your development environment:

- **GitHub Codespaces**: A cloud-based, fully managed development environment accessible directly from your browser or VS Code.
- **Local Devcontainer**: A Docker-based development environment running locally on your machine using Visual Studio Code.

Both options provide consistent, isolated environments with all necessary tools and dependencies pre-installed. We do not recommend setting up a local development environment manually.

## Option 1: GitHub Codespaces (Preferred)

GitHub Codespaces provides a fully managed, cloud-hosted development environment accessible directly from your browser or VS Code. It includes all necessary dependencies and tools pre-installed, allowing you to start contributing immediately.

### Getting Started with Codespaces

1. Navigate to the repository on GitHub.
2. Click the green **Code** button, then select **Open with Codespaces** > **New codespace**.
3. GitHub will automatically build and launch your Codespace environment. This may take a few minutes.
4. Once ready, you can edit, run, debug, and test the code directly in your browser or connect to it from your local VS Code.

For more information, see the [GitHub Codespaces documentation](https://github.com/features/codespaces).

### Configuring Git Credentials in Codespaces

GitHub Codespaces automatically configures your Git credentials based on your GitHub account. No additional configuration is typically required.

### Managing Secrets in Codespaces

Use [Codespaces Secrets](https://docs.github.com/en/codespaces/managing-your-codespaces/managing-secrets-for-your-codespaces) to securely store sensitive information such as credentials or environment variables.  Codespace secrets are automatically projected into your codespace as environment variables.

## Option 2: Local Devcontainer

A local devcontainer provides a consistent, isolated development environment on your local machine using Docker and Visual Studio Code.

### Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed in VS Code.

### Opening the Devcontainer Locally

1. Clone or fork this repository to your local machine.
2. Open VS Code and press `F1` to open the command palette. Type `Remote-Containers: Open Folder in Container…` and select it.
3. Browse to the cloned repository folder and click **Open**.
4. VS Code will reload and build the devcontainer image. This may take a few minutes.
5. Once ready, you'll see **Dev Container: Power Platform Terraform Provider Development** in the lower-left corner of VS Code. Open a new terminal (`Ctrl+Shift+``) to confirm you're inside the container.
6. You can now edit, run, debug, and test the code. Changes will reflect both in the container and your local file system.

For more information, see the [VS Code Devcontainer documentation](https://code.visualstudio.com/docs/devcontainers/containers).

### Configuring Git Credentials in Local Devcontainer

Verify or configure your Git credentials within the devcontainer terminal:

- Verify Git username and email:

```bash
git config --list
```

- Set or update your Git username and email:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@address"
```

If accessing the container shell outside VS Code, run:

```bash
export SSH_AUTH_SOCK=$(ls -t /tmp/vscode-ssh-auth* | head -1)
export REMOTE_CONTAINERS_IPC=$(ls -t /tmp/vscode-remote-containers-ipc* | head -1)
```

For more details, see [sharing Git credentials with your container](https://code.visualstudio.com/remote/advancedcontainers/sharing-git-credentials).

### Using Terminal

If you prefer to use you operating system terminal instead of VSCode you can run the following command:

```bash
docker exec -u vscode -w /workspaces/terraform-provider-power-platform -it <your_docker_container_name_goes_here> bash -c "exec bash"
```

## Power Platform Prerequisites

Regardless of your chosen environment, you'll need access to a Power Platform tenant, licenses, and appropriate credentials.

### Tenant Setup

Ensure you have access to a tenant where you can create and delete Power Platform environments and resources. Follow the [bootstrap readme](https://github.com/microsoft/power-platform-terraform-quickstarts/blob/main/bootstrap/README.md) from our quickstarts repository to set up your tenant.

### Licensing
Verify that you have the necessary licenses for the users and services you want to create, adding users to the power platform environments requires licensing, you can verify your licenses following the [View license consumption for Power Apps and Power Automate (preview)](https://learn.microsoft.com/power-platform/admin/view-license-consumption-issues)

### Credentials and Authentication

Refer to the [provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication) for detailed instructions on authenticating to the provider.  When getting started we recommend authenticating with a user-context using `az login` command from az cli and `use_cli` provider configuration.

#### Service Principals

While we recommend getting started using user-context authentication, most production scenarios will use service principals. Some APIs behave differently when called in a user-context vs service principal context so we highly recommend testing your contributions with a service principal.

Configure service principal credentials securely using environment variables:

- In **Codespaces**, leverage [Codespaces Secrets](https://docs.github.com/en/codespaces/managing-your-codespaces/managing-secrets-for-your-codespaces).
- In **Local Devcontainer**, avoid using local `.env` files to prevent accidental exposure of sensitive information.

Credentials can be passed as Terraform variables (`TF_VAR_*`) or provider-specific environment variables (`POWER_PLATFORM_*`). See the [provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication) for more details.

## Building and Running the Provider in VSCode

Open bash terminal inside VSCode and execute the following commands:

```bash
# you should already be in this directory, but just in case
cd /workspaces/terraform-provider-power-platform

# Build and install the provider's binary locally
make install

# Navigate to a folder that contains *.tf files and run below
cd examples/data-sources/powerplatform_environments

# Run terraform to validate that provider is functioning
terraform plan
```

> [!NOTE]
> You cannot run `terraform init` when using dev overrides. `terraform init` will validate the versions and provider source, while `terraform plan` will skip those validations when `dev overrides` is part of your config. You can simply run `terraforn plan` and `terrafirn apply` when working in devcontainer.

> [!TIP]
> Because when working locally the `terraform init` command will not work, if you need additional terraform providers from terraform registry, all of them have to be added locally to the devcontainer in order to `terraform plan` and `terraform apply` work. You can add you missing terraform providers by adding them `.devcontainer/features/acceptance_test_dependencies/main.tf` and rebuild the devcontainer.

## Debugging provider in VSCode

1. Open VSCode with the root folder as the parent of this ReadMe
1. Click On Run and Debug (F5)
1. Copy `TF_REATTACH_PROVIDERS` value in the Debug Console
1. Run `export TF_REATTACH_PROVIDERS=<value>` with the value copied from the above step
1. Add breakpoints
1. `cd` to a parent folder where main.tf exists
1. Run `terraform apply`

## Debugging network calls

This devcontainer has `mitmproxy` preinstalled.  There are 3 versions available: `mitmproxy` a command line UI, `mitmweb` a web based UI, and `mitmdump` a headless command to dump network traffic to file.  For example, first run:

```bash
mitmproxy
```

Then in a **separate terminal** you can either launch acceptance tests with:

```bash
make acctest TEST=<test_name> USE_PROXY=1
```

or run an example with:

```bash
HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 terraform apply
```

Flip back to your `mitmproxy` terminal to inspect the network traffic.

## Testing Guidelines

Quality tests are essential to ensure our provider works reliably across different environments and scenarios. As a contributor, you'll be expected to write comprehensive tests for any code you submit.

Our testing approach includes both unit tests and acceptance tests:

- **Unit tests**: Fast-executing tests that verify individual components in isolation
- **Acceptance tests**: End-to-end tests that validate provider behavior against actual Power Platform APIs

The detailed [testing guidelines](/devdocs/testing_guidelines.md) will help you understand our testing patterns, best practices, and expectations. We've designed them to make writing tests straightforward while ensuring they're thorough enough to catch potential issues. You'll find examples of how to use our testing frameworks, how to properly mock external dependencies using `httpmock`, and how to structure your tests for readability and maintainability.

Testing is a collaborative effort - if you're unsure about testing approaches for specific scenarios, don't hesitate to ask questions in your PR or discussions.

## Dependencies

Managing dependencies carefully is crucial for maintaining both the security and reliability of our provider. Each external dependency introduces potential risks, from security vulnerabilities to licensing issues, and contributes to our software supply chain - the complete set of components and processes that go into our software.

To ensure transparency and thoughtful consideration, **any Pull Request that adds a new dependency must include an Architectural Decision Record (ADR)** explaining:

- Why the dependency is necessary
- Alternative options considered
- Security implications
- Maintenance considerations
- Licensing compatibility

This process helps us maintain a clean dependency tree and reduces potential attack vectors or unexpected behavior in production environments.

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```sh
go get github.com/author/dependency
make deps
```

### Updating Dependencies

Open a terminal in your devcontainer and type `make deps`

## User Documentation

User documentation markdown files in [/docs](/docs/) are auto-generated by running `make userdocs` which uses the [terraform plugin docs tool](https://github.com/hashicorp/terraform-plugin-docs). It parses information about providers, resources, and data sources directly from the Go code. By eliminating the need to manually maintain separate documentation files, we reduce duplication, keep docs in sync with the schema, and ensure the Terraform Registry always reflects the latest code.  **Do not manually edit the markdown files in [`/docs`](/docs/)**

If you need to edit documentation edit the following sources:

1. **Schema Definitions in `/internal`**  
   - Each resource and data source has a schema defined in the Go code under the [`/internal/`](/internal/) directory. The schema typically includes fields like `Name`, `Description`, and required/optional indicators.
   - `terraform-plugin-docs` automatically looks for these fields and uses the `MarkdownDescription` metadata in those schema blocks to build out the documentation for each resource or data source.
   - When possible include links in `MarkdownDescription` to official Microsoft Power Platform Admin Center documenation for the equivalent feature in PPAC.
   - Terraform provider documentation is served from the Terraform Registry website, which uses a specialized markdown rendering engine. Unlike GitHub markdown, the Registry renderer has specific syntax requirements for features like callouts (which use `-> NOTE:` syntax instead of the GitHub style blockquotes).

1. **Template Files in `/templates`**  
   - In addition to reading the schema, the tool merges content from the [template files](/templates/). These template files often contain the main structure of the user documentation, placeholders for resource attributes, usage examples, and any provider-level guidance or usage instructions.
   - If you need to provide custom instructions, additional examples, or overviews, this is done in the templates (rather than in the automatically-generated docs).

1. **Examples in `/examples`**
   - Data Source examples in `examples/data-sources/{your_data_source}/` should have `data-source.tf`, `outputs.tf`, and optionally `variables.tf` files
   - Resources examples in `examples/resources/{your_data_source}/` should have `resource.tf`, `outputs.tf` files, and optionally `variables.tf` files.  Including a `import.sh` will signal that your resource is compatible with the `import` command.

**Before committing, open a terminal in your devcontainer and type `make precommit`** which internally calls `make userdocs`, `make lint`, and `make unittest` to regenerate documentation, runs linters, and runs unit tests.

### Developer Documentation

Developer documentation in this repository exists to help other contributors understand design decisions, architecture, and development practices. Unlike user documentation which is auto-generated, developer documentation is maintained manually.

Developer documentation primarily lives in two places:

- **[/devdocs](/devdocs/)**: Contains detailed guidelines on specific aspects of development
- **[DEVELOPER.md](/DEVELOPER.md)**: This file, which provides an overview of the development process

You should update or add developer documentation when:

- Adding new features with non-trivial implementation details
- Changing existing behavior that other developers should be aware of
- Making architectural decisions that impact the codebase
- Adding/improving development workflows or tools

## Preparing a PR

Before submitting a pull request, ensure your contribution meets all the requirements in our [Pull Request Checklist](/CONTRIBUTING.md#pull-request-checklist). This checklist covers essential aspects such as implementation quality, test coverage, documentation, and schema descriptions.

When your changes are ready:

1. Verify all tests pass with `make unittest` and `make acctest`
2. Ensure documentation is updated by running `make userdocs`
3. Create a changelog entry with `changie new`
4. Write a clear PR description explaining your changes, their purpose, and any limitations
5. Reference any related issues using GitHub's referencing syntax (#issue-number)

PRs that follow these guidelines will move through the review process more smoothly. Maintainers may request changes to ensure your contribution aligns with the project's quality standards and design principles.
