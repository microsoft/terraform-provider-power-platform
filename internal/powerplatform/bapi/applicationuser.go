package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"

	//"fmt"
	"net/http"
	"net/url"
	"strings"

	//"time"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetApplicationUser(ctx context.Context, environmentName string) ([]models.ApplicationUserDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerapps.com",
		Path:   "/providers/Microsoft.PowerApps/apis",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	values.Add("showApisWithToS", "true")
	values.Add("hideDlpExemptApis", "true")
	values.Add("showAllDlpEnforceableApis", "true")
	values.Add("$filter", "environment eq '~Default'")
	apiUrl.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	ApplicationUserArray := models.ApplicationUserDtoArray{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&ApplicationUserArray)
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/PowerPlatform.Governance/v1/applicationuser/metadata/unblockable", //verify this path
	}
	request, err = http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	unblockableApplicationUserArray := []models.UnblockableApplicationUserDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&unblockableApplicationUserArray)
	if err != nil {
		return nil, err
	}

	for inx, applicationUser := range ApplicationUserArray.Value {
		for _, unblockableApplicationUser := range unblockableApplicationUserArray {
			if applicationUser.Id == unblockableApplicationUser.Id {
				ApplicationUserArray.Value[inx].Properties.Unblockable = unblockableApplicationUser.Metadata.Unblockable
			}
		}
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/PowerPlatform.Governance/v1/applicationUser/metadata/virtual", //verify this path
	}
	request, err = http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	virtualApplicationUserArray := []models.VirtualApplicationUserDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&virtualApplicationUserArray)
	if err != nil {
		return nil, err
	}
	for _, virutualApplicationUser := range virtualApplicationUserArray {
		ApplicationUserArray.Value = append(ApplicationUserArray.Value, models.ApplicationUserDto{
			Id:   virutualApplicationUser.Id,
			Name: virutualApplicationUser.Metadata.Name,
			Type: virutualApplicationUser.Metadata.Type,
			Properties: models.ApplicationUserPropertiesDto{
				DisplayName: virutualApplicationUser.Metadata.DisplayName,
				Unblockable: false,
				Tier:        "Built-in",
				Publisher:   "Microsoft",
				Description: "",
			},
		})
	}

	for inx, applicationUser := range ApplicationUserArray.Value {
		nameSplit := strings.Split(applicationUser.Id, "/")
		ApplicationUserArray.Value[inx].Name = nameSplit[len(nameSplit)-1]
	}

	return ApplicationUserArray.Value, nil
}
