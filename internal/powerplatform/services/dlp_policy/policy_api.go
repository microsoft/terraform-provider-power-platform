package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewDlpPolicyClient(bapi api.BapiClientInterface) DlpPolicyClient {
	return DlpPolicyClient{
		BapiApiClient: bapi,
	}
}

type DlpPolicyClient struct {
	BapiApiClient api.BapiClientInterface
}

func (client *DlpPolicyClient) GetPolicies(ctx context.Context) ([]DlpPolicyModelDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiApiClient.GetBase().GetConfig().Urls.BapiUrl,
		Path:   "providers/PowerPlatform.Governance/v2/policies",
	}
	policiesArray := DlpPolicyDefinitionDtoArray{}
	_, err := client.BapiApiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policiesArray)
	if err != nil {
		return nil, err
	}

	policies := make([]DlpPolicyModelDto, 0)
	for _, policy := range policiesArray.Value {

		apiUrl := &url.URL{
			Scheme: "https",
			Host:   client.BapiApiClient.GetBase().GetConfig().Urls.BapiUrl,
			Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.PolicyDefinition.Name),
		}
		policy := DlpPolicyDto{}
		_, err := client.BapiApiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)
		if err != nil {
			return nil, err
		}
		v, err := covertDlpPolicyToPolicyModelDto(policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, *v)
	}
	return policies, nil
}
