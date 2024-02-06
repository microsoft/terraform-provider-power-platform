package cli

import (
	"context"
	"fmt"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

func SilentLogin(ctx context.Context, scope []string, account public.Account, cache cache.ExportReplace) ([]public.AuthResult, error) {
	publicClient, err := createPublicClient(&account.Realm, cache)
	if err != nil {
		return nil, err
	}

	var authResults []public.AuthResult
	for _, scope := range scope {
		var authResult public.AuthResult
		authResult, err = publicClient.AcquireTokenSilent(ctx, []string{scope}, public.WithTenantID(account.Realm), public.WithSilentAccount(account))
		if err != nil {
			return nil, err
		}
		authResults = append(authResults, authResult)
	}
	return authResults, nil
}

func DeviceCodeLogin(ctx context.Context, tenantId *string, scope []string, cache cache.ExportReplace) ([]public.AuthResult, error) {
	publicClient, err := createPublicClient(tenantId, cache)
	if err != nil {
		return nil, err
	}

	var authResults []public.AuthResult
	for _, scope := range scope {
		var authResult public.AuthResult
		var deviceCode public.DeviceCode
		if tenantId == nil || *tenantId == "" {
			deviceCode, err = publicClient.AcquireTokenByDeviceCode(ctx, []string{scope})

		} else {
			deviceCode, err = publicClient.AcquireTokenByDeviceCode(ctx, []string{scope}, public.WithTenantID(*tenantId))
		}

		if err != nil {
			return nil, err
		}
		fmt.Printf("To sign in, use a web browser to open the page https://aka.ms/devicecode and enter the code %s to authenticate.\n", deviceCode.Result.UserCode)

		authResult, err = deviceCode.AuthenticationResult(ctx)
		if err != nil {
			return nil, err
		}

		authResults = append(authResults, authResult)
	}
	return authResults, nil
}

func UserPassLogin(ctx context.Context, tenantId, username, password *string, scope []string, cache cache.ExportReplace) ([]public.AuthResult, error) {
	publicClient, err := createPublicClient(tenantId, cache)
	if err != nil {
		return nil, err
	}

	var authResults []public.AuthResult
	for _, scope := range scope {
		var authResult public.AuthResult
		if tenantId == nil || *tenantId == "" {
			authResult, err = publicClient.AcquireTokenByUsernamePassword(ctx, []string{scope}, *username, *password)

		} else {
			authResult, err = publicClient.AcquireTokenByUsernamePassword(ctx, []string{scope}, *username, *password, public.WithTenantID(*tenantId))
		}
		if err != nil {
			return nil, err
		}
		authResults = append(authResults, authResult)
	}
	return authResults, nil
}

func InteractiveLogin(ctx context.Context, tenantId *string, scope []string, cache cache.ExportReplace) ([]public.AuthResult, error) {

	fmt.Println("Starting interactive login flow")

	publicClient, err := createPublicClient(tenantId, cache)
	if err != nil {
		return nil, err
	}

	fmt.Println("Created public client")

	var authResults []public.AuthResult
	for _, scope := range scope {
		var authResult public.AuthResult
		if tenantId == nil || *tenantId == "" {

			fmt.Println("Acquiring token interactive without tenant id")

			authResult, err = publicClient.AcquireTokenInteractive(ctx, []string{scope})

		} else {

			fmt.Println("Acquiring token interactive with tenant id")

			authResult, err = publicClient.AcquireTokenInteractive(ctx, []string{scope}, public.WithTenantID(*tenantId))
		}

		if err != nil {
			return nil, err
		}

		fmt.Println("Acquired token interactive")

		authResults = append(authResults, authResult)
	}
	return authResults, nil
}

func createPublicClient(tenantId *string, cache cache.ExportReplace) (*public.Client, error) {
	var err error = nil
	var publicClient public.Client
	if tenantId == nil {
		publicClient, err = public.New(
			constants.CLIENT_ID,
			public.WithCache(cache))
	} else {
		publicClient, err = public.New(
			constants.CLIENT_ID,
			public.WithAuthority(constants.OAUTH_AUTHORITY_URL+*tenantId),
			public.WithCache(cache))
	}

	if err != nil {
		return nil, err
	}
	return &publicClient, nil
}
