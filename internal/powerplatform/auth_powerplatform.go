package powerplatform

import (
	"context"
	"time"
)

var _ PowerPlatformAuthInterface = &PowerPlatformAuthImplementation{}

type PowerPlatformAuthInterface interface {
	IsTokenExpiredOrEmpty() bool
	RefreshToken() (string, error)

	AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error)
	AuthenticateClientSecret(ctx context.Context, tenantId, applicationid, secret string) (string, error)
}

type PowerPlatformAuthImplementation struct {
	Token       string
	TokenExpiry time.Time
}

func (client *PowerPlatformAuthImplementation) IsTokenExpiredOrEmpty() bool {
	if client.Token == "" {
		return true
	} else {
		return time.Now().After(client.TokenExpiry)
	}
}

func (client *PowerPlatformAuthImplementation) RefreshToken() (string, error) {
	//todo implement token refresh
	panic("[RefreshToken] not implemented")
}

func (client *PowerPlatformAuthImplementation) AuthenticateUserPass(ctx context.Context, tenantId, username, password string) (string, error) {
	//todo implement
	panic("[AuthenticateUserPass] not implemented")
}

func (client *PowerPlatformAuthImplementation) AuthenticateClientSecret(ctx context.Context, tenantId, applicationId, secret string) (string, error) {
	//todo implement
	panic("[AuthenticateClientSecret] not implemented")
}
