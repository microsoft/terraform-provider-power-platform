package powerplatform

import (
	"context"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWER_PLATFORM_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWER_PLATFORM_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWER_PLATFORM_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"powerplatform_environment":                 resourceEnvironment(),
			"powerplatform_solution":                    resourceSolution(),
			"powerplatform_data_loss_prevention_policy": resourceDataLossPreventionPolicy(),
			"powerplatform_user":                        resourceUser(),
			"powerplatform_package":                     resourcePackage(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"powerplatform_powerapps":                     dataSourcePowerApps(),
			"powerplatform_environments":                  dataSourceEnvironments(),
			"powerplatform_data_loss_prevention_policies": dataSourceDataLossPreventionPolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// TODO implement OAuth
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var host *string
	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	if (username != "") && (password != "") {
		c, error := powerplatform.NewClient(host, &username, &password)
		if error != nil {
			return nil, diag.FromErr(error)
		}
		return c, diags
	}

	c, error := powerplatform.NewClient(host, nil, nil)
	if error != nil {
		return nil, diag.FromErr(error)
	}
	return c, diags

}
