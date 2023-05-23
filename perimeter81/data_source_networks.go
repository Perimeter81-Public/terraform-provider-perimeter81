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
		return appendErrorDiags(diags, "Unable to get Networks", err)
	}

	newNetworks := flattenNetworksData(networks.Networks)
	if err := d.Set("networks", newNetworks); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
