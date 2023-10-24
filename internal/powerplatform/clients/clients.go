package clients

import (
	powerplatform "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type ProviderClient struct {
	Config           *common.ProviderConfig
	BapiApi          *BapiClient
	DataverseApi     *DataverseClient
	PowerPlatformApi *PowerPlatoformApiClient
}

type BapiClient struct {
	Auth   *powerplatform.BapiAuth
	Client *powerplatform.BapiClientApi
}

func NewBapiClient(auth *powerplatform.BapiAuth, client *powerplatform.BapiClientApi) *BapiClient {
	return &BapiClient{
		Auth:   auth,
		Client: client,
	}
}

type DataverseClient struct {
	Auth   *powerplatform.DataverseAuth
	Client *powerplatform.DataverseClientApi
}

func NewDataverseClient(auth *powerplatform.DataverseAuth, client *powerplatform.DataverseClientApi) *DataverseClient {
	return &DataverseClient{
		Auth:   auth,
		Client: client,
	}
}

type PowerPlatoformApiClient struct {
	Auth   *powerplatform.PowerPlatformAuth
	Client *powerplatform.PowerPlatformClientApi
}

func NewPowerPlatformApiClient(auth *powerplatform.PowerPlatformAuth, client *powerplatform.PowerPlatformClientApi) *PowerPlatoformApiClient {
	return &PowerPlatoformApiClient{
		Auth:   auth,
		Client: client,
	}
}
