# Community and Support Documentation

## FAQ (Frequently Asked Questions)

### Q: Can I develop on this provider without using the VS Code Dev Container?

**A:** It's not recommended to develop this provider outside the provided Dev Container. GitHub Codespaces is the preferred environment due to consistency, controlled dependencies, and reproducibility. Using the Dev Container locally is acceptable, but manually setting up a local environment can introduce risks such as mismatched dependencies, version conflicts, and environmental inconsistencies that could complicate development and troubleshooting.

### Q: After building the provider locally, terraform init fails or doesn’t use my changes.

**A:** When using a dev build, configure Terraform to use the local plugin binary. The repository provides a `.terraformrc` file with a dev override for the provider. If you are in the Dev Container, this is likely already set up. Instead of `terraform init`, simply run `terraform plan`. Terraform will skip the usual provider download and use your locally installed provider binary. However, if your examples also include other providers such as `azapi` or `azurerm`, running `terraform init` may still be necessary. You can safely ignore any errors related specifically to the Power Platform provider when executing `terraform init`. If Terraform still attempts to download the provider, verify the dev override in `.terraformrc` in your home or working directory, and ensure you've executed `make install`.

### Q: Where can I find usage examples or documentation for end users?

**A:** End user documentation is available on the [Terraform Registry: Power Platform Provider Docs](https://registry.terraform.io/providers/microsoft/power-platform/latest/docs). It includes detailed resource and data source documentation, examples, and authentication configuration. Additional resources include an [official Microsoft Learn article](https://learn.microsoft.com/en-us/business-applications/playbook/enterprise-solutions/power-platform-terraform-provider/) and the [QuickStarts repository](https://github.com/microsoft/power-platform-terraform-quickstarts) containing real-world Terraform configurations.

### Q: I wrote an acceptance test, but it’s failing due to permissions or timeouts.

**A:** Ensure your test environment has appropriate permissions. For instance, testing environment creation requires admin credentials. Timeouts can be adjusted via Terraform Plugin SDK’s testing framework. Ensure no leftover resources from previous tests are causing conflicts. For persistent issues, open a discussion or issue marked related to testing on GitHub.

### Q: The provider build or tests are failing on my machine. What can I do?

**A:** Update your code (`main` branch), run `make deps`, and rebuild your Dev Container if necessary. Use verbose testing (`-v`) for clearer output. Ensure environment variables are properly configured, especially in the Dev Container. Check file permissions, running the provided `chown` commands if necessary. If issues persist, consult or file GitHub issues.

## Developer Support Channels

The primary channel for support and communication is the project's GitHub repository:

- **[GitHub Issues](https://github.com/microsoft/terraform-provider-power-platform/issues):** Report bugs or request features by providing detailed logs, Terraform configurations, and reproduction steps.
- **[GitHub Discussions](https://github.com/microsoft/terraform-provider-power-platform/discussions):** Ideal for Q&A, design topics, implementation guidance, or feedback.
- **[Security Issues](https://github.com/microsoft/terraform-provider-power-platform/security/policy):** Report privately following the security reporting guidelines.

Microsoft Tech Community and Terraform community forums are secondary resources. There is currently no public Slack or Teams channel dedicated to this project.

Always adhere to the project's [Code of Conduct](https://github.com/microsoft/terraform-provider-power-platform/blob/main/CODE_OF_CONDUCT.md), maintaining respectful and patient communication.

## Working with Maintainers

When engaging maintainers (e.g., PR code reviews), follow these best practices:

- **Be Responsive:** Address feedback promptly to expedite merges.
- **Follow Guidelines:** Adhere to the PR templates and guidelines provided.
- **Testing Evidence:** Provide clear test results or log snippets indicating successful tests.
- **Understand Priorities:** Accept feedback constructively, aligning with the project roadmap or considering maintaining a fork if your contribution is niche.

Support the community by answering questions, reviewing PRs, and mentoring newcomers.

