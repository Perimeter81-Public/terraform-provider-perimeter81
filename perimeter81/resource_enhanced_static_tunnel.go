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
resourceEnhancedStaticTunnel Setup the Enhanced Static Tunnel Resource CRUD operations

@return &schema.Resource
*/
func resourceEnhancedStaticTunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnhancedStaticTunnelCreate,
		ReadContext:   resourceEnhancedStaticTunnelRead,
		UpdateContext: resourceEnhancedStaticTunnelUpdate,
		DeleteContext: resourceEnhancedStaticTunnelDelete,
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
				Description: "The ID of the enhanced network this static tunnel belongs to.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target region ID within the enhanced network.",
			},
			"tunnel_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the static IPSec tunnel.",
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
			"auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authentication type for the tunnel. Use 'psk' for pre-shared key or 'cert' for certificate.",
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description for the static tunnel.",
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
			"p81_gateway_subnets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of Check Point SASE gateway subnet CIDR blocks.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remote_gateway_subnets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of remote gateway subnet CIDR blocks.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			StateContext: resourceEnhancedStaticTunnelImportState,
		},
	}
}

/*
resourceEnhancedStaticTunnelImportState Import an enhanced static tunnel by its ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceEnhancedStaticTunnelImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceEnhancedStaticTunnelRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import enhanced static tunnel: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
flattenIPSecPhaseConfigV23 converts a phase schema block into an IPSecPhaseConfigV23 SDK model.
*/
func flattenIPSecPhaseConfigV23(phaseList []interface{}) perimeter81Sdk.IPSecPhaseConfigV23 {
	if len(phaseList) == 0 {
		return perimeter81Sdk.IPSecPhaseConfigV23{}
	}
	phaseMap := phaseList[0].(map[string]interface{})
	return perimeter81Sdk.IPSecPhaseConfigV23{
		Auth:              flattenStringsArrayData(phaseMap["auth"].([]interface{})),
		Encryption:        flattenStringsArrayData(phaseMap["encryption"].([]interface{})),
		KeyExchangeMethod: flattenStringsArrayData(phaseMap["key_exchange_method"].([]interface{})),
	}
}

/*
flattenIPSecPhaseConfigV23ToMap converts an IPSecPhaseConfigV23 SDK model to a Terraform-compatible map.
*/
func flattenIPSecPhaseConfigV23ToMap(phase perimeter81Sdk.IPSecPhaseConfigV23) []interface{} {
	phaseMap := map[string]interface{}{
		"auth":                phase.Auth,
		"encryption":          phase.Encryption,
		"key_exchange_method": phase.KeyExchangeMethod,
	}
	return []interface{}{phaseMap}
}

