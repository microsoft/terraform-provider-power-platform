package powerplatform

import (
	"context"
	"fmt"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataLossPreventionPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDataLossPreventionPolicyRead,
		CreateContext: resourceDataLossPreventionPolicyCreate,
		UpdateContext: resourceDataLossPreventionPolicyUpdate,
		DeleteContext: resourceDataLossPreventionPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource to manage Data Loss Prevention policies.",

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"created_by": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_by": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"e_tag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					string("AllEnvironments"),
					string("ExceptEnvironments"),
					string("OnlyEnvironments"),
				}, false),
			},
			"default_connectors_classification": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					string("General"),
					string("Confidential"),
					string("Blocked"),
				}, false),
			},
			"environment": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"connector_group": {
				Type:      schema.TypeList,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"classification": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ValidateFunc: StringInSlice([]string{
								string("General"),
								string("Confidential"),
								string("Blocked"),
							}, false),
						},
						"connector": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
										Optional: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceDataLossPreventionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	policy, err := client.GetPolicy(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if policy == nil {
		return diag.FromErr(fmt.Errorf("policy '%s' not found", d.Id()))
	}

	d.Set("name", policy.Name)
	d.Set("display_name", policy.DisplayName)
	d.Set("created_by", policy.CreatedBy)
	d.Set("created_time", policy.CreatedTime)
	d.Set("last_modified_by", policy.LastModifiedBy)
	d.Set("last_modified_time", policy.LastModifiedTime)
	d.Set("e_tag", policy.ETag)
	d.Set("environment_type", policy.EnvironmentType)
	d.Set("default_connectors_classification", policy.DefaultConnectorsClassification)

	policyMap := make(map[string]interface{})
	policyMap["environment"] = make([]map[string]interface{}, 0)
	policyMap["connector_group"] = make([]map[string]interface{}, 0)

	for _, environment := range policy.Environments {
		environmentMap := make(map[string]interface{})
		environmentMap["name"] = environment.Name
		environmentMap["id"] = environment.Id
		environmentMap["type"] = environment.Type
		policyMap["environment"] = append(policyMap["environment"].([]map[string]interface{}), environmentMap)
	}
	d.Set("environment", policyMap["environment"])

	for _, connectorGroup := range policy.ConnectorGroups {
		connectorGroupMap := make(map[string]interface{})
		connectorGroupMap["classification"] = connectorGroup.Classification
		connectorGroupMap["connector"] = make([]map[string]interface{}, 0)

		for _, connector := range connectorGroup.Connectors {
			connectorMap := make(map[string]interface{})
			connectorMap["id"] = connector.Id
			connectorMap["name"] = connector.Name
			connectorMap["type"] = connector.Type
			connectorGroupMap["connector"] = append(connectorGroupMap["connector"].([]map[string]interface{}), connectorMap)
		}
		policyMap["connector_group"] = append(policyMap["connector_group"].([]map[string]interface{}), connectorGroupMap)
	}
	d.Set("connector_group", policyMap["connector_group"])

	d.SetId(policy.Name)
	return diags
}

func resourceDataLossPreventionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	policyToCreate := powerplatform.DlpPolicy{
		DisplayName:                     d.Get("display_name").(string),
		EnvironmentType:                 d.Get("environment_type").(string),
		DefaultConnectorsClassification: d.Get("default_connectors_classification").(string),
		Environments:                    []powerplatform.DlpEnvironment{},
		ConnectorGroups:                 []powerplatform.DlpConnectorGroups{},
	}
	policyToCreate.Environments = make([]powerplatform.DlpEnvironment, 0)
	policyToCreate.ConnectorGroups = make([]powerplatform.DlpConnectorGroups, 0)

	for _, environment := range d.Get("environment").([]interface{}) {
		policyToCreate.Environments = append(policyToCreate.Environments, powerplatform.DlpEnvironment{
			Name: environment.(map[string]interface{})["name"].(string),
		})
	}

	for _, connectorGroup := range d.Get("connector_group").([]interface{}) {

		connectorGroupMap := powerplatform.DlpConnectorGroups{
			Classification: connectorGroup.(map[string]interface{})["classification"].(string),
			Connectors:     []powerplatform.DlpConnector{},
		}

		for _, connector := range connectorGroup.(map[string]interface{})["connector"].([]interface{}) {
			connectorGroupMap.Connectors = append(connectorGroupMap.Connectors, powerplatform.DlpConnector{
				Id:   connector.(map[string]interface{})["id"].(string),
				Name: connector.(map[string]interface{})["name"].(string),
				Type: connector.(map[string]interface{})["type"].(string),
			})
		}
		policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, connectorGroupMap)
	}

	policy, error := client.CreatePolicy(policyToCreate)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId(policy.Name)

	d.Set("name", policy.Name)
	d.Set("display_name", policy.DisplayName)
	d.Set("created_by", policy.CreatedBy)
	d.Set("created_time", policy.CreatedTime)
	d.Set("last_modified_by", policy.LastModifiedBy)
	d.Set("last_modified_time", policy.LastModifiedTime)
	d.Set("e_tag", policy.ETag)
	d.Set("environment_type", policy.EnvironmentType)
	d.Set("default_connectors_classification", policy.DefaultConnectorsClassification)

	policyMap := make(map[string]interface{})
	policyMap["environment"] = make([]map[string]interface{}, 0)
	policyMap["connector_group"] = make([]map[string]interface{}, 0)

	for _, environment := range policy.Environments {
		environmentMap := make(map[string]interface{})
		environmentMap["name"] = environment.Name
		environmentMap["id"] = environment.Id
		environmentMap["type"] = environment.Type
		policyMap["environment"] = append(policyMap["environment"].([]map[string]interface{}), environmentMap)
	}
	d.Set("environment", policyMap["environment"])

	for _, connectorGroup := range policy.ConnectorGroups {
		connectorGroupMap := make(map[string]interface{})
		connectorGroupMap["classification"] = connectorGroup.Classification
		connectorGroupMap["connector"] = make([]map[string]interface{}, 0)

		for _, connector := range connectorGroup.Connectors {
			connectorMap := make(map[string]interface{})
			connectorMap["id"] = connector.Id
			connectorMap["name"] = connector.Name
			connectorMap["type"] = connector.Type
			connectorGroupMap["connector"] = append(connectorGroupMap["connector"].([]map[string]interface{}), connectorMap)
		}
		policyMap["connector_group"] = append(policyMap["connector_group"].([]map[string]interface{}), connectorGroupMap)
	}
	d.Set("connector_group", policyMap["connector_group"])

	return diags
}

func resourceDataLossPreventionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)
	var diags diag.Diagnostics

	//there is no point of checking if a given attribute has changes
	//as API requires all attributes to be sent in the request

	policyToUpdate := powerplatform.DlpPolicy{
		Name:                            d.Id(),
		DisplayName:                     d.Get("display_name").(string),
		EnvironmentType:                 d.Get("environment_type").(string),
		DefaultConnectorsClassification: d.Get("default_connectors_classification").(string),
		Environments:                    []powerplatform.DlpEnvironment{},
		ConnectorGroups:                 []powerplatform.DlpConnectorGroups{},
	}
	policyToUpdate.Environments = make([]powerplatform.DlpEnvironment, 0)
	policyToUpdate.ConnectorGroups = make([]powerplatform.DlpConnectorGroups, 0)

	for _, environment := range d.Get("environment").([]interface{}) {
		policyToUpdate.Environments = append(policyToUpdate.Environments, powerplatform.DlpEnvironment{
			Name: environment.(map[string]interface{})["name"].(string),
			Id:   environment.(map[string]interface{})["id"].(string),
			Type: environment.(map[string]interface{})["type"].(string),
		})
	}

	for _, connectorGroup := range d.Get("connector_group").([]interface{}) {

		connectorGroupMap := powerplatform.DlpConnectorGroups{
			Classification: connectorGroup.(map[string]interface{})["classification"].(string),
			Connectors:     []powerplatform.DlpConnector{},
		}
		for _, connector := range connectorGroup.(map[string]interface{})["connector"].([]interface{}) {
			connectorGroupMap.Connectors = append(connectorGroupMap.Connectors, powerplatform.DlpConnector{
				Id:   connector.(map[string]interface{})["id"].(string),
				Name: connector.(map[string]interface{})["name"].(string),
				Type: connector.(map[string]interface{})["type"].(string),
			})
		}
		policyToUpdate.ConnectorGroups = append(policyToUpdate.ConnectorGroups, connectorGroupMap)
	}

	policy, error := client.UpdatePolicy(d.Id(), policyToUpdate)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId(policy.Name)
	d.Set("name", policy.Name)
	d.Set("display_name", policy.DisplayName)
	d.Set("created_by", policy.CreatedBy)
	d.Set("created_time", policy.CreatedTime)
	d.Set("last_modified_by", policy.LastModifiedBy)
	d.Set("last_modified_time", policy.LastModifiedTime)
	d.Set("e_tag", policy.ETag)
	d.Set("environment_type", policy.EnvironmentType)
	d.Set("default_connectors_classification", policy.DefaultConnectorsClassification)

	policyMap := make(map[string]interface{})
	policyMap["environment"] = make([]map[string]interface{}, 0)
	policyMap["connector_group"] = make([]map[string]interface{}, 0)

	for _, environment := range policy.Environments {
		environmentMap := make(map[string]interface{})
		environmentMap["name"] = environment.Name
		environmentMap["id"] = environment.Id
		environmentMap["type"] = environment.Type
		policyMap["environment"] = append(policyMap["environment"].([]map[string]interface{}), environmentMap)
	}
	d.Set("environment", policyMap["environment"])

	for _, connectorGroup := range policy.ConnectorGroups {
		connectorGroupMap := make(map[string]interface{})
		connectorGroupMap["classification"] = connectorGroup.Classification
		connectorGroupMap["connector"] = make([]map[string]interface{}, 0)

		for _, connector := range connectorGroup.Connectors {
			connectorMap := make(map[string]interface{})
			connectorMap["id"] = connector.Id
			connectorMap["name"] = connector.Name
			connectorMap["type"] = connector.Type
			connectorGroupMap["connector"] = append(connectorGroupMap["connector"].([]map[string]interface{}), connectorMap)
		}
		policyMap["connector_group"] = append(policyMap["connector_group"].([]map[string]interface{}), connectorGroupMap)
	}
	d.Set("connector_group", policyMap["connector_group"])

	return diags
}

func resourceDataLossPreventionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	error := client.DeletePolicy(d.Id())
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId("")
	return diags
}
