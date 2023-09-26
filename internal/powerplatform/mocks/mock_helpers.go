package powerplatform_mocks

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
)

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
