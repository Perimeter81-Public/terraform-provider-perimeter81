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
resourceIpsecRedundant Setup the IpSec-Redundant Resource CRUD operations

@return &schema.Resource
*/
func resourceIpsecRedundant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsecRedundantCreate,
		ReadContext:   resourceIpsecRedundantRead,
		UpdateContext: resourceIpsecRedundantUpdate,
		DeleteContext: resourceIpsecRedundantDelete,
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
			"tunnel_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"advanced_settings": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					}},
			},
			"shared_settings": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					}},
			},
			"tunnel1": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passphrase": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tunnel_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"p81_gwinternal_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_gwinternal_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_public_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_asn": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					}},
			},
			"tunnel2": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passphrase": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
						"tunnel_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"p81_gwinternal_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_gwinternal_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_public_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_asn": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					}},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsecRedundantImportState,
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
func resourceIpsecRedundantImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	// get the network and tunnel id and validate
	if len(strings.Split(d.Id(), "-")) != 2 {
		return nil, fmt.Errorf("could not import tunnel without provider the network_id and the tunnel_id in format network_id-tunnel_id\n")
	}

	diagnostics := resourceIpsecRedundantRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import ipsec redundant tunnel: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceIpsecRedundantCreate Create a Ipsec Redundant Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/

func resourceIpsecRedundantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)

	if ctx == nil {
		ctx = context.Background()
	}

	// get the tunnel data from the terraform resource and flatten what need to be flattened for the api
	networkId := d.Get("network_id").(string)
	regionId := d.Get("region_id").(string)
	tunnelName := d.Get("tunnel_name").(string)
	tunnel1Data := d.Get("tunnel1").([]interface{})[0].(map[string]interface{})
	tunnel2Data := d.Get("tunnel2").([]interface{})[0].(map[string]interface{})
	gatewayId1 := tunnel1Data["gateway_id"].(string)
	passphrase1 := tunnel1Data["passphrase"].(string)
	p81GWinternalIP1 := tunnel1Data["p81_gwinternal_ip"].(string)
	remoteGWinernalIP1 := tunnel1Data["remote_gwinternal_ip"].(string)
	remotePublicIP1 := tunnel1Data["remote_public_ip"].(string)
	remoteAsn1 := tunnel1Data["remote_asn"].(float64)
	remoteId1 := tunnel1Data["remote_id"].(string)
	gatewayId2 := tunnel2Data["gateway_id"].(string)
	passphrase2 := tunnel2Data["passphrase"].(string)
	p81GWinternalIP2 := tunnel2Data["p81_gwinternal_ip"].(string)
	remoteGWinernalIP2 := tunnel2Data["remote_gwinternal_ip"].(string)
	remotePublicIP2 := tunnel2Data["remote_public_ip"].(string)
	remoteAsn2 := tunnel2Data["remote_asn"].(float64)
	remoteId2 := tunnel2Data["remote_id"].(string)
	sharedSettingsData := d.Get("shared_settings").([]interface{})[0].(map[string]interface{})
	p81GatewaySubnets := flattenStringsArrayData(sharedSettingsData["p81_gateway_subnets"].([]interface{}))
	remoteGatewaySubnets := flattenStringsArrayData(sharedSettingsData["remote_gateway_subnets"].([]interface{}))
	advancedSettingsData := d.Get("advanced_settings").([]interface{})[0].(map[string]interface{})
	keyExchange := advancedSettingsData["key_exchange"].(string)
	dpdTimeout := advancedSettingsData["dpd_timeout"].(string)
	dpdDelay := advancedSettingsData["dpd_delay"].(string)
	lifetime := advancedSettingsData["lifetime"].(string)
	ikeLifeTime := advancedSettingsData["ike_life_time"].(string)
	phase1Data := advancedSettingsData["phase1"].([]interface{})[0].(map[string]interface{})
	phase2Data := advancedSettingsData["phase2"].([]interface{})[0].(map[string]interface{})
	authPhase1 := flattenStringsArrayData(phase1Data["auth"].([]interface{}))
	authPhase2 := flattenStringsArrayData(phase2Data["auth"].([]interface{}))
	encryptionPhase1 := flattenStringsArrayData(phase1Data["encryption"].([]interface{}))
	encryptionPhase2 := flattenStringsArrayData(phase2Data["encryption"].([]interface{}))
	dhPhase1 := flattenIntsArrayData(phase1Data["dh"].([]interface{}))
	dhPhase2 := flattenIntsArrayData(phase2Data["dh"].([]interface{}))

	// create the payload for the api
	ipSecRedundantBody := perimeter81Sdk.CreateIpSecRedundantPayload{
		RegionID:   regionId,
		TunnelName: tunnelName,
		Tunnel1: &perimeter81Sdk.IpSecRedundantTunnelPayload{
			Passphrase:        passphrase1,
			GatewayID:         gatewayId1,
			P81GWinternalIP:   p81GWinternalIP1,
			RemoteGWinernalIP: remoteGWinernalIP1,
			RemotePublicIP:    remotePublicIP1,
			RemoteASN:         remoteAsn1,
			RemoteID:          remoteId1,
		},
		Tunnel2: &perimeter81Sdk.IpSecRedundantTunnelPayload{
			Passphrase:        passphrase2,
			GatewayID:         gatewayId2,
			P81GWinternalIP:   p81GWinternalIP2,
			RemoteGWinernalIP: remoteGWinernalIP2,
			RemotePublicIP:    remotePublicIP2,
			RemoteASN:         remoteAsn2,
			RemoteID:          remoteId2,
		},
		SharedSettings: &perimeter81Sdk.IpSecSharedSettings{
			P81GatewaySubnets:    p81GatewaySubnets,
			RemoteGatewaySubnets: remoteGatewaySubnets,
		},

		AdvancedSettings: &perimeter81Sdk.IpSecAdvancedSettings{
			KeyExchange: keyExchange,
			IkeLifeTime: ikeLifeTime,
			Lifetime:    lifetime,
			DpdTimeout:  dpdTimeout,
			DpdDelay:    dpdDelay,
			Phase1: &perimeter81Sdk.IpSecPhase{
				Auth:       authPhase1,
				Encryption: encryptionPhase1,
				Dh:         dhPhase1,
			},
			Phase2: &perimeter81Sdk.IpSecPhase{
				Auth:       authPhase2,
				Encryption: encryptionPhase2,
				Dh:         dhPhase2,
			},
		},
	}

	// create the ipsec-redundant tunnel using the client sdk and check for errors
	status, _, err := client.IPSecRedundantApi.CreateIPSecRedundantTunnel(ctx, ipSecRedundantBody, networkId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create ipsec-redundant tunnel", err)
	}

	// get the status id of the ipsec-redundant tunnel creation
	var ipSecRedundantTunnelId string
	statusId := getIdFromUrl(status.StatusUrl)

	// check the status of the ipsec-redundant tunnel creation
	for {
		// check the status of the network that contains the ipsec-redundant tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the network status is completed, get the ipsec-redundant tunnel id and break the loop
		if networkStatus.Completed {
			baseTunnelBody := perimeter81Sdk.BaseTunnelValues{
				RegionID:   regionId,
				GatewayID:  gatewayId1,
				TunnelName: tunnelName,
			}
			ipSecRedundantTunnelId, diags = getRedundantTunnelId(ctx, networkId, baseTunnelBody, *client, diags)
			if ipSecRedundantTunnelId == "" {
				return diags
			}
			break
		}
		// delay for 20 seconds before checking the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId(ipSecRedundantTunnelId)

	return resourceIpsecRedundantRead(ctx, d, m)
}

