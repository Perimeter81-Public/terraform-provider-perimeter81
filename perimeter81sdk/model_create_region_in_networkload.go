/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

type CreateRegionInNetworkload struct {
	// cpRegion ID.
	CpRegionId string `json:"cpRegionId"`
	// Desired number of instances in region.
	InstanceCount int32 `json:"instanceCount"`
	// Create the gateway as disabled if true.
	Idle bool `json:"idle"`
}