/*
resourceEnhancedStaticTunnelCreate Create an Enhanced Static IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedStaticTunnelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	regionId := d.Get("region_id").(string)
	tunnelName := d.Get("tunnel_name").(string)
	keyExchange := d.Get("key_exchange").(string)
	ikeLifeTime := d.Get("ike_life_time").(string)
	lifetime := d.Get("lifetime").(string)
	dpdDelay := d.Get("dpd_delay").(string)
	dpdTimeout := d.Get("dpd_timeout").(string)
	p81GatewaySubnets := flattenStringsArrayData(d.Get("p81_gateway_subnets").([]interface{}))
	remoteGatewaySubnets := flattenStringsArrayData(d.Get("remote_gateway_subnets").([]interface{}))
	phase1 := flattenIPSecPhaseConfigV23(d.Get("phase1").([]interface{}))
	phase2 := flattenIPSecPhaseConfigV23(d.Get("phase2").([]interface{}))

	payload := perimeter81Sdk.StaticTunnelCreate{
		RegionID:             regionId,
		TunnelName:           tunnelName,
		KeyExchange:          keyExchange,
		IkeLifeTime:          ikeLifeTime,
		Lifetime:             lifetime,
		DpdDelay:             dpdDelay,
		DpdTimeout:           dpdTimeout,
		P81GatewaySubnets:    p81GatewaySubnets,
		RemoteGatewaySubnets: remoteGatewaySubnets,
		Phase1:               phase1,
		Phase2:               phase2,
	}

	if v, ok := d.GetOk("remote_public_ip"); ok {
		s := v.(string)
		payload.RemotePublicIP = &s
	}
	if v, ok := d.GetOk("remote_id"); ok {
		s := v.(string)
		payload.RemoteID = &s
	}
	if v, ok := d.GetOk("auth_type"); ok {
		s := v.(string)
		payload.AuthType = &s
	}
	if v, ok := d.GetOk("passphrase"); ok {
		s := v.(string)
		payload.Passphrase = &s
	}
	if v, ok := d.GetOk("customer_root_ca"); ok {
		s := v.(string)
		payload.CustomerRootCA = &s
	}
	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		payload.Description = &s
	}
	peakBandwidth := int32(d.Get("peak_bandwidth").(int))
	payload.PeakBandwidth = &peakBandwidth

	status, _, err := client.EnhancedTunnelsAPI.CreateStaticTunnel(ctx, networkId).StaticTunnelCreate(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Enhanced Static Tunnel", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var tunnelId string
	for {
		var networkStatus perimeter81Sdk.AsyncOperationStatus
		networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			d.Partial(true)
			return diags
		}
		if networkStatus.GetCompleted() {
			tunnelId = getIdFromUrl(networkStatus.Result.GetResource())
			break
		}
		time.Sleep(60 * time.Second)
	}

	d.SetId(tunnelId)
	return resourceEnhancedStaticTunnelRead(ctx, d, m)
}

/*
resourceEnhancedStaticTunnelRead Read an Enhanced Static IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedStaticTunnelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	tunnelId := d.Id()

	tunnelData, _, err := client.EnhancedTunnelsAPI.GetStaticTunnel(ctx, networkId, tunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Enhanced Static Tunnel", err)
	}

	if err := d.Set("tunnel_name", tunnelData.TunnelName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel tunnel_name", err)
	}
	if err := d.Set("key_exchange", tunnelData.KeyExchange); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel key_exchange", err)
	}
	if err := d.Set("ike_life_time", tunnelData.IkeLifeTime); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel ike_life_time", err)
	}
	if err := d.Set("lifetime", tunnelData.Lifetime); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel lifetime", err)
	}
	if err := d.Set("dpd_delay", tunnelData.DpdDelay); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel dpd_delay", err)
	}
	if err := d.Set("dpd_timeout", tunnelData.DpdTimeout); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel dpd_timeout", err)
	}
	if err := d.Set("remote_public_ip", tunnelData.RemotePublicIP); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel remote_public_ip", err)
	}
	if err := d.Set("remote_id", tunnelData.RemoteID); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel remote_id", err)
	}
	if err := d.Set("auth_type", tunnelData.AuthType); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel auth_type", err)
	}
	if err := d.Set("p81_gateway_subnets", tunnelData.P81GatewaySubnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel p81_gateway_subnets", err)
	}
	if err := d.Set("remote_gateway_subnets", tunnelData.RemoteGatewaySubnets); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel remote_gateway_subnets", err)
	}
	if err := d.Set("phase1", flattenIPSecPhaseConfigV23ToMap(tunnelData.Phase1)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel phase1", err)
	}
	if err := d.Set("phase2", flattenIPSecPhaseConfigV23ToMap(tunnelData.Phase2)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Static Tunnel phase2", err)
	}

	return diags
}

/*
resourceEnhancedStaticTunnelUpdate Update an Enhanced Static IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedStaticTunnelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	tunnelId := d.Id()

	payload := perimeter81Sdk.StaticTunnelUpdate{}

	if v, ok := d.GetOk("tunnel_name"); ok {
		s := v.(string)
		payload.TunnelName = &s
	}
	if v, ok := d.GetOk("remote_public_ip"); ok {
		s := v.(string)
		payload.RemotePublicIP = &s
	}
	if v, ok := d.GetOk("remote_id"); ok {
		s := v.(string)
		payload.RemoteID = &s
	}
	if v, ok := d.GetOk("auth_type"); ok {
		s := v.(string)
		payload.AuthType = &s
	}
	if v, ok := d.GetOk("passphrase"); ok {
		s := v.(string)
		payload.Passphrase = &s
	}
	if v, ok := d.GetOk("customer_root_ca"); ok {
		s := v.(string)
		payload.CustomerRootCA = &s
	}
	if v, ok := d.GetOk("description"); ok {
		s := v.(string)
		payload.Description = &s
	}
	if v, ok := d.GetOk("key_exchange"); ok {
		s := v.(string)
		payload.KeyExchange = &s
	}
	if v, ok := d.GetOk("ike_life_time"); ok {
		s := v.(string)
		payload.IkeLifeTime = &s
	}
	if v, ok := d.GetOk("lifetime"); ok {
		s := v.(string)
		payload.Lifetime = &s
	}
	if v, ok := d.GetOk("dpd_delay"); ok {
		s := v.(string)
		payload.DpdDelay = &s
	}
	if v, ok := d.GetOk("dpd_timeout"); ok {
		s := v.(string)
		payload.DpdTimeout = &s
	}
	if v, ok := d.GetOk("p81_gateway_subnets"); ok {
		payload.P81GatewaySubnets = flattenStringsArrayData(v.([]interface{}))
	}
	if v, ok := d.GetOk("remote_gateway_subnets"); ok {
		payload.RemoteGatewaySubnets = flattenStringsArrayData(v.([]interface{}))
	}
	if v := d.Get("phase1").([]interface{}); len(v) > 0 {
		phase1 := flattenIPSecPhaseConfigV23(v)
		payload.Phase1 = &phase1
	}
	if v := d.Get("phase2").([]interface{}); len(v) > 0 {
		phase2 := flattenIPSecPhaseConfigV23(v)
		payload.Phase2 = &phase2
	}
	peakBandwidth := int32(d.Get("peak_bandwidth").(int))
	payload.PeakBandwidth = &peakBandwidth

	_, _, err := client.EnhancedTunnelsAPI.UpdateStaticTunnel(ctx, networkId, tunnelId).StaticTunnelUpdate(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to update Enhanced Static Tunnel", err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceEnhancedStaticTunnelRead(ctx, d, m)
}

/*
resourceEnhancedStaticTunnelDelete Delete an Enhanced Static IPSec Tunnel.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceEnhancedStaticTunnelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	tunnelId := d.Id()

	status, _, err := client.EnhancedTunnelsAPI.DeleteStaticTunnel(ctx, networkId, tunnelId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete Enhanced Static Tunnel", err)
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
