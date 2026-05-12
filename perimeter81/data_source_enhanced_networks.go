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
dataSourceEnhancedNetworks Query all Enhanced Networks

@return &schema.Resource
*/
func dataSourceEnhancedNetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnhancedNetworksRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of enhanced networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the enhanced network.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the enhanced network.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet CIDR block of the enhanced network.",
						},
						"dns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DNS of the enhanced network.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access type of the enhanced network.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of tags associated with the enhanced network.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the default enhanced network.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tenant ID that owns this enhanced network.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation timestamp of the enhanced network.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceEnhancedNetworksRead Use the SDK to query all Enhanced Networks.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceEnhancedNetworksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networks, _, err := client.EnhancedNetworksAPI.GetEnhancedNetworks(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Enhanced Networks", err)
	}

	networksData := flattenEnhancedNetworksData(networks)
	if err := d.Set("networks", networksData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Enhanced Networks data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenEnhancedNetworksData flattens a list of EnhancedNetwork SDK models to a Terraform-compatible list.
*/
func flattenEnhancedNetworksData(networks []perimeter81Sdk.EnhancedNetwork) []interface{} {
	if networks == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(networks))
	for i, network := range networks {
		networkMap := map[string]interface{}{
			"id":          network.Id,
			"name":        network.Name,
			"subnet":      network.Subnet,
			"dns":         network.Dns,
			"access_type": network.AccessType,
			"tags":        network.Tags,
			"is_default":  network.IsDefault,
			"tenant_id":   network.TenantId,
			"created_at":  network.CreatedAt.String(),
		}
		result[i] = networkMap
	}
	return result
}
