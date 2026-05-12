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
resourceEnhancedNetwork Setup the Enhanced Network Resource CRUD operations

@return &schema.Resource
*/
func resourceEnhancedNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnhancedNetworkCreate,
		ReadContext:   resourceEnhancedNetworkRead,
		UpdateContext: resourceEnhancedNetworkUpdate,
		DeleteContext: resourceEnhancedNetworkDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp of the last update to this resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the enhanced network.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet CIDR block for the enhanced network. Cannot be changed after creation. Allowed private subnet ranges are 10.0.0.0/12-22, 172.16.0.0/12-22, 192.168.0.0/16-22, or 198.18.0.0/15-22.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of tags to associate with the enhanced network.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of regions to deploy the enhanced network in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"harmony_sase_region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Check Point SASE region ID. Retrieve available IDs from the enhanced_regions data source.",
						},
						"scale_units": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The number of scale units for the region. Higher values provide greater throughput and connection capacity. Defaults to 1.",
						},
						"idle": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether the region gateway is disabled for users. Defaults to true.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the created region.",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEnhancedNetworkImportState,
		},
	}
}

/*
resourceEnhancedNetworkImportState Import an enhanced network by its ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceEnhancedNetworkImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceEnhancedNetworkRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import enhanced network: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceEnhancedNetworkCreate Create an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	name := d.Get("name").(string)
	subnet := d.Get("subnet").(string)
	tags := flattenStringsArrayData(d.Get("tags").([]interface{}))

	regionItems := d.Get("region").([]interface{})
	regionPayloads := make([]perimeter81Sdk.EnhancedRegionCreate, len(regionItems))
	for i, regionItem := range regionItems {
		regionMap := regionItem.(map[string]interface{})
		harmonySaseRegionId := regionMap["harmony_sase_region_id"].(string)
		scaleUnits := int32(regionMap["scale_units"].(int))
		idle := regionMap["idle"].(bool)
		regionPayloads[i] = perimeter81Sdk.EnhancedRegionCreate{
			HarmonySaseRegionId: harmonySaseRegionId,
			ScaleUnits:          &scaleUnits,
			Idle:                &idle,
		}
	}

	networkPayload := perimeter81Sdk.DeployEnhancedNetworkNetwork{
		Name:   name,
		Subnet: &subnet,
		Tags:   tags,
	}
	deployPayload := perimeter81Sdk.DeployEnhancedNetwork{
		Network: networkPayload,
		Regions: regionPayloads,
	}

	status, _, err := client.EnhancedNetworksAPI.CreateEnhancedNetwork(ctx).DeployEnhancedNetwork(deployPayload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Enhanced Network", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var networkId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			networks, _, listErr := client.EnhancedNetworksAPI.GetEnhancedNetworks(ctx).Execute()
			if listErr != nil {
				d.Partial(true)
				return appendErrorDiags(diags, "Unable to create Enhanced Network", listErr)
			}
			for _, networkData := range networks {
				if networkData.Name == name {
					d.SetId(networkData.Id)
					return resourceEnhancedNetworkRead(ctx, d, m)
				}
			}
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			networkId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(networkId)
	return resourceEnhancedNetworkRead(ctx, d, m)
}

/*
resourceEnhancedNetworkRead Read an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Id()
	networkData, _, err := client.EnhancedNetworksAPI.GetEnhancedNetwork(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Enhanced Network", err)
	}

	if err := d.Set("name", networkData.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Network name", err)
	}
	if err := d.Set("subnet", networkData.Subnet); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Network subnet", err)
	}
	if err := d.Set("tags", networkData.Tags); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Network tags", err)
	}

	return diags
}

/*
resourceEnhancedNetworkUpdate Update an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	if d.HasChanges("name", "tags") {
		networkId := d.Id()
		name := d.Get("name").(string)
		tags := flattenStringsArrayData(d.Get("tags").([]interface{}))

		updateNetwork := perimeter81Sdk.EnhancedNetworkUpdateNetwork{
			Name: &name,
			Tags: tags,
		}
		updatePayload := perimeter81Sdk.EnhancedNetworkUpdate{
			Network: &updateNetwork,
		}
		_, _, err := client.EnhancedNetworksAPI.UpdateEnhancedNetwork(ctx, networkId).EnhancedNetworkUpdate(updatePayload).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update Enhanced Network", err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceEnhancedNetworkRead(ctx, d, m)
}

/*
resourceEnhancedNetworkDelete Delete an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Id()
	status, _, err := client.EnhancedNetworksAPI.DeleteEnhancedNetwork(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete Enhanced Network", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId("")
	return diags
}
