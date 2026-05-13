package checkpointsase

import (
	"context"
	"fmt"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
resourceObjectAddresses Setup the Object Addresses Resource CRUD operations

@return &schema.Resource
*/
func resourceObjectAddresses() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an address object in Check Point SASE's shared object library. " +
			"Address objects are reusable references to a single IP, a list of IPs, a CIDR " +
			"block, or an FQDN; they're typically referenced from firewall policy rules " +
			"and service definitions. Use `checkpointsase_object_services` for the parallel " +
			"service-object resource.",
		CreateContext: resourceObjectAddressesCreate,
		ReadContext:   resourceObjectAddressesRead,
		UpdateContext: resourceObjectAddressesUpdate,
		DeleteContext: resourceObjectAddressesDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Timestamp of the last update to this resource.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Display name of the address object. Must be 3–100 characters.",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the address object.",
			},
			"value_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Category of the `value` list. Must be `ip` (single IP), `list` (multiple IPs), `cidr` (single CIDR block), or `fqdn` (single domain name).",
				ValidateFunc: validation.StringInSlice([]string{"ip", "list", "cidr", "fqdn"}, false),
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "Has no effect on the v2.3 server. The Public API hardcodes `ipv4` server-side and strips this field from both request and response. Will be removed in a future major release.",
				Description: "IP version (e.g. `ipv4`). Not transmitted to or returned by the v2.3 server — values you set here are silently discarded.",
			},
			"value": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Address values. Shape depends on `value_type`: exactly 1 element for `ip` / `cidr` / `fqdn`, 1+ elements for `list`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceObjectAddressesImportState,
		},
	}
}

/*
resourceOpenvpnImportState Import gateways
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectAddressesImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceObjectAddressesRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import object Addresses: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceObjectAddressesCreate Create a Object Addresses
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectAddressesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the object services data from the terraform resource data and flatten what need to be flattened for the api
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	valueType := d.Get("value_type").(string)
	value := flattenStringsArrayData(d.Get("value").([]interface{}))

	objectAddressesPayload := perimeter81Sdk.ObjectsAddressObj{
		Name:        name,
		Description: &description,
		ValueType:   valueType,
		Value:       value,
	}
	// create the Object Addresses and check for errors
	objectAddresses, _, err := client.ObjectsAddressesAPI.PostObjectsAddresses(ctx).ObjectsAddressObj(objectAddressesPayload).Execute()

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Object Addresses", err)
	}

	d.SetId(objectAddresses.GetId())
	return resourceObjectAddressesRead(ctx, d, m)
}

/*
resourceObjectAddressesRead Read a Object Addresses
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectAddressesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the object addresses and check for errors
	objectsAddresses, _, err := client.ObjectsAddressesAPI.GetObjectsAddresses(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find object addresses", err)
	}
	currentObjectAddresses := getCurrentObjectAddressesInArray(objectsAddresses, d.Id())

	if err := d.Set("name", currentObjectAddresses.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses name", err)
	}
	if err := d.Set("description", currentObjectAddresses.GetDescription()); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses description", err)
	}
	if err := d.Set("value_type", currentObjectAddresses.ValueType); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses value_type", err)
	}
	if err := d.Set("value", currentObjectAddresses.Value); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses value", err)
	}

	return diags
}

/*
resourceObjectAddressesUpdate Update a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectAddressesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	if d.HasChanges("value", "description", "name", "ip_version", "value_type") {

		// get the object addresses data from the terraform resource data and flatten what need to be flattened for the api
		objectAddressesId := d.Id()
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		valueType := d.Get("value_type").(string)
		value := flattenStringsArrayData(d.Get("value").([]interface{}))

		// prepare the object addresses data for the api service
		updateObjectAddressesPayload := perimeter81Sdk.ObjectsAddressObj{
			Name:        name,
			Description: &description,
			ValueType:   valueType,
			Value:       value,
		}
		//update the object addresses and check for errors
		_, _, err := client.ObjectsAddressesAPI.PutObjectsAddresses(ctx, objectAddressesId).ObjectsAddressObj(updateObjectAddressesPayload).Execute()
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update object addresses", err)
		}
	}

	return resourceObjectAddressesRead(ctx, d, m)
}

/*
resourceObjectAddressesDelete Delete a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectAddressesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// delete the object Addresses and check for errors
	_, err := client.ObjectsAddressesAPI.DeleteObjectsAddresses(ctx, d.Id()).Execute()

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete object addresses", err)
	}

	d.SetId("")
	return nil
}
