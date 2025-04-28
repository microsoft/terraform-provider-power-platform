---
page_title: Security Considerations
subcategory: Security
description: Security model, threat boundaries, and hardening guidance for the Terraform Power Platform provider.
---

# Security Considerations

_Last reviewed: 2025-04-26_

## Threat model

| Aspect        | Consideration                                                                                                   |
|---------------|-----------------------------------------------------------------------------------------------------------------|
| **Assets**    | Power Platform resources (environments, connections, Copilot Studio assets) managed by this Terraform provider |
| **Adversary** | Anyone able to submit malicious Terraform configuration, intercept network traffic, or inspect local logs       |
| **Assumptions** | Terraform CLI/Cloud worker runs on a patched host; the Go 1.22+ runtime is uncompromised                      |

## Trust boundaries

Terraform CLI / Cloud → **provider RPC** → Provider → **HTTPS (TLS 1.2+)** → Power Platform REST APIs

* Configuration files cross from user space into the provider process.  
* OAuth tokens traverse TLS-protected channels to Power Platform APIs.  
* No additional privilege-escalation path is expected.

## Security design goals

The provider is **designed to**:

* **Avoid logging secrets** – client secrets, certificates, and tokens remain in memory and are redacted from diagnostics.  
* **Enforce TLS 1.2+ with certificate verification** for every HTTP request; no plaintext or downgraded transport path exists.  
* **Validate input against the defined Terraform schema** – Rejects unexpected keys, type mismatches, and malformed strings based on schema rules before API calls.  
* **Ship verifiable artifacts** – each release includes a SHA-256 checksum file and detached GPG signature.  
* **Request appropriate OAuth scopes** – Infers scopes based on target API endpoints, often using `.default` to leverage permissions granted during application consent.

## Limits / out of scope

* Terraform state may contain plaintext IDs or connection strings; securing state storage (e.g., encrypted remote backend) is outside the provider.  
* The provider does not conceal values shown by `terraform plan` if a field is intentionally marked `sensitive = false`.  
* If a tenant admin grants excessive API permissions to the service principal, the provider could use them.
* Securing credentials (secrets, certificates, passwords) used by the provider is the user's responsibility.
* The provider executes the defined configuration; it does not analyze the user's intent or prevent potentially harmful (but valid) configurations.
* Monitoring Power Platform or Entra ID audit logs for security events related to provider actions is the user's responsibility.

## Hardening recommendations

* Run Terraform from a CI worker with limited privileges and secured logs.  
* Store state in an encrypted remote backend (Terraform Cloud, Azure Storage with SSE, etc.).  
* Prefer Entra workload-identity federation over long-lived client secrets.  
* Rotate credentials regularly and monitor audit logs for unexpected API calls.
* Apply the principle of least privilege to the Entra ID application registration or managed identity used by the provider. Grant only the necessary Power Platform and Azure API permissions required for the resources being managed.
* Use static code analysis tools (e.g., `tfsec`, `checkov`) to scan Terraform configurations for potential security issues before applying changes.
* Regularly review Terraform state files for any sensitive data that might have been unintentionally stored.
