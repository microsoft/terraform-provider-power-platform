---
page_title: "Continuous Access Evaluation (CAE) - Power Platform Provider"
subcategory: "Authentication"
description: |-
  Guide on using Continuous Access Evaluation (CAE) with the Power Platform Terraform Provider.
---

# Continuous Access Evaluation (CAE)

This guide explains how to use Continuous Access Evaluation (CAE) with the Power Platform Terraform Provider.

## What is Continuous Access Evaluation?

Continuous Access Evaluation (CAE) is a security feature that enables real-time security policy enforcement for access tokens. 
With CAE enabled, authentication tokens can be immediately invalidated when security policy changes occur, such as:

- User account termination or disablement
- Password changes or resets
- Changes to conditional access policies
- Changes to network location policies

Without CAE, security policy changes might not take effect until an access token expires (typically 1 hour), creating a 
potential security gap where a user might retain access despite security policy changes.

## Enabling CAE

To enable CAE in the Power Platform Provider, add the `enable_continuous_access_evaluation` option to your provider configuration:

```terraform
provider "powerplatform" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  
  # Enable Continuous Access Evaluation
  enable_continuous_access_evaluation = true
}
```

When CAE is enabled, the provider will:

1. Request CAE-enabled tokens from Microsoft Entra ID (Azure AD)
2. Handle CAE challenge responses by detecting security policy violations
3. Report clear error messages when security policies block access

## How CAE Works

When you enable CAE, the provider instructs the Azure Identity library to request tokens with CAE support. These tokens 
contain additional claims that allow Microsoft Entra ID to revoke them in real-time when security policies change.

If a security policy violation occurs (e.g., the user's account is disabled), API requests made with the token will 
receive a CAE challenge response. The provider detects these responses and surfaces a clear error message indicating 
that access was denied due to a security policy change.

## When to Use CAE

CAE is recommended for all production environments as it significantly improves security posture by ensuring that 
security policies are enforced in near real-time. It's especially valuable when:

- Managing resources in high-security environments
- Implementing Zero Trust security architectures
- Complying with regulatory requirements that mandate immediate access revocation
- Running automation where timely security policy enforcement is critical

## Considerations

- CAE is supported for all authentication methods in the provider (client secret, certificate, CLI, managed identity, OIDC)
- There is no performance impact from enabling CAE for normal operations
- When a security policy violation occurs, the provider will report a clear error message indicating the cause

## Learn More

For more information about Continuous Access Evaluation, see the Microsoft documentation:

- [Continuous Access Evaluation in Microsoft Entra ID](https://learn.microsoft.com/entra/identity/conditional-access/concept-continuous-access-evaluation)
- [How to use Continuous Access Evaluation enabled APIs in your applications](https://learn.microsoft.com/en-us/entra/identity/conditional-access/developer/app-developer-guide)
