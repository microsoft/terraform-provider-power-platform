# Power Platform Terraform Provider

The Power Platform Terraform Provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/).

See our [Roadmap](./ROADMAP.md) for the project's direction over the next 12 months.

> [!CAUTION]
> Bugs or errors in Infrastructure-as-Code (IaC) software could lead to service interruptions or data loss. We strongly recommend backing up your data and testing thoroughly in non-production environments before using any feature in production. Your feedback is valuable to us, so please share any issues or suggestions you encounter via [GitHub issues](https://github.com/microsoft/terraform-provider-power-platform/issues).

Some resources and data sources are made available as a preview. Preview features may have restricted or limited functionality. Future updates could include breaking changes; however, we adhere to [Semantic Versioning](https://semver.org/) to clearly communicate these changes.

The following resources are in **preview**:

- powerplatform_analytics_data_exports
- powerplatform_copilot_studio_application_insights
- powerplatform_environment (only when creating developer environment types)
- powerplatform_environment_group_rule_set
- powerplatform_environment_wave
- powerplatform_tenant_capacity

## Using the Provider

The [user documentation](https://microsoft.github.io/terraform-provider-power-platform) contains information about how to install, configure, and use the provider to manage Power Platform resources and data sources. More advances examples together with bootstrap script can be found in the [Quick Starts Repository](https://github.com/microsoft/power-platform-terraform-quickstarts).

## Contributing

Refer to the [Contributing Guide](/CONTRIBUTING.md) to learn about the different types of contributions you can make to the repo.

For developers interested in contributing to the provider, we offer comprehensive documentation to guide you through the development process:

- [Developer Guide](DEVELOPER.md): Main guide covering development environment setup, workflow, and contributing practices
- [Schema Guidelines](/devdocs/schema_guidelines.md): Best practices for designing resource and data source schemas
- [Testing Guidelines](/devdocs/testing_guidelines.md): Approaches for writing robust unit and acceptance tests
- [Security Guidelines](/devdocs/security_guidelines.md): Security best practices and vulnerability reporting
- [Observability Guidelines](/devdocs/observability_guidelines.md): Implementing proper logging and telemetry
- [Feature Request Guidelines](/devdocs/feature_request_guidelines.md): Criteria for determining which features belong in the provider
- [Release Guidelines](/devdocs/release_guidelines.md): Process for releasing new provider versions
- [Community Support](/devdocs/community_support.md): How to engage with the community and get help

These resources will help you understand our development philosophy, technical requirements, and contribution processes.

## Security

For information about reporting security concerns see the [Security](SECURITY.md) documentation.

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8714/badge)](https://www.bestpractices.dev/projects/8714)

## Data Collection

The software may collect information about you and your use of the software and send it to Microsoft. Microsoft may use this information to provide services and improve our products and services. You may turn off the telemetry as described in the repository. There are also some features in the software that may enable you and Microsoft to collect data from users of your applications. If you use these features, you must comply with applicable law, including providing appropriate notices to users of your applications together with a copy of Microsoft’s privacy statement. Our privacy statement is located at <https://go.microsoft.com/fwlink/?LinkID=824704>. You can learn more about data collection and use in the help documentation and our privacy statement. Your use of the software operates as your consent to these practices.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general). Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.
