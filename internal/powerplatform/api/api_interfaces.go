package powerplatform_api

import (
	"context"

	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

type ApiClientInterface interface {
	SetAuth(auth AuthBaseOperationInterface)
	GetConfig() common.ProviderConfig

	InitializeBase(ctx context.Context) (string, error)
	ExecuteBase(ctx context.Context, token, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error)
}

type BapiClientInterface interface {
	GetBase() ApiClientInterface
	SetDataverseClient(dataverseClient DataverseClientInterface)
	Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error)

	GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error)
	GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error)
	CreateEnvironment(ctx context.Context, environment models.EnvironmentCreateDto) (*models.EnvironmentDto, error)
	UpdateEnvironment(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error)
	DeleteEnvironment(ctx context.Context, environmentId string) error

	GetPowerApps(ctx context.Context, environmentId string) ([]models.PowerAppBapi, error)

	GetConnectors(ctx context.Context) ([]models.ConnectorDto, error)
	GetPolicies(ctx context.Context) ([]models.DlpPolicyModel, error)
	GetPolicy(ctx context.Context, name string) (*models.DlpPolicyModel, error)
	DeletePolicy(ctx context.Context, name string) error
	UpdatePolicy(ctx context.Context, name string, policyToUpdate models.DlpPolicyModel) (*models.DlpPolicyModel, error)
	CreatePolicy(ctx context.Context, policyToCreate models.DlpPolicyModel) (*models.DlpPolicyModel, error)
}

type DataverseClientInterface interface {
	Initialize(ctx context.Context, environmentUrl string) (string, error)
	SetBapiClient(bapiClient BapiClientInterface)
	Execute(ctx context.Context, environmentUrl, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error)

	GetTableData(ctx context.Context, environmentId, tableName, odataQuery string, responseObj interface{}) error

	GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error)
	CreateSolution(ctx context.Context, environmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error)
	GetSolution(ctx context.Context, environmentId string, solutionName string) (*models.SolutionDto, error)
	DeleteSolution(ctx context.Context, environmentId string, solutionName string) error

	GetDefaultCurrencyForEnvironment(ctx context.Context, environmentId string) (*models.TransactionCurrencyDto, error)
}

type PowerPlatformClientApiInterface interface {
	GetBase() ApiClientInterface
	Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error)

	GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error)
}
