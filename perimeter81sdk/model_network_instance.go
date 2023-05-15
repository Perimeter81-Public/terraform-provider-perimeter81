/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

type NetworkInstance struct {
	// ID of the network.
	Network string `json:"network"`
	// ID of the network region.
	Region       string `json:"region"`
	InstanceType string `json:"instanceType"`
	ImageType    string `json:"imageType"`
	ImageVersion string `json:"imageVersion"`
	// Unique ID.
	ResourceId string `json:"resourceId"`
	Dns        string `json:"dns"`
	Ip         string `json:"ip"`
	// List of network tunnels.
	Tunnels []NetworkTunnel `json:"tunnels"`
	// Unique ID.
	Id string `json:"id"`
	// ID of the tenant.
	TenantId string `json:"tenantId"`
	// The date when this record was created.
	CreatedAt string `json:"createdAt"`
	// The date of last update of the record.
	UpdatedAt string `json:"updatedAt,omitempty"`
}
