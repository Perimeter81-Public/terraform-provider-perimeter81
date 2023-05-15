/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

type CreateWireguardTunnelPayload struct {
	// Region ID
	RegionID string `json:"regionID"`
	// Gatwway ID
	GatewayID      string   `json:"gatewayID"`
	TunnelName     string   `json:"tunnelName"`
	RemoteEndpoint string   `json:"remoteEndpoint,omitempty"`
	RemoteSubnets  []string `json:"remoteSubnets,omitempty"`
}
