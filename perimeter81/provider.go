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
				DefaultFunc: schema.EnvDefaultFunc("CHECKPOINT_SASE_API_KEY", nil),
				Description: descriptions["api_key"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASE_URL", perimeter81Sdk.BaseURLUS),
				Description: descriptions["base_url"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sase_network":                 resourceNetwork(),
			"sase_wireguard":               resourceWireguard(),
			"sase_openvpn":                 resourceOpenvpn(),
			"sase_ipsec_single":            resourceIpsecSingle(),
			"sase_ipsec_redundant":         resourceIpsecRedundant(),
			"sase_gateway":                 resourceGateway(),
			"sase_object_services":         resourceObjectServices(),
			"sase_object_addresses":        resourceObjectAddresses(),
			"sase_enhanced_network":        resourceEnhancedNetwork(),
			"sase_enhanced_region":         resourceEnhancedRegion(),
			"sase_enhanced_static_tunnel":  resourceEnhancedStaticTunnel(),
			"sase_enhanced_dynamic_tunnel": resourceEnhancedDynamicTunnel(),
			"sase_enhanced_route_table":    resourceEnhancedRouteTable(),
			"sase_application":             resourceApplication(),
			"sase_firewall_policy":         resourceFirewallPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sase_networks":                dataSourceNetworks(),
			"sase_standard_networks":       dataSourceStandardNetworks(),
			"sase_all_networks":            dataSourceAllNetworks(),
			"sase_regions":                 dataSourceRegions(),
			"sase_object_services":         dataSourceObjectServices(),
			"sase_object_addresses":        dataSourceObjectAddresses(),
			"sase_enhanced_networks":       dataSourceEnhancedNetworks(),
			"sase_enhanced_regions":        dataSourceEnhancedRegions(),
			"sase_applications":            dataSourceApplications(),
			"sase_route_table":             dataSourceRouteTable(),
			"sase_enhanced_route_table":    dataSourceEnhancedRouteTable(),
			"sase_network_health":          dataSourceNetworkHealth(),
			"sase_enhanced_network_health": dataSourceEnhancedNetworkHealth(),
			"sase_enhanced_tunnels":        dataSourceEnhancedTunnels(),
			"sase_customer_certificates":   dataSourceCustomerCertificates(),
			"sase_status":                  dataSourceStatus(),
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

	// Initialize the Check Point Check Point SASE client SDK
	var client interface{}
	if apiKey != "" {
		client = perimeter81Sdk.NewAPIClient(perimeter81Sdk.NewConfiguration(apiKey, baseUrl))
	}

	// check if the client is initialized correctly
	if client == nil {
		log.Println("[ERROR] Initializing Check Point SASE client is not completed")
		return nil, nil
	}
	log.Println("[INFO] Initializing Check Point SASE client")

	return client, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key":  "The API key for the Check Point SASE Public API.",
		"base_url": "The base URL for the Check Point SASE REST API. Defaults to the US endpoint if not set.",
	}
}
