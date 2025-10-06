# Power Platform Terraform Provider Roadmap

This roadmap outlines the planned direction for the Power Platform Terraform Provider over the next 12 months. It serves as a high-level guide for contributors and users, providing transparency into where the project is headed while maintaining flexibility to adapt to community needs and evolving requirements.

## Key Themes and Planned Direction

| Key Theme | Planned Direction (12-month horizon) |
|-----------|--------------------------------------|
| **Quality & Reliability** | Prioritize fixes for issues raised by users in the GitHub repository, improve error messages and diagnostics to enhance the user experience, and host periodic community discussions to gather feedback on feature priorities. Streamline the contribution process and expand documentation for new contributors. |
| **Security Improvements** | Add support for service principal authentication for resources as they become available in the APIs. Strengthen binary-signing and SBOM (Software Bill of Materials) practices in our release pipeline, integrate stricter static and dynamic analysis tools into CI/CD workflows, and address any newly reported CVEs promptly. Update release process to newer Terraform registry release process. |
| **AI & Copilot Studio Support** | We would consider adding resources and data sources to support this space. |
| **Governance & Environment Management** | We would consider extending or adding new resources and data sources with new features to support environment governance, policies, or large-scale environment management. Examples include disaster recovery settings on environments. |

## Out of Scope for This Period

To maintain focus and deliver quality improvements, the following areas are **not** planned for the next 12 months:

- **Major API Breaking Changes**: While we follow semantic versioning, we will avoid unnecessary breaking changes that would disrupt existing user configurations
- **Experimental Power Platform Services**: Features that are in private preview or experimental stages within Power Platform will not be prioritized until they reach public preview or general availability
- **Cross-Cloud Integration**: Integration with non-Microsoft cloud services or hybrid scenarios beyond what Power Platform natively supports
- **Legacy Power Platform Features**: Deprecated or legacy Power Platform features that Microsoft is phasing out will not receive new provider support

## How This Roadmap Works

- **Flexibility**: This roadmap provides direction but remains flexible to adapt to urgent security fixes, critical bugs, or significant community requests
- **Community Input**: We welcome feedback and suggestions through [GitHub Discussions](https://github.com/microsoft/terraform-provider-power-platform/discussions) and issues
- **Regular Updates**: This roadmap will be reviewed and updated quarterly to reflect progress and any necessary adjustments
- **Semantic Versioning**: All changes will follow [semantic versioning](https://semver.org/) principles to clearly communicate the impact of updates

## Contributing

Interested in contributing to these initiatives? Check out our [Contributing Guide](CONTRIBUTING.md) and [Developer Guide](DEVELOPER.md) to get started. We particularly welcome contributions in areas aligned with this roadmap.

For questions about this roadmap or to suggest changes, please open a [GitHub Discussion](https://github.com/microsoft/terraform-provider-power-platform/discussions).