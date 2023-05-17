/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

type Region struct {
	CountryCode   string `json:"countryCode,omitempty"`
	ContinentCode string `json:"continentCode,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Name          string `json:"name,omitempty"`
	ClassName     string `json:"className,omitempty"`
	ObjectName    string `json:"objectName,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
	Id            string `json:"id,omitempty"`
}