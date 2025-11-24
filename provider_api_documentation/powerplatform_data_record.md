---
page_title: "powerplatform_data_record Resource - powerplatform"
description: |-
  The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.
---

# powerplatform_data_record (Resource)

The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.

Data Record is a special type of a resources, that allows creation of any type Dataverese table record. The syntax for working with `data_record` resource is simmilar to raw WebAPI HTTP requests that this record uses:

- [WebAPI overview - Power Platform | Microsoft Learn](https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/overview)

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud)                                                              |
| ------------------- | ----------- | -------------------------------------------------------------------------------------------- |
| Create              | `POST`      | `https://<environment_host>/api/data/<dataverse_api_version>/<entity_collection_name>`         |
| Read                | `GET`       | `https://<environment_host>/api/data/<dataverse_api_version>/<entity_collection_name>(<record_id>)` |
| Update              | `PATCH`     | `https://<environment_host>/api/data/<dataverse_api_version>/<entity_collection_name>(<record_id>)` |
| Delete              | `DELETE`    | `https://<environment_host>/api/data/<dataverse_api_version>/<entity_collection_name>(<record_id>)` |

## Attribute Mapping

The `powerplatform_data_record` resource does not have a fixed attribute mapping, as it is a generic resource that can be used to manage any type of Dataverse record. The attributes are dynamically determined based on the `columns` attribute of the resource.

The `columns` attribute is a map of key-value pairs, where the key is the logical name of the attribute and the value is the value of the attribute. The value can be a string, number, boolean, or a map for lookups.

For lookup attributes, the value should be a map with the following keys:

- `table_logical_name`: The logical name of the table to which the lookup is pointing.
- `data_record_id`: The ID of the record to which the lookup is pointing.

For relations, the value should be a list of maps, where each map has the following keys:

- `table_logical_name`: The logical name of the table to which the relation is pointing.
- `data_record_id`: The ID of the record to which the relation is pointing.

### Example API Response

An example of the API response can be found [here](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/data_record/tests/resource/Validate_Create/get_contact_00000000-0000-0000-0000-000000000010.json).
