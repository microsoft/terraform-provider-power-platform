
## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.bap.microsoft.com/api/environments/{organizationId}/features?geo={geo}` |
| Update              | `POST`      | `https://api.bap.microsoft.com/api/environments/{organizationId}/features/{featureName}/enable?geo={geo}` |

## Attribute Mapping

| Resource Attribute | API Response JSON Field |
| ------------------ | ----------------------- |
| `id`               | composed as `{environment_id}/{feature_name}` |
| `environment_id`   | from Terraform configuration                    |
| `feature_name`     | from Terraform configuration                    |
| `state`            | `appsUpgradeState`                              |

### Example API Response

Examples of API responses used by this resource can be found in the test fixtures:

- [`environment_wave/test/resource/EnvironmentWaveResource_Create/get_organizations.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_wave/test/resource/EnvironmentWaveResource_Create/get_organizations.json)
- [`environment_wave/test/resource/EnvironmentWaveResource_Create/get_features_upgrading.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_wave/test/resource/EnvironmentWaveResource_Create/get_features_upgrading.json)
- [`environment_wave/test/resource/EnvironmentWaveResource_Create/get_features_enabled.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/environment_wave/test/resource/EnvironmentWaveResource_Create/get_features_enabled.json)
