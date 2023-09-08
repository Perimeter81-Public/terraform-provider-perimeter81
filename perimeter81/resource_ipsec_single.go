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
resourceIpsecSingle Setup the IpSec-Signle Resource CRUD operations

@return &schema.Resource
*/
func resourceIpsecSingle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecSingleCreate,
		ReadContext:   resourceIpsecSingleRead,
		UpdateContext: resourceIpsecSingleUpdate,
		DeleteContext: resourceIpsecSingleDelete,
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
			"key_exchange": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ike_life_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lifetime": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dpd_delay": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dpd_timeout": {
				Type:     schema.TypeString,
				Required: true,
			},
			"passphrase": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"remote_public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_id": {
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
			"p81_gateway_subnets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_gateway_subnets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"phase1": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encryption": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dh": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					}},
			},
			"phase2": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encryption": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dh": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					}},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecSingleImportState,
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
func resourceIpsecSingleImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// get the network and tunnel id and validate
	if len(strings.Split(d.Id(), "-")) != 2 {
		return nil, fmt.Errorf("could not import tunnel without provider the network_id and the tunnel_id in format network_id-tunnel_id\n")
	}

	diagnostics := resourceIpsecSingleRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import ipsec single tunnel: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceIpsecSingleCreate Create a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecSingleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the ipsec single data from the terraform resource data and flatten what need to be flattened for the api
	networkId := d.Get("network_id").(string)
	regionId := d.Get("region_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	tunnelName := d.Get("tunnel_name").(string)
	keyExchange := d.Get("key_exchange").(string)
	remotePublicIP := d.Get("remote_public_ip").(string)
	passphrase := d.Get("passphrase").(string)
	dpdTimeout := d.Get("dpd_timeout").(string)
	dpdDelay := d.Get("dpd_delay").(string)
	lifetime := d.Get("lifetime").(string)
	ikeLifeTime := d.Get("ike_life_time").(string)
	p81GatewaySubnets := flattenStringsArrayData(d.Get("p81_gateway_subnets").([]interface{}))
	remoteGatewaySubnets := flattenStringsArrayData(d.Get("remote_gateway_subnets").([]interface{}))
	phase1Data := d.Get("phase1").([]interface{})[0].(map[string]interface{})
	phase2Data := d.Get("phase2").([]interface{})[0].(map[string]interface{})
	authPhase1 := flattenStringsArrayData(phase1Data["auth"].([]interface{}))
	authPhase2 := flattenStringsArrayData(phase2Data["auth"].([]interface{}))
	encryptionPhase1 := flattenStringsArrayData(phase1Data["encryption"].([]interface{}))
	encryptionPhase2 := flattenStringsArrayData(phase2Data["encryption"].([]interface{}))
	dhPhase1 := flattenIntsArrayData(phase1Data["dh"].([]interface{}))
	dhPhase2 := flattenIntsArrayData(phase2Data["dh"].([]interface{}))

	// create the ipsec single payload
	phase1 := perimeter81Sdk.IpSecPhase{
		Auth:       authPhase1,
		Encryption: encryptionPhase1,
		Dh:         dhPhase1,
	}
	phase2 := perimeter81Sdk.IpSecPhase{
		Auth:       authPhase2,
		Encryption: encryptionPhase2,
		Dh:         dhPhase2,
	}
	ipSecSingleBody := perimeter81Sdk.CreateIpSecSinglePayload{
		RegionID:             regionId,
		GatewayID:            gatewayId,
		TunnelName:           tunnelName,
		KeyExchange:          keyExchange,
		RemotePublicIP:       remotePublicIP,
		Lifetime:             lifetime,
		IkeLifeTime:          ikeLifeTime,
		Passphrase:           passphrase,
		DpdTimeout:           dpdTimeout,
		DpdDelay:             dpdDelay,
		P81GatewaySubnets:    p81GatewaySubnets,
		Phase1:               &phase1,
		Phase2:               &phase2,
		RemoteGatewaySubnets: remoteGatewaySubnets,
	}

	// create the ipsec single tunnel and check for errors
	status, _, err := client.IPSecSingleApi.CreateIPSecSingleTunnel(ctx, ipSecSingleBody, networkId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create IpsecSingle tunnel", err)
	}

	// get the status id of the ipsec-single tunnel creation
	var ipSecSingleTunnelId string
	statusId := getIdFromUrl(status.StatusUrl)

	// check the status of the ipsec-redundant tunnel creation
	for {
		// check the status of the network that contains the ipsec-redundant tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the network status is completed, get the ipsec-single tunnel id and break the loop
		if networkStatus.Completed {
			baseTunnelBody := perimeter81Sdk.BaseTunnelValues{
				RegionID:   regionId,
				GatewayID:  gatewayId,
				TunnelName: tunnelName,
			}
			ipSecSingleTunnelId, diags = getTunnelId(ctx, networkId, baseTunnelBody, *client, diags)
			if ipSecSingleTunnelId == "" {
				return diags
			}
			break
		}
		// delay for 20 seconds before checking the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId(ipSecSingleTunnelId)

	return resourceIpsecSingleRead(ctx, d, m)
}

/*
resourceIpsecSingleRead Read a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecSingleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the ipsec-single tunnel id and the network id
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

	// get the ipsec-single tunnel and check for errors
	tunnel, _, err := client.IPSecSingleApi.GetIPSecSingleTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read ipsec-single tunnel", err)
	}
	// set the ipsec-single tunnel computed data
	if err := d.Set("created_at", tunnel.CreatedAt.String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set createdAt", err)
	}
	if err := d.Set("updated_at", tunnel.UpdatedAt.String()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set updatedAt", err)
	}
	if err := d.Set("remote_id", tunnel.RemoteID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set remoteId", err)
	}
	if err := d.Set("key_exchange", tunnel.KeyExchange); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set keyExchange", err)
	}
	if err := d.Set("remote_public_ip", tunnel.RemotePublicIP); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Remote Public IP", err)
	}
	if err := d.Set("ike_life_time", tunnel.IkeLifeTime); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set ike lifetime", err)
	}
	if err := d.Set("lifetime", tunnel.LifeTime); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set lifetime", err)
	}
	if err := d.Set("dpd_delay", tunnel.DpdDelay); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Dpd delay", err)
	}
	if err := d.Set("dpd_timeout", tunnel.DpdTimeout); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Dpd timeout", err)
	}
	if err := d.Set("passphrase", tunnel.Passphrase); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set passphrase", err)
	}
	if err := d.Set("network_id", networkId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set networkId", err)
	}
	if err := d.Set("region_id", tunnel.RegionID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set regionId", err)
	}
	if err := d.Set("gateway_id", tunnel.GatewayID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set gatewayId", err)
	}
	if err := d.Set("tunnel_name", tunnel.TunnelName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set name", err)
	}
	if err := d.Set("p81_gateway_subnets", tunnel.P81GatewaySubnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set p81 gateway subnets", err)
	}
	if err := d.Set("remote_gateway_subnets", tunnel.RemoteGatewaySubnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set remote subnets", err)
	}
	if err := d.Set("phase1", flattenPhasesData(tunnel.Phase1)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set phase1", err)
	}
	if err := d.Set("phase2", flattenPhasesData(tunnel.Phase2)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set phase2", err)
	}
	d.SetId(tunnelId)
	return diags
}

/*
resourceIpsecSingleUpdate Update a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecSingleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// check if the ipsec single data has changes
	if d.HasChanges("key_exchange", "remote_public_ip", "remote_id", "passphrase", "dpd_timeout", "dpd_delay", "lifetime", "ike_life_time", "p81_gateway_subnets", "remote_gateway_subnets", "phase1", "phase2") {
		// get the ipsec-single tunnel id and the network id and data from the terraform resource data and flatten the data that need to be flattened
		tunnelId := d.Id()
		networkId := d.Get("network_id").(string)
		keyExchange := d.Get("key_exchange").(string)
		remotePublicIP := d.Get("remote_public_ip").(string)
		passphrase := d.Get("passphrase").(string)
		dpdTimeout := d.Get("dpd_timeout").(string)
		dpdDelay := d.Get("dpd_delay").(string)
		lifetime := d.Get("lifetime").(string)
		ikeLifeTime := d.Get("ike_life_time").(string)
		remoteId := d.Get("remote_id").(string)
		p81GatewaySubnets := flattenStringsArrayData(d.Get("p81_gateway_subnets").([]interface{}))
		remoteGatewaySubnets := flattenStringsArrayData(d.Get("remote_gateway_subnets").([]interface{}))
		phase1Data := d.Get("phase1").([]interface{})[0].(map[string]interface{})
		phase2Data := d.Get("phase2").([]interface{})[0].(map[string]interface{})
		authPhase1 := flattenStringsArrayData(phase1Data["auth"].([]interface{}))
		authPhase2 := flattenStringsArrayData(phase2Data["auth"].([]interface{}))
		encryptionPhase1 := flattenStringsArrayData(phase1Data["encryption"].([]interface{}))
		encryptionPhase2 := flattenStringsArrayData(phase2Data["encryption"].([]interface{}))
		dhPhase1 := flattenIntsArrayData(phase1Data["dh"].([]interface{}))
		dhPhase2 := flattenIntsArrayData(phase2Data["dh"].([]interface{}))

		// create the ipsec-single tunnel body
		phase1 := perimeter81Sdk.IpSecPhase{
			Auth:       authPhase1,
			Encryption: encryptionPhase1,
			Dh:         dhPhase1,
		}
		phase2 := perimeter81Sdk.IpSecPhase{
			Auth:       authPhase2,
			Encryption: encryptionPhase2,
			Dh:         dhPhase2,
		}

		ipSecSingleDetails := perimeter81Sdk.IpSecSingleDetails{
			KeyExchange:          keyExchange,
			RemotePublicIP:       remotePublicIP,
			Passphrase:           passphrase,
			DpdTimeout:           dpdTimeout,
			DpdDelay:             dpdDelay,
			Lifetime:             lifetime,
			IkeLifeTime:          ikeLifeTime,
			P81GatewaySubnets:    p81GatewaySubnets,
			RemoteGatewaySubnets: remoteGatewaySubnets,
			Phase1:               &phase1,
			Phase2:               &phase2,
			RemoteID:             remoteId,
		}

		// update the ipsec-single tunnel and check for errors
		status, _, err := client.IPSecSingleApi.UpdateIPSecSingleTunnel(ctx, ipSecSingleDetails, networkId, tunnelId)
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update ipsec-single Tunnel", err)
		}

		// get the status id from the status url
		statusId := getIdFromUrl(status.StatusUrl)

		// check the status of the ipsec-single tunnel and check for errors
		for {
			// check the status of the ipsec-single tunnel and check for errors
			networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
			if err != nil {
				d.Partial(true)
				return diags
			}
			// if the ipsec-single tunnel status is completed break the loop
			if networkStatus.Completed {
				break
			}
			// sleep for 20 seconds
			time.Sleep(20 * time.Second)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceIpsecSingleRead(ctx, d, m)
}

/*
resourceIpsecSingleDelete Delete a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecSingleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the ipsec-single tunnel id and the network id from the terraform resource data
	tunnelId := d.Id()
	networkId := d.Get("network_id").(string)

	// delete the ipsec-single tunnel and check for errors
	status, _, err := client.IPSecSingleApi.DeleteIPSecSingleTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete ipsec-single tunnel", err)
	}

	// get the status id from the status url
	statusId := getIdFromUrl(status.StatusUrl)
	// check the status of the ipsec-single tunnel and check for errors
	for {
		// check the status of the ipsec-single tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the ipsec-single tunnel status is completed break the loop
		if networkStatus.Completed {
			break
		}
		// sleep for 20 seconds
		time.Sleep(20 * time.Second)
	}
	d.SetId("")
	return diags
}
