package perimeter81

import (
	"context"
	"fmt"
	"strings"
	perimeter81Sdk "terraform-provider-perimeter81/perimeter81sdk"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpregionid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instancecount": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"idle": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}
	network := d.Get("network").([]interface{})[0].(map[string]interface{})
	name := network["name"].(string)
	tags := flattenTagsData(network["tags"].([]interface{}))
	regions := flattenRegionsData(d.Get("region").([]interface{}))

	CreateNetworkPayload := perimeter81Sdk.CreateNetworkPayload{
		Name: name,
		Tags: tags,
	}
	DeployNetworkPayload := perimeter81Sdk.DeployNetworkPayload{
		Network: &CreateNetworkPayload,
		Regions: regions,
	}

	status, _, err := client.NetworksApi.NetworksControllerV2NetworkCreate(ctx, DeployNetworkPayload)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Network",
			Detail:   err.Error(),
		})
	}

	statusUrl := strings.Split(status.StatusUrl, "/")
	statusId := statusUrl[len(statusUrl)-1]

	var networkId string
	for {
		networkStatus, _, err := client.NetworksApi.NetworksControllerV2Status(ctx, statusId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Network Status",
				Detail:   err.Error(),
			})
		}
		if networkStatus.Completed {
			networkUrl := strings.Split(networkStatus.Result.Resource, "/")
			networkId = networkUrl[len(networkUrl)-1]
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(networkId)

	return diags
}
func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}
	networkId := d.Id()
	network, _, err := client.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkId)

	if err != nil {
		return diag.FromErr(err)
	}

	CreateNetworkPayload := perimeter81Sdk.CreateNetworkPayload{
		Name:   network.Name,
		Tags:   network.Tags,
		Subnet: network.Subnet,
	}

	n := flattenNetworkData([]perimeter81Sdk.CreateNetworkPayload{CreateNetworkPayload})
	if err := d.Set("network", n); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}

	if d.HasChange("network") {
		networkId := d.Id()
		network := d.Get("network").([]interface{})[0].(map[string]interface{})
		name := network["name"].(string)
		tags := flattenTagsData(network["tags"].([]interface{}))
		networkDto := perimeter81Sdk.BaseNetworkDto{Name: name, Tags: tags}
		updateNetworkDto := perimeter81Sdk.UpdateNetworkDto{Network: &networkDto}
		_, _, err := client.NetworksApi.NetworksControllerV2NetworkUpdate(ctx, updateNetworkDto, networkId)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update Network",
				Detail:   err.Error(),
			})
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceNetworkRead(ctx, d, m)
}
func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}
	networkId := d.Id()
	_, _, err := client.NetworksApi.NetworksControllerV2NetworkDelete(ctx, networkId)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func flattenTagsData(tags []interface{}) []string {
	tagsData := make([]string, len(tags))
	for i, v := range tags {
		tagsData[i] = fmt.Sprint(v)
	}
	return tagsData
}

func flattenRegionsData(regionItems []interface{}) []perimeter81Sdk.CreateRegionInNetworkload {
	if regionItems != nil {
		regions := make([]perimeter81Sdk.CreateRegionInNetworkload, len(regionItems))

		for i, regionItem := range regionItems {
			region := perimeter81Sdk.CreateRegionInNetworkload{}

			region.CpRegionId = regionItem.(map[string]interface{})["cpregionid"].(string)
			region.InstanceCount = int32(regionItem.(map[string]interface{})["instancecount"].(int))
			region.Idle = regionItem.(map[string]interface{})["idle"].(bool)
			regions[i] = region
		}

		return regions
	}

	return make([]perimeter81Sdk.CreateRegionInNetworkload, 0)
}

func flattenNetworkData(networkItems []perimeter81Sdk.CreateNetworkPayload) []interface{} {
	if networkItems != nil {
		networks := make([]interface{}, len(networkItems))

		for i, networkItem := range networkItems {
			network := make(map[string]interface{})

			network["name"] = networkItem.Name
			network["tags"] = networkItem.Tags
			network["subnet"] = networkItem.Subnet
			networks[i] = network
		}

		return networks
	}

	return make([]interface{}, 0)
}
