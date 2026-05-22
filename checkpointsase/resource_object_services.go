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
resourceObjectServices Setup the Object Services Resource CRUD operations

@return &schema.Resource
*/
func resourceObjectServices() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a service object in Check Point SASE's shared object library. " +
			"Service objects are reusable references to one or more transport-layer " +
			"protocol + port combinations; they're typically referenced from firewall " +
			"policy rules. Use `checkpointsase_object_addresses` for the parallel " +
			"address-object resource. " +
			"**ICMP**: the v2.3 API also supports `protocol = \"icmp\"`, but this " +
			"provider does not yet expose the corresponding `protocolOptions` payload, " +
			"so only `tcp` and `udp` are usable here.",
		CreateContext: resourceObjectServicesCreate,
		ReadContext:   resourceObjectServicesRead,
		UpdateContext: resourceObjectServicesUpdate,
		DeleteContext: resourceObjectServicesDelete,
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
				Description:  "Display name of the service object. Must be 3–100 characters.",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description of the service object.",
			},
			"protocols": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of protocol+port combinations covered by this service object. At least one entry is required.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Transport protocol. Must be `tcp` or `udp`.",
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
						},
						"value_type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Shape of the `value` list. Must be `single` (one port), `range` (exactly two ports, low–high), or `list` (multiple discrete ports).",
							ValidateFunc: validation.StringInSlice([]string{"single", "range", "list"}, false),
						},
						"value": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Port numbers. Shape depends on `value_type`: 1 element for `single`, 2 elements (start, end) for `range`, 1+ for `list`. Each value must be a valid port (1–65535).",
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IsPortNumber,
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

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	protocolsPayload := flattenProtocolsData(d.Get("protocols").([]interface{}))

	createObjectsServicesPayload := perimeter81Sdk.ObjectsServicesRequestObj{
		Name:        name,
		Description: &description,
		Protocols:   protocolsPayload,
	}
	newObjectServices, _, err := client.ObjectsServicesAPI.PostObjectsServices(ctx).ObjectsServicesRequestObj(createObjectsServicesPayload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Object Services", err)
	}

	d.SetId(newObjectServices.GetId())
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
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	// BUG-21 fix preserved: look up by id, not by name. On terraform import
	// only d.Id() is seeded, so a by-name lookup would panic; a by-name
	// lookup also misbehaves if the service is renamed server-side. The
	// SDK list now returns the flat protocols structure (BUG-17 SDK fix in
	// P81-123406) so a single GET + id match is all we need.
	objectsServices, _, err := client.ObjectsServicesAPI.GetObjectsServices(ctx).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to fetch object services", err)
	}
	var match *perimeter81Sdk.ObjectsServicesResponseObj
	for i := range objectsServices.Data {
		if objectsServices.Data[i].Id != nil && *objectsServices.Data[i].Id == d.Id() {
			match = &objectsServices.Data[i]
			break
		}
	}
	if match == nil {
		// Resource gone server-side — let terraform schedule a recreate.
		d.SetId("")
		return diags
	}

	if err := d.Set("name", match.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services name", err)
	}
	if err := d.Set("description", match.Description); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services description", err)
	}
	if err := d.Set("protocols", flattenObjectServicesProtocols(match.Protocols)); err != nil {
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
		objectServicesId := d.Id()
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		protocolsPayload := flattenProtocolsData(d.Get("protocols").([]interface{}))
		updateObjectServicesPayload := perimeter81Sdk.ObjectsServicesRequestObj{
			Name:        name,
			Description: &description,
			Protocols:   protocolsPayload,
		}
		if _, _, err := client.ObjectsServicesAPI.PutObjectsServices(ctx, objectServicesId).ObjectsServicesRequestObj(updateObjectServicesPayload).Execute(); err != nil {
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
	_, err := client.ObjectsServicesAPI.DeleteObjectsServices(ctx, d.Id()).Execute()

	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to delete ipsec-single tunnel", err)
	}

	d.SetId("")
	return nil
}
