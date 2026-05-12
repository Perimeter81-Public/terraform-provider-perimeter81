package checkpointsase

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceNetwork Setup the IpSec-Signle Resource CRUD operations

@return &schema.Resource
*/
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
						"dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpregion_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idle": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: ResourceNetworkImportState,
		},
	}
}

/*
ResourceNetworkImportState Import gateways
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func ResourceNetworkImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceNetworkRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import network: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceNetworkCreate Create a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the network data from the resource data and flatten what need to be flattened
	network := d.Get("network").([]interface{})[0].(map[string]interface{})
	name := network["name"].(string)
	tags := flattenStringsArrayData(network["tags"].([]interface{}))
	subnet := network["subnet"].(string)
	regions := flattenRegionsData(d.Get("region").([]interface{}))

	// convert region configs to CreateRegionInNetworkPayload
	regionPayloads := make([]perimeter81Sdk.CreateRegionInNetworkPayload, len(regions))
	for i, r := range regions {
		regionPayloads[i] = perimeter81Sdk.CreateRegionInNetworkPayload{
			HarmonySaseRegionId: r.CpRegionId,
			Idle:                r.Idle,
		}
	}

	// create the network payload
	CreateNetworkPayload := perimeter81Sdk.CreateNetworkPayload{
		Name:   name,
		Tags:   tags,
		Subnet: &subnet,
	}
	DeployNetworkPayload := perimeter81Sdk.DeployNetworkPayload{
		Network: CreateNetworkPayload,
		Regions: regionPayloads,
	}
	// create the network and check for errors
	status, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkCreate(ctx).DeployNetworkPayload(DeployNetworkPayload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Network", err)
	}

	// get the status id from the status url
	statusId := getIdFromUrl(status.GetStatusUrl())
	var networkId string
	// check the status of the network creation
	for {
		// check the status of the network creation and check for errors
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			networks, _, err := client.StandardNetworksAPI.StandardGetNetworks(ctx).Execute()
			if err != nil {
				d.Partial(true)
				return appendErrorDiags(diags, "Unable to Create Network", err)
			}
			for _, networkData := range networks {
				if networkData.Name == name {
					d.SetId(networkData.Id)
					return resourceNetworkRead(ctx, d, m)
				}
			}
			d.Partial(true)
			return diags
		}
		// if the network creation is completed, get the network id and break the loop
		if networkStatus.GetCompleted() {
			networkId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		// sleep for 60 seconds and check the status again
		time.Sleep(60 * time.Second)
	}

	d.SetId(networkId)

	return resourceNetworkRead(ctx, d, m)
}

/*
resourceNetworkRead Read a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the network id from the resource data
	networkId := d.Id()
	// get the network data and check for errors
	networkData, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Network", err)
	}

	// get the regions data and check for errors
	regionsData, _, err := client.RegionsAPI.StandardNetworksControllerV2GetRegions(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get CpRegions", err)
	}

	// handle regions terraform import
	regions := flattenRegionsData(d.Get("region").([]interface{}))
	regions = importRegions(networkData, regionsData, regions)

	// flatten the regions data and set the network region infos
	setNetworkRegionInfos(regionsData, networkData, regions)
	CreateNetworkPayload := perimeter81Sdk.CreateNetworkPayload{
		Name:   networkData.Name,
		Tags:   networkData.Tags,
		Subnet: &networkData.Subnet,
	}
	// set the network data and the regions data
	network := flattenNetworkData([]perimeter81Sdk.CreateNetworkPayload{CreateNetworkPayload})
	regions = setDefaultGatewayIpForRegions(regions, networkData)
	if err := d.Set("network", network); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Network data", err)
	}
	if err := d.Set("region", flattenNetworkRegions(regions)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set regions data", err)
	}

	return diags
}

/*
resourceNetworkUpdate Upadte a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// check if the network has changed
	if d.HasChange("network") {

		// get the network data from the resource data and flatten what need to be flattened
		networkId := d.Id()
		network := d.Get("network").([]interface{})[0].(map[string]interface{})
		name := network["name"].(string)
		tags := flattenStringsArrayData(network["tags"].([]interface{}))

		// creata the update payload
		networkDto := perimeter81Sdk.BaseNetworkDto{Name: &name, Tags: tags}
		updateNetworkDto := perimeter81Sdk.UpdateNetworkDto{Network: networkDto}
		// update the network and check for errors
		_, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkUpdate(ctx, networkId).UpdateNetworkDto(updateNetworkDto).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update Network", err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	// check if the regions has changed
	if d.HasChange("region") {
		// get the network id from the resource data
		networkId := d.Id()
		// get the old and new regions data and flatten what need to be flattened
		oldRegoins, newRegions := d.GetChange("region")
		oldRegionsFlattened := flattenRegionsData(oldRegoins.([]interface{}))
		newRegionsFlattened := flattenRegionsData(newRegions.([]interface{}))
		// Add new region to the network if any and check for errors
		statusId, cpRegionId, err := resourceRegionCreate(ctx, networkId, oldRegionsFlattened, newRegionsFlattened, d, client)
		if err != nil {
			return appendErrorDiags(diags, "Unable to update network region "+cpRegionId, err)
		}
		// Delete old region from the network if any and check for errors
		var RegionId string
		statusId, RegionId, err = resourceRegionDelete(ctx, networkId, oldRegionsFlattened, newRegionsFlattened, d, client, statusId)
		if err != nil {
			return appendErrorDiags(diags, "Unable to delete network region "+RegionId, err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
		// wait for the network to be updated with the regions (if a status id is available)
		if statusId != "" {
			for {
				// check the network status and check for errors
				var networkStatus perimeter81Sdk.AsyncOperationStatus
				networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
				if err != nil {
					d.Partial(true)
					return diags
				}
				// if the network status is completed break the loop
				if networkStatus.GetCompleted() {
					break
				}
				// wait for 60 seconds and check again
				time.Sleep(60 * time.Second)
			}
		}
	}

	return resourceNetworkRead(ctx, d, m)
}

/*
resourceNetworkDelete Delete a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the network id from the resource data
	networkId := d.Id()
	// delete the network and check for errors (synchronous operation — returns AsyncOperationResult, no status URL)
	_, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkDelete(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete network", err)
	}
	d.SetId("")
	return diags
}
