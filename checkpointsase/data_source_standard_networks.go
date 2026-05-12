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
dataSourceStandardNetworks Query all Standard Networks via the /v2.3/networks/standard endpoint.
This is an explicit alias that makes it clear the data source returns standard networks only.

@return &schema.Resource
*/
func dataSourceStandardNetworks() *schema.Resource {
	return &schema.Resource{
		Description: "List all standard networks in Check Point SASE. Uses the /v2.3/networks/standard endpoint.",
		ReadContext: dataSourceStandardNetworksRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of standard networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the network.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tags associated with the network.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"dns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DNS name of the network.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet CIDR block of the network.",
						},
						"accesstype": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access type of the network.",
						},
						"isdefault": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the default network.",
						},
						"tenantid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tenant ID that owns this network.",
						},
						"createdat": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation timestamp.",
						},
						"updatedat": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last update timestamp.",
						},
					},
				},
			},
		},
	}
}

func dataSourceStandardNetworksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networks, _, err := client.StandardNetworksAPI.StandardGetNetworks(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get standard networks", err)
	}

	flatNetworks := make([]interface{}, len(networks))
	for i, n := range networks {
		network := make(map[string]interface{})
		network["id"] = n.GetId()
		network["name"] = n.GetName()
		network["tags"] = n.GetTags()
		network["dns"] = n.GetDns()
		network["subnet"] = n.GetSubnet()
		network["accesstype"] = n.GetAccessType()
		network["isdefault"] = n.GetIsDefault()
		network["tenantid"] = n.GetTenantId()
		network["createdat"] = n.GetCreatedAt().String()
		network["updatedat"] = n.GetUpdatedAt().String()
		flatNetworks[i] = network
	}

	if err := d.Set("networks", flatNetworks); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
