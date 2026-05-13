package checkpointsase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
resourceIpsecRedundant Setup the IpSec-Redundant Resource CRUD operations

@return &schema.Resource
*/
func resourceIpsecRedundant() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an active/standby IPsec redundant tunnel pair for a " +
			"`checkpointsase_network`. Two tunnels (`tunnel1` + `tunnel2`) terminate at " +
			"distinct remote endpoints for failover; `shared_settings` (gateway subnets) " +
			"and `advanced_settings` (IKE/IPSec parameters, phase1/phase2 proposals) " +
			"apply to both tunnels uniformly. " +
			"**This resource has no in-place update path** — every attribute change " +
			"forces full replacement (destroy + recreate). Updating in place will be " +
			"supported in a future version.",
		CreateContext: resourceIpsecRedundantCreate,
		ReadContext:   resourceIpsecRedundantRead,
		UpdateContext: resourceIpsecRedundantUpdate,
		DeleteContext: resourceIpsecRedundantDelete,
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
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the standard network the tunnel pair belongs to.",
			},
			"tunnel_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Display name for the redundant tunnel pair.",
			},
			"advanced_settings": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "IKE/IPSec parameters and phase1/phase2 proposals shared by both tunnels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_exchange": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "IKE version for key exchange. Must be `ikev1` or `ikev2`.",
							ValidateFunc: validation.StringInSlice([]string{"ikev1", "ikev2"}, false),
						},
						"ike_life_time": {
							Type:     schema.TypeString,
							Required: true,
							Description: "IKE lifetime as a `<int><unit>` duration string, e.g. `28800s`, `480m`, or `8h`. " +
								"Server-enforced ranges: `s` 10–86400, `m` 1–1440, `h` 1–24.",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+[smh]$`),
								"must be a duration with unit `s`, `m`, or `h` (e.g. `28800s`, `480m`, `8h`)"),
						},
						"lifetime": {
							Type:     schema.TypeString,
							Required: true,
							Description: "IPSec SA lifetime as a `<int><unit>` duration string, e.g. `3600s`, `60m`, or `1h`. " +
								"Server-enforced ranges: `s` 10–86400, `m` 1–1440, `h` 1–24.",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d+[smh]$`),
								"must be a duration with unit `s`, `m`, or `h` (e.g. `3600s`, `60m`, `1h`)"),
						},
						"dpd_delay": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dead peer detection delay interval, formatted `<int>s`. Allowed range is `5s`–`60s`.",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([5-9]|[1-5]\d|60)s$`),
								"must be a duration like `5s`–`60s`"),
						},
						"dpd_timeout": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dead peer detection timeout, formatted `<int>s`. Allowed range is `5s`–`60s`.",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([5-9]|[1-5]\d|60)s$`),
								"must be a duration like `5s`–`60s`"),
						},
						"phase1": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Phase 1 (IKE) IPSec proposal lists.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 1 authentication algorithms (e.g. `[\"sha256\"]`).",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"encryption": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 1 encryption algorithms (e.g. `[\"aes-cbc-256\"]`).",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dh": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 1 Diffie-Hellman group numbers (e.g. `[14]` for MODP2048).",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								}},
						},
						"phase2": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Phase 2 (ESP/IPSec) proposal lists.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 2 authentication algorithms.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"encryption": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 2 encryption algorithms.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dh": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of phase 2 Diffie-Hellman group numbers.",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								}},
						},
					}},
			},
			"shared_settings": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet routing settings shared by both tunnels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"p81_gateway_subnets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Check Point SASE gateway subnet CIDR blocks reachable through either tunnel.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remote_gateway_subnets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Remote-side subnet CIDR blocks reachable through either tunnel.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					}},
			},
			"tunnel1": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Primary tunnel endpoint configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passphrase": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Required:    true,
							Description: "Pre-shared key for this tunnel (8–64 characters).",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the SASE gateway that terminates this tunnel locally.",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Optional remote tunnel ID. Computed if not supplied.",
						},
						"tunnel_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The server-assigned tunnel ID. Computed.",
						},
						"p81_gwinternal_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Check Point SASE gateway internal IP on this tunnel.",
						},
						"remote_gwinternal_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote gateway internal IP on this tunnel.",
						},
						"remote_public_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote gateway public IP on this tunnel.",
						},
						"remote_asn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote peer's BGP ASN as a string (e.g. `\"65010\"`).",
						},
					}},
			},
			"tunnel2": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Standby tunnel endpoint configuration. Same shape as `tunnel1`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passphrase": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Required:    true,
							Description: "Pre-shared key for this tunnel (8–64 characters).",
						},
						"tunnel_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The server-assigned tunnel ID. Computed.",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the SASE gateway that terminates this tunnel locally.",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Optional remote tunnel ID. Computed if not supplied.",
						},
						"p81_gwinternal_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Check Point SASE gateway internal IP on this tunnel.",
						},
						"remote_gwinternal_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote gateway internal IP on this tunnel.",
						},
						"remote_public_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote gateway public IP on this tunnel.",
						},
						"remote_asn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The remote peer's BGP ASN as a string (e.g. `\"65010\"`).",
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
	ctx = context.Background()

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
	remoteId1 := tunnel1Data["remote_id"].(string)
	gatewayId2 := tunnel2Data["gateway_id"].(string)
	passphrase2 := tunnel2Data["passphrase"].(string)
	p81GWinternalIP2 := tunnel2Data["p81_gwinternal_ip"].(string)
	remoteGWinernalIP2 := tunnel2Data["remote_gwinternal_ip"].(string)
	remotePublicIP2 := tunnel2Data["remote_public_ip"].(string)
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
	remoteId1Value := perimeter81Sdk.StringAsRemoteID(&remoteId1)
	remoteId2Value := perimeter81Sdk.StringAsRemoteID(&remoteId2)
	ipSecRedundantBody := perimeter81Sdk.CreateIPSecRedundantPayload{
		RegionID:   regionId,
		TunnelName: tunnelName,
		Tunnel1: perimeter81Sdk.IPSecRedundantTunnelPayload{
			Passphrase:         passphrase1,
			GatewayID:          gatewayId1,
			P81GWInternalIP:    p81GWinternalIP1,
			RemoteGWInternalIP: remoteGWinernalIP1,
			RemotePublicIP:     remotePublicIP1,
			RemoteASN:          perimeter81Sdk.RemoteASN{},
			RemoteID:           remoteId1Value,
		},
		Tunnel2: perimeter81Sdk.IPSecRedundantTunnelPayload{
			Passphrase:         passphrase2,
			GatewayID:          gatewayId2,
			P81GWInternalIP:    p81GWinternalIP2,
			RemoteGWInternalIP: remoteGWinernalIP2,
			RemotePublicIP:     remotePublicIP2,
			RemoteASN:          perimeter81Sdk.RemoteASN{},
			RemoteID:           remoteId2Value,
		},
		SharedSettings: perimeter81Sdk.IPSecSharedSettingsCreate{
			P81GatewaySubnets:    p81GatewaySubnets,
			RemoteGatewaySubnets: remoteGatewaySubnets,
			P81ASN:               perimeter81Sdk.RemoteASN{},
		},
		AdvancedSettings: perimeter81Sdk.IPSecAdvancedSettings{
			KeyExchange: keyExchange,
			IkeLifeTime: ikeLifeTime,
			Lifetime:    lifetime,
			DpdTimeout:  dpdTimeout,
			DpdDelay:    dpdDelay,
			Phase1: perimeter81Sdk.IPSecPhaseConfig{
				Auth:       authPhase1,
				Encryption: encryptionPhase1,
				Dh:         dhPhase1,
			},
			Phase2: perimeter81Sdk.IPSecPhaseConfig{
				Auth:       authPhase2,
				Encryption: encryptionPhase2,
				Dh:         dhPhase2,
			},
		},
	}
	// create the ipsec-redundant tunnel using the client sdk and check for errors
	status, _, err := client.IPSecRedundantAPI.StandardCreateIPSecRedundantTunnel(ctx, networkId).CreateIPSecRedundantPayload(ipSecRedundantBody).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create ipsec-redundant tunnel", err)
	}

	// get the status id of the ipsec-redundant tunnel creation
	var ipSecRedundantTunnelId string
	statusId := getIdFromUrl(status.GetStatusUrl())

	// check the status of the ipsec-redundant tunnel creation
	for {
		// check the status of the network that contains the ipsec-redundant tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the network status is completed, get the ipsec-redundant tunnel id and break the loop
		if networkStatus.GetCompleted() {
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
	ctx = context.Background()

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
	tunnel, _, err := client.IPSecRedundantAPI.StandardGetIPSecRedundantTunnel(ctx, networkId, tunnelId).Execute()

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read ipsec-redundant tunnel", err)
	}
	if err := d.Set("network_id", networkId); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set networkId", err)
	}
	if err := d.Set("region_id", tunnel.GetRegionID()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set regionId", err)
	}
	if err := d.Set("tunnel_name", tunnel.GetTunnelName()); err != nil {
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
	if len(ids) != 1 {
		if err := d.Set("tunnel1", flattenTunnelData(tunnel.Tunnel1)); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set tunnel1", err)
		}
		if err := d.Set("tunnel2", flattenTunnelData(tunnel.Tunnel2)); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set tunnel2", err)
		}
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
	ctx = context.Background()

	// get the ipsec-redundant tunnel id and the network id from the terraform resource data
	tunnelId := d.Id()
	networkId := d.Get("network_id").(string)

	// delete the ipsec-redundant tunnel using the client sdk and check for errors
	status, _, err := client.IPSecRedundantAPI.StandardDeleteIPSecRedundantTunnel(ctx, networkId, tunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete ipsec-redundant tunnel", err)
	}

	// get the status id of the ipsec-redundant tunnel deletion
	statusId := getIdFromUrl(status.GetStatusUrl())
	for {
		// check the status of the network that contains the ipsec-redundant tunnel and check for errors
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		// if the network status is completed, break the loop
		if networkStatus.GetCompleted() {
			break
		}
		// delay for 20 seconds before checking the status again
		time.Sleep(20 * time.Second)
	}
	d.SetId("")
	return diags
}
