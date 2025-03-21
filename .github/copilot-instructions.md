# Terraform Provider Power Platform Instructions

Use `make install` to compile the code
Use `make lint` to run the linter
Use `make unittest` to run all unit tests.
The TEST parameter for unit tests and acceptance tests (e.g. `make unittest TEST=TestPrefix`) can be used to run test functions that match `^(TestAcc|TestUnit)TestPrefix`
Only run acceptance tests with user consent, and always specify a test prefix when running `make acctest TEST=TestPrefix`
When running acceptance tests you don't need to specify `TF_ACC=1` because `make acctest` already does that
Always run `make` commands in `/workspaces/terraform-provider-power-platform` working directory
Don't ever run `terraform init`
Do not direcly edit the files under `/docs` because they are auto-generated from MarkdownDescription on schemas using `make userdocs`
test files should have `_test` appended to their package name
To run an example, `cd` to its working directory and run `terraform apply -auto-approve`

When creating Schema for the resource or datasource, use MarkdownDescription and never use Description attribute
Methods that are not used outside the namespace scope, should be kept private.
Helper methods that covert DTO to model and model to DTO should be in models.go file
The DTO structures should always have `Dto` suffix and be in models.go file
Use tflog.Debug for logging unless there is something really important (tflog.Info) or an error/warning
Comments on methods should provide information about how to use it, its parameters, and expected results. Omit comments that don't substantially improve the readability of the code.
The functions to convert Dto and model objects should be always named `convertDtoToModel` and `convertModelToDto` and should be placed in models.go file

## Testing

When writing unit tests for resources you must register mock responders for every step of the process:

- Test steps will call the `Create`, `Read`, `Update`, and `Delete` methods in the resource.  All the API calls made in those functions need to be mocked for each time the operation is called.
- First test step will call `Create` then `Read` methods. Test expects that JSON from read include the changes applied in create
- Subsequent test steps will call `Read`, then `Update`, then `Read`. First read should match what was read at the end of previous step.  The final read should include the changes applied in the update step.
- Steps that delete a resource (or omit a previously configured resource from config) will call `Delete` then `Read`.  Unless otherwise specified the final read will return a 404
- Tests that import state should be a single test step and only call `Read` using the `ImportStateId`
- All the JSON response for unit test should be stored in .json files.
  - Files should be placed in a folder with a name coresponding to the Unit Test that is being used. Folder name should ommit `UnitTest` in its name.
  - Each Unit Test folder with .json files be stored at `services\{service_name}\test\resource` or `services\{service_name}\datasource` with all other resource and/or datasource .go files.
  - The .json file name should consist of the mock request method (`get`, 'post', `delete`) followed by `_` and name of the returned mock object name or action.
  - The file names have to be sensible without empty spaces and special characters.
