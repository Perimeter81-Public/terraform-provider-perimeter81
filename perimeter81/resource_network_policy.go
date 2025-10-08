package perimeter81

import (
	"context"
	"fmt"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceNetworkPolicy Setup the Network Policy Resource CRUD operations

@return &schema.Resource
*/
func resourceNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPolicyCreate,
		ReadContext:   resourceNetworkPolicyRead,
		UpdateContext: resourceNetworkPolicyUpdate,
		DeleteContext: resourceNetworkPolicyDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ID this policy belongs to",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the policy is enabled",
			},
			"allowed": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the policy is allow or deny",
			},
			"policy_rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Array of policy rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy rule ID",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the policy rule",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the rule is enabled",
						},
						"allowed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the rule allows or denies traffic",
						},
						"sources": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Source users, groups, and addresses",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"users": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of user IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of group IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"addresses": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of address object IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"destinations": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Destination users, groups, and addresses",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"users": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of user IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"groups": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of group IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"addresses": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of address object IDs",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"services": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of service IDs",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceNetworkPolicyImportState,
		},
	}
}

/*
resourceNetworkPolicyImportState Import network policy
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceNetworkPolicyImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// The ID for import should be the network_id
	networkId := d.Id()
	d.Set("network_id", networkId)
	d.SetId(networkId)

	diagnostics := resourceNetworkPolicyRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import network policy: %s, %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceNetworkPolicyCreate Create a Network Policy (actually updates since there's one per network)
  - @param ctx context.Context
  - @param d *schema.ResourceData
  - @param m interface{}

@return diag.Diagnostics
*/
func resourceNetworkPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For network policy, create is the same as update since there's one policy per network
	return resourceNetworkPolicyUpdate(ctx, d, m)
}

/*
resourceNetworkPolicyRead Read a Network Policy
  - @param ctx context.Context
  - @param d *schema.ResourceData
  - @param m interface{}

@return diag.Diagnostics
*/
func resourceNetworkPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	// Get the network policy
	firewallPolicy, _, err := client.FirewallPolicyApi.GetFirewallPolicy(ctx, networkId)
	if err != nil {
		return appendErrorDiags(diags, "Unable to read network policy", err)
	}

	// Set the resource ID to the network ID (since there's one policy per network)
	d.SetId(networkId)

	if err := d.Set("enabled", firewallPolicy.Enabled); err != nil {
		return appendErrorDiags(diags, "Unable to set enabled", err)
	}
	if err := d.Set("allowed", firewallPolicy.Allowed); err != nil {
		return appendErrorDiags(diags, "Unable to set allowed", err)
	}

	// Flatten and set policy rules
	policyRules := flattenNetworkPolicyRules(firewallPolicy.PolicyRules)
	if err := d.Set("policy_rules", policyRules); err != nil {
		return appendErrorDiags(diags, "Unable to set policy_rules", err)
	}

	return diags
}

/*
resourceNetworkPolicyUpdate Update a Network Policy
  - @param ctx context.Context
  - @param d *schema.ResourceData
  - @param m interface{}

@return diag.Diagnostics
*/
func resourceNetworkPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)
	enabled := d.Get("enabled").(bool)
	allowed := d.Get("allowed").(bool)
	policyRulesRaw := d.Get("policy_rules").([]interface{})

	// Expand policy rules
	policyRules := expandNetworkPolicyRules(policyRulesRaw)

	// Create the network policy payload
	firewallPolicyPayload := perimeter81Sdk.FirewallPolicy{
		Enabled:     enabled,
		Allowed:     allowed,
		PolicyRules: policyRules,
	}

	// Update the network policy
	_, _, err := client.FirewallPolicyApi.UpdateFirewallPolicy(ctx, firewallPolicyPayload, networkId)
	if err != nil {
		return appendErrorDiags(diags, "Unable to update network policy", err)
	}

	// Set the resource ID
	d.SetId(networkId)
	d.Set("last_updated", time.Now().Format(time.RFC3339))

	return resourceNetworkPolicyRead(ctx, d, m)
}

/*
resourceNetworkPolicyDelete Delete a Network Policy (sets to disabled state)
  - @param ctx context.Context
  - @param d *schema.ResourceData
  - @param m interface{}

@return diag.Diagnostics
*/
func resourceNetworkPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	// Since we can't actually delete the policy (one per network), we disable it
	firewallPolicyPayload := perimeter81Sdk.FirewallPolicy{
		Enabled:     false,
		Allowed:     true,
		PolicyRules: []perimeter81Sdk.FirewallPolicyRule{},
	}

	_, _, err := client.FirewallPolicyApi.UpdateFirewallPolicy(ctx, firewallPolicyPayload, networkId)
	if err != nil {
		return appendErrorDiags(diags, "Unable to disable network policy", err)
	}

	d.SetId("")
	return diags
}

