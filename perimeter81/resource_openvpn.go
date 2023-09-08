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
resourceOpenvpn Setup the Openvpn Resource CRUD operations

@return &schema.Resource
*/
func resourceOpenvpn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenvpnCreate,
		ReadContext:   resourceOpenvpnRead,
		UpdateContext: resourceOpenvpnUpdate,
		DeleteContext: resourceOpenvpnDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
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
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
	status, _, err := client.OpenVPNApi.CreateOpenVPNTunnel(ctx, baseTunnelBody, networkId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Openvpn tunnel", err)
	}

	// get the status id from the status url
	var openvpnTunnelId string
	statusId := getIdFromUrl(status.StatusUrl)

	// check the status of the tunnel creation
	for {
		// check the status of the tunnel creation and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			return diags
		}
		// if the status is completed, get the tunnel id and break the loop
		if networkStatus.Completed {
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
	tunnel, _, err := client.OpenVPNApi.GetOpenVPNTunnel(ctx, networkId, tunnelId)
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
	if err := d.Set("access_key_id", tunnel.AccessKeyId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set access key id", err)
	}
	if err := d.Set("secret_access_key", tunnel.SecretAccessKey); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set secret access key", err)
	}
	if err := d.Set("type", tunnel.Type_); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set secret access key", err)
	}
	if err := d.Set("updated_at", tunnel.UpdatedAt.String()); err != nil {
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
		status, _, err := client.OpenVPNApi.UpdateOpenVPNTunnel(ctx, networkId, tunnelId)
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update openvpn Tunnel", err)
		}

		// get the status id from the status url
		statusId := getIdFromUrl(status.StatusUrl)

		// check the status of the tunnel update
		for {
			// check the status of the tunnel update and check for errors
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
	status, _, err := client.OpenVPNApi.DeleteOpenVPNTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete openvpn tunnel", err)
	}

	// get the status id from the status url
	statusId := getIdFromUrl(status.StatusUrl)
	// check the status of the tunnel deletion
	for {
		// check the status of the tunnel deletion and check for errors
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
