package checkpointsase

import (
	"context"
	"strconv"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
dataSourceEnhancedTunnels Query all tunnels in an enhanced network

@return &schema.Resource
*/
func dataSourceEnhancedTunnels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnhancedTunnelsRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the enhanced network to retrieve tunnels for.",
			},
			"items_total": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The total number of tunnels in the enhanced network.",
			},
			"page": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The current page number of the paginated result.",
			},
			"total_page": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The total number of pages in the paginated result.",
			},
			"tunnels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of tunnels in the enhanced network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the enhanced tunnel.",
						},
						"tunnel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the tunnel.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the target region for this tunnel.",
						},
						"ha_tunnel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enhanced dynamic tunnel group ID (or tunnel ID for static tunnels).",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The authentication type for the tunnel ('psk' for pre-shared key, 'cert' for certificate).",
						},
						"key_exchange": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IKE version used for key exchange.",
						},
						"ike_life_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IKE phase 1 lifetime.",
						},
						"lifetime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IKE phase 2 lifetime.",
						},
						"dpd_delay": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Dead Peer Detection delay interval.",
						},
						"dpd_timeout": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Dead Peer Detection timeout.",
						},
						"dpd_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The action taken when Dead Peer Detection triggers.",
						},
						"remote_public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remote gateway's public IP address.",
						},
						"remote_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remote gateway ID.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An optional description for the tunnel.",
						},
						"routing_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The routing type for the tunnel.",
						},
						"peak_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The expected peak throughput of the tunnel communication in Mbps.",
						},
						"p81_gateway_subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of Check Point SASE gateway subnets.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"remote_gateway_subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of remote gateway subnets.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

/*
dataSourceEnhancedTunnelsRead Use the SDK to query all tunnels in an enhanced network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceEnhancedTunnelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	response, _, err := client.EnhancedTunnelsAPI.GetEnhancedRegionTunnelsPerNetwork(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Enhanced Tunnels", err)
	}

	if err := d.Set("items_total", float64(response.GetItemsTotal())); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Tunnels items_total", err)
	}
	if err := d.Set("page", float64(response.GetPage())); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Tunnels page", err)
	}
	if err := d.Set("total_page", float64(response.GetTotalPage())); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Tunnels total_page", err)
	}

	tunnels := flattenEnhancedTunnelsData(response.GetData())
	if err := d.Set("tunnels", tunnels); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Tunnels data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenEnhancedTunnelsData flattens a list of EnhancedTunnel SDK models to a Terraform-compatible list.
  - @param tunnels []perimeter81Sdk.EnhancedTunnel - the tunnels to flatten

@return []interface{} - the flattened tunnels
*/
func flattenEnhancedTunnelsData(tunnels []perimeter81Sdk.EnhancedTunnel) []interface{} {
	if tunnels == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(tunnels))
	for i, tunnel := range tunnels {
		tunnelMap := map[string]interface{}{
			"id":                     tunnel.GetId(),
			"tunnel_name":            tunnel.GetTunnelName(),
			"region_id":              tunnel.GetRegionID(),
			"ha_tunnel_id":           tunnel.GetHaTunnelID(),
			"auth_type":              tunnel.GetAuthType(),
			"key_exchange":           tunnel.GetKeyExchange(),
			"ike_life_time":          tunnel.GetIkeLifeTime(),
			"lifetime":               tunnel.GetLifetime(),
			"dpd_delay":              tunnel.GetDpdDelay(),
			"dpd_timeout":            tunnel.GetDpdTimeout(),
			"dpd_action":             tunnel.GetDpdAction(),
			"remote_public_ip":       tunnel.GetRemotePublicIP(),
			"remote_id":              tunnel.GetRemoteID(),
			"description":            tunnel.GetDescription(),
			"routing_type":           string(tunnel.GetRoutingType()),
			"peak_bandwidth":         int(tunnel.GetPeakBandwidth()),
			"p81_gateway_subnets":    tunnel.GetP81GatewaySubnets(),
			"remote_gateway_subnets": tunnel.GetRemoteGatewaySubnets(),
		}
		result[i] = tunnelMap
	}
	return result
}
