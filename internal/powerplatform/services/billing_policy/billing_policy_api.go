package billing_policy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	powerplatformapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/ppapi"
)

type BillingPolicyClientInterface interface {
	GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error)
	GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error)
	CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error)
}

type BillingPolicyClient struct {
	ppapi powerplatformapi.PowerPlatformClientApiInterface
}

const (
	API_VERSION = "2022-03-01-preview"
)

func NewBillingPolicyClient(ctx context.Context) BillingPolicyClientInterface {
	return &BillingPolicyClient{}
}

func (client *BillingPolicyClient) GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerplatform.com",
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", id),
	}

	policy := BillingPolicyDto{}
	_, err := client.ppapi.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *BillingPolicyClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerplatform.com",
		Path:   "/licensing/billingPolicies",
	}

	policies := []BillingPolicyDto{}
	_, err := client.ppapi.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policies)

	return policies, err
}

func (client *BillingPolicyClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}

func (client *BillingPolicyClient) UpdateBillingPolicy(ctx context.Context, policyToUpdate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}

func (client *BillingPolicyClient) DeleteBillingPolicy(ctx context.Context, id string) error {
	return nil
}
