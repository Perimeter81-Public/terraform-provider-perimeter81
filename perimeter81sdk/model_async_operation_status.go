/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

type AsyncOperationStatus struct {
	Completed bool                  `json:"completed,omitempty"`
	Result    *AsyncOperationResult `json:"result,omitempty"`
}
