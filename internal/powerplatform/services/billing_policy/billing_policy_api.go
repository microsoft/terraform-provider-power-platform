package billing_policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

type BillingPolicyClientInterface interface {
	GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error)
	GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error)
	CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error)
}

type BillingPolicyClient struct {
	ppapi *http.Client
}

const (
	API_VERSION = "2022-03-01-preview"
)

func NewBillingPolicyClient(ctx context.Context) BillingPolicyClientInterface {
	return &BillingPolicyClient{
		ppapi: clients.NewPowerPlatformApiClient(ctx, []int{http.StatusOK}),
	}
}

func (client *BillingPolicyClient) GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerplatform.com",
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", id),
	}

	resp, err := client.ppapi.Get(apiUrl.String())
	if err != nil {
		return nil, err
	}
	
	policy := BillingPolicyDto{}
	err = json.NewDecoder(resp.Body).Decode(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (client *BillingPolicyClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	return nil, nil
}

func (client *BillingPolicyClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}
