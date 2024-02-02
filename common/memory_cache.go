package common

import (
	"context"
	"fmt"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	public "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	constants "github.com/microsoft/terraform-provider-power-platform/constants"
)

type MemoryCache struct {
	content string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		content: "{}",
	}
}

var _ ExportReplaceCacheExtension = &MemoryCache{}

func (c *MemoryCache) GetAccounts(ctx context.Context) ([]public.Account, error) {
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

func (c *MemoryCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	if len(c.content) != 0 {
		err := cache.Unmarshal([]byte(c.content))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MemoryCache) Export(ctx context.Context, cache cache.Marshaler, hints cache.ExportHints) error {
	contentBytes, err := cache.Marshal()
	if err != nil {
		return err
	}
	c.content = string(contentBytes)
	return nil
}
