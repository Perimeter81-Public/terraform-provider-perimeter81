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
dataSourceObjectServices Query all ObjectServices

@return &schema.Resource
*/
func dataSourceObjectServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObjectServicesRead,
		Schema: map[string]*schema.Schema{
			"object_services": {
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
						"protocols": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								}},
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

func dataSourceObjectServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}

	// call the api and check if there is an error
	objectServices, _, err := client.ObjectsServicesApi.GetObjectsServices(ctx)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get object services", err)
	}

	// flatten the data so it fit the terraform schema and set the terraform resource data
	newObjectServices := flattenObjectServicesData(objectServices.Data)
	if err := d.Set("object_services", newObjectServices); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
