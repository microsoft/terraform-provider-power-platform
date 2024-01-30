package cli

import (
	"context"
	"fmt"

	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	common "github.com/microsoft/terraform-provider-power-platform/common"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

func Login(ctx context.Context, tenantId, username, password *string, useDeviceCode *bool, cache *common.AuthenticationCache) {

}

func TokenMode(ctx context.Context, tenantId, username, password, scope *string, cache *common.AuthenticationCache, isListAccountsMode *bool) {
	if *tenantId == "" {
		fmt.Println("Error: '--tenantid' parameter is required.")
		return
	}
	if *username == "" {
		fmt.Println("Error: '--username' parameter is required.")
		return
	}
	if *scope == "" {
		fmt.Println("Error: '--scope' parameter is required.")
		return
	}
	accounts, err := cache.GetAccounts(ctx, *tenantId)
	if err != nil {
		fmt.Println(err)
		return
	}

	searchedAccountInx := -1
	for inx, account := range accounts {
		if account.PreferredUsername == *username {
			searchedAccountInx = inx
			break
		}
	}
	if searchedAccountInx == -1 {
		fmt.Printf("Error: Account '%s' not found in cache.\nPlease log in first.", *username)
		return
	}

	authResult, err := SilentLogin(*tenantId, *scope, accounts[searchedAccountInx], cache)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Token:  %s\nExpires: %s", authResult.AccessToken, authResult.ExpiresOn)
}

func LoginMode(ctx context.Context, tenantId, username, password, scope *string, cache *common.AuthenticationCache) {
	if *tenantId == "" {
		fmt.Println("Error: '--tenantid' parameter is required.")
		return
	}
	for _, scope := range constants.REQUIRED_SCOPES {
		var err error = nil
		if *username != "" && *password != "" {
			err = UserPassLogin(*tenantId, *username, *password, scope, cache)
		} else {
			err = InteractiveLogin(*tenantId, scope, cache)
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ListAccountMode(ctx context.Context, tenantId *string, cache *common.AuthenticationCache) ([]public.Account, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority("https://login.microsoftonline.com/"+*tenantId+"/"), public.WithCache(cache))
	if err != nil {
		return nil, err
	}

	accounts, err := publicClient.Accounts(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return accounts, nil
}
