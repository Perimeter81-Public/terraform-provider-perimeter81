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
dataSourceApplications Query all Applications

@return &schema.Resource
*/
func dataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationsRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of applications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the application.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the application.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the application (e.g., 'http', 'https', 'rdp').",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the network associated with this application.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host address of the application.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the application is currently enabled.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation timestamp of the application.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceApplicationsRead Use the SDK to query all Applications.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceApplicationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	response, _, err := client.ApplicationAPI.GetApplications(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Applications", err)
	}

	applicationsData := flattenApplicationsData(response.GetData())
	if err := d.Set("applications", applicationsData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Applications data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenApplicationsData flattens a list of ApplicationsListObject SDK models to a Terraform-compatible list.
*/
func flattenApplicationsData(apps []perimeter81Sdk.ApplicationsListObject) []interface{} {
	if apps == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(apps))
	for i, app := range apps {
		appMap := map[string]interface{}{
			"id":         app.Id,
			"name":       app.Name,
			"type":       app.Type,
			"host":       app.Host,
			"network_id": app.Network.Id,
			"enabled":    app.GetEnabled(),
			"created_at": app.GetCreatedAt(),
		}
		result[i] = appMap
	}
	return result
}
