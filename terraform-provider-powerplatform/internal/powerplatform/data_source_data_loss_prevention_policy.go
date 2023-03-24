package powerplatform

import (
	"context"
	"strconv"
	"time"

	powerplatform "github.com/microsoft/terraform-provider-powerplatform/internal/powerplatform/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDataLossPreventionPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataLossPreventionPolicyRead,
		Schema: map[string]*schema.Schema{
			"policies": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
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
							Computed: true,
						},
						"default_connectors_classification": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"environments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connector_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"classification": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"connectors": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": &schema.Schema{
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDataLossPreventionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*powerplatform.Client)

	var diags diag.Diagnostics

	policies, err := client.GetPolicies()
	if err != nil {
		return diag.FromErr(err)
	}

	all := make([]map[string]interface{}, 0)
	for _, policy := range policies {
		policyMap := make(map[string]interface{})
		policyMap["name"] = policy.Name
		policyMap["display_name"] = policy.DisplayName
		policyMap["created_by"] = policy.CreatedBy
		policyMap["created_time"] = policy.CreatedTime
		policyMap["last_modified_by"] = policy.LastModifiedBy
		policyMap["last_modified_time"] = policy.LastModifiedTime
		policyMap["e_tag"] = policy.ETag
		policyMap["environment_type"] = policy.EnvironmentType
		policyMap["default_connectors_classification"] = policy.DefaultConnectorsClassification
		policyMap["environments"] = make([]map[string]interface{}, 0)
		policyMap["connector_groups"] = make([]map[string]interface{}, 0)

		for _, environment := range policy.Environments {
			environmentMap := make(map[string]interface{})
			environmentMap["name"] = environment.Name
			environmentMap["id"] = environment.Id
			environmentMap["type"] = environment.Type
			policyMap["environments"] = append(policyMap["environments"].([]map[string]interface{}), environmentMap)
		}

		for _, connectorGroup := range policy.ConnectorGroups {
			connectorGroupMap := make(map[string]interface{})
			connectorGroupMap["classification"] = connectorGroup.Classification
			connectorGroupMap["connectors"] = make([]map[string]interface{}, 0)

			for _, connector := range connectorGroup.Connectors {
				connectorMap := make(map[string]interface{})
				connectorMap["id"] = connector.Id
				connectorMap["name"] = connector.Name
				connectorMap["type"] = connector.Type
				connectorGroupMap["connectors"] = append(connectorGroupMap["connectors"].([]map[string]interface{}), connectorMap)
			}
			policyMap["connector_groups"] = append(policyMap["connector_groups"].([]map[string]interface{}), connectorGroupMap)
		}
		all = append(all, policyMap)
	}

	if err := d.Set("policies", all); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
