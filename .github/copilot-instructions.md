Use `make install` to compile the code
Use `make lint` to run the linter
Use `make unittest` to run all unit tests.
The TEST parameter for unit tests and acceptance tests (e.g. `make unittest TEST=TestPrefix`) can be used to run test functions that match `^(TestAcc|TestUnit)TestPrefix`
Only run acceptance tests with user consent, and always specify a test prefix when running `make acctest TEST=TestPrefix`
When running acceptance tests you don't need to specify `TF_ACC=1` because `make acctest` already does that
Always run `make` commands in `/workspaces/terraform-provider-power-platform` working directory
Don't ever run `terraform init`
test files should have `_test` appended to their package name
To run an example, `cd` to its working directory and run `terraform apply -auto-approve`
When writing unit tests for resources you must register mock responders for every step of the process
- Test steps will call the `Create`, `Read`, `Update`, and `Delete` methods in the resource.  All the API calls made in those functions need to be mocked for each time the operation is called.
- First test step will call `Create` then `Read` methods. Test expects that JSON from read include the changes applied in create
- Subsequent test steps will call `Read`, then `Update`, then `Read`. First read should match what was read at the end of previous step.  The final read should include the changes applied in the update step.
- Steps that delete a resource (or omit a previously configured resource from config) will call `Delete` then `Read`.  Unless otherwise specified the final read will return a 404
- Tests that import state should be a single test step and only call `Read` using the `ImportStateId`
