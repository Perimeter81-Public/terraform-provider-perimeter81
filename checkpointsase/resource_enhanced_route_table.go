package checkpointsase

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
resourceEnhancedRouteTable Setup the Enhanced Route Table Resource CRUD operations

@return &schema.Resource
*/
func resourceEnhancedRouteTable() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a route-table entry for a `checkpointsase_enhanced_network`. " +
			"A route directs traffic for the specified `subnets` through a tunnel " +
			"(or list of dynamic tunnels). " +
			"Use `type = \"static\"` with `tunnel_id` for static tunnels, or " +
			"`type = \"dynamic\"` with `tunnel_ids` for dynamic tunnels. " +
			"The selected tunnel(s) must not already have a route table attached. " +
			"**`network_id`, `type`, and the chosen tunnel field are immutable** — " +
			"changing any of them forces resource replacement.",
		CreateContext: resourceEnhancedRouteTableCreate,
		ReadContext:   resourceEnhancedRouteTableRead,
		UpdateContext: resourceEnhancedRouteTableUpdate,
		DeleteContext: resourceEnhancedRouteTableDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp of the last update to this resource.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the enhanced network this route table entry belongs to.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The route type. Must be either `static` or `dynamic`. " +
					"Use with `tunnel_id` for static, or `tunnel_ids` for dynamic.",
				ValidateFunc: validation.StringInSlice([]string{"static", "dynamic"}, false),
			},
			"tunnel_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"tunnel_ids"},
				Description:   "The static tunnel ID. Required when type is `static`. Mutually exclusive with `tunnel_ids`. The selected static tunnel must not already have a route table.",
			},
			"tunnel_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"tunnel_id"},
				Description:   "The list of dynamic tunnel IDs. Required when type is `dynamic`. Mutually exclusive with `tunnel_id`. The selected dynamic tunnels must not already have a route table.",
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"subnets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of subnet CIDR blocks for the route table entry.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"propagated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the route is propagated automatically.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEnhancedRouteTableImportState,
		},
	}
}

/*
resourceEnhancedRouteTableImportState Import an enhanced route table entry by its ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceEnhancedRouteTableImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceEnhancedRouteTableRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import enhanced route table: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceEnhancedRouteTableCreate Create an Enhanced Route Table entry.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRouteTableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	routeType := d.Get("type").(string)
	subnets := flattenStringsArrayData(d.Get("subnets").([]interface{}))

	var status *perimeter81Sdk.AsyncOperationResponse
	var err error

	switch routeType {
	case "static":
		tunnelId, ok := d.GetOk("tunnel_id")
		if !ok {
			return appendErrorDiags(diags, "tunnel_id is required for static route type", fmt.Errorf("tunnel_id must be set when type is 'static'"))
		}
		payload := perimeter81Sdk.EnhancedRouteTableStaticCreate{
			TunnelId: tunnelId.(string),
			Subnets:  subnets,
		}
		status, _, err = client.EnhancedRouteTablesAPI.CreateStaticRoute(ctx, networkId).EnhancedRouteTableStaticCreate(payload).Execute()
	case "dynamic":
		tunnelIds := flattenStringsArrayData(d.Get("tunnel_ids").([]interface{}))
		if len(tunnelIds) == 0 {
			return appendErrorDiags(diags, "tunnel_ids is required for dynamic route type", fmt.Errorf("tunnel_ids must be non-empty when type is 'dynamic'"))
		}
		payload := perimeter81Sdk.EnhancedRouteTableDynamicCreate{
			TunnelIds: tunnelIds,
			Subnets:   subnets,
		}
		status, _, err = client.EnhancedRouteTablesAPI.CreateDynamicRoute(ctx, networkId).EnhancedRouteTableDynamicCreate(payload).Execute()
	default:
		return appendErrorDiags(diags, "Invalid route type", fmt.Errorf("type must be 'static' or 'dynamic', got: %s", routeType))
	}

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Enhanced Route Table entry", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var routeId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			routeId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(routeId)
	return resourceEnhancedRouteTableRead(ctx, d, m)
}

/*
resourceEnhancedRouteTableRead Read an Enhanced Route Table entry.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	routeId := d.Id()

	routeData, _, err := client.EnhancedRouteTablesAPI.GetRouteEntry(ctx, networkId, routeId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Enhanced Route Table entry", err)
	}

	if err := d.Set("subnets", routeData.Subnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Route Table subnets", err)
	}
	if err := d.Set("tunnel_ids", routeData.TunnelIds); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Route Table tunnel_ids", err)
	}
	if err := d.Set("propagated", routeData.Propagated); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Route Table propagated", err)
	}

	return diags
}

/*
resourceEnhancedRouteTableUpdate Update an Enhanced Route Table entry.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRouteTableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	if d.HasChange("subnets") {
		networkId := d.Get("network_id").(string)
		routeId := d.Id()
		subnets := flattenStringsArrayData(d.Get("subnets").([]interface{}))

		payload := perimeter81Sdk.RouteTableUpdate{
			Subnets: subnets,
		}

		_, _, err := client.EnhancedRouteTablesAPI.UpdateRouteEntry(ctx, networkId, routeId).RouteTableUpdate(payload).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update Enhanced Route Table entry", err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceEnhancedRouteTableRead(ctx, d, m)
}

/*
resourceEnhancedRouteTableDelete Delete an Enhanced Route Table entry.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedRouteTableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	routeId := d.Id()

	status, _, err := client.EnhancedRouteTablesAPI.DeleteRouteEntry(ctx, networkId, routeId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete Enhanced Route Table entry", err)
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
