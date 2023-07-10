package perimeter81

import (
	"context"
	"log"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
Provider Set up the provider schema

@return &schema.Provider
*/
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PERIMETER81_API_KEY", nil),
				Description: descriptions["api_key"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASE_URL", nil),
				Description: descriptions["base_url"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"perimeter81_network":         resourceNetwork(),
			"perimeter81_wireguard":       resourceWireguard(),
			"perimeter81_openvpn":         resourceOpenvpn(),
			"perimeter81_ipsec_single":    resourceIpsecSingle(),
			"perimeter81_ipsec_redundant": resourceIpsecRedundant(),
			"perimeter81_gateway":         resourceGateway(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"perimeter81_networks": dataSourceNetworks(),
			"perimeter81_regions":  dataSourceRegions(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

/*
providerConfigure Intialize the provider client SDK configuration
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data

@return interface{} - the terraform meta data that contains the client, and diag.Diagnostics
*/
func providerConfigure(con context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	// Get the api key and base url from the provider schema
	apiKey := d.Get("api_key").(string)
	baseUrl := d.Get("base_url").(string)

	// Initialize the perimeter81 client sdk
	var client interface{}
	if apiKey != "" {
		client = perimeter81Sdk.NewAPIClient(perimeter81Sdk.NewConfiguration(apiKey, baseUrl))
	}

	// check if the client is initialized correctly
	if client == nil {
		log.Println("[ERROR] Initializing perimeter81 client is not completed")
		return nil, nil
	}
	log.Println("[INFO] Initializing perimeter81 client")

	return client, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key":  "The Api key for the Preimeter81 Public API.",
		"base_url": "The base url for the rest api.",
	}
}
