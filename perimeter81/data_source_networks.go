package perimeter81

import (
	"context"
	"strconv"
	perimeter81Sdk "terraform-provider-perimeter81/perimeter81sdk"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworksRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"applications": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accesstype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isdefault": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tenantid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"createdat": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updatedat": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regions": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"network": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dns": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tenantid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"createdat": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"updatedat": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instances": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"network": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"dns": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"resourceid": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"ip": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"tenantid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"imageversion": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"imagetype": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"instancetype": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"createdat": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"updatedat": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"tunnels": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"instance": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"interfacename": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"leftallowedip": {
																Type:     schema.TypeList,
																Computed: true,
																Optional: true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"leftendpoint": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"network": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"region": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"requestconfigtoken": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"id": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"tenantid": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"createdat": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"updatedat": {
																Type:     schema.TypeString,
																Optional: true,
															},
														}}},
											},
										}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}

	networks, _, err := client.NetworksApi.GetNetworks(ctx)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Network",
			Detail:   err.Error(),
		})
	}

	newNetworks := flattenNetworksData(networks.Networks)
	if err := d.Set("networks", newNetworks); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenNetworksData(networkItems []perimeter81Sdk.Network) []interface{} {
	if networkItems != nil {
		networks := make([]interface{}, len(networkItems))
		for i, serverItem := range networkItems {
			network := make(map[string]interface{})
			network["name"] = serverItem.Name
			network["id"] = serverItem.Id
			network["tags"] = serverItem.Tags
			network["subnet"] = serverItem.Subnet
			network["dns"] = serverItem.Dns
			network["accesstype"] = serverItem.AccessType
			network["isdefault"] = serverItem.IsDefault
			network["tenantid"] = serverItem.TenantId
			network["createdat"] = serverItem.CreatedAt
			network["updatedat"] = serverItem.UpdatedAt
			network["regions"] = flattenNetworkRegionsData(serverItem.Regions)
			networks[i] = network
		}
		return networks
	}
	return make([]interface{}, 0)
}

func flattenNetworkRegionsData(regionItems []perimeter81Sdk.NetworkRegion) []interface{} {
	if regionItems != nil {
		regions := make([]interface{}, len(regionItems))
		for i, regionItem := range regionItems {
			region := make(map[string]interface{})
			region["network"] = regionItem.Network
			region["dns"] = regionItem.Dns
			region["name"] = regionItem.Name
			region["tenantid"] = regionItem.TenantId
			region["createdat"] = regionItem.CreatedAt
			region["updatedat"] = regionItem.UpdatedAt
			region["id"] = regionItem.Id
			region["instances"] = flattenNetworkInstancesData(regionItem.Instances)
			regions[i] = region
		}
		return regions
	}
	return make([]interface{}, 0)
}

func flattenNetworkInstancesData(instanceItems []perimeter81Sdk.NetworkInstance) []interface{} {
	if instanceItems != nil {
		instances := make([]interface{}, len(instanceItems))
		for i, instanceItem := range instanceItems {
			instance := make(map[string]interface{})
			instance["network"] = instanceItem.Network
			instance["dns"] = instanceItem.Dns
			instance["tenantid"] = instanceItem.TenantId
			instance["createdat"] = instanceItem.CreatedAt
			instance["updatedat"] = instanceItem.UpdatedAt
			instance["resourceid"] = instanceItem.ResourceId
			instance["ip"] = instanceItem.Ip
			instance["id"] = instanceItem.Id
			instance["imageversion"] = instanceItem.ImageVersion
			instance["imagetype"] = instanceItem.ImageType
			instance["region"] = instanceItem.Region
			instance["instancetype"] = instanceItem.InstanceType
			instance["tunnels"] = flattenNetworkTunnelsData(instanceItem.Tunnels)
			instances[i] = instance
		}
		return instances
	}

	return make([]interface{}, 0)
}

func flattenNetworkTunnelsData(tunnelItems []perimeter81Sdk.NetworkTunnel) []interface{} {
	if tunnelItems != nil {
		tunnels := make([]interface{}, len(tunnelItems))
		for i, tunnelItem := range tunnelItems {
			tunnel := make(map[string]interface{})
			tunnel["instance"] = tunnelItem.Instance
			tunnel["interfacename"] = tunnelItem.InterfaceName
			tunnel["leftallowedip"] = tunnelItem.LeftAllowedIP
			tunnel["leftendpoint"] = tunnelItem.LeftEndpoint
			tunnel["network"] = tunnelItem.Network
			tunnel["region"] = tunnelItem.Region
			tunnel["requestconfigtoken"] = tunnelItem.RequestConfigToken
			tunnel["type"] = tunnelItem.Type_
			tunnel["id"] = tunnelItem.Id
			tunnel["tenantid"] = tunnelItem.TenantId
			tunnel["createdat"] = tunnelItem.CreatedAt
			tunnel["updatedat"] = tunnelItem.UpdatedAt
			tunnels[i] = tunnel
		}
		return tunnels
	}

	return make([]interface{}, 0)
}
