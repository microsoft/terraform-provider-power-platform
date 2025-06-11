---
page_title: "Authenticating to Power Platform using a Service Principal and a Client Secret/Certificate"
subcategory: "Authentication"
description: |-
  <no value>
---


# Authenticating to Power Platform using a Service Principal and a Client Secret/Certificate

The Power Platform provider can use a Service Principal with Client Secret to authenticate to Power Platform services.

1. [Create an app registration for the Power Platform Terraform Provider](app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. Configure the provider to use a Service Principal with a Client Secret with either environment variables or using Terraform variables

    ```terraform
    provider "powerplatform" {
      client_id     = var.client_id
      tenant_id     = var.tenant_id
      client_secret = var.client_secret
    }
    ```

## Authenticating to Power Platform using Service Principal and certificate

1. [Create an app registration for the Power Platform Terraform Provider](app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. Generate a certificate using openssl or other tools

    ```bash
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365
    ```

1. Merge public and private part of the certificate files together

    Using linux shell

    ```bash
    cat *.pem > cert+key.pem
    ```

    Using Powershell

    ```powershell
    Get-Content .\cert.pem, .\key.pem | Set-Content cert+key.pem
    ```

1. Generate pkcs12 file

    ```bash
    openssl pkcs12 -export -out cert.pkcs12 -in cert+key.pem
    ```

1. Add public part of the certificate (`cert.pem` file) to the app registration
1. Store your key.pem and the password used to generate in a safe place
1. Configure the provider to use certificate with the following code:

    ```terraform
    provider "powerplatform" {
      client_id     = var.client_id
      tenant_id     = var.tenant_id
      client_certificate_file_path = "${path.cwd}/cert.pkcs12"
      client_certificate_password  = var.cert_pass
    }
    ```
