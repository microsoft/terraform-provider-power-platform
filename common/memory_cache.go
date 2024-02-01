package common

import (
	"context"

	cache "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
)

type MemoryCache struct {
	cache string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

// func (c *MemoryCache) GetAccounts(ctx context.Context, tenantId string) ([]public.Account, error) {
// 	publicClient, err := public.New(constants.CLIENT_ID, public.WithAuthority("https://login.microsoftonline.com/"+tenantId+"/"), public.WithCache(c))
// 	if err != nil {
// 		return nil, err
// 	}

// 	accounts, err := publicClient.Accounts(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	return accounts, nil
// }

func (c *MemoryCache) Replace(ctx context.Context, cache cache.Unmarshaler, hints cache.ReplaceHints) error {
	if len(c.cache) != 0 {
		err := cache.Unmarshal([]byte(c.cache))
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
	c.cache = string(contentBytes)
	return nil
}
