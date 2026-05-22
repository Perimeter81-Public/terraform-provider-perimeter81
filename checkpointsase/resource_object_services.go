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

	// BUG-17 workaround: bypass the SDK's broken protocols serialization
	// (it never sets the `protocol` field, so the server saves the service
	// with empty protocols). Send the flat JSON shape directly.
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	rawProtocols := hclProtocolsToRaw(d.Get("protocols").([]interface{}))

	id, err := createRawObjectService(ctx, client, name, description, rawProtocols)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Object Services", err)
	}
	if id == "" {
		// Server returned an unexpected (or 204-style) shape — fall back to
		// listing and finding by name. Best-effort.
		rawList, listErr := fetchRawObjectServices(ctx, client)
		if listErr != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Created Object Service but couldn't recover its id", listErr)
		}
		for _, svc := range rawList {
			if svc.Name == name {
				id = svc.Id
				break
			}
		}
	}
	if id == "" {
		d.Partial(true)
		return appendErrorDiags(diags, "Created Object Service but server response had no id", fmt.Errorf("empty id"))
	}

	d.SetId(id)
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

	// BUG-21 fix (and BUG-17 protocols workaround in the same path): look up
	// the service by id from the raw GET. The previous implementation looked
	// up by `d.Get("name")`, which (1) panicked on `terraform import` because
	// the import handler only seeds `d.Id()`, and (2) silently misbehaved if
	// the service was renamed server-side. Looking up by id also lets us
	// drop the broken SDK call (which couldn't deserialize protocols) — the
	// raw fetch already returns id/name/description/protocols.
	rawServices, rawErr := fetchRawObjectServices(ctx, client)
	if rawErr != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to fetch object services", rawErr)
	}
	var rawMatch *rawObjectService
	for i := range rawServices {
		if rawServices[i].Id == d.Id() {
			rawMatch = &rawServices[i]
			break
		}
	}
	if rawMatch == nil {
		// Resource gone server-side — let terraform schedule a recreate.
		d.SetId("")
		return diags
	}

	if err := d.Set("name", rawMatch.Name); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services name", err)
	}
	if err := d.Set("description", rawMatch.Description); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set object services description", err)
	}
	if err := d.Set("protocols", rawProtocolsToTerraform(rawMatch.Protocols)); err != nil {
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
		// BUG-17 workaround: same flat-shape PUT as the CREATE workaround.
		objectServicesId := d.Id()
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		rawProtocols := hclProtocolsToRaw(d.Get("protocols").([]interface{}))

		if err := updateRawObjectService(ctx, client, objectServicesId, name, description, rawProtocols); err != nil {
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
