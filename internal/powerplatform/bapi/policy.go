package powerplatform_bapi

import (
	"context"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetPolicies(ctx context.Context) ([]models.DlpPolicyDto, error) {
	return nil, nil
}

func (client *ApiClient) GetPolicy(ctx context.Context, name string) (*models.DlpPolicyDto, error) {
	return nil, nil
}

func (client *ApiClient) DeletePolicy(ctx context.Context, name string) error {
	return nil
}

func (client *ApiClient) UpdatePolicy(ctx context.Context, name string, PolicyDtoToUpdate models.DlpPolicyDto) (*models.DlpPolicyDto, error) {
	return nil, nil
}

func (client *ApiClient) CreatePolicy(ctx context.Context, PolicyDtoToCreate models.DlpPolicyDto) (*models.DlpPolicyDto, error) {
	return nil, nil
}
