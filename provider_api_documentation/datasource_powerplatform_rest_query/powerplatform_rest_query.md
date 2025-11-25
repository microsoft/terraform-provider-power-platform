# `powerplatform_rest_query` (Data Source)

This data source is used to execute arbitrary HTTP requests against Power Platform or Dataverse REST APIs and return the raw response body.

## API Endpoints

Because this data source allows you to specify any absolute URL, HTTP method, headers, and body, there is no single fixed endpoint. Instead, the request is sent to the URL you provide in the `url` argument.

| Data Source Argument | HTTP Method / Endpoint Behavior |
| -------------------- | -------------------------------- |
| `method`             | The HTTP method used for the request (for example, `GET`, `POST`, `PATCH`, `DELETE`). |
| `url`                | The absolute URL of the API call to execute. |
| `headers`            | Optional list of HTTP headers that will be sent with the request. |
| `body`               | Optional request body that will be sent with the request. |
| `expected_http_status` | Optional list of HTTP status codes that are considered successful. If the response status code does not match any of these, the data source will return an error. |

## Attribute Mapping

| Data Source Attribute | API Response JSON Field |
| --------------------- | ----------------------- |
| `output.body`         | Entire HTTP response body as a string (may contain JSON, XML, or plain text depending on the target API). |

### Example API Response

An example of the API response body captured by this data source (for a `WhoAmI` call against Dataverse Web API) can be found in the test fixture [`rest/tests/datasource/Web_Apis_WhoAmI/get_whoami.json`](https://github.com/microsoft/terraform-provider-power-platform/blob/main/internal/services/rest/tests/datasource/Web_Apis_WhoAmI/get_whoami.json).
