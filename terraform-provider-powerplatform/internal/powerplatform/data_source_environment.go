package powerplatform

import (
	"context"
	"strconv"
	"time"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			// magodo: No need for this "environment" layer, just keep the schema as same as the managed resource counterpart
			"environments": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"environment_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_data_service_database_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"language_name": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							//ForceNew: true, // magodo: data source attribute don't require ForceNew
						},
						"currency_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							//ForceNew: true, // magodo: data source attribute don't require ForceNew
						},
					},
				},
			},
		},
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	envs, err := client.GetEnvironments()
	if err != nil {
		return diag.FromErr(err)
	}

	all := make([]map[string]interface{}, 0)
	for _, env := range envs {
		envMap := make(map[string]interface{})
		envMap["environment_name"] = env.EnvironmentName
		envMap["display_name"] = env.DisplayName
		envMap["location"] = env.Location
		envMap["environment_type"] = env.EnvironmentType
		envMap["common_data_service_database_type"] = env.CommonDataServiceDatabaseType
		envMap["organization_id"] = env.OrganizationId
		envMap["security_group_id"] = env.SecurityGroupId
		envMap["language_name"] = env.LanguageName
		envMap["currency_name"] = env.CurrencyName
		all = append(all, envMap)
	}

	if err := d.Set("environments", all); err != nil {
		return diag.FromErr(err)
	}

	// magodo: you can use the environment name as its ID, just as how the managed resource counterpart does
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
