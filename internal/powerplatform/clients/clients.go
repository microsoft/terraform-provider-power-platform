package clients

import (
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	dvapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/dataverse"
	ppapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/ppapi"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type ProviderClient struct {
	Config           *common.ProviderConfig
	BapiApi          *BapiClient
	DataverseApi     *DataverseClient
	PowerPlatformApi *PowerPlatoformApiClient
}

type BapiClient struct {
	Auth   bapi.BapiAuthInterface
	Client api.BapiClientInterface
}

type DataverseClient struct {
	Auth   dvapi.DataverseAuthInterface
	Client api.DataverseClientInterface
}

type PowerPlatoformApiClient struct {
	Auth   ppapi.PowerPlatformAuthInterface
	Client api.PowerPlatformClientApiInterface
}
