package powerplatform

import (
	"context"
	"strconv"
	"time"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePowerApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePowerAppsRead,
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"apps": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePowerAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	envName := d.Get("environment_name").(string)

	apps, err := client.GetPowerApps(envName)
	if err != nil {
		return diag.FromErr(err)
	}

	all := make([]map[string]interface{}, 0)
	for _, app := range apps {
		appMap := make(map[string]interface{})
		appMap["environment_name"] = app.EnvironmentName
		appMap["display_name"] = app.DisplayName
		appMap["name"] = app.Name
		appMap["created_time"] = app.CreatedTime
		all = append(all, appMap)
	}

	if err := d.Set("apps", all); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
