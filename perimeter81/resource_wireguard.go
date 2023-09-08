package perimeter81

import (
	"context"
	"fmt"
	"strings"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceWireguard Setup the Wireguard tunnel Resource CRUD operations

@return &schema.Resource
*/
func resourceWireguard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWireguardCreate,
		ReadContext:   resourceWireguardRead,
		UpdateContext: resourceWireguardUpdate,
		DeleteContext: resourceWireguardDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_subnets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceWireguardImportState,
		},
	}
}

/*
resourceWireguardImportState Import gateways
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceWireguardImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// get the network and tunnel id and validate
	if len(strings.Split(d.Id(), "-")) != 2 {
		return nil, fmt.Errorf("could not import tunnel without provider the network_id and the tunnel_id in format network_id-tunnel_id\n")
	}

	diagnostics := resourceWireguardRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import wireguard: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceWireguardCreate Create a Wireguard tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceWireguardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the resource data from the terraform resource and flatten what needs to be flattened
	networkId := d.Get("network_id").(string)
	remoteEndpoint := d.Get("remote_endpoint").(string)
	regionId := d.Get("region_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	tunnelName := d.Get("tunnel_name").(string)
	remoteSubnets := flattenStringsArrayData(d.Get("remote_subnets").([]interface{}))

	// create the wireguard tunnel payload
	wireguardBody := perimeter81Sdk.CreateWireguardTunnelPayload{
		RegionID:       regionId,
		RemoteEndpoint: remoteEndpoint,
		GatewayID:      gatewayId,
		TunnelName:     tunnelName,
		RemoteSubnets:  remoteSubnets,
	}

	// create the wireguard tunnel and check for errors
	status, _, err := client.WireguardApi.CreateWireguardTunnel(ctx, wireguardBody, networkId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Wireguard tunnel", err)
	}

	// get the status id from the status url
	var wireguardTunnelId string
	statusId := getIdFromUrl(status.StatusUrl)

	// check the status of the wireguard tunnel creation
	for {
		// check the status of the wireguard tunnel creation and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the status is completed, get the tunnel id and break the loop
		if networkStatus.Completed {
			baseTunnelBody := perimeter81Sdk.BaseTunnelValues{
				RegionID:   regionId,
				GatewayID:  gatewayId,
				TunnelName: tunnelName,
			}
			wireguardTunnelId, diags = getTunnelId(ctx, networkId, baseTunnelBody, *client, diags)
			if wireguardTunnelId == "" {
				return diags
			}
			break
		}
		// sleep for 20 seconds and check the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId(wireguardTunnelId)

	return diags
}

/*
resourceWireguardRead Read a Wireguard tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceWireguardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the tunnel id and the network id from the terraform resource
	ids := strings.Split(d.Id(), "-")
	var networkId string
	var tunnelId string
	if len(ids) == 1 {
		tunnelId = d.Id()
		networkId = d.Get("network_id").(string)
	} else {
		networkId = ids[0]
		tunnelId = ids[1]
	}

	// get the wireguard tunnel and check for errors
	tunnel, _, err := client.WireguardApi.GetWireguardTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read wireguard tunnel", err)
	}
	// set the resource computed data
	if err := d.Set("remote_endpoint", tunnel.RemoteEndpoint); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set remoteendpoint", err)
	}
	if err := d.Set("network_id", networkId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set networkId", err)
	}
	if err := d.Set("gateway_id", tunnel.GatewayID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set gatewayId", err)
	}
	if err := d.Set("region_id", tunnel.RegionID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set regionId", err)
	}
	if err := d.Set("tunnel_name", tunnel.TunnelName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set tunnelName", err)
	}
	if err := d.Set("created_at", tunnel.CreatedAt.String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set createdAt", err)
	}
	if err := d.Set("updated_at", tunnel.UpdatedAt.String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set updatedAt", err)
	}
	if err := d.Set("remote_subnets", tunnel.RemoteSubnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set remotesubnets", err)
	}
	d.SetId(tunnelId)

	return diags
}

/*
resourceWireguardUpdate Update a Wireguard tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceWireguardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// check if the remote endpoint or the remote subnets have changed
	if d.HasChanges("remote_endpoint", "remote_subnets") {

		// get the resource data from the terraform resource and flatten what needs to be flattened
		tunnelId := d.Id()
		networkId := d.Get("network_id").(string)
		remoteEndpoint := d.Get("remote_endpoint").(string)
		remoteSubnets := flattenStringsArrayData(d.Get("remote_subnets").([]interface{}))
		// create the wireguard tunnel update payload
		wireguardDetails := perimeter81Sdk.WireGuradDetails{
			RemoteEndpoint: remoteEndpoint,
			RemoteSubnets:  remoteSubnets,
		}
		// update the wireguard tunnel and check for errors
		status, _, err := client.WireguardApi.UpdateWireguardTunnel(ctx, wireguardDetails, networkId, tunnelId)
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update wireguard Tunnel", err)
		}

		// get the status id from the status url
		statusId := getIdFromUrl(status.StatusUrl)
		// check the status of the wireguard tunnel update
		for {
			// check the status of the wireguard tunnel update and check for errors
			networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
			if err != nil {
				d.Partial(true)
				return diags
			}
			// if the status is completed, break the loop
			if networkStatus.Completed {
				break
			}
			// sleep for 20 seconds and check the status again
			time.Sleep(20 * time.Second)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceWireguardRead(ctx, d, m)
}

/*
resourceWireguardDelete Delete a Wireguard tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceWireguardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the tunnel id and the network id from the terraform resource
	tunnelId := d.Id()
	networkId := d.Get("network_id").(string)

	// delete the wireguard tunnel and check for errors
	status, _, err := client.WireguardApi.DeleteWireguardTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete wireguard tunnel", err)
	}

	// get the status id from the status url
	statusId := getIdFromUrl(status.StatusUrl)
	// check the status of the wireguard tunnel deletion
	for {
		// check the status of the wireguard tunnel deletion and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the status is completed, break the loop
		if networkStatus.Completed {
			break
		}
		// sleep for 20 seconds and check the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId("")
	return diags
}
