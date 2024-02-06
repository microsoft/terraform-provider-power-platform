# Terraform Power Platform Provider CLI

This is a CLI for the Terraform Power Platform Provider. It is used to authenticate with the Power Platform and generate a refresh and access token for use with the Terraform Power Platform Provider with `use_cli="true"` provider configuration option.

## Usage

### Logging in options

1. Logging in interactively (a popup window will be opened in the browser multiple times, please login with the same account/tenant for all of them)

    ```bash
    terraform-provider-power-platform login --tenant <tenant_id>
    ```

1. Logging using username/password

    ```bash
    terraform-provider-power-platform login --tenant <tenantid>  --user <username> --password <password>
    ```

1. Logging with device code

    ```bash
    terraform-provider-power-platform --tenant <tenantid> --use-device-code
    ```

### Listing accounts saved in cache

After you have logged in, you can list the accounts saved in the cache by running the following command. When using Terraform Power Platform Provider with `use_cli="true"` provider configuration option, it is important to note that **first account from the list will be used in the provider**.

```bash
terraform-provider-power-platform account --list
```

### Removing accounts saved in cache

```bash
terraform-provider-power-platform account --clear
```

### Setting default authentication account

If there is more than single account saved in the cache, you can change which one will be use for authentication by default by running the following command:

```bash
terraform-provider-power-platform account --set <index>
```

The index is the number of the account in the list when using `terraform-provider-power-platform account --list` command.

### Getting access token for a given scope

After you have logged in, you can get the access token for a given scope by running the following command:

```bash
terraform-provider-power-platform account --get-access-token --scope <scope>
```

## Important Notes

- For windows platforms, CLI will store the cache in `%APPDATA%\Microsoft\Terraform Power Platform Provider\terraform_power_platform_cache.dat` and **will be** encrypted using [Data Protection API](https://en.wikipedia.org/wiki/Data_Protection_API) in context of the current user. This means that the cache will be available only for the current user on the current machine.

- For non-windows platforms, CLI will store the cache in `/home/<user>/.local/share/Microsoft/TerraformPowerPlatformProvider/terraform_power_platform_cache.dat` and **will not be** encrypted. Only `chmod 600` will be applied to the file.

- When using Terraform Power Platform Provider with `use_cli="true"` provider configuration option, it is important to note that if you login with CLI using different accounts/tenants only **first account from the list will be used in the provider**. If you want to use different account/tenant, you need to remove the cache file and login again.
