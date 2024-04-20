package perimeter81

import (
	"context"
	"fmt"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceObjectAddresses Setup the Object Addresses Resource CRUD operations

@return &schema.Resource
*/
func resourceObjectAddresses() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectAddressesCreate,
		ReadContext:   resourceObjectAddressesRead,
		UpdateContext: resourceObjectAddressesUpdate,
		DeleteContext: resourceObjectAddressesDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
	ipVersion := d.Get("ip_version").(string)
	value := flattenStringsArrayData(d.Get("value").([]interface{}))

	objectAddressesPayload := perimeter81Sdk.ObjectsAddressObj{
		Name:        name,
		Description: description,
		ValueType:   valueType,
		IpVersion:   ipVersion,
		Value:       value,
	}
	// create the Object Addresses and check for errors
	objectAddresses, _, err := client.ObjectsAddressesApi.PostObjectsAddresses(ctx, objectAddressesPayload)

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Object Addresses", err)
	}

	d.SetId(objectAddresses.Id)
	return resourceObjectServicesRead(ctx, d, m)
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
	objectsAddresses, _, err := client.ObjectsAddressesApi.GetObjectsAddresses(ctx)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find object addresses", err)
	}
	currentObjectAddresses := getCurrentObjectAddressesInArray(&objectsAddresses, d.Id())

	if err := d.Set("name", currentObjectAddresses.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses name", err)
	}
	if err := d.Set("description", currentObjectAddresses.Description); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses description", err)
	}
	if err := d.Set("value_type", currentObjectAddresses.ValueType); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses value_type", err)
	}
	if err := d.Set("ip_version", currentObjectAddresses.IpVersion); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object addresses ip_version", err)
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
		ipVersion := d.Get("ip_version").(string)
		value := flattenStringsArrayData(d.Get("value").([]interface{}))

		// prepare the object addresses data for the api service
		updateObjectAddressesPayload := perimeter81Sdk.ObjectsAddressObj{
			Name:        name,
			Description: description,
			ValueType:   valueType,
			IpVersion:   ipVersion,
			Value:       value,
		}
		//update the object addresses and check for errors
		_, _, err := client.ObjectsAddressesApi.PutObjectsAddresses(ctx, updateObjectAddressesPayload, objectAddressesId)
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
	_, err := client.ObjectsAddressesApi.DeleteObjectsAddresses(ctx, d.Id())

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete object addresses", err)
	}

	d.SetId("")
	return nil
}
