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
dataSourceStatus Query the Check Point SASE API status

@return &schema.Resource
*/
func dataSourceStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStatusRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the Check Point SASE API.",
			},
		},
	}
}

/*
dataSourceStatusRead Use the SDK to query the API status.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceStatusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	status, _, err := client.NetworksAPI.GetStatus(ctx).Execute()
	if err != nil {
		// The /v2.3/status endpoint returns plain text "Ok" which the SDK
		// may fail to JSON-decode. If the body contains "Ok", treat it as success.
		if apiErr, ok := err.(*perimeter81Sdk.GenericOpenAPIError); ok {
			body := string(apiErr.Body())
			if body == "Ok" || body == "\"Ok\"" {
				status = "Ok"
			} else {
				d.Partial(true)
				return appendErrorDiags(diags, "Unable to get API status", err)
			}
		} else {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to get API status", err)
		}
	}

	if err := d.Set("status", status); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set API status", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
