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

Regardless of your chosen environment, you'll need access to a Power Platform tenant and appropriate credentials.

### Tenant Setup

Ensure you have access to a tenant where you can create and delete Power Platform environments and resources. Follow the [bootstrap readme](https://github.com/microsoft/power-platform-terraform-quickstarts/blob/main/bootstrap/README.md) from our quickstarts repository to set up your tenant.

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

## Running Acceptance Tests

To run all acceptance tests

```bash
make acctest
```

To run single acceptance test

```bash
TF_ACC=1 go test -v ./... -run TestAcc<test_name>
```

## Running Unit Tests

To run all unit tests

```bash
make unittest
```

To run single unit test

```bash
TF_ACC=0 go test -v ./... -run TestUnit<test_name>
```

> [!NOTE]
> The tests require permissions on the folders, these permissions are assigned when creating your container. If you have permission problems when running the unit tests, you can rebuild your development container or run the following commands again to assign the permissions to the necessary folders.

```bash
sudo chown -R vscode /workspaces/terraform-provider-power-platform/
sudo chown -R vscode /go/pkg
```

## Writing Tests

All the test for a given resource/datasource are located in `/internal/<resource/datasource_name>_test.go` file. When writing a new feature you should try to create [happy path](https://en.wikipedia.org/wiki/Happy_path) test(s) for you feature covering create, read and deletion of your new feature. For updates you should cover not only update of all properties but situation when a force recreate of a resource is requried (if you have such propeties in you resource).

### Writing Unit Tests

Unit test are created by mocking HTTP request, some of the often used HTTP mocks encapsulated in `ActivateEnvironmentHttpMocks` function, so that you don't have to write them for every test. When implementing new mocks, the mokcked response json files should be located in `/internal/services/<your_service_name>/test/<resource_or_datasource>/<name_of_the_unit_test>` folder

> [!TIP]
> When creating mocked json responses you can resuse the exising one by **duplicating** then into you `<name_of_the_unit_test>` folder.

> [!CAUTION]
> Your mocked json response file should not contain any Personally Identifiable Information such as tenantid, usernames, phone numbers, emails, addresses etc. You should anonymize that data.

### Writing Acceptance Tests

Each acceptance test is a copy of an unit test from tested use case perspective. That means for a given unit test we should have an acceptance test that validates the same use case but against a real infrastructure.

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```sh
go get github.com/author/dependency
make deps
```

## Updating Dependencies

Open a terminal in your devcontainer and type

```sh
make deps
```

## Updating Documentation

User documentation markdown files in [/docs](/docs/) are auto-generated by the [terraform plugin docs tool](https://github.com/hashicorp/terraform-plugin-docs).

> [!IMPORTANT]
> Do not manually edit the markdown files in [/docs](/docs/). If you need to edit documentation edit the following sources:

- schema information in the provider, resource, and data-source golang files that are in [/internal](/internal/)
- [template files](templates/)

```sh
make userdocs
```

User documentation is temporarily served on GitHub Pages which requires the [pages.yml GitHub workflow](/.github/workflows/pages.yml) to transform /docs markdown files into a static website.  Once this provider is published to the Terraform registry, documentation will be hosted on the registry instead.

## Making a Release

> [!TIP]
> In you development work flow, you don't have to release the provider in order to test it locally, instread you can use the devcontainer and keep installing it locally in there by using `make install` command.

Our releases use [semantic versioning](https://semver.org/).

Given a version number MAJOR.MINOR.PATCH, increment the:

- MAJOR version when you make incompatible API changes
- MINOR version when you add functionality in a backward compatible manner
- PATCH version when you make backward compatible bug fixes

Use the `preview` extension to the MAJOR.MINOR.PATCH format for preview release such as `v0.7.0-preview`.

### Using the CLI

As a last PR to `main` branch before new release, create documentation using [Changie](https://github.com/miniscruff/changie):

``` bash
changie batch 1.0.0-preview
```

to release, use the `git tag` command on the `main` branch (ensure main is up to date) and then push that release back to origin.

``` bash
git tag -a v1.0.0-preview -m "v1.0.0-preview"
git push origin v1.0.0-preview
```

Once the release is pushed to the repo, the [release.yml](/.github/workflows/release.yml) action pipeline will detect the new tag and create a draft release. After the build completes, you can publish the release if it looks good.

## Developer work flow

Once you decide to contribute back to this reposity by fixing a bug or adding a feature you work flow will be as follows:

1. Fork this repository and open in locally
1. Start working in devcontainer on your changes. (commands: `make install`, `terraform plan`, `terraform apply`)
    - Completly new feature should be located in a new `/internal/services/<new_service_name>` folder.
1. Add and/or update unit and accaptance tests. Tests for new feature should be created in new resource/datasource_test.go file (commands: `make unittest`, `make acctest`)
    - When working on a bug remember to add a new unit and acceptance test(s) covering your use case if that test does not exist yet.
    - When working on a new feature add unit and acceptance tests covering [happy path](https://en.wikipedia.org/wiki/Happy_path) for your feature, ideally also some edge cases. If you feature enhances existing resource/datasource, add/change validation of your new properties in all tests that use that resource/datasource
1. Create/Update examples in `/examples/...` folder(s)
    - When working on enhacement remeber to add new enhacement properties to all existing examples using that resource/datasource, especially if it is a requried property.
    - When creating new resource/datasource, create new examples showcasing how to use it.
1. Regenerate the docs. (commands: `make docs`)
1. Raise a pull request from your fork back the this repository

> [!NOTE]
> If your use case requries testing outside local devcontainer like for example running it from a Github action, then you will need to create a realease from your fork repository.
