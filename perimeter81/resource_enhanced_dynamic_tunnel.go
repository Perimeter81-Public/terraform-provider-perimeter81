package perimeter81

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceEnhancedDynamicTunnel Setup the Enhanced Dynamic Tunnel Resource CRUD operations

@return &schema.Resource
*/
func resourceEnhancedDynamicTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnhancedDynamicTunnelCreate,
		ReadContext:   resourceEnhancedDynamicTunnelRead,
		UpdateContext: resourceEnhancedDynamicTunnelUpdate,
		DeleteContext: resourceEnhancedDynamicTunnelDelete,
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
				Description: "The ID of the enhanced network this dynamic tunnel belongs to.",
			},
			"tunnel_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the dynamic IPSec tunnel.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description for the dynamic tunnel.",
			},
			"tunnel": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of individual tunnel endpoints for this dynamic tunnel group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The enhanced region ID for this tunnel endpoint.",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication type for this tunnel endpoint ('psk' or 'cert').",
						},
						"passphrase": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Pre-shared key for tunnel authentication (8-64 characters). Required when auth_type is 'psk'.",
						},
						"customer_root_ca": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customer root certificate authority. Required when auth_type is 'cert'.",
						},
						"remote_public_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The remote gateway public IP address.",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The remote gateway ID.",
						},
						"p81_gw_internal_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Check Point SASE gateway internal IP address.",
						},
						"remote_gw_internal_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The remote gateway internal IP address.",
						},
					},
				},
			},
			"p81_gateway_subnets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of Check Point SASE gateway subnet CIDR blocks (shared settings).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"remote_gateway_subnets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of remote gateway subnet CIDR blocks (shared settings).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"peak_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Expected peak throughput of the tunnel communication in Mbps. Defaults to 1000.",
			},
			"key_exchange": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IKE version for key exchange (e.g., 'ikev2').",
			},
			"ike_life_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IKE lifetime value (e.g., '28800s', '480m', '8h').",
			},
			"lifetime": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPSec SA lifetime value (e.g., '3600s', '60m', '1h').",
			},
			"dpd_delay": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Dead peer detection delay interval (e.g., '30s').",
			},
			"dpd_timeout": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Dead peer detection timeout value (e.g., '60s').",
			},
			"phase1": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Phase 1 (IKE) IPSec configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 1 authentication algorithms.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"encryption": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 1 encryption algorithms.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"key_exchange_method": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 1 key exchange methods (Diffie-Hellman groups).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"phase2": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Phase 2 (ESP/IPSec) configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 2 authentication algorithms.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"encryption": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 2 encryption algorithms.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"key_exchange_method": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of phase 2 key exchange methods (Diffie-Hellman groups).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEnhancedDynamicTunnelImportState,
		},
	}
}

/*
resourceEnhancedDynamicTunnelImportState Import an enhanced dynamic tunnel by its ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceEnhancedDynamicTunnelImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceEnhancedDynamicTunnelRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import enhanced dynamic tunnel: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
flattenDynamicTunnelDetails converts a list of tunnel schema blocks into []DynamicTunnelDetails SDK models.
*/
func flattenDynamicTunnelDetails(tunnelItems []interface{}) []perimeter81Sdk.DynamicTunnelDetails {
	tunnels := make([]perimeter81Sdk.DynamicTunnelDetails, len(tunnelItems))
	for i, item := range tunnelItems {
		tunnelMap := item.(map[string]interface{})
		regionId := tunnelMap["region_id"].(string)
		detail := perimeter81Sdk.DynamicTunnelDetails{
			RegionID: regionId,
		}
		if v, ok := tunnelMap["auth_type"].(string); ok && v != "" {
			detail.AuthType = &v
		}
		if v, ok := tunnelMap["passphrase"].(string); ok && v != "" {
			detail.Passphrase = &v
		}
		if v, ok := tunnelMap["customer_root_ca"].(string); ok && v != "" {
			detail.CustomerRootCA = &v
		}
		if v, ok := tunnelMap["remote_public_ip"].(string); ok && v != "" {
			detail.RemotePublicIP = &v
		}
		if v, ok := tunnelMap["remote_id"].(string); ok && v != "" {
			detail.RemoteID = &v
		}
		if v, ok := tunnelMap["p81_gw_internal_ip"].(string); ok && v != "" {
			detail.P81GWInternalIP = &v
		}
		if v, ok := tunnelMap["remote_gw_internal_ip"].(string); ok && v != "" {
			detail.RemoteGWInternalIP = &v
		}
		tunnels[i] = detail
	}
	return tunnels
}

