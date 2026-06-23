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
dataSourceEnhancedRegions Query all available Check Point SASE regions that support enhanced networks

@return &schema.Resource
*/
func dataSourceEnhancedRegions() *schema.Resource {
	return &schema.Resource{
		Description: "List the cloud regions available for deploying " +
			"`checkpointsase_enhanced_network` resources. Each entry provides the " +
			"`id` to pass as `harmony_sase_region_id` when declaring an enhanced " +
			"network's region.",
		ReadContext: dataSourceEnhancedRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of available Check Point SASE regions for enhanced networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the Check Point SASE region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Check Point SASE region.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the Check Point SASE region.",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO 3166-1 alpha-2 country code for the region.",
						},
						"continent_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO 3166-1 alpha-2 continent code for the region.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceEnhancedRegionsRead Use the SDK to query all available Check Point SASE regions for enhanced networks.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceEnhancedRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	regions, _, err := client.EnhancedRegionsAPI.EnhancedNetworksControllerV2GetRegions(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Enhanced Regions", err)
	}

	regionsData := flattenHarmonySaseRegions(regions)
	if err := d.Set("regions", regionsData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Regions data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenHarmonySaseRegions flattens a list of HarmonySaseRegion SDK models to a Terraform-compatible list.
*/
func flattenHarmonySaseRegions(regions []perimeter81Sdk.HarmonySaseRegion) []interface{} {
	if regions == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(regions))
	for i, region := range regions {
		regionMap := map[string]interface{}{
			"id":             region.Id,
			"name":           region.Name,
			"display_name":   region.DisplayName,
			"country_code":   region.CountryCode,
			"continent_code": region.ContinentCode,
		}
		result[i] = regionMap
	}
	return result
}
