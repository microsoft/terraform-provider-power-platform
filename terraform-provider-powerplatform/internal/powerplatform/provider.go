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

	host := d.Get("host").(string)
	if host == "" {
		return nil, diag.Errorf(`"host" is not specified`)
	}

	c, err := powerplatform.NewClient(host, username, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags

}
