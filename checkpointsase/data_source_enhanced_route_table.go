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
dataSourceEnhancedRouteTable Query the enhanced network route table for a given network

@return &schema.Resource
*/
func dataSourceEnhancedRouteTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnhancedRouteTableRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the enhanced network to fetch the route table for.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of route table entries for the given enhanced network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route entry ID.",
						},
						"tunnel_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tunnel IDs associated with this route entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of subnet CIDR blocks for this route entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"propagated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the route is propagated automatically.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceEnhancedRouteTableRead Use the SDK to query the enhanced network route table.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceEnhancedRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	response, _, err := client.EnhancedRouteTablesAPI.GetEnhancedRouteTable(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Enhanced Route Table", err)
	}

	routesData := flattenEnhancedRouteTableData(response.GetData())
	if err := d.Set("routes", routesData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Route Table data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenEnhancedRouteTableData flattens a list of EnhancedRouteTable SDK models to a Terraform-compatible list.
*/
func flattenEnhancedRouteTableData(routes []perimeter81Sdk.EnhancedRouteTable) []interface{} {
	if routes == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		routeMap := map[string]interface{}{
			"id":         route.Id,
			"tunnel_ids": route.TunnelIds,
			"subnets":    route.Subnets,
			"propagated": route.Propagated,
		}
		result[i] = routeMap
	}
	return result
}
