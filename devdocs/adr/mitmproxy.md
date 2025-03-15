# ADR: Use of mitmproxy in the DevContainer for Network Debugging

## Status

Proposed

## Context

Debugging network traffic in the development environment is challenging. In the context of the Terraform Provider for Power Platform, rapid identification of issues related to API calls and server responses is essential. There are two significant upcoming benefits:

- Capturing network traffic and converting the requests/responses into mock data for unit tests.
- Simplifying the debugging process when building new resources or data sources by providing clear insights into server responses.

## Decision

We will integrate mitmproxy into the devcontainer. This setup allows developers to observe, intercept, and debug network traffic directly in the development environment. The primary focus will be on:

- Enabling comprehensive visibility into API interactions.
- Providing future capabilities to record and transform captured traffic into reusable mock scenarios for unit testing.
- Offering a robust mechanism to troubleshoot issues arising from unpredictable server responses.

## Consequences

- **Enhanced Debugging:** Developers benefit from real-time insights into network traffic, which improves the speed and accuracy of debugging.
- **Test Data Generation:** The recorded traffic has the potential to serve as a foundation for generating mock responses, thereby enhancing the unit testing framework.
- **Development Overhead:** Initial setup and maintenance of mitmproxy may introduce complexity. However, this is offset by the long-term benefits in diagnostics and test automation.
- **Security Considerations:** Care must be taken to ensure that sensitive information does not leak when capturing and storing network data.
- **Provider Code Stability:** No changes have been made to the provider code to support mitmproxy, ensuring that there is no test-only code path.

This ADR should be reviewed periodically. Future improvements may involve automating the conversion of captured data into unit test fixtures, thereby further streamlining the development process.
