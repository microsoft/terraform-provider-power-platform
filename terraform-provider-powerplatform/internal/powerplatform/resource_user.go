package powerplatform

import (
	"context"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource to manage users in environments",

		Schema: map[string]*schema.Schema{
			"application_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_roles": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
			"is_app_user": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"aad_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"user_principal_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//ForceNew: true,
				//todo add validation that is the user is not an app user
				//for that to recrate because DV sets this values itself from null
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func flattenTypeList(typeSet interface{}) []string {
	typeList := typeSet.(*schema.Set).List()
	flattenList := []string{}
	for _, v := range typeList {
		if v != "" {
			flattenList = append(flattenList, v.(string))
		}
	}
	return flattenList
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	//TODO in case of app user the aadID is the enterprise object id not application object id
	//we're not getting the correct aadID from terraform Azure AD provider
	var environment_name = d.Get("environment_name").(string)
	var aadId = d.Id()

	user, err := client.ReadUser(environment_name, aadId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.AadObjectId)
	d.Set("environment", environment_name)
	d.Set("security_roles", user.SecurityRoles)
	d.Set("is_app_user", user.IsApplicationUser)
	d.Set("aad_id", user.AadObjectId)
	d.Set("user_principal_name", user.DomainName)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("is_disabled", user.IsDisabled)
	d.Set("application_id", user.ApplicationId)

	return diags
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	var environment_name = d.Get("environment_name").(string)
	var aadId = d.Get("aad_id").(string)

	userToCreate := powerplatform.User{
		IsApplicationUser: d.Get("is_app_user").(bool),
		AadObjectId:       aadId,
		DomainName:        d.Get("user_principal_name").(string),
		FirstName:         d.Get("first_name").(string),
		LastName:          d.Get("last_name").(string),
		SecurityRoles:     flattenTypeList(d.Get("security_roles")),
		//IsDisabled:    d.Get("is_disabled").(bool),
		ApplicationId: d.Get("application_id").(string),
	}

	user, err := client.CreateUser(environment_name, userToCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.AadObjectId)
	d.Set("environment", environment_name)
	d.Set("security_roles", user.SecurityRoles)
	d.Set("is_app_user", user.IsApplicationUser)
	d.Set("aad_id", user.AadObjectId)
	d.Set("user_principal_name", user.DomainName)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("is_disabled", user.IsDisabled)
	d.Set("application_id", user.ApplicationId)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	var environment_name = d.Get("environment_name").(string)
	var aadId = d.Id()

	userToUpdate := powerplatform.User{
		IsApplicationUser: d.Get("is_app_user").(bool),
		AadObjectId:       aadId,
		DomainName:        d.Get("user_principal_name").(string),
		FirstName:         d.Get("first_name").(string),
		LastName:          d.Get("last_name").(string),
		SecurityRoles:     flattenTypeList(d.Get("security_roles")),
		//IsDisabled:        d.Get("is_disabled").(bool),
		ApplicationId: d.Get("application_id").(string),
	}

	user, err := client.UpdateUser(environment_name, userToUpdate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.AadObjectId)
	d.Set("environment", environment_name)
	d.Set("security_roles", user.SecurityRoles)
	d.Set("is_app_user", user.IsApplicationUser)
	d.Set("aad_id", user.AadObjectId)
	d.Set("user_principal_name", user.DomainName)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("is_disabled", user.IsDisabled)
	d.Set("application_id", user.ApplicationId)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	environment := d.Get("environment_name").(string)
	aadId := d.Id()

	err := client.DeleteUser(environment, aadId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
