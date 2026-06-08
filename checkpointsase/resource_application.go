package checkpointsase

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
resourceApplication Setup the Application Resource CRUD operations.
Note: there is no update or delete endpoint — all fields are ForceNew.

@return &schema.Resource
*/
func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an Application in Check Point SASE. " +
			"**All attributes are immutable**: any change to a field on this resource " +
			"forces full replacement (destroy + re-create), not in-place update. " +
			"**`terraform destroy` only removes the resource from state.** " +
			"The Harmony SASE v2.3 API does not expose a delete endpoint for " +
			"applications, so the application continues to exist on the server. " +
			"Delete it manually via the Infinity Portal if needed.",
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		DeleteContext: resourceApplicationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The application name.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The application type. The v2.3 API supports creating " +
					"applications of these types only: `http`, `https`, `rdp`. " +
					"(Existing `ssh` and `vnc` applications can be read but not " +
					"created through the Public API.)",
				ValidateFunc: validation.StringInSlice([]string{"http", "https", "rdp"}, false),
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ID to associate with this application.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The application host address.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "The application port number (1–65535).",
				ValidateFunc: validation.IsPortNumber,
			},
			"users": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of user IDs allowed to access this application.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of group IDs allowed to access this application.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImportState,
		},
	}
}

/*
resourceApplicationImportState Import an application by its ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceApplicationImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	diagnostics := resourceApplicationRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import application: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
buildApplicationHostNullable builds a NullableOneOfFixedHostIdpHost from a host string value.
*/
func buildApplicationHostNullable(host string) perimeter81Sdk.NullableOneOfFixedHostIdpHost {
	hostValue := perimeter81Sdk.StringAsFixedHostValue(&host)
	fixedHost := perimeter81Sdk.FixedHost{
		Source: "fixed",
		Value:  hostValue,
	}
	var hostInterface interface{} = fixedHost
	return *perimeter81Sdk.NewNullableOneOfFixedHostIdpHost(&hostInterface)
}

/*
buildApplicationPortNullable builds a NullableOneOfFixedPortIdpPort from a port int32 value.
*/
func buildApplicationPortNullable(port int32) perimeter81Sdk.NullableOneOfFixedPortIdpPort {
	fixedPort := perimeter81Sdk.FixedPort{
		Source: "fixed",
		Value:  port,
	}
	var portInterface interface{} = fixedPort
	return *perimeter81Sdk.NewNullableOneOfFixedPortIdpPort(&portInterface)
}

/*
resourceApplicationCreate Create an Application.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	appName := d.Get("name").(string)
	appType := d.Get("type").(string)
	networkId := d.Get("network").(string)
	host := d.Get("host").(string)
	port := int32(d.Get("port").(int))
	users := flattenStringsArrayData(d.Get("users").([]interface{}))
	groups := flattenStringsArrayData(d.Get("groups").([]interface{}))

	hostNullable := buildApplicationHostNullable(host)
	portNullable := buildApplicationPortNullable(port)

	var payload perimeter81Sdk.CreateApplicationRequest

	switch appType {
	case "http":
		httpApp := perimeter81Sdk.HttpCreateApplication{
			Name:       appName,
			Type:       appType,
			Network:    networkId,
			Host:       hostNullable,
			Port:       portNullable,
			Users:      users,
			Groups:     groups,
			Headers:    map[string]interface{}{},
			Attributes: perimeter81Sdk.HttpAttributes{},
		}
		payload = perimeter81Sdk.CreateApplicationRequest{
			HttpCreateApplication: &httpApp,
		}
	case "https":
		httpsApp := perimeter81Sdk.HttpsCreateApplication{
			Name:       appName,
			Type:       appType,
			Network:    networkId,
			Host:       hostNullable,
			Port:       portNullable,
			Users:      users,
			Groups:     groups,
			Headers:    map[string]interface{}{},
			Attributes: perimeter81Sdk.HttpsAttributes{},
		}
		payload = perimeter81Sdk.CreateApplicationRequest{
			HttpsCreateApplication: &httpsApp,
		}
	case "rdp":
		rdpApp := perimeter81Sdk.RdpCreateApplication{
			Name:       appName,
			Type:       appType,
			Network:    networkId,
			Host:       hostNullable,
			Port:       portNullable,
			Users:      users,
			Groups:     groups,
			Attributes: perimeter81Sdk.RdpAttributes{},
			Auth:       perimeter81Sdk.ApplicationAuth{AuthEnabled: false},
		}
		payload = perimeter81Sdk.CreateApplicationRequest{
			RdpCreateApplication: &rdpApp,
		}
	default:
		return appendErrorDiags(diags, "Unsupported application type", fmt.Errorf("type must be 'http', 'https', or 'rdp', got: %s", appType))
	}

	status, _, err := client.ApplicationAPI.CreateApplication(ctx).CreateApplicationRequest(payload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to create Application", err)
	}

	statusId := getIdFromUrl(status.GetStatusUrl())
	var applicationId string
	for {
		appStatus, _, statusErr := client.ApplicationAPI.GetApplicationStatus(ctx, statusId).Execute()
		if statusErr != nil {
			d.Partial(true)
			return appendErrorDiags(diags, "Unable to get Application status", statusErr)
		}
		if appStatus.GetCompleted() {
			if appStatus.Result != nil {
				applicationId = getIdFromUrl(appStatus.Result.GetResource())
			}
			if applicationId == "" {
				// Async result didn't carry a resource URL. Fall back to
				// listing applications and finding by name.
				appName := d.Get("name").(string)
				resp, _, lerr := client.ApplicationAPI.GetApplications(ctx).Execute()
				if lerr == nil && resp != nil {
					for _, a := range resp.Data {
						if a.Name == appName {
							applicationId = a.Id
							break
						}
					}
				}
				if applicationId == "" {
					d.Partial(true)
					return appendErrorDiags(diags, "Unable to extract Application id post-Create",
						fmt.Errorf("async status completed but result.resource was empty and list-by-name found no match for name=%s", appName))
				}
			}
			break
		}
		time.Sleep(30 * time.Second)
	}

	d.SetId(applicationId)
	return resourceApplicationRead(ctx, d, m)
}

/*
resourceApplicationRead Read an Application.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	applicationId := d.Id()
	appData, _, err := client.ApplicationAPI.GetApplicationById(ctx, applicationId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Application", err)
	}

	var appName, appType string
	if appData.HttpApplication != nil {
		appName = appData.HttpApplication.Name
		appType = appData.HttpApplication.Type
	} else if appData.HttpsApplication != nil {
		appName = appData.HttpsApplication.Name
		appType = appData.HttpsApplication.Type
	} else if appData.RdpApplication != nil {
		appName = appData.RdpApplication.Name
		appType = appData.RdpApplication.Type
	}

	if err := d.Set("name", appName); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Application name", err)
	}
	if err := d.Set("type", appType); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Application type", err)
	}

	return diags
}

/*
resourceApplicationDelete is a no-op since there is no delete endpoint in the v2.3 API.
The resource is removed from state only.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No delete endpoint exists for applications in the v2.3 API.
	// Remove from state only.
	d.SetId("")
	return nil
}
