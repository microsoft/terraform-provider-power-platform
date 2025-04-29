# Title

Potential infinite loop in the `AddDataverseToEnvironment` method.

##

`/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

## Problem

The `AddDataverseToEnvironment` method contains a `for {}` loop that polls for the provisioning state of a Dataverse until it reaches the `Succeeded` state. However, there is no explicit timeout or maximum retry limit in place, which could lead to an infinite loop in case the provisioning state never reaches the desired condition or gets stuck in another state like `LinkedDatabaseProvisioning`.

## Impact

This can result in resource exhaustion, as the loop could endlessly occupy system resources and block program execution. It also makes the application less reliable, as there is no fallback mechanism to handle scenarios where the desired provisioning state is never reached. Severity: High.

## Location

The problem occurs in the following method:

```go
func (client *Client) AddDataverseToEnvironment(ctx context.Context, environmentId string, environmentCreateLinkEnvironmentMetadata createLinkEnvironmentMetadataDto) (*EnvironmentDto, error)
```

## Code Issue

```go
	for {
		lifecycleEnv := EnvironmentDto{}
		lifecycleResponse, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, &lifecycleEnv)
		if err != nil {
			return nil, err
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, "Dataverse Creation Operation State: '"+lifecycleEnv.Properties.ProvisioningState+"'")
		tflog.Debug(ctx, "Dataverse Creation Operation HTTP Status: '"+lifecycleResponse.HttpResponse.Status+"'")

		if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
			return &lifecycleEnv, nil
		} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
			return &lifecycleEnv, errors.New("dataverse creation failed. provisioning state: " + lifecycleEnv.Properties.ProvisioningState)
		}
	}
```

## Fix

Introduce a timeout mechanism using a context with a deadline or a retry counter to break out of the loop if the desired state is not achieved within a reasonable timeframe.

### Fix Using Retry Counter

```go
	retryCount := 0
	maxRetries := 10 // Specify a suitable maximum retry count.
	for retryCount < maxRetries {
		lifecycleEnv := EnvironmentDto{}
		lifecycleResponse, err := client.Api.Execute(ctx, nil, "GET", locationHeader, nil, nil, []int{http.StatusOK, http.StatusAccepted}, &lifecycleEnv)
		if err != nil {
			return nil, err
		}

		err = client.Api.SleepWithContext(ctx, retryAfter)
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, "Dataverse Creation Operation State: '"+lifecycleEnv.Properties.ProvisioningState+"'")
		tflog.Debug(ctx, "Dataverse Creation Operation HTTP Status: '"+lifecycleResponse.HttpResponse.Status+"'")

		if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
			return &lifecycleEnv, nil
		} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
			return &lifecycleEnv, errors.New("dataverse creation failed. provisioning state: " + lifecycleEnv.Properties.ProvisioningState)
		}

		retryCount++
	}

	return nil, errors.New("max retries reached while waiting for dataverse provisioning state")
```

This fix ensures that the loop will terminate after a certain number of retries, preventing infinite execution and improving application reliability.
