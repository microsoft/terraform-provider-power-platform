package powerplatform_mocks

import (
	context "context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"
	powerplatform_bapi "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/bapi"
)

// TODO Remove
func NewUnitTestsMockClientInterface(t *testing.T) *MockClientInterface {
	ctrl := gomock.NewController(t)
	return NewMockClientInterface(ctrl)
}

// TODO remove
func DoUnitTestsBasicAuth(clientMock *MockClientInterface) {
	a := &powerplatform.AuthResponse{
		AuthHash: "auth_hash_placeholder",
	}
	clientMock.EXPECT().DoBasicAuth(gomock.Any(), gomock.Any(), gomock.Any()).Return(a, nil).AnyTimes()
}

func NewUnitTestsMockApiClientInterface(t *testing.T) *MockApiClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockApiClientInterface(ctrl)

	authResponse := powerplatform_bapi.AuthResponse{
		Token: "mock_token",
	}

	clientMock.EXPECT().DoAuthClientSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tenantId, applicationId, clientSecret string) (*powerplatform_bapi.AuthResponse, error) {
		return &authResponse, nil
	}).AnyTimes()
	clientMock.EXPECT().DoAuthUsernamePassword(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tenantId, username, password string) (*powerplatform_bapi.AuthResponse, error) {
		return &authResponse, nil
	}).AnyTimes()

	return clientMock
}
