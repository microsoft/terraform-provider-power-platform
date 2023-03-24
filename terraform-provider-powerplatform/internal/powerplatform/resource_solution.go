package powerplatform

import (
	"context"
	"io/ioutil"
	"strings"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSolution() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceSolutionCreate,
		ReadContext:   resourceSolutionRead,
		UpdateContext: resourceSolutionUpdate,
		DeleteContext: resourceSolutionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource to manage solutions.",

		Schema: map[string]*schema.Schema{
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"solution_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"solution_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"solution_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"settings_file": &schema.Schema{
				Type: schema.TypeString,
				//Optional: true,
				//ForceNew: true,
				Required: true,
				//TODO when content of the file changes we should update solution
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) == strings.ToLower(new)
				},
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"is_managed": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

// TODO we should validate zip checksum. Name may be the same but content may be different
// TODO we should use checksum of the file as Id together with solution name
// https://github.com/hashicorp/terraform-provider-local/blob/c5ac21491b93e549bb698b5fa881759372f31b8a/internal/provider/resource_local_file.go#L100
func resourceSolutionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	var solutionName = d.Get("solution_name").(string)
	var environmentName = d.Get("environment_name").(string)

	solutions, err := client.ReadSolutions(environmentName)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool = false
	for _, solution := range solutions {
		if solution.SolutionName == solutionName {
			d.Set("solution_name", solution.SolutionName)

			//TODO test a case when solution version changes
			d.Set("solution_version", solution.SolutionVersion)
			d.Set("is_managed", solution.IsManaged)
			d.Set("display_name", solution.DisplayName)
			d.Set("environment_name", environmentName)
			d.Set("solution_file", d.Get("solution_file"))
			d.Set("settings_file", d.Get("settings_file"))
			d.SetId("_" + environmentName + "_" + solution.SolutionName)
			found = true
		}
	}
	//TODO test a case when solution is not found
	if !found {
		d.SetId("")
		//return diag.Errorf("solution %s not found", solutionName)
	}
	return diags
}

func resourceSolutionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	solutionToCreate := powerplatform.Solution{
		EnvironmentName: d.Get("environment_name").(string),
		File:            d.Get("solution_file").(string),
		SettingsFile:    d.Get("settings_file").(string),
		SolutionName:    d.Get("solution_name").(string),
	}

	solutionContent, err := ioutil.ReadFile(solutionToCreate.File)
	if err != nil {
		return diag.FromErr(err)
	}

	settingsContent := make([]byte, 0)
	if solutionToCreate.SettingsFile != "" {
		settingsContent, err = ioutil.ReadFile(solutionToCreate.SettingsFile)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	solution, err := client.CreateSolutions(solutionToCreate, solutionContent, settingsContent)
	if err != nil {
		return diag.FromErr(err)
	}
	if solution.SolutionName == "" {
		return diag.Errorf("solution name is empty")
	}

	d.SetId("_" + solution.EnvironmentName + "_" + solution.SolutionName)
	d.Set("environment_name", solution.EnvironmentName)
	d.Set("solution_name", solution.SolutionName)
	d.Set("solution_version", solution.SolutionVersion)
	d.Set("is_managed", solution.IsManaged)
	d.Set("display_name", solution.DisplayName)
	d.Set("solution_file", d.Get("solution_file"))
	d.Set("settings_file", d.Get("settings_file"))

	return diags
}

func resourceSolutionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSolutionCreate(ctx, d, m)
}

// We can't really delete a solution, because that will remove all the data,
// so we just remove the resource from the state
// TODO solution should be under environment, so we can delete it whe environment changes
func resourceSolutionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
