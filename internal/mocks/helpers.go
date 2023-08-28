package powerplatform_mocks

import (
	context "context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
)

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
