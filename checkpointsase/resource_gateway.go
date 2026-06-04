package checkpointsase

import (
	"context"
	"fmt"
	"strings"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceGateway Setup the gateway Resource CRUD operations

@return &schema.Resource
*/
func resourceGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the gateway pool of a single region within a `checkpointsase_network`. " +
			"Each `gateways` block declares one gateway (named, with an `idle` flag). " +
			"The resource is keyed by the composite `<network_id>-<region_id>` for import. " +
			"**`network_id` and `region_id` are immutable** — changing either would orphan the " +
			"managed gateway list from its region and is not supported.",
		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceGatewayRead,
		UpdateContext: resourceGatewayUpdate,
		DeleteContext: resourceGatewayDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp of the last update to this resource.",
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The ID of the network's region within which to manage the gateway pool. " +
					"This is the network-region ID returned by `checkpointsase_network.region.region_id`, " +
					"not the cloud region ID (`cpregion_id`).",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the standard network whose gateways this resource manages.",
			},
			"gateways": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of gateways to provision in the region. Order is not significant.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The gateway name. Must be unique within the region.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID assigned to the gateway by the server.",
						},
						"dns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DNS hostname assigned to the gateway.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address assigned to the gateway.",
						},
						"idle": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the gateway is idle (disabled for user traffic). Set to `false` to make the gateway active.",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceGatewayImportState,
		},
	}
}

