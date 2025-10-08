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
dataSourceObjectAddresses Query all ObjectAddresses

@return &schema.Resource
*/
func dataSourceObjectAddresses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectAddressesRead,
		Schema: map[string]*schema.Schema{
			"object_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
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
	objectAddresses, _, err := client.ObjectsAddressesApi.GetObjectsAddresses(ctx)
	// a, _ := json.Marshal(objectAddresses)
	// return append(diags, diag.Diagnostic{
	// 	Severity: diag.Error,
	// 	Summary:  "error",
	// 	Detail:   string(a),
	// })
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
