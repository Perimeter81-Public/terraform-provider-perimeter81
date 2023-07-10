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
dataSourceRegions Query all Regions

@return &schema.Resource
*/
func dataSourceRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"continent_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"class_name": {
							Type:     schema.TypeString,
							Computed: true,
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
	if ctx == nil {
		ctx = context.Background()
	}

	// call the api and check if there is an error
	regionsData, _, err := client.RegionsApi.NetworksControllerV2GetRegions(ctx)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Regions data", err)
	}

	// flatten the data so it fit the terraform schema and set the terraform resource data
	if err := d.Set("regions", flattenRegions(regionsData.Regions)); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
