package perimeter81

import (
	"context"
	"strconv"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
dataSourceNetworkHealth Query the health status of a standard network

@return &schema.Resource
*/
func dataSourceNetworkHealth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkHealthRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the standard network to retrieve health status for.",
			},
			"health_checks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of health check results for gateways and tunnels in the standard network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of health check component ('gateway' or 'tunnel').",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health status of the component ('passing', 'critical', or 'unknown').",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the network this health check belongs to.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance (gateway or tunnel) being checked.",
						},
						"tunnel_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the tunnel (only populated for tunnel type health checks).",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceNetworkHealthRead Use the SDK to query the health of a standard network.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceNetworkHealthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	healthResponse, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2GetNetworkHealth(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Network health status", err)
	}

	healthChecks := flattenStandardHealthChecks(healthResponse.GetData())
	if err := d.Set("health_checks", healthChecks); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Network health data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenStandardHealthChecks flattens a list of StandardHealthCheck SDK models to a Terraform-compatible list.
  - @param checks []perimeter81Sdk.StandardHealthCheck - the health checks to flatten

@return []interface{} - the flattened health checks
*/
func flattenStandardHealthChecks(checks []perimeter81Sdk.StandardHealthCheck) []interface{} {
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
			"instance_id": meta.GetInstanceId(),
			"tunnel_name": meta.GetTunnelName(),
		}
		result[i] = checkMap
	}
	return result
}
