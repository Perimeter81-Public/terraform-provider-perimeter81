package perimeter81

import (
	"context"
	"fmt"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceObjectServices Setup the Object Services Resource CRUD operations

@return &schema.Resource
*/
func resourceObjectServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectServicesCreate,
		ReadContext:   resourceObjectServicesRead,
		UpdateContext: resourceObjectServicesUpdate,
		DeleteContext: resourceObjectServicesDelete,
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
		Importer: &schema.ResourceImporter{
			StateContext: resourceObjectServicesImportState,
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
func resourceObjectServicesImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceObjectServicesRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import object services: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceObjectServicesCreate Create a Object Services
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectServicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the object services data from the terraform resource data and flatten what need to be flattened for the api
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	protocols := flattenProtocolsData(d.Get("protocols").([]interface{}))
	protocolsPayload := make([]perimeter81Sdk.ObjectsServicesProtocolRequestObj, len(protocols))

	for index, protocol := range protocols {
		protocolsPayload[index].ObjectServiceProtocolTcpudp = protocol
	}

	CreateObjectsServicesPayload := perimeter81Sdk.ObjectsServicesRequestObj{
		Name:        name,
		Description: description,
		Protocols:   protocolsPayload,
	}
	// return diags
	// create the Object Services and check for errors
	newObjectServices, _, err := client.ObjectsServicesApi.PostObjectsServices(ctx, CreateObjectsServicesPayload)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Object Services", err)
	}

	d.SetId(newObjectServices.Id)
	return resourceObjectServicesRead(ctx, d, m)
}

/*
resourceObjectServicesRead Read a Object Services
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectServicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// get the object services and check for errors
	objectsServices, _, err := client.ObjectsServicesApi.GetObjectsServices(ctx)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Network", err)
	}
	currentObjectServices := getCurrentObjectServicesInArray(&objectsServices, d.Id())

	if err := d.Set("name", currentObjectServices.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services name", err)
	}
	if err := d.Set("description", currentObjectServices.Description); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services description", err)
	}
	if err := d.Set("protocols", flattenObjectServicesProtocols(currentObjectServices.Protocols)); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services protocol data", err)
	}

	return diags
}

/*
resourceObjectServicesUpdate Update a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectServicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	if d.HasChanges("name", "description", "protocols") {

		// get the object services data from the terraform resource data and flatten what need to be flattened for the api
		objectServicesId := d.Id()
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		protocols := flattenProtocolsData(d.Get("protocols").([]interface{}))

		// prepare the object services data for the api service
		protocolsPayload := make([]perimeter81Sdk.ObjectsServicesProtocolRequestObj, len(protocols))
		for index, protocol := range protocols {
			protocolsPayload[index].ObjectServiceProtocolTcpudp = protocol
		}
		updateObjectServicesPayload := perimeter81Sdk.ObjectsServicesRequestObj{
			Id:          objectServicesId,
			Name:        name,
			Description: description,
			Protocols:   protocolsPayload,
		}
		//update the object services and check for errors
		_, _, err := client.ObjectsServicesApi.PutObjectsServices(ctx, updateObjectServicesPayload, objectServicesId)
		if err != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to update object services", err)
		}
	}

	return resourceObjectServicesRead(ctx, d, m)
}

/*
resourceObjectServicesDelete Delete a Ipsec single Tunnel
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceObjectServicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// delete the object services and check for errors
	_, err := client.ObjectsServicesApi.DeleteObjectsServices(ctx, d.Id())

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete ipsec-single tunnel", err)
	}

	d.SetId("")
	return nil
}