/*
resourceGatewayImportState Import gateways
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceGatewayImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	var diagnostics diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()
	// get the network and region id and validate
	ids := strings.Split(d.Id(), "-")
	if len(ids) != 2 {
		return nil, fmt.Errorf("could not import gateways without provider the network_id and the region_id in format network_id-region_id\n")
	}

	// call the api and check if there is an error
	networkData, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, ids[0]).Execute()
	if err != nil {
		diagnostics = appendErrorDiags(diagnostics, "Unable to find Network", err)
	}

	// get the gateways that are available inside that region and validate
	gateways := getGatewaysInArray(ids[1], networkData)
	if len(gateways) == 0 {
		return nil, fmt.Errorf("could not import gateways please make sure that the netwrok_id and the region_id are correct\n")
	}
	newGateways := make([]GatewayConfig, 0)
	for _, gateway := range gateways {
		newGateways = append(newGateways, GatewayConfig{
			Idle: false,
			Id:   gateway.Id,
			Name: "$" + gateway.Id + "$",
			Dns:  gateway.Dns,
			Ip:   gateway.Ip,
		})
	}
	// set the gateway and ids after getting the gateway id to the resource data
	if err := d.Set("gateways", flattenGateways(newGateways)); err != nil {
		return nil, fmt.Errorf("Unable to set Gateway data after import\n")
	}
	if err := d.Set("network_id", ids[0]); err != nil {
		return nil, fmt.Errorf("Unable to set network_id after import\n")
	}
	if err := d.Set("region_id", ids[1]); err != nil {
		return nil, fmt.Errorf("Unable to set region_id after import\n")
	}
	d.SetId(ids[0] + "-" + ids[1])
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import gateways: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceGatewayCreate Create a gateway
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the gateways data from the resource data

	gateways := flattenGatewaysData(d.Get("gateways").([]interface{}))
	network_id := d.Get("network_id").(string)
	region_id := d.Get("region_id").(string)

	if check, name := checkGatewayDuplicatesInArray(gateways); check {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Gateway", fmt.Errorf("gateway name %s is duplicated", name))
	}

	// add the gateways to the region and check for errors
	diags, err := addGatewayToRegion(ctx, client, gateways, network_id, region_id, diags)
	if err != nil {
		d.Partial(true)
		return diags
	}

	// set the gateway after getting the gateway id to the resource data
	if err := d.Set("gateways", flattenGateways(gateways)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Gateway data", err)
	}
	d.SetId(network_id + "-" + region_id)
	return diags
}

/*
resourceGatewayRead Read a gateway
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	regionId := d.Get("region_id").(string)

	networkData, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Network for gateway read", err)
	}

	serverGateways := getGatewaysInArray(regionId, networkData)
	if len(serverGateways) == 0 {
		// Region/network no longer has gateways — drop from state.
		d.SetId("")
		return diags
	}

	// Map server response by ID for lookup.
	serverByID := make(map[string]perimeter81Sdk.NetworkInstance, len(serverGateways))
	for _, g := range serverGateways {
		serverByID[g.Id] = g
	}

	// Reconcile with prior state. The API doesn't return `name` / `idle`,
	// so preserve those from the user's HCL/state and refresh only the
	// server-computed fields (id/dns/ip).
	prior := flattenGatewaysData(d.Get("gateways").([]interface{}))
	refreshed := make([]GatewayConfig, 0, len(prior))
	for _, p := range prior {
		if srv, ok := serverByID[p.Id]; ok {
			refreshed = append(refreshed, GatewayConfig{
				Id:   srv.Id,
				Dns:  srv.Dns,
				Ip:   srv.Ip,
				Name: p.Name, // preserve — API doesn't carry it
				Idle: p.Idle, // preserve — API doesn't carry it
			})
			delete(serverByID, p.Id)
		}
	}
	// Adopt-style: do NOT surface gateways that exist server-side but are
	// not in our prior state. Standard networks auto-provision an initial
	// gateway at network-create time; surfacing it here would cause every
	// post-apply plan to flag the auto-gateway as extraneous and would
	// destroy it on apply. To bring an unmanaged gateway under this
	// resource, use `terraform import 'checkpointsase_gateway.X'
	// <network_id>-<gateway_id>`.
	if err := d.Set("gateways", flattenGateways(refreshed)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Gateway data", err)
	}
	return diags
}

/*
resourceGatewayUpdate Update a gateway
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()
	// check if the region_id or network_id is changed
	if d.HasChanges("region_id", "network_id") {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to change network_id or region_id", fmt.Errorf("region_id and network_id cannot be updated"))
	}
	// check if the gateways is changed
	if d.HasChange("gateways") {
		// get the old and new gateways and get the gateways info
		oldGateways, newGateways := d.GetChange("gateways")
		network_id := d.Get("network_id").(string)
		region_id := d.Get("region_id").(string)

		// flatten the gateways data to match the schema
		oldGatewaysFlattened := flattenGatewaysData(oldGateways.([]interface{}))
		newGatewaysFlattened := flattenGatewaysData(newGateways.([]interface{}))

		// handle the name change
		pass := false
		if len(oldGatewaysFlattened) == len(newGatewaysFlattened) {
			for _, gateway := range oldGatewaysFlattened {
				if gateway.Name != "$"+gateway.Id+"$" {
					pass = true
					break
				}

			}
		}
		if pass || len(oldGatewaysFlattened) != len(newGatewaysFlattened) {
			if check, name := checkGatewayDuplicatesInArray(newGatewaysFlattened); check {
				d.Partial(true)
				return appendErrorDiags(diags, "Unable to create Gateway", fmt.Errorf("gateway name %s is duplicated", name))
			}
			// get the gateways to be added and add them to the region and check for errors
			gateways := getNewGateway(oldGatewaysFlattened, newGatewaysFlattened)
			diags, err := addGatewayToRegion(ctx, client, gateways, network_id, region_id, diags)
			if err != nil {
				return diags
			}

			// add the id to the new gateways after being created
			for index, gateway := range newGatewaysFlattened {
				if gateway.Id == "" {
					for _, newGateway := range gateways {
						if gateway.Name == newGateway.Name {
							newGatewaysFlattened[index].Id = newGateway.Id
						}
					}
				}
			}
			// get the gateways to be deleted and delete them from the region and check for errors
			gateways = getGatewayToBeDeleted(oldGatewaysFlattened, newGatewaysFlattened)
			diags, err = deleteGatewayFromRegion(ctx, client, gateways, network_id, region_id, diags)
			if err != nil {
				return diags
			}
			// set the gateway after getting the gateway id to the resource data
			if err := d.Set("gateways", flattenGateways(newGatewaysFlattened)); err != nil {
				return appendErrorDiags(diags, "Unable to set Gateway data", err)
			}
		} else {
			// set the gateway after getting the gateway id to the resource data
			if err := d.Set("gateways", flattenGateways(newGatewaysFlattened)); err != nil {
				return appendErrorDiags(diags, "Unable to set Gateway data", err)
			}
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return diags
}

/*
resourceGatewayDelete Delete a gateway
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()
	// get the gateways data from the resource data
	gateways := flattenGatewaysData(d.Get("gateways").([]interface{}))
	network_id := d.Get("network_id").(string)
	region_id := d.Get("region_id").(string)

	// delete the gateways from the region and check for errors
	diags, err := deleteGatewayFromRegion(ctx, client, gateways, network_id, region_id, diags)
	if err != nil {
		d.Partial(true)
		return diags
	}
	d.SetId("")
	return diags
}
