package checkpointsase

import (
	"context"
	"log"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

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
			"checkpointsase_network":                 resourceNetwork(),
			"checkpointsase_wireguard":               resourceWireguard(),
			"checkpointsase_openvpn":                 resourceOpenvpn(),
			"checkpointsase_ipsec_single":            resourceIpsecSingle(),
			"checkpointsase_ipsec_redundant":         resourceIpsecRedundant(),
			"checkpointsase_gateway":                 resourceGateway(),
			"checkpointsase_object_services":         resourceObjectServices(),
			"checkpointsase_object_addresses":        resourceObjectAddresses(),
			"checkpointsase_enhanced_network":        resourceEnhancedNetwork(),
			"checkpointsase_enhanced_region":         resourceEnhancedRegion(),
			"checkpointsase_enhanced_static_tunnel":  resourceEnhancedStaticTunnel(),
			"checkpointsase_enhanced_dynamic_tunnel": resourceEnhancedDynamicTunnel(),
			"checkpointsase_enhanced_route_table":    resourceEnhancedRouteTable(),
			"checkpointsase_application":             resourceApplication(),
			"checkpointsase_firewall_policy":         resourceFirewallPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"checkpointsase_networks":                dataSourceNetworks(),
			"checkpointsase_standard_networks":       dataSourceStandardNetworks(),
			"checkpointsase_all_networks":            dataSourceAllNetworks(),
			"checkpointsase_regions":                 dataSourceRegions(),
			"checkpointsase_object_services":         dataSourceObjectServices(),
			"checkpointsase_object_addresses":        dataSourceObjectAddresses(),
			"checkpointsase_enhanced_networks":       dataSourceEnhancedNetworks(),
			"checkpointsase_enhanced_regions":        dataSourceEnhancedRegions(),
			"checkpointsase_applications":            dataSourceApplications(),
			"checkpointsase_route_table":             dataSourceRouteTable(),
			"checkpointsase_enhanced_route_table":    dataSourceEnhancedRouteTable(),
			"checkpointsase_network_health":          dataSourceNetworkHealth(),
			"checkpointsase_enhanced_network_health": dataSourceEnhancedNetworkHealth(),
			"checkpointsase_enhanced_tunnels":        dataSourceEnhancedTunnels(),
			"checkpointsase_customer_certificates":   dataSourceCustomerCertificates(),
			"checkpointsase_status":                  dataSourceStatus(),
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
