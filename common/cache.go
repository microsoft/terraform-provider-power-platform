package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

type AuthenticationCache struct {
	FileProtectData *FileProtectData
}

func NewAuthenticationCache() *AuthenticationCache {
	return &AuthenticationCache{
		FileProtectData: &FileProtectData{},
	}
}

func (c *AuthenticationCache) GetCacheFilePath() (string, error) {
	dir, err := c.FileProtectData.GetOrCreateCacheFileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, constants.MSAL_CACHE_FILE_NAME), nil
}

func (c *AuthenticationCache) GetAccounts(ctx context.Context, tenantId string) ([]public.Account, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority("https://login.microsoftonline.com/"+tenantId+"/"), public.WithCache(c))
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

func (c *AuthenticationCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	cacheFilePath, err := c.GetCacheFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		encryptedData, err := c.FileProtectData.Encrypt([]byte("{}"))
		if err != nil {
			return err
		}
		os.WriteFile(cacheFilePath, encryptedData, 0600)
	}
	contentBytes, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return err
	}

	decryptedData, err := c.FileProtectData.Decrypt(contentBytes)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = cache.Unmarshal(decryptedData)
	if err != nil {
		return err
	}
	return nil
}

func (c *AuthenticationCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) error {
	cacheFilePath, err := c.GetCacheFilePath()
	if err != nil {
		return err
	}

	contentBytes, err := cache.Marshal()
	if err != nil {
		return err
	}

	encryptedData, err := c.FileProtectData.Encrypt(contentBytes)
	if err != nil {
		return err
	}

	os.WriteFile(cacheFilePath, encryptedData, 0600)

	return nil
}
