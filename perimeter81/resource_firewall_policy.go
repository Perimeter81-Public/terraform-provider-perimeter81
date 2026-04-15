package perimeter81

import (
	"context"
	"fmt"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceFirewallPolicy Setup the Firewall Policy Resource CRUD operations.
Firewall policies are auto-created with each network, so this resource "adopts"
the existing policy. There is no Create or Delete operation — only Read and Update.

@return &schema.Resource
*/
func resourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallPolicyCreate,
		ReadContext:   resourceFirewallPolicyRead,
		UpdateContext: resourceFirewallPolicyUpdate,
		DeleteContext: resourceFirewallPolicyDelete,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the network whose firewall policy to manage.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the firewall policy is enabled.",
			},
			"allowed": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the default policy action is allow (true) or drop (false).",
			},
			"trace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the policy is traced.",
			},
			"policy_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of firewall policy rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique ID of the policy rule.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the policy rule.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether this rule is enabled.",
						},
						"allowed": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether this rule allows (true) or denies (false) the traffic.",
						},
						"services": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of service object IDs to match in this rule.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceFirewallPolicyImportState,
		},
	}
}

/*
resourceFirewallPolicyImportState Import a firewall policy by its network ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return []*schema.ResourceData, error
*/
func resourceFirewallPolicyImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Use the resource ID as the network_id
	if err := d.Set("network_id", d.Id()); err != nil {
		return nil, fmt.Errorf("could not set network_id: %s", err)
	}
	diagnostics := resourceFirewallPolicyRead(ctx, d, m)
	if diagnostics.HasError() {
		for _, diagnostic := range diagnostics {
			if diagnostic.Severity == diag.Error {
				return nil, fmt.Errorf("could not import firewall policy: %s, \n %s", diagnostic.Summary, diagnostic.Detail)
			}
		}
	}
	return []*schema.ResourceData{d}, nil
}

/*
resourceFirewallPolicyCreate "adopts" the existing firewall policy for the given network by reading its current
state and setting the resource ID to the network_id. Then it applies the desired configuration via update.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceFirewallPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	// Read existing policy to get its ID
	policyData, _, err := client.FirewallPolicyAPI.GetFirewallPolicy(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read Firewall Policy for adoption", err)
	}

	// Use the network_id as the resource ID (policy ID is stored separately)
	d.SetId(networkId)
	if err := d.Set("enabled", policyData.Enabled); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Firewall Policy enabled", err)
	}
	if err := d.Set("allowed", policyData.Allowed); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Firewall Policy allowed", err)
	}

	// Now apply the desired configuration
	return resourceFirewallPolicyUpdate(ctx, d, m)
}

/*
resourceFirewallPolicyRead Read a Firewall Policy by network ID.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceFirewallPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	policyData, _, err := client.FirewallPolicyAPI.GetFirewallPolicy(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to find Firewall Policy", err)
	}

	if err := d.Set("enabled", policyData.Enabled); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Firewall Policy enabled", err)
	}
	if err := d.Set("allowed", policyData.Allowed); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Firewall Policy allowed", err)
	}

	policyRules := flattenFirewallPolicyRules(policyData.PolicyRules)
	if err := d.Set("policy_rules", policyRules); err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to set Firewall Policy rules", err)
	}

	return diags
}

/*
flattenFirewallPolicyRules converts a list of FirewallPolicyRule SDK models to a Terraform-compatible list.
*/
func flattenFirewallPolicyRules(rules []perimeter81Sdk.FirewallPolicyRule) []interface{} {
	if rules == nil {
		return make([]interface{}, 0)
	}
	result := make([]interface{}, len(rules))
	for i, rule := range rules {
		ruleMap := map[string]interface{}{
			"name":     rule.Name,
			"enabled":  rule.Enabled,
			"allowed":  rule.Allowed,
			"services": rule.Services,
		}
		if rule.Id != nil {
			ruleMap["id"] = *rule.Id
		} else {
			ruleMap["id"] = ""
		}
		result[i] = ruleMap
	}
	return result
}

/*
resourceFirewallPolicyUpdate Update the Firewall Policy configuration.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceFirewallPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	ctx = context.Background()

	networkId := d.Get("network_id").(string)

	// Read current policy to get the policy ID
	policyData, _, err := client.FirewallPolicyAPI.GetFirewallPolicy(ctx, networkId).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to read Firewall Policy for update", err)
	}

	enabled := d.Get("enabled").(bool)
	allowed := d.Get("allowed").(bool)

	// Build policy rules from schema
	policyRulesRaw := d.Get("policy_rules").([]interface{})
	policyRules := make([]perimeter81Sdk.FirewallPolicyRule, len(policyRulesRaw))
	for i, ruleRaw := range policyRulesRaw {
		ruleMap := ruleRaw.(map[string]interface{})
		rule := perimeter81Sdk.FirewallPolicyRule{
			Name:    ruleMap["name"].(string),
			Enabled: ruleMap["enabled"].(bool),
			Allowed: ruleMap["allowed"].(bool),
			// Sources and Destinations are not yet managed by this resource —
			// keep them as zero values to leave them unchanged.
			Sources:      perimeter81Sdk.SourcesAndDestinations{},
			Destinations: perimeter81Sdk.SourcesAndDestinations{},
		}
		if v, ok := ruleMap["id"].(string); ok && v != "" {
			rule.Id = &v
		}
		if v, ok := ruleMap["services"].([]interface{}); ok {
			rule.Services = flattenStringsArrayData(v)
		}
		policyRules[i] = rule
	}

	updatePayload := perimeter81Sdk.FirewallPolicy{
		Id:          policyData.Id,
		Enabled:     enabled,
		Allowed:     allowed,
		PolicyRules: policyRules,
	}

	if v, ok := d.GetOk("trace"); ok {
		trace := v.(bool)
		updatePayload.Trace = &trace
	}

	_, _, err = client.FirewallPolicyAPI.UpdateFirewallPolicy(ctx, networkId).FirewallPolicy(updatePayload).Execute()
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to update Firewall Policy", err)
	}

	return resourceFirewallPolicyRead(ctx, d, m)
}

/*
resourceFirewallPolicyDelete is a no-op since firewall policies cannot be deleted —
they are created automatically with each network. The resource is simply removed from state.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client

@return diag.Diagnostics
*/
func resourceFirewallPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Firewall policies cannot be deleted — they are auto-created with the network.
	// Remove from state only.
	d.SetId("")
	return nil
}
