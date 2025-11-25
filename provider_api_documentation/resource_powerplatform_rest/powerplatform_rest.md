# `powerplatform_rest`

This resource is used to execute custom HTTP operations against Power Platform or Dataverse Web APIs as part of a Terraform-managed lifecycle. It supports separate `create`, `read`, `update`, and `destroy` operations, each defined as an independent web API call.

## API Endpoints

Because this resource allows you to specify arbitrary URLs and methods per operation, there is no single fixed endpoint. Instead, each operation executes the URL and HTTP method that you configure.

| Resource Block | HTTP Method / Endpoint Behavior |
| -------------- | -------------------------------- |
| `create`       | Executes the HTTP request defined in the `create` block (for example, `POST` to a Dataverse Web API entity URL). |
| `read`         | Executes the HTTP request defined in the `read` block to retrieve the current state. |
| `update`       | Executes the HTTP request defined in the `update` block whenever Terraform detects changes that require an update. |
| `destroy`      | Executes the HTTP request defined in the `destroy` block during resource deletion. |

## Attribute Mapping

| Resource Attribute                      | API Request / Response Field |
| --------------------------------------- | ----------------------------- |
| `id`                                    | Terraform-only identifier (timestamp-based string). |
| `create.scope`                          | Authentication scope for the `create` HTTP request. |
| `create.method`                         | HTTP method used for the `create` request. |
| `create.url`                            | Absolute URL used for the `create` request. |
| `create.body`                           | Request body sent during the `create` operation. |
| `create.headers[*].name`                | Header name sent with the `create` request. |
| `create.headers[*].value`               | Header value sent with the `create` request. |
| `create.expected_http_status`           | List of HTTP status codes considered successful for the `create` operation. |
| `read.*`, `update.*`, `destroy.*`       | Same as the `create.*` attributes, but applied to their respective operations. |
| `output.body`                           | Entire HTTP response body from the last executed operation as a string (may contain JSON, XML, or plain text). |

### Example API Response

An example of the HTTP response body captured by this resource (for a `POST` that creates an `account` record via Dataverse Web API) can be found in the test fixture [`rest/tests/resource/Web_Api_Validate_Create/post_account.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/rest/tests/resource/Web_Api_Validate_Create/post_account.json).
