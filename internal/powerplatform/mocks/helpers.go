package powerplatform_mocks

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
)

// func NewUnitTestsMockApiClientInterface(t *testing.T) *MockApiClientInterface {
// 	ctrl := gomock.NewController(t)
// 	clientMock := NewMockApiClientInterface(ctrl)

// 	// authResponse := powerplatform_bapi.AuthResponse{
// 	// 	Token: "mock_token",
// 	// }

// 	// clientMock.EXPECT().DoAuthClientSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tenantId, applicationId, clientSecret string) (*powerplatform_bapi.AuthResponse, error) {
// 	// 	return &authResponse, nil
// 	// }).AnyTimes()
// 	// clientMock.EXPECT().DoAuthUsernamePassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tenantId, username, password string) (*powerplatform_bapi.AuthResponse, error) {
// 	// 	return &authResponse, nil
// 	// }).AnyTimes()

// 	return clientMock
// }

func NewUnitTestsMockBapiClientInterface(t *testing.T) *MockBapiClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockBapiClientInterface(ctrl)

	return clientMock
}

func NewUnitTestMockDataverseClientInterface(t *testing.T) *MockDataverseClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockDataverseClientInterface(ctrl)

	return clientMock
}

func NewUnitTestMockPowerPlatformClientInterface(t *testing.T) *MockPowerPlatformClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockPowerPlatformClientInterface(ctrl)

	return clientMock
}
