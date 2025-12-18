# Roadmap

This document outlines the high-level direction of the Terraform Provider for Power Platform over the next 12 months. It is intended to give contributors and users visibility into planned focus areas without committing to specific dates or deliverables.

> [!NOTE]
> This roadmap reflects our current intentions and may evolve based on community feedback, changing priorities, or new opportunities. For detailed progress, see the [issue tracker](https://github.com/microsoft/terraform-provider-power-platform/issues), [past releases](https://github.com/microsoft/terraform-provider-power-platform/releases) and [changelog](./CHANGELOG.md).

## Quarterly Milestones

### Q1 2025 (Jan–Mar): Quality & Reliability

Focus on making the provider more stable and trustworthy.

- Increase unit test coverage to at least 80% across all services
- Reduce acceptance tests and improve test isolation
- Establish a more predictable release schedule
- Keep user documentation accurate and up-to-date

### Q2 2025 (Apr–Jun): Security Hardening

Focus on strengthening security practices.

- Continue signing all releases with GPG keys
- Improve SBOM (Software Bill of Materials) practices for supply chain transparency
- Integrate stricter static and dynamic analysis tools into CI/CD
- Address any reported CVEs
- Meet and exceed [OpenSSF Best Practices](https://www.bestpractices.dev/projects/8714) requirements

### Q3 2025 (Jul–Sep): AI & Copilot Studio Support

Focus on Power Platform's expanding AI capabilities.

- Add resources and data sources for AI-powered features
- Expand support for Copilot Studio configuration
- Provide clear examples for using Terraform with Power Platform AI features

### Q4 2025 (Oct–Dec): Community-Driven Improvements

Focus on user experience and community engagement.

- Prioritize fixes for issues raised by users
- Improve error messages for better troubleshooting
- Host periodic GitHub Discussions to gather feedback
- Make it easier for new contributors to get started

## How to Provide Feedback

We welcome your input! Here's how you can help shape the roadmap:

- **GitHub Issues**: [Open an issue](https://github.com/microsoft/terraform-provider-power-platform/issues/new) for bug reports or feature requests.
- **GitHub Discussions**: Join the conversation in [Discussions](https://github.com/microsoft/terraform-provider-power-platform/discussions) to share ideas or ask questions.
- **Pull Requests**: Contributions are welcome! See the [Contributing Guide](./CONTRIBUTING.md) to get started.

## Related Documents

- [Contributing Guide](./CONTRIBUTING.md) – How to contribute to the project
- [Developer Guide](./DEVELOPER.md) – Development setup and workflow
- [Feature Request Guidelines](devdocs/feature_request_guidelines.md) – Criteria for feature inclusion
- [Security](./SECURITY.md) – How to report security concerns
- [Changelog](./CHANGELOG.md) – History of changes

---

_Last updated on December 2025_
