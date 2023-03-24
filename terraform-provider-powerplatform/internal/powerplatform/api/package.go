package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *Client) CreatePackage(environmentName string, packageName string, file string, settings string) (*PackageDeploy, error) {

	packageContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content, err := json.Marshal(PackageDeploy{
		PackageName:     packageName,
		PackageSettings: settings,
		PackageContent:  packageContent,
	})

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/environments/%s/packages", client.HostURL, environmentName),
		bytes.NewReader(content))

	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("no body returned")
	}

	packageDeploy := PackageDeploy{}
	err = json.Unmarshal(body, &packageDeploy)
	if err != nil {
		return nil, err
	}

	return &packageDeploy, nil
}
