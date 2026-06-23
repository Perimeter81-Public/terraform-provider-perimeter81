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
dataSourceRouteTable Query the standard network route table for a given network

@return &schema.Resource
*/
func dataSourceRouteTable() *schema.Resource {
	return &schema.Resource{
		Description: "Read the route table of a single standard `checkpointsase_network`. " +
			"Returns one entry per route with `subnets`, interface name, and propagation flag. " +
			"Use `network_id` from `checkpointsase_standard_networks` to look up the right network.",
		ReadContext: dataSourceRouteTableRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the standard network to fetch the route table for.",
			},
			"routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of route table entries for the given standard network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route entry ID.",
						},
						"interface_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The route table interface name.",
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
dataSourceRouteTableRead Use the SDK to query the standard network route table.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	routes, _, err := client.RouteTableAPI.StandardGetRouteTable(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Route Table", err)
	}

	routesData := flattenRouteTableData(routes)
	if err := d.Set("routes", routesData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Route Table data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenRouteTableData flattens a list of GetRouteTable200ResponseInner SDK models to a Terraform-compatible list.
*/
func flattenRouteTableData(routes []perimeter81Sdk.GetRouteTable200ResponseInner) []interface{} {
	if routes == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(routes))
	for i, route := range routes {
		routeMap := map[string]interface{}{
			"id":             route.GetId(),
			"interface_name": route.GetInterfaceName(),
			"subnets":        route.Subnets,
			"propagated":     route.GetPropagated(),
		}
		result[i] = routeMap
	}
	return result
}
