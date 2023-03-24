package powerplatform

import (
	"context"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource to manage environments.",

		Schema: map[string]*schema.Schema{
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"common_data_service_database_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"language_name": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"currency_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			//TODO: move to some settings > features subnode
			"is_custom_controls_in_canvas_apps_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	envToCreate := powerplatform.EnvironmentCreate{
		DisplayName:                         d.Get("display_name").(string),
		Location:                            d.Get("location").(string),
		EnvironmentType:                     d.Get("environment_type").(string),
		LanguageName:                        d.Get("language_name").(int),
		CurrencyName:                        d.Get("currency_name").(string),
		IsCustomControlsInCanvasAppsEnabled: d.Get("is_custom_controls_in_canvas_apps_enabled").(bool),
	}

	env, err := client.CreateEnvironment(envToCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(env.EnvironmentName)
	d.Set("environment_name", env.EnvironmentName)
	d.Set("display_name", env.DisplayName)
	d.Set("url", env.Url)
	d.Set("domain", env.Domain)
	d.Set("location", env.Location)
	d.Set("environment_type", env.EnvironmentType)
	d.Set("common_data_service_database_type", env.CommonDataServiceDatabaseType)
	d.Set("organization_id", env.OrganizationId)
	d.Set("security_group_id", env.SecurityGroupId)
	d.Set("environment_type", env.EnvironmentType)
	d.Set("language_name", env.LanguageName)
	d.Set("currency_name", env.CurrencyName)
	d.Set("IsCustomControlsInCanvasAppsEnabled", env.IsCustomControlsInCanvasAppsEnabled)

	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	env, err := client.GetEnvironment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("environment_name", env.EnvironmentName)
	d.Set("display_name", env.DisplayName)
	d.Set("location", env.Location)
	d.Set("url", env.Url)
	d.Set("domain", env.Domain)
	d.Set("environment_type", env.EnvironmentType)
	d.Set("common_data_service_database_type", env.CommonDataServiceDatabaseType)
	d.Set("organization_id", env.OrganizationId)
	d.Set("security_group_id", env.SecurityGroupId)
	d.Set("environment_type", env.EnvironmentType)
	d.Set("language_name", env.LanguageName)
	d.Set("currency_name", env.CurrencyName)
	d.Set("IsCustomControlsInCanvasAppsEnabled", env.IsCustomControlsInCanvasAppsEnabled)
	d.SetId(env.EnvironmentName)

	return diags
}

// todo support security_group_id updates
func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//todo validate security_group_id updates
	if d.HasChange("display_name") {
		client := m.(*powerplatform.Client)

		envToUpdate := powerplatform.Environment{
			EnvironmentName:                     d.Id(),
			DisplayName:                         d.Get("display_name").(string),
			OrganizationId:                      d.Get("organization_id").(string),
			Url:                                 d.Get("url").(string),
			Domain:                              d.Get("domain").(string),
			SecurityGroupId:                     d.Get("security_group_id").(string),
			Location:                            d.Get("location").(string),
			LanguageName:                        d.Get("language_name").(int),
			IsCustomControlsInCanvasAppsEnabled: d.Get("is_custom_controls_in_canvas_apps_enabled").(bool),
		}

		err := client.UpdateEnvironment(d.Id(), envToUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	err := client.DeleteEnvironment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
