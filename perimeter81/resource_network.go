package perimeter81

import (
	"context"
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
	tags := flattenStringsArrayData(network["tags"].([]interface{}))
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
		return appendErrorDiags(diags, "Unable to create Network", err)
	}

	statusId := getIdFromUrl(status.StatusUrl)

	var networkId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, _ = checkNetworkStatus(ctx, statusId, *client, diags)
		if networkStatus.Completed {
			networkId = getIdFromUrl(networkStatus.Result.Resource)
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
	networkData, _, err := client.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkId)

	if err != nil {
		return appendErrorDiags(diags, "Unable to find Network", err)
	}

	CreateNetworkPayload := perimeter81Sdk.CreateNetworkPayload{
		Name:   networkData.Name,
		Tags:   networkData.Tags,
		Subnet: networkData.Subnet,
	}

	network := flattenNetworkData([]perimeter81Sdk.CreateNetworkPayload{CreateNetworkPayload})
	if err := d.Set("network", network); err != nil {
		return appendErrorDiags(diags, "Unable to set Network data", err)
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
		tags := flattenStringsArrayData(network["tags"].([]interface{}))

		networkDto := perimeter81Sdk.BaseNetworkDto{Name: name, Tags: tags}
		updateNetworkDto := perimeter81Sdk.UpdateNetworkDto{Network: &networkDto}
		_, _, err := client.NetworksApi.NetworksControllerV2NetworkUpdate(ctx, updateNetworkDto, networkId)
		if err != nil {
			return appendErrorDiags(diags, "Unable to update Network", err)
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
		return appendErrorDiags(diags, "Unable to delete network", err)
	}
	d.SetId("")
	return diags
}
