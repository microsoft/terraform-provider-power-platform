package cli

import (
	"context"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

func SilentLogin(tenantId, scope string, account public.Account, cache cache.ExportReplace) (*public.AuthResult, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority("https://login.microsoftonline.com/"+tenantId+"/"), public.WithCache(cache))
	if err != nil {
		return nil, err
	}

	authResult, err := publicClient.AcquireTokenSilent(context.Background(), []string{
		scope,
	},
		public.WithTenantID(tenantId),
		public.WithSilentAccount(account))

	if err != nil {
		return nil, err
	}

	return &authResult, nil
}

func UserPassLogin(tenantId, username, password, scope string, cache cache.ExportReplace) error {
	publicClient, err := public.New(
		constants.CLIENT_ID,
		public.WithAuthority("https://login.microsoftonline.com/"+tenantId+"/"),
		public.WithCache(cache))

	if err != nil {
		return err
	}

	_, err = publicClient.AcquireTokenByUsernamePassword(
		context.Background(),
		[]string{
			scope,
		},
		username,
		password,
		public.WithTenantID(tenantId),
	)
	if err != nil {
		return err
	}
	return nil
}

func InteractiveLogin(tenantId, scope string, cache cache.ExportReplace) error {
	publicClient, err := public.New(
		constants.CLIENT_ID,
		public.WithAuthority("https://login.microsoftonline.com/"+tenantId+"/"),
		public.WithCache(cache))

	if err != nil {
		return err
	}

	_, err = publicClient.AcquireTokenInteractive(
		context.Background(),
		[]string{
			scope,
		},
		public.WithTenantID(tenantId),
	)
	if err != nil {
		return err
	}
	return nil
}
