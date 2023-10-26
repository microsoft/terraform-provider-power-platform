# Developer Guide

The Terraform Provider for Power Platform extends Terraform's capabilities to allow Terraform to manage Power Platform infrastructure and services.  The provider is built on the modern [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and NOT on the the older Terraform SDK.  Ensure that you are referencing the correct [Plugin Framework documentation](https://developer.hashicorp.com/terraform/plugin/framework) when developing for this provider.

If you want to contribute to the provider, refer to the [Contributing Guide](/CONTRIBUTING.md) which can help you learn about the different types of contributions you can make to the repo.  The following documentation will help developers get setup and prepared to make code contributions to the repo.

## Devcontainer

If you want to contribute to this project, you can use the devcontainer feature in Visual Studio Code to create a consistent and isolated development environment. A devcontainer is a Docker container that has all the tools and dependencies needed to work with the codebase. You can open any folder inside the container and use VS Code’s full feature set, including IntelliSense, code navigation, debugging, and extensions.

## Developer Requirements

To use the devcontainer in this repo, you need to have the following prerequisites:

- [Docker](https://www.docker.com/products/docker-desktop/)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed in VS Code.

## Opening the Devcontainer

Once you have the prerequisites, you can follow these steps to open the repo in a devcontainer:

1. Clone or fork this repo to your local machine.
1. Open VS Code and press F1 to open the command palette. Type “Remote-Containers: Open Folder in Container…” and select it.
1. Browse to the folder where you cloned or forked the repo and click “Open”.
1. VS Code will reload and start building the devcontainer image. This may take a few minutes depending on your network speed and the size of the image.
1. When the devcontainer is ready, you will see “Dev Container: Power Platform Terraform Provider Development” in the lower left corner of the VS Code status bar. You can also open a new terminal (Ctrl+Shift+`) and see that you are inside the container.
1. You can now edit, run, debug, and test the code as if you were on your local machine. Any changes you make will be reflected in the container and in your local file system.

Note: To work with the repository you will need to verify or configure your GIT credentials, you can do it as follows in the dev Container terminal:

- Verify Git user name and email:

```bash
git config --list
```

You should see your username and email listed, if they do not appear or you want to change them you must
establish them following the step below, (to quit the "git config" mode type "q").

- Change or set your Git username and email in the Dev Container:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@address"
```

Note: if you logging to docker container's shell outside the VS Code, in order to work with git repository, run the following commands:

```bash
export SSH_AUTH_SOCK=$(ls -t /tmp/vscode-ssh-auth* | head -1)
export REMOTE_CONTAINERS_IPC=$(ls -t /tmp/vscode-remote-containers-ipc* | head -1)
```

For more information about devcontainers, you can check out the [devcontainer documentation](https://code.visualstudio.com/docs/devcontainers/containers) and [sharing Git credentials with your container](https://code.visualstudio.com/remote/advancedcontainers/sharing-git-credentials).

## Power Platform Prerequisites

### Tenant

Developers should have access to a tenant where they can create and delete Power Platform environments and other resources.

TODO: more information about tenant and permissions needed by the dev user

### Credentials

See the [provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication) on getting a service principal or user account credentials configured.

### Environment Variables

Use environment variables to configure the provider to use your chosen credentials.  You may either pass credentials as terraform variables (via `TF_VAR_*` environment variables) or by using the provider's own environment variables (`POWER_PLATFORM_*`).  See the [provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication) for more information on configuring credentials for the provider.

## Running Provider locally in VSCode (linux)

Open bash terminal inside VS Code and execute the following commands:

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

Note: You cannot run `terraform init` when using dev overrides. `terraform init` will validate the versions and provider source, while `terraform plan` will skip those validations when `dev overrides` is part of your config.

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

Note: The tests require permissions on the folders, these permissions are assigned when creating your container.
If you have permission problems when running the unit tests, you can rebuild your development container
or run the following commands again to assign the permissions to the necessary folders.

```bash
sudo chown -R vscode /workspaces/terraform-provider-power-platform/
sudo chown -R vscode /go/pkg
```

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

Do not manually edit the markdown files in [/docs](/docs/). If you need to edit documentation edit the following sources:

- schema information in the provider, resource, and data-source golang files that are in [/internal](/internal/)
- [template files](templates/)

```sh
make userdocs
```

User documentation is temporarily served on GitHub Pages which requires the [pages.yml GitHub workflow](/.github/workflows/pages.yml) to transform /docs markdown files into a static website.  Once this provider is published to the Terraform registry, documentation will be hosted on the registry instead.