/*
expandNetworkPolicyRules Expand policy rules from Terraform schema to SDK type
*/
func expandNetworkPolicyRules(rulesRaw []interface{}) []perimeter81Sdk.FirewallPolicyRule {
	if len(rulesRaw) == 0 {
		return []perimeter81Sdk.FirewallPolicyRule{}
	}

	rules := make([]perimeter81Sdk.FirewallPolicyRule, len(rulesRaw))
	for i, ruleRaw := range rulesRaw {
		ruleMap := ruleRaw.(map[string]interface{})

		rule := perimeter81Sdk.FirewallPolicyRule{
			Name:    ruleMap["name"].(string),
			Enabled: ruleMap["enabled"].(bool),
			Allowed: ruleMap["allowed"].(bool),
		}

		// Set ID if it exists
		if id, ok := ruleMap["id"].(string); ok && id != "" {
			rule.Id = id
		}

		// Expand sources
		if sourcesRaw, ok := ruleMap["sources"].([]interface{}); ok && len(sourcesRaw) > 0 {
			sourcesMap := sourcesRaw[0].(map[string]interface{})
			sources := &perimeter81Sdk.SourcesAndDestinations{}
			
			if users, ok := sourcesMap["users"].([]interface{}); ok && len(users) > 0 {
				sources.Users = flattenStringsArrayData(users)
			}
			if groups, ok := sourcesMap["groups"].([]interface{}); ok && len(groups) > 0 {
				sources.Groups = flattenStringsArrayData(groups)
			}
			if addresses, ok := sourcesMap["addresses"].([]interface{}); ok && len(addresses) > 0 {
				sources.Addresses.Addresses = flattenStringsArrayData(addresses)
			}
			
			rule.Sources = sources
		}

		// Expand destinations
		if destinationsRaw, ok := ruleMap["destinations"].([]interface{}); ok && len(destinationsRaw) > 0 {
			destinationsMap := destinationsRaw[0].(map[string]interface{})
			destinations := &perimeter81Sdk.SourcesAndDestinations{}
			
			if users, ok := destinationsMap["users"].([]interface{}); ok && len(users) > 0 {
				destinations.Users = flattenStringsArrayData(users)
			}
			if groups, ok := destinationsMap["groups"].([]interface{}); ok && len(groups) > 0 {
				destinations.Groups = flattenStringsArrayData(groups)
			}
			if addresses, ok := destinationsMap["addresses"].([]interface{}); ok && len(addresses) > 0 {
				destinations.Addresses.Addresses = flattenStringsArrayData(addresses)
			}
			
			rule.Destinations = destinations
		}

		// Expand services
		if servicesRaw, ok := ruleMap["services"].([]interface{}); ok {
			rule.Services = flattenStringsArrayData(servicesRaw)
		}

		rules[i] = rule
	}

	return rules
}

/*
flattenNetworkPolicyRules Flatten policy rules from SDK type to Terraform schema
*/
func flattenNetworkPolicyRules(rules []perimeter81Sdk.FirewallPolicyRule) []interface{} {
	if len(rules) == 0 {
		return []interface{}{}
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		ruleMap := map[string]interface{}{
			"id":      rule.Id,
			"name":    rule.Name,
			"enabled": rule.Enabled,
			"allowed": rule.Allowed,
		}

		// Flatten sources
		sources := map[string]interface{}{
			"users":     []string{},
			"groups":    []string{},
			"addresses": []string{},
		}
		if rule.Sources != nil {
			if rule.Sources.Users != nil {
				sources["users"] = rule.Sources.Users
			}
			if rule.Sources.Groups != nil {
				sources["groups"] = rule.Sources.Groups
			}
			if len(rule.Sources.Addresses.Addresses) > 0 {
				sources["addresses"] = rule.Sources.Addresses.Addresses
			}
		}
		ruleMap["sources"] = []interface{}{sources}

		// Flatten destinations
		destinations := map[string]interface{}{
			"users":     []string{},
			"groups":    []string{},
			"addresses": []string{},
		}
		if rule.Destinations != nil {
			if rule.Destinations.Users != nil {
				destinations["users"] = rule.Destinations.Users
			}
			if rule.Destinations.Groups != nil {
				destinations["groups"] = rule.Destinations.Groups
			}
			if len(rule.Destinations.Addresses.Addresses) > 0 {
				destinations["addresses"] = rule.Destinations.Addresses.Addresses
			}
		}
		ruleMap["destinations"] = []interface{}{destinations}

		// Flatten services
		if rule.Services != nil {
			ruleMap["services"] = rule.Services
		} else {
			ruleMap["services"] = []string{}
		}

		flattenedRules[i] = ruleMap
	}

	return flattenedRules
}