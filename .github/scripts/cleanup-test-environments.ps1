# Deletes every environment whose display name starts with "Test", which is the prefix the
# acceptance tests give the environments they create. An acceptance test does not always destroy
# its environment: a create that outlives its timeout or a delete that does not finish leaves it
# behind in the tenant, so the workflow runs this before and after the tests.
# Requires an authenticated pac CLI profile.

$output = pac admin list --name "Test" --json
if ($LASTEXITCODE -eq 0 -and $output) {
  try {
    $environmentsList = $output | ConvertFrom-Json
    $environmentsList | ForEach-Object -Parallel {
        if ($_.DisplayName.StartsWith("Test")){
            Write-Output $_.DisplayName
            pac admin delete -env $_.EnvironmentID
        }
    } -ThrottleLimit 15
  } catch {
    Write-Output "Failed to parse environments list: $_"
    Write-Output "Raw output: $output"
  }
} else {
  Write-Output "Failed to list environments or no output received. Exit code: $LASTEXITCODE"
  Write-Output "Output: $output"
}
