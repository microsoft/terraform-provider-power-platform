package auth

import "context"

type PowerPlatformApiTokenProvider struct {
	accessToken  string
	refreshToken string
	expires_in   int64
	tokenType    string
}

func NewPowerPlatformApiTokenProvider(ctx context.Context) *PowerPlatformApiTokenProvider {
	return &PowerPlatformApiTokenProvider{
		
	}
}
