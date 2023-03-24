package powerplatform

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePackage() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourcePackageCreate,
		ReadContext:   resourcePackageRead,
		UpdateContext: resourcePackageUpdate,
		DeleteContext: resourcePackageDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: schema.ImportStatePassthroughContext,
		// },

		Description: "Resource to manage packages.",

		Schema: map[string]*schema.Schema{
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"package_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_settings": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) == strings.ToLower(new)
				},
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"import_logs": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePackageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	//d.Set("package_settings", d.Get("package_settings"))
	return diags
}

func resourcePackageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	environmentName := d.Get("environment_name").(string)
	packageName := d.Get("package_name").(string)
	packageFile := d.Get("package_file").(string)
	settings := d.Get("package_settings").(string)

	packageDeploy, err := client.CreatePackage(environmentName, packageName, packageFile, settings)
	if err != nil {
		return diag.FromErr(err)
	}

	currentTime := time.Now().UnixNano()
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", currentTime)))
	timeHash := h.Sum32()

	d.SetId(fmt.Sprint(timeHash))
	d.Set("import_logs", packageDeploy.ImportLogs)
	d.Set("package_name", packageDeploy.PackageName)
	d.Set("package_settings", packageDeploy.PackageSettings)
	d.Set("package_file", packageFile)

	return diags
}

func resourcePackageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePackageCreate(ctx, d, m)
}

func resourcePackageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//we can't delete a package, so we just remove it from the state
	d.SetId("")
	return diags
}
