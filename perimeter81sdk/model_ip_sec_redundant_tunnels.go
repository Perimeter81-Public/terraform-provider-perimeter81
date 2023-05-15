/*
 * Perimeter81 Public API
 *
 * The YAML for Preimeter81 Public API.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package perimeter81sdk

import (
	"time"
)

type IpSecRedundantTunnels struct {
	TunnelName       string                 `json:"tunnelName,omitempty"`
	RegionID         string                 `json:"regionID,omitempty"`
	Tunnel1          *IpSecRedundantTunnel  `json:"tunnel1,omitempty"`
	Tunnel2          *IpSecRedundantTunnel  `json:"tunnel2,omitempty"`
	SharedSettings   *IpSecSharedSettings   `json:"sharedSettings,omitempty"`
	AdvancedSettings *IpSecAdvancedSettings `json:"advancedSettings,omitempty"`
	CreatedAt        time.Time              `json:"createdAt,omitempty"`
	UpdatedAt        time.Time              `json:"updatedAt,omitempty"`
}
