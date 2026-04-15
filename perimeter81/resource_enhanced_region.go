package perimeter81

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceEnhancedRegion Setup the Enhanced Region Resource CRUD operations

@return &schema.Resource
*/
func resourceEnhancedRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnhancedRegionCreate,
		ReadContext:   resourceEnhancedRegionRead,
		UpdateContext: resourceEnhancedRegionUpdate,
		DeleteContext: resourceEnhancedRegionDelete,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the enhanced network to add this region to.",
			},
			"harmony_sase_region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Description: "Whether the region gateway is disabled for users. Defaults to true.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEnhancedRegionImportState,
		},
	}
}

/*
resourceEnhancedRegionImportState Import an enhanced region by its composite ID (networkId/regionId).
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceEnhancedRegionImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceEnhancedRegionRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import enhanced region: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceEnhancedRegionCreate Create an Enhanced Region in an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRegionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	harmonySaseRegionId := d.Get("harmony_sase_region_id").(string)
	scaleUnits := int32(d.Get("scale_units").(int))
	idle := d.Get("idle").(bool)

	payload := perimeter81Sdk.EnhancedRegionCreate{
		HarmonySaseRegionId: harmonySaseRegionId,
		ScaleUnits:          &scaleUnits,
		Idle:                &idle,
	}

	status, _, err := client.EnhancedRegionsAPI.CreateEnhancedRegion(ctx, networkId).EnhancedRegionCreate(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Enhanced Region", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var regionId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			regionId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(regionId)
	return resourceEnhancedRegionRead(ctx, d, m)
}

/*
resourceEnhancedRegionRead Read an Enhanced Region.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	regionId := d.Id()

	regionData, _, err := client.EnhancedRegionsAPI.GetEnhancedRegion(ctx, networkId, regionId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Enhanced Region", err)
	}

	if err := d.Set("scale_units", regionData.ScaleUnits); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Region scale_units", err)
	}

	return diags
}

/*
resourceEnhancedRegionUpdate Update an Enhanced Region's scale units.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRegionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	if !d.HasChange("scale_units") {
		return resourceEnhancedRegionRead(ctx, d, m)
	}

	networkId := d.Get("network_id").(string)
	regionId := d.Id()

	oldVal, newVal := d.GetChange("scale_units")
	oldScaleUnits := int32(oldVal.(int))
	newScaleUnits := int32(newVal.(int))

	if newScaleUnits > oldScaleUnits {
		// Increase scale units one unit at a time
		unitsToAdd := newScaleUnits - oldScaleUnits
		payload := perimeter81Sdk.ScaleUnitsOperation{
			UnitType:        "standard",
			ScaleUnitsCount: unitsToAdd,
		}

		status, _, err := client.EnhancedRegionsAPI.IncreaseScaleUnit(ctx, networkId, regionId).ScaleUnitsOperation(payload).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to increase Enhanced Region scale units", err)
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
	} else if newScaleUnits < oldScaleUnits {
		// Reduce scale units
		unitsToRemove := oldScaleUnits - newScaleUnits
		payload := perimeter81Sdk.ScaleUnitsOperation{
			UnitType:        "standard",
			ScaleUnitsCount: unitsToRemove,
		}

		status, _, err := client.EnhancedRegionsAPI.ReduceScaleUnit(ctx, networkId, regionId).ScaleUnitsOperation(payload).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to reduce Enhanced Region scale units", err)
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
	}

	return resourceEnhancedRegionRead(ctx, d, m)
}

/*
resourceEnhancedRegionDelete Delete an Enhanced Region from an Enhanced Network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRegionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	regionId := d.Id()

	status, _, err := client.EnhancedRegionsAPI.DeleteEnhancedRegion(ctx, networkId, regionId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete Enhanced Region", err)
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
