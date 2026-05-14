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
dataSourceObjectAddresses Query all ObjectAddresses

@return &schema.Resource
*/
func dataSourceObjectAddresses() *schema.Resource {
	return &schema.Resource{
		Description: "List all address objects in Check Point SASE's shared object library. " +
			"Use `checkpointsase_object_addresses` (the resource) to manage individual entries.",
		ReadContext: dataSourceObjectAddressesRead,
		Schema: map[string]*schema.Schema{
			"object_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of address objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the address object.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the address object.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the address object.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Category of the `value` list. One of `ip`, `list`, `cidr`, `fqdn`.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "Not transmitted to or returned by the v2.3 server (hardcoded to `ipv4` and stripped from responses). Field retained for backward compatibility.",
							Description: "IP version. Has no effect on the v2.3 server.",
						},
						"value": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Address values. Shape depends on `value_type`.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

/*
dataSourceObjectServicesRead Use the SDK to query all ObjectServices
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client
@return diag.Diagnostics
*/

func dataSourceObjectAddressesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}

	// call the api and check if there is an error
	objectAddresses, _, err := client.ObjectsAddressesAPI.GetObjectsAddresses(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get object addresses", err)
	}

	// flatten the data so it fit the terraform schema and set the terraform resource data
	newObjectServices := flattenObjectAddressesData(objectAddresses.Data)
	if err := d.Set("object_addresses", newObjectServices); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
