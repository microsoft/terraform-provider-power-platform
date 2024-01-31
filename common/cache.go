package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

type AuthenticationCache struct {
	FileProtectData *FileProtectData
	CacheContent    *CacheContent
}

func NewAuthenticationCache() *AuthenticationCache {
	return &AuthenticationCache{
		FileProtectData: &FileProtectData{},
		CacheContent:    &CacheContent{},
	}
}

type CacheContent struct {
	DefaultAccount   string `json:"default_account"`
	MsalCacheContent string `json:"msal_cache_content"`
}

func (c *AuthenticationCache) GetCacheFilePath() (string, error) {
	dir, err := c.FileProtectData.GetOrCreateCacheFileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, constants.MSAL_CACHE_FILE_NAME), nil
}

func (c *AuthenticationCache) GetAccounts(ctx context.Context) ([]public.Account, error) {
	publicClient, err := public.New(constants.CLIENT_ID, public.WithCache(c))
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

func (c *AuthenticationCache) GetDefaultAccount(ctx context.Context) (*public.Account, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.PreferredUsername == c.CacheContent.DefaultAccount {
			return &account, nil
		}
	}
	return nil, nil
}

func (c *AuthenticationCache) SetDefaultAccount(ctx context.Context, account public.Account) error {
	c.CacheContent.DefaultAccount = account.PreferredUsername
	contentBytes, err := json.Marshal(c.CacheContent)
	if err != nil {
		return err
	}
	return c.writeProtectedFile(ctx, contentBytes)
}

func (c *AuthenticationCache) writeProtectedFile(ctx context.Context, contentBytes []byte) error {

	cacheFilePath, err := c.GetCacheFilePath()
	if err != nil {
		return err
	}

	encryptedData, err := c.FileProtectData.Encrypt(contentBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFilePath, encryptedData, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *AuthenticationCache) readProtectedFile(ctx context.Context) ([]byte, error) {

	cacheFilePath, err := c.GetCacheFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, nil
	}

	contentBytes, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	decryptedData, err := c.FileProtectData.Decrypt(contentBytes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(decryptedData, c.CacheContent)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

func (c *AuthenticationCache) DeleteFile(ctx context.Context) error {
	cacheFilePath, err := c.GetCacheFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(cacheFilePath)
}

func (c *AuthenticationCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {

	contentByes, err := c.readProtectedFile(ctx)
	if err != nil {
		return err
	}

	if contentByes == nil {
		emptyContent := []byte("{}")
		contentByes, err = c.FileProtectData.Encrypt(emptyContent)
		if err != nil {
			return err
		}
		err = c.writeProtectedFile(ctx, contentByes)
		if err != nil {
			return err
		}
		json.Unmarshal(emptyContent, c.CacheContent)
	}

	err = json.Unmarshal(contentByes, c.CacheContent)
	if err != nil {
		return err
	}

	err = cache.Unmarshal([]byte(c.CacheContent.MsalCacheContent))
	if err != nil {
		return err
	}
	return nil
}

func (c *AuthenticationCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) error {

	msalContentBytes, err := cache.Marshal()
	if err != nil {
		return err
	}
	c.CacheContent.MsalCacheContent = string(msalContentBytes)

	contentBytes, err := json.Marshal(c.CacheContent)
	if err != nil {
		return err
	}

	err = c.writeProtectedFile(ctx, contentBytes)
	if err != nil {
		return err
	}
	return nil
}