/*
resourceEnhancedDynamicTunnelCreate Create an Enhanced Dynamic IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedDynamicTunnelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	tunnelName := d.Get("tunnel_name").(string)
	p81GatewaySubnets := flattenStringsArrayData(d.Get("p81_gateway_subnets").([]interface{}))
	remoteGatewaySubnets := flattenStringsArrayData(d.Get("remote_gateway_subnets").([]interface{}))
	peakBandwidth := int32(d.Get("peak_bandwidth").(int))
	keyExchange := d.Get("key_exchange").(string)
	ikeLifeTime := d.Get("ike_life_time").(string)
	lifetime := d.Get("lifetime").(string)
	dpdDelay := d.Get("dpd_delay").(string)
	dpdTimeout := d.Get("dpd_timeout").(string)
	phase1 := flattenIPSecPhaseConfigV23(d.Get("phase1").([]interface{}))
	phase2 := flattenIPSecPhaseConfigV23(d.Get("phase2").([]interface{}))
	tunnels := flattenDynamicTunnelDetails(d.Get("tunnel").([]interface{}))

	sharedSettings := perimeter81Sdk.EnhancedIPSecSharedSettingsCreate{
		P81GatewaySubnets:    p81GatewaySubnets,
		RemoteGatewaySubnets: remoteGatewaySubnets,
		PeakBandwidth:        &peakBandwidth,
		Features:             perimeter81Sdk.NetworkFeaturesCreate{},
		LeftASN:              perimeter81Sdk.RemoteASN{},
	}

	advancedSettings := perimeter81Sdk.IPSecAdvancedSettingsV23{
		KeyExchange: keyExchange,
		IkeLifeTime: ikeLifeTime,
		Lifetime:    lifetime,
		DpdDelay:    dpdDelay,
		DpdTimeout:  dpdTimeout,
		Phase1:      phase1,
		Phase2:      phase2,
	}

	payload := perimeter81Sdk.DynamicTunnelCreate{
		TunnelName:       tunnelName,
		Tunnels:          tunnels,
		SharedSettings:   sharedSettings,
		AdvancedSettings: advancedSettings,
	}

	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		payload.Description = &s
	}

	status, _, err := client.EnhancedTunnelsAPI.CreateDynamicTunnel(ctx, networkId).DynamicTunnelCreate(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Enhanced Dynamic Tunnel", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var dynamicTunnelId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			dynamicTunnelId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(dynamicTunnelId)
	return resourceEnhancedDynamicTunnelRead(ctx, d, m)
}

/*
resourceEnhancedDynamicTunnelRead Read an Enhanced Dynamic IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedDynamicTunnelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	dynamicTunnelId := d.Id()

	tunnelsData, _, err := client.EnhancedTunnelsAPI.GetDynamicTunnel(ctx, networkId, dynamicTunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Enhanced Dynamic Tunnel", err)
	}

	if len(tunnelsData) > 0 {
		tunnel := tunnelsData[0]
		if err := d.Set("tunnel_name", tunnel.TunnelName); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel tunnel_name", err)
		}
		if err := d.Set("key_exchange", tunnel.KeyExchange); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel key_exchange", err)
		}
		if err := d.Set("ike_life_time", tunnel.IkeLifeTime); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel ike_life_time", err)
		}
		if err := d.Set("lifetime", tunnel.Lifetime); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel lifetime", err)
		}
		if err := d.Set("dpd_delay", tunnel.DpdDelay); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel dpd_delay", err)
		}
		if err := d.Set("dpd_timeout", tunnel.DpdTimeout); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel dpd_timeout", err)
		}
		if err := d.Set("phase1", flattenIPSecPhaseConfigV23ToMap(tunnel.Phase1)); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel phase1", err)
		}
		if err := d.Set("phase2", flattenIPSecPhaseConfigV23ToMap(tunnel.Phase2)); err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to set Enhanced Dynamic Tunnel phase2", err)
		}
	}

	return diags
}

/*
resourceEnhancedDynamicTunnelUpdate Update an Enhanced Dynamic IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedDynamicTunnelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	dynamicTunnelId := d.Id()
	tunnelName := d.Get("tunnel_name").(string)

	payload := perimeter81Sdk.DynamicTunnelUpdate{
		TunnelName: tunnelName,
	}

	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		payload.Description = &s
	}

	_, _, err := client.EnhancedTunnelsAPI.UpdateDynamicTunnel(ctx, networkId, dynamicTunnelId).DynamicTunnelUpdate(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to update Enhanced Dynamic Tunnel", err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceEnhancedDynamicTunnelRead(ctx, d, m)
}

/*
resourceEnhancedDynamicTunnelDelete Delete an Enhanced Dynamic IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedDynamicTunnelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	dynamicTunnelId := d.Id()

	status, _, err := client.EnhancedTunnelsAPI.DeleteDynamicTunnel(ctx, networkId, dynamicTunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete Enhanced Dynamic Tunnel", err)
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