/*
resourceIpsecRedundantRead Read a Ipsec Redundant Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecRedundantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)

	if ctx == nil {
		ctx = context.Background()
	}

	// get the ipsec-redundant tunnel id and the network id from the terraform resource data
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

	// get the ipsec-redundant tunnel using the client sdk and check for errors
	tunnel, _, err := client.IPSecRedundantApi.GetIPSecRedundantTunnel(ctx, networkId, tunnelId)

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read ipsec-redundant tunnel", err)
	}
	if err := d.Set("network_id", networkId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set networkId", err)
	}
	if err := d.Set("region_id", tunnel.RegionID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set regionId", err)
	}
	if err := d.Set("tunnel_name", tunnel.TunnelName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set tunnel name", err)
	}
	if err := d.Set("advanced_settings", flattenAdvancedSettingsData(tunnel.AdvancedSettings)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set advanced settings", err)
	}
	if err := d.Set("shared_settings", flattenSharedSettingsData(tunnel.SharedSettings)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set shared settings", err)
	}
	if err := d.Set("tunnel1", flattenTunnelData(tunnel.Tunnel1)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set tunnel1", err)
	}
	if err := d.Set("tunnel2", flattenTunnelData(tunnel.Tunnel2)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set tunnel2", err)
	}
	d.SetId(tunnelId)
	return diags
}

/*
resourceIpsecRedundantUpdate Update a Ipsec Redundant Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecRedundantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return appendErrorDiags(diags, "Unable to delete ipsec-redundant tunnel", fmt.Errorf("ipsec-redundant tunnel update is not available yet"))
}

/*
resourceIpsecRedundantDelete Delete a Ipsec Redundant Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceIpsecRedundantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)

	if ctx == nil {
		ctx = context.Background()
	}

	// get the ipsec-redundant tunnel id and the network id from the terraform resource data
	tunnelId := d.Id()
	networkId := d.Get("network_id").(string)

	// delete the ipsec-redundant tunnel using the client sdk and check for errors
	status, _, err := client.IPSecRedundantApi.DeleteIPSecRedundantTunnel(ctx, networkId, tunnelId)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete ipsec-redundant tunnel", err)
	}

	// get the status id of the ipsec-redundant tunnel deletion
	statusId := getIdFromUrl(status.StatusUrl)
	for {
		// check the status of the network that contains the ipsec-redundant tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the network status is completed, break the loop
		if networkStatus.Completed {
			break
		}
		// delay for 20 seconds before checking the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId("")
	return diags
}
