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
resourceOpenvpn Setup the Openvpn Resource CRUD operations

@return &schema.Resource
*/
func resourceOpenvpn() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an OpenVPN client tunnel attached to one gateway of a " +
			"`checkpointsase_network`. OpenVPN tunnels are credentialed: on creation " +
			"the server returns a one-time `secret_access_key` (read it from state — " +
			"it's not retrievable from the API later). " +
			"**`version` is a credential-rotation trigger**, not a real version: " +
			"change its integer value to call the server's update endpoint and rotate " +
			"the tunnel's credentials. The actual numeric value is opaque — what " +
			"matters is that it changes from the previous run. " +
			"**`network_id`, `region_id`, `gateway_id`, and `tunnel_name` are " +
			"immutable** — changing any of them forces resource replacement.",
		CreateContext: resourceOpenvpnCreate,
		ReadContext:   resourceOpenvpnRead,
		UpdateContext: resourceOpenvpnUpdate,
		DeleteContext: resourceOpenvpnDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp of the last update to this resource.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the network's region. Returned by `checkpointsase_network.region.region_id`.",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Credential-rotation trigger. Increment (or change) this integer to trigger a server-side rotation of `access_key_id` / `secret_access_key`. The numeric value itself has no meaning beyond change detection.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the standard network the tunnel belongs to.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the SASE gateway that terminates this tunnel locally.",
			},
			"tunnel_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Display name for the OpenVPN tunnel.",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server-assigned credential ID for the OpenVPN client. Rotated when `version` changes.",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Server-assigned credential secret for the OpenVPN client. Returned on create and on each rotation; the API does not allow re-fetching this value later, so the terraform state is the only durable copy.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tunnel type (always `openvpn` server-side).",
			},
			"created_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp when the tunnel was created (server-assigned).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp when the tunnel was last updated server-side.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenvpnImportState,
		},
	}
}

/*
resourceOpenvpnImportState Import gateways
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceOpenvpnImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// get the network and tunnel id and validate
	if len(strings.Split(d.Id(), "-")) != 2 {
		return nil, fmt.Errorf("could not import tunnel without provider the network_id and the tunnel_id in format network_id-tunnel_id\n")
	}

	diagnostics := resourceOpenvpnRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import openvpn tunnel: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceOpenvpnCreate Create a Openvpn tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceOpenvpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the data from the resource data
	networkId := d.Get("network_id").(string)
	regionId := d.Get("region_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	tunnelName := d.Get("tunnel_name").(string)

	// create the tunnel payload
	baseTunnelBody := perimeter81Sdk.BaseTunnelValues{
		RegionID:   regionId,
		GatewayID:  gatewayId,
		TunnelName: tunnelName,
	}
	// create the tunnel and check for errors
	status, _, err := client.OpenVPNAPI.StandardCreateOpenVPNTunnel(ctx, networkId).BaseTunnelValues(baseTunnelBody).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Openvpn tunnel", err)
	}

	// get the status id from the status url
	var openvpnTunnelId string
	statusId := getIdFromUrl(status.GetStatusUrl())

	// check the status of the tunnel creation
	for {
		// check the status of the tunnel creation and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			return diags
		}
		// if the status is completed, get the tunnel id and break the loop
		if networkStatus.GetCompleted() {
			openvpnTunnelId, diags = getTunnelId(ctx, networkId, baseTunnelBody, *client, diags)
			if openvpnTunnelId == "" {
				return diags
			}
			break
		}
		// sleep for 20 seconds and check the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId(openvpnTunnelId)

	return resourceOpenvpnRead(ctx, d, m)
}

/*
resourceOpenvpnRead Read a Openvpn tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceOpenvpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the tunnel id and the network id from the resource data
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

	// get the tunnel and check for errors
	tunnel, _, err := client.OpenVPNAPI.StandardGetOpenVPNTunnel(ctx, networkId, tunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read openvpn tunnel", err)
	}
	// set the resource computed data
	if err := d.Set("tunnel_name", tunnel.TunnelName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set tunnel_name", err)
	}
	if err := d.Set("network_id", networkId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to network id", err)
	}
	if err := d.Set("region_id", tunnel.RegionID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to region id", err)
	}
	if err := d.Set("gateway_id", tunnel.GatewayID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set gateway id", err)
	}
	if err := d.Set("access_key_id", tunnel.GetAccessKeyId()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set access key id", err)
	}
	if err := d.Set("secret_access_key", tunnel.GetSecretAccessKey()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set secret access key", err)
	}
	if err := d.Set("type", tunnel.GetType()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set type", err)
	}
	if err := d.Set("updated_at", tunnel.GetUpdatedAt().String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set updated at", err)
	}
	if err := d.Set("created_at", tunnel.CreatedAt.String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set created at", err)
	}
	d.SetId(tunnelId)
	return diags
}

/*
resourceOpenvpnUpdate Update a Openvpn tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceOpenvpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// check if the version has changed
	if d.HasChange("version") {

		// get the tunnel id and the network id from the resource data
		tunnelId := d.Id()
		networkId := d.Get("network_id").(string)
		// update the tunnel and check for errors
		status, _, err := client.OpenVPNAPI.StandardUpdateOpenVPNTunnel(ctx, networkId, tunnelId).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update openvpn Tunnel", err)
		}

		// get the status id from the status url
		statusId := getIdFromUrl(status.GetStatusUrl())

		// check the status of the tunnel update
		for {
			// check the status of the tunnel update and check for errors
			networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
			if err != nil {
				d.Partial(true)
				return diags
			}
			// if the status is completed, break the loop
			if networkStatus.GetCompleted() {
				break
			}
			// sleep for 20 seconds and check the status again
			time.Sleep(20 * time.Second)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceOpenvpnRead(ctx, d, m)
}

/*
resourceOpenvpnDelete Delete a Openvpn tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceOpenvpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the tunnel id and the network id from the resource data
	tunnelId := d.Id()
	networkId := d.Get("network_id").(string)

	// delete the tunnel and check for errors
	status, _, err := client.OpenVPNAPI.StandardDeleteOpenVPNTunnel(ctx, networkId, tunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete openvpn tunnel", err)
	}

	// get the status id from the status url
	statusId := getIdFromUrl(status.GetStatusUrl())
	// check the status of the tunnel deletion
	for {
		// check the status of the tunnel deletion and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the status is completed, break the loop
		if networkStatus.GetCompleted() {
			break
		}
		// sleep for 20 seconds and check the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId("")
	return diags
}
