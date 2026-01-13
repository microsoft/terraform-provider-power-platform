# Roadmap

This document outlines the high-level direction of the Terraform Provider for Power Platform. It is intended to give contributors and users visibility into planned focus areas and help guide prioritization of the backlog without committing to specific dates or deliverables.

> [!NOTE]
> This roadmap reflects our current intentions and may evolve based on community feedback, changing priorities, or new opportunities. For detailed progress, see the [issue tracker](https://github.com/microsoft/terraform-provider-power-platform/issues), [past releases](https://github.com/microsoft/terraform-provider-power-platform/releases) and [changelog](./CHANGELOG.md).

## Priorities

The following priorities guide how we triage issues, plan work, and accept contributions. They are listed in approximate order of importance.

### Quality & Reliability

Making the provider stable and trustworthy is our top priority.

- Increase unit test coverage to at least 80% across all services
- Reduce acceptance test runtime and improve test isolation
- Establish a more predictable release schedule
- Keep user documentation accurate and up-to-date

### Security Hardening

Maintaining strong security practices is critical.

- Address any reported CVEs promptly
- Continue signing all releases with GPG keys
- Improve SBOM (Software Bill of Materials) practices for supply chain transparency
- Integrate stricter static and dynamic analysis tools into CI/CD
- Meet and exceed [OpenSSF Best Practices](https://www.bestpractices.dev/projects/8714) requirements

### Public API Migration

Prefer public Power Platform APIs over internal or undocumented APIs.

- Migrate existing resources to use public Power Platform APIs where available
- Prioritize new resources that can be built on stable, documented APIs
- Reduce reliance on internal APIs that may change without notice

### Release Process Improvements

Align with Terraform ecosystem best practices.

- Evaluate and adopt updated release processes for the Terraform Registry/marketplace
- Improve automation around releases and changelog generation

### AI & Copilot Studio Support

Support Power Platform's expanding AI capabilities.

- Add resources and data sources for AI-powered features
- Expand support for Copilot Studio configuration
- Provide clear examples for using Terraform with Power Platform AI features

### Community-Driven Improvements

Improve user experience and community engagement.

- Prioritize fixes for issues raised by users
- Improve error messages for better troubleshooting
- Host periodic GitHub Discussions to gather feedback
- Make it easier for new contributors to get started

## Not a Priority

The following items are explicitly out of scope or deprioritized. This helps set expectations and keeps our focus on higher-impact work.

- **Parity with every Power Platform Admin Center feature** – We focus on the most impactful Terraform use cases rather than replicating every UI feature.
- **Supporting deprecated or legacy APIs** – We prioritize modern, supported APIs and will not invest in deprecated endpoints.

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

_Last updated on January 2026_
