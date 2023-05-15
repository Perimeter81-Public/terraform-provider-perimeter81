package perimeter81

import (
	"context"
	"log"
	perimeter81Sdk "terraform-provider-perimeter81/perimeter81sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"perimeter81_network": resourceNetwork(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"perimeter81_networks": dataSourceNetworks(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(con context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)

	var client interface{}
	if apiKey != "" {
		client = perimeter81Sdk.NewAPIClient(perimeter81Sdk.NewConfiguration(apiKey))
	}
	// var diags diag.Diagnostics
	// diags = append(diags, diag.Diagnostic{
	// 	Severity: diag.Error,
	// 	Summary:  "Unable to create HashiCups client",
	// 	Detail:   client,
	// })

	if client == nil {
		log.Println("[ERROR] Initializing postmark client is not completed")
		return nil, nil
	}
	log.Println("[INFO] Initializing postmark client")

	return client, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_key": "The Api key for the Preimeter81 Public API.",
	}
}
