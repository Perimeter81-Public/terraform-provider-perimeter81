package perimeter81

import (
	"context"
	"strconv"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
dataSourceNetworks Query all Networks

@return &schema.Resource
*/
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

/*
dataSourceNetworksRead Use the SDK to query all Networks
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param d *schema.ResourceData - the terraform resource data
  - @param m interface{} - the terraform meta data that contains the client
@return diag.Diagnostics
*/

func dataSourceNetworksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// intialize the client and the context if not exists
	var diags diag.Diagnostics
	client := m.(*perimeter81Sdk.APIClient)
	if ctx == nil {
		ctx = context.Background()
	}

	// call the api and check if there is an error
	networks, _, err := client.NetworksApi.GetNetworks(ctx)
	if err != nil {
		d.Partial(true)
		return appendErrorDiags(diags, "Unable to get Networks", err)
	}

	// flatten the data so it fit the terraform schema and set the terraform resource data
	newNetworks := flattenNetworksData(networks)
	if err := d.Set("networks", newNetworks); err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
