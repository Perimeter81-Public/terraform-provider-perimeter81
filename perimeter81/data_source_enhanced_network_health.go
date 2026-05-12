package perimeter81

import (
	"context"
	"strconv"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
dataSourceEnhancedNetworkHealth Query the health status of an enhanced network

@return &schema.Resource
*/
func dataSourceEnhancedNetworkHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnhancedNetworkHealthRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the enhanced network to retrieve health status for.",
			},
			"health_checks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of health check results for tunnels in the enhanced network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of health check component (always 'tunnel' for enhanced networks).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health status of the tunnel ('passing', 'critical', or 'unknown').",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the enhanced network this health check belongs to.",
						},
						"tunnel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the tunnel being checked.",
						},
						"tunnel_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the tunnel being checked.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the region where the tunnel is deployed.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceEnhancedNetworkHealthRead Use the SDK to query the health of an enhanced network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceEnhancedNetworkHealthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	healthResponse, _, err := client.EnhancedNetworksAPI.GetEnhancedNetworkHealth(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Enhanced Network health status", err)
	}

	healthChecks := flattenEnhancedHealthChecks(healthResponse.GetData())
	if err := d.Set("health_checks", healthChecks); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Network health data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenEnhancedHealthChecks flattens a list of EnhancedHealthCheck SDK models to a Terraform-compatible list.
  - @param checks []perimeter81Sdk.EnhancedHealthCheck - the health checks to flatten

@return []interface{} - the flattened health checks
*/
func flattenEnhancedHealthChecks(checks []perimeter81Sdk.EnhancedHealthCheck) []interface{} {
	if checks == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(checks))
	for i, check := range checks {
		meta := check.GetMeta()
		checkMap := map[string]interface{}{
			"type":        string(check.GetType()),
			"status":      string(check.GetStatus()),
			"network_id":  meta.GetNetworkId(),
			"tunnel_name": meta.GetTunnelName(),
			"tunnel_id":   meta.GetTunnelId(),
			"region_id":   meta.GetRegionId(),
		}
		result[i] = checkMap
	}
	return result
}
