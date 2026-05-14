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
dataSourceObjectServices Query all ObjectServices

@return &schema.Resource
*/
func dataSourceObjectServices() *schema.Resource {
	return &schema.Resource{
		Description: "List all service objects in Check Point SASE's shared object library. " +
			"Use `checkpointsase_object_services` (the resource) to manage individual entries.",
		ReadContext: dataSourceObjectServicesRead,
		Schema: map[string]*schema.Schema{
			"object_services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of service objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the service object.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the service object.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the service object.",
						},
						"protocols": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of protocol+port combinations covered by this service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Transport protocol (`tcp`, `udp`, or `icmp`).",
									},
									"value_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Shape of `value`: `single`, `range`, or `list`.",
									},
									"value": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Port numbers. Shape depends on `value_type`.",
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
	objectServices, _, err := client.ObjectsServicesAPI.GetObjectsServices(ctx).Execute()
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
