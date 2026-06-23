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
dataSourceCustomerCertificates Query all customer certificates for enhanced networks

@return &schema.Resource
*/
func dataSourceCustomerCertificates() *schema.Resource {
	return &schema.Resource{
		Description: "List customer-uploaded TLS certificates available to enhanced networks. " +
			"Used when configuring `auth_type = \"cert\"` on enhanced static or dynamic " +
			"tunnels — the certificate IDs returned here can be referenced from the " +
			"tunnel's `customer_root_ca` attribute. " +
			"Returns an empty list (no error) when the tenant has no enhanced networks.",
		ReadContext: dataSourceCustomerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of customer certificates associated with enhanced networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the customer certificate.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the customer certificate.",
						},
						"expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date of the customer certificate in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

/*
dataSourceCustomerCertificatesRead Use the SDK to query all customer certificates.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func dataSourceCustomerCertificatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	certificates, resp, err := client.EnhancedNetworksAPI.EnhancedNetworksControllerV23GetNetworkCustomerCertificate(ctx).Execute()
	if err != nil {
		// 409 Conflict means the tenant doesn't have enhanced networks or the feature is not enabled.
		// Return empty list instead of failing.
		if resp != nil && resp.StatusCode == 409 {
			certificates = nil
		} else {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to get Customer Certificates", err)
		}
	}

	certData := flattenCustomerCertificates(certificates)
	if err := d.Set("certificates", certData); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Customer Certificates data", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

/*
flattenCustomerCertificates flattens a list of NetworkCustomerCertificateResponseInner SDK models to a Terraform-compatible list.
  - @param certificates []perimeter81Sdk.NetworkCustomerCertificateResponseInner - the certificates to flatten

@return []interface{} - the flattened certificates
*/
func flattenCustomerCertificates(certificates []perimeter81Sdk.NetworkCustomerCertificateResponseInner) []interface{} {
	if certificates == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(certificates))
	for i, cert := range certificates {
		certMap := map[string]interface{}{
			"id":           cert.GetId(),
			"display_name": cert.GetDisplayName(),
			"expires_at":   cert.GetExpiresAt().Format(time.RFC3339),
		}
		result[i] = certMap
	}
	return result
}
