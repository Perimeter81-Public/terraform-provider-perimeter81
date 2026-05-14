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
dataSourceRegions Query all Regions

@return &schema.Resource
*/
func dataSourceRegions() *schema.Resource {
	return &schema.Resource{
		Description: "List the cloud regions available for deploying " +
			"`checkpointsase_network` (standard) resources. Each entry provides " +
			"the `id` to pass as `cpregion_id` when declaring a network's region. " +
			"For enhanced networks use `checkpointsase_enhanced_regions` instead.",
		ReadContext: dataSourceRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of available Check Point SASE cloud regions for standard networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server-side object identifier (internal).",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region ID. Used as `cpregion_id` on `checkpointsase_network.region`.",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ISO 3166-1 alpha-2 country code for the region.",
						},
						"continent_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ISO 3166-1 alpha-2 continent code for the region.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable region name.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal name of the region.",
						},
						"class_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server-side class identifier (internal).",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceRegionsRead Use the SDK to query all Regions
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// call the api and check if there is an error
	regionsData, _, err := client.RegionsAPI.StandardNetworksControllerV2GetRegions(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Regions data", err)
	}

	// flatten the data so it fit the terraform schema and set the terraform resource data
	if err := d.Set("regions", flattenRegions(regionsData)); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
