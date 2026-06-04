package checkpointsase

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// parseASNString converts a user-supplied ASN string (e.g. "65010") to the
// SDK's RemoteASN type. Invalid input returns 0; the API validator catches
// out-of-range values. Used by resources whose HCL schema declares the ASN
// as a string (historical reasons) — newer resources use TypeInt directly.
func parseASNString(s string) perimeter81Sdk.RemoteASN {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return perimeter81Sdk.RemoteASN(int32(n))
}

/*
flattenStringsArrayData flatten string array data
  - @param strs []interface{} - the strings that need to be flattened

@return []string - the flattened strings
*/
func flattenStringsArrayData(strs []interface{}) []string {
	strsData := make([]string, len(strs))
	for index, str := range strs {
		strsData[index] = fmt.Sprint(str)
	}
	return strsData
}

/*
flattenIntsArrayData flatten ints array data
  - @param ints []interface{} - the ints that need to be flattened

@return []int32 - the flattened ints
*/
func flattenIntsArrayData(ints []interface{}) []int32 {
	intsData := make([]int32, len(ints))
	for index, ele := range ints {
		intsData[index] = int32(ele.(int))
	}
	return intsData
}

/*
getIdFromUrl split the url and get the last element(id)
  - @param url string - the url that need to be splited to get the id

@return string - the id
*/
func getIdFromUrl(url string) string {
	urlSplited := strings.Split(url, "/")
	return urlSplited[len(urlSplited)-1]
}

/*
flattenRegions flatten regions data
  - @param regionsDate []perimeter81Sdk.Region - the regions that need to be flattened

@return []interface{} - the flattened regions
*/
func flattenRegions(regionsDate []perimeter81Sdk.Region) []interface{} {
	if regionsDate != nil {
		regions := make([]interface{}, len(regionsDate))

		for i, regionData := range regionsDate {
			region := make(map[string]interface{})

			region["country_code"] = regionData.GetCountryCode()
			region["continent_code"] = regionData.GetContinentCode()
			region["display_name"] = regionData.GetDisplayName()
			region["name"] = regionData.GetName()
			region["class_name"] = regionData.GetClassName()
			region["object_id"] = regionData.GetId()
			region["id"] = regionData.GetId()
			regions[i] = region
		}

		return regions
	}

	return make([]interface{}, 0)
}

/*
flattenRegionsData flatten regions data
  - @param regionItems []interface{} - the regions that need to be flattened

@return []perimeter81Sdk.StandardNetworkRegionConfig - the flattened regions
*/
func flattenRegionsData(regionItems []interface{}) []StandardNetworkRegionConfig {
	if regionItems != nil {
		regions := make([]StandardNetworkRegionConfig, len(regionItems))

		for i, regionItem := range regionItems {
			region := StandardNetworkRegionConfig{}

			region.CpRegionId = regionItem.(map[string]interface{})["cpregion_id"].(string)
			region.Idle = regionItem.(map[string]interface{})["idle"].(bool)
			region_id := regionItem.(map[string]interface{})["region_id"]
			if region_id != nil {
				region.RegionID = region_id.(string)
			}
			regions[i] = region
		}

		return regions
	}

	return make([]StandardNetworkRegionConfig, 0)
}

/*
flattenProtocolsData flatten Protocols data
  - @param protocolItems []interface{} - the protocols that need to be flattened

@return []perimeter81Sdk.ObjectsServicesProtocolRequestObj - the flattened protocols
*/
func flattenProtocolsData(protocolItems []interface{}) []perimeter81Sdk.ObjectsServicesProtocolRequestObj {
	if protocolItems == nil {
		return make([]perimeter81Sdk.ObjectsServicesProtocolRequestObj, 0)
	}
	protocols := make([]perimeter81Sdk.ObjectsServicesProtocolRequestObj, len(protocolItems))
	for i, protocolItem := range protocolItems {
		m := protocolItem.(map[string]interface{})
		entry := perimeter81Sdk.ObjectsServicesProtocolRequestObj{
			Protocol:  m["protocol"].(string),
			ValueType: m["value_type"].(string),
			Value:     flattenIntsArrayData(m["value"].([]interface{})),
		}
		protocols[i] = entry
	}
	return protocols
}

// StandardNetworkRegionConfig holds the internal representation of a standard network region
// used for tracking region create/delete operations.
type StandardNetworkRegionConfig struct {
	// CpRegionId is the region ID used in the get-regions endpoint (cpregion / harmony-sase region id).
	CpRegionId string
	// RegionID is the ID of the created region inside the network.
	RegionID string
	// Idle indicates whether the gateway should be created as disabled.
	Idle bool
	// Name is the display name of the region.
	Name string
	// Dns is the DNS of the region.
	Dns string
	// DefaultGatewayIp is the IP of the default gateway.
	DefaultGatewayIp string
}

// GatewayConfig holds the internal representation of a gateway.
type GatewayConfig struct {
	Name string
	Idle bool
	Id   string
	Dns  string
	Ip   string
}

/*
flattenGatewaysData flatten gateways data
  - @param gatewaysItems []interface{} - the gateways data that need to be flattened

@return []GatewayConfig - the flattened gateways
*/
func flattenGatewaysData(gatewaysItems []interface{}) []GatewayConfig {
	if gatewaysItems != nil {
		gateways := make([]GatewayConfig, len(gatewaysItems))

		for i, gatewayItem := range gatewaysItems {
			gateway := GatewayConfig{}

			gateway.Name = gatewayItem.(map[string]interface{})["name"].(string)
			gateway.Idle = gatewayItem.(map[string]interface{})["idle"].(bool)
			id := gatewayItem.(map[string]interface{})["id"]
			if id != nil {
				gateway.Id = id.(string)
			}
			gateways[i] = gateway
		}

		return gateways
	}

	return make([]GatewayConfig, 0)
}

/*
flattenGateways flatten gateways data
  - @param gatewaysItems []GatewayConfig - the gateways that need to be flattened

@return []interface{} - the flattened gateways data
*/
func flattenGateways(gatewaysItems []GatewayConfig) []interface{} {
	if gatewaysItems != nil {
		gateways := make([]interface{}, len(gatewaysItems))

		for i, gatewayItems := range gatewaysItems {
			gateway := make(map[string]interface{})

			gateway["name"] = gatewayItems.Name
			gateway["idle"] = gatewayItems.Idle
			gateway["id"] = gatewayItems.Id
			gateway["dns"] = gatewayItems.Dns
			gateway["ip"] = gatewayItems.Ip
			gateways[i] = gateway
		}
		return gateways
	}
	return make([]interface{}, 0)
}

/*
flattenNetworkData flatten network data
  - @param networkItems []perimeter81Sdk.CreateNetworkPayload - the network data that need to be flattened

@return []interface{} - the flattened network data
*/
func flattenNetworkData(networkItems []perimeter81Sdk.CreateNetworkPayload) []interface{} {
	if networkItems != nil {
		networks := make([]interface{}, len(networkItems))

		for i, networkItem := range networkItems {
			network := make(map[string]interface{})

			network["name"] = networkItem.Name
			network["tags"] = networkItem.Tags
			network["subnet"] = networkItem.GetSubnet()
			networks[i] = network
		}

		return networks
	}

	return make([]interface{}, 0)
}

/*
flattenNetworkRegions flatten network regions
  - @param regionItems []StandardNetworkRegionConfig - the network regions that need to be flattened

@return []interface{} - the flattened  network regions
*/
func flattenNetworkRegions(regionItems []StandardNetworkRegionConfig) []interface{} {
	if regionItems != nil {
		regions := make([]interface{}, len(regionItems))

		for i, regionItem := range regionItems {
			region := make(map[string]interface{})

			region["cpregion_id"] = regionItem.CpRegionId
			region["region_id"] = regionItem.RegionID
			region["idle"] = regionItem.Idle
			region["name"] = regionItem.Name
			region["dns"] = regionItem.Dns
			region["default_gateway_ip"] = regionItem.DefaultGatewayIp
			regions[i] = region
		}

		return regions
	}

	return make([]interface{}, 0)
}

/*
flattenObjectServicesProtocols flatten object services protocols
  - @param protocolItems []perimeter81Sdk.ObjectsServicesProtocolResponseObj - the object services protocols that need to be flattened

@return []interface{} - the flattened  network regions
*/
func flattenObjectServicesProtocols(protocolItems []perimeter81Sdk.ObjectsServicesProtocolResponseObj) []interface{} {
	if protocolItems == nil {
		return make([]interface{}, 0)
	}
	protocols := make([]interface{}, len(protocolItems))
	for i, protocolItem := range protocolItems {
		protocols[i] = map[string]interface{}{
			"protocol":   protocolItem.Protocol,
			"value_type": protocolItem.ValueType,
			"value":      protocolItem.Value,
		}
	}
	return protocols
}

/*
flattenRegionData flatten network
  - @param networkItems []perimeter81Sdk.CreateNetworkPayload - the network that need to be flattened

@return []interface{} - the flattened  network data
*/
func flattenRegionData(networkItems []perimeter81Sdk.CreateNetworkPayload) []interface{} {
	if networkItems != nil {
		networks := make([]interface{}, len(networkItems))

		for i, networkItem := range networkItems {
			network := make(map[string]interface{})

			network["name"] = networkItem.Name
			network["tags"] = networkItem.Tags
			network["subnet"] = networkItem.GetSubnet()
			networks[i] = network
		}

		return networks
	}

	return make([]interface{}, 0)
}

/*
flattenNetworksData flatten networks data
  - @param networkItems []perimeter81Sdk.Network - the networks that need to be flattened

@return []interface{} - the flattened  networks data
*/
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
			network["createdat"] = serverItem.CreatedAt.String()
			if serverItem.UpdatedAt != nil {
				network["updatedat"] = serverItem.UpdatedAt.String()
			} else {
				network["updatedat"] = ""
			}
			network["regions"] = flattenNetworkRegionsData(serverItem.Regions)
			networks[i] = network
		}
		return networks
	}
	return make([]interface{}, 0)
}

/*
flattenNetworkRegionsData flatten network regions data
  - @param regionItems []perimeter81Sdk.NetworkRegion - the network regions that need to be flattened

@return []interface{} - the flattened  network regions data
*/
func flattenNetworkRegionsData(regionItems []perimeter81Sdk.NetworkRegion) []interface{} {
	if regionItems != nil {
		regions := make([]interface{}, len(regionItems))
		for i, regionItem := range regionItems {
			region := make(map[string]interface{})
			region["network"] = regionItem.Network
			region["dns"] = regionItem.Dns
			region["name"] = regionItem.Name
			region["tenantid"] = regionItem.TenantId
			region["createdat"] = regionItem.CreatedAt.String()
			if regionItem.UpdatedAt != nil {
				region["updatedat"] = regionItem.UpdatedAt.String()
			} else {
				region["updatedat"] = ""
			}
			region["id"] = regionItem.Id
			region["instances"] = flattenNetworkInstancesData(regionItem.Instances)
			regions[i] = region
		}
		return regions
	}
	return make([]interface{}, 0)
}

/*
flattenNetworkInstancesData flatten network instances data
  - @param instanceItems []perimeter81Sdk.NetworkInstance - the network instances that need to be flattened

@return []interface{} - the flattened  network instances data
*/
func flattenNetworkInstancesData(instanceItems []perimeter81Sdk.NetworkInstance) []interface{} {
	if instanceItems != nil {
		instances := make([]interface{}, len(instanceItems))
		for i, instanceItem := range instanceItems {
			instance := make(map[string]interface{})
			instance["network"] = instanceItem.Network
			instance["dns"] = instanceItem.Dns
			instance["tenantid"] = instanceItem.TenantId
			instance["createdat"] = instanceItem.CreatedAt.String()
			if instanceItem.UpdatedAt != nil {
				instance["updatedat"] = instanceItem.UpdatedAt.String()
			} else {
				instance["updatedat"] = ""
			}
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

/*
flattenNetworkTunnelsData flatten network tunnels data
  - @param tunnelItems []perimeter81Sdk.NetworkTunnel - the network tunnels that need to be flattened

@return []interface{} - the flattened  network tunnels data
*/
func flattenNetworkTunnelsData(tunnelItems []perimeter81Sdk.NetworkTunnel) []interface{} {
	if tunnelItems != nil {
		tunnels := make([]interface{}, len(tunnelItems))
		for i, tunnelItem := range tunnelItems {
			tunnel := make(map[string]interface{})
			// NetworkTunnel is a union type - extract base fields from whichever variant is set
			if tunnelItem.NetworkTunnelWireguard != nil {
				wg := tunnelItem.NetworkTunnelWireguard
				tunnel["instance"] = wg.Instance
				tunnel["interfacename"] = wg.InterfaceName
				tunnel["leftallowedip"] = wg.LeftAllowedIP
				tunnel["leftendpoint"] = wg.LeftEndpoint
				tunnel["network"] = wg.Network
				tunnel["region"] = wg.Region
				tunnel["requestconfigtoken"] = wg.RequestConfigToken
				tunnel["type"] = wg.Type
				tunnel["id"] = wg.Id
				tunnel["tenantid"] = wg.TenantId
				tunnel["createdat"] = wg.CreatedAt.String()
				if wg.UpdatedAt != nil {
					tunnel["updatedat"] = wg.UpdatedAt.String()
				} else {
					tunnel["updatedat"] = ""
				}
			} else if tunnelItem.NetworkTunnelIpsecSingle != nil {
				t := tunnelItem.NetworkTunnelIpsecSingle
				tunnel["instance"] = t.Instance
				tunnel["interfacename"] = t.InterfaceName
				tunnel["leftallowedip"] = []string{}
				tunnel["leftendpoint"] = ""
				tunnel["network"] = t.Network
				tunnel["region"] = t.Region
				tunnel["requestconfigtoken"] = ""
				tunnel["type"] = t.Type
				tunnel["id"] = t.Id
				tunnel["tenantid"] = t.TenantId
				tunnel["createdat"] = t.CreatedAt.String()
				if t.UpdatedAt != nil {
					tunnel["updatedat"] = t.UpdatedAt.String()
				} else {
					tunnel["updatedat"] = ""
				}
			} else if tunnelItem.NetworkTunnelIpsecRedundant != nil {
				t := tunnelItem.NetworkTunnelIpsecRedundant
				tunnel["instance"] = t.Instance
				tunnel["interfacename"] = t.InterfaceName
				tunnel["leftallowedip"] = []string{}
				tunnel["leftendpoint"] = ""
				tunnel["network"] = t.Network
				tunnel["region"] = t.Region
				tunnel["requestconfigtoken"] = ""
				tunnel["type"] = t.Type
				tunnel["id"] = t.Id
				tunnel["tenantid"] = t.TenantId
				tunnel["createdat"] = t.CreatedAt.String()
				if t.UpdatedAt != nil {
					tunnel["updatedat"] = t.UpdatedAt.String()
				} else {
					tunnel["updatedat"] = ""
				}
			} else if tunnelItem.NetworkTunnelOpenvpn != nil {
				t := tunnelItem.NetworkTunnelOpenvpn
				tunnel["instance"] = t.Instance
				tunnel["interfacename"] = t.InterfaceName
				tunnel["leftallowedip"] = []string{}
				tunnel["leftendpoint"] = ""
				tunnel["network"] = t.Network
				tunnel["region"] = t.Region
				tunnel["requestconfigtoken"] = ""
				tunnel["type"] = t.Type
				tunnel["id"] = t.Id
				tunnel["tenantid"] = t.TenantId
				tunnel["createdat"] = t.CreatedAt.String()
				if t.UpdatedAt != nil {
					tunnel["updatedat"] = t.UpdatedAt.String()
				} else {
					tunnel["updatedat"] = ""
				}
			}
			tunnels[i] = tunnel
		}
		return tunnels
	}

	return make([]interface{}, 0)
}

/*
flattenPhasesData flatten Phases date
  - @param phasesItem *perimeter81Sdk.IPSecPhaseConfig - the phase config that need to be flattened

@return []interface{} - the flattened  phases data
*/
func flattenPhasesData(phasesItem *perimeter81Sdk.IPSecPhaseConfig) []interface{} {
	if phasesItem != nil {
		phase := make([]interface{}, 1)
		phaseData := make(map[string]interface{})
		phaseData["auth"] = phasesItem.Auth
		phaseData["encryption"] = phasesItem.Encryption
		phaseData["dh"] = phasesItem.Dh
		phase[0] = phaseData
		return phase
	}

	return make([]interface{}, 0)
}

/*
flattenAdvancedSettingsData flatten Advanced Settings date
  - @param advancedSettingsItem *IPSecAdvancedSettings - the advanced settings that need to be flattened

@return []interface{} - the flattened advanced settings data
*/
func flattenAdvancedSettingsData(advancedSettingsItem *perimeter81Sdk.IPSecAdvancedSettings) []interface{} {
	if advancedSettingsItem != nil {
		advancedSettings := make([]interface{}, 1)
		advancedSettingsData := make(map[string]interface{})
		advancedSettingsData["key_exchange"] = advancedSettingsItem.KeyExchange
		advancedSettingsData["ike_life_time"] = advancedSettingsItem.IkeLifeTime
		advancedSettingsData["lifetime"] = advancedSettingsItem.Lifetime
		advancedSettingsData["dpd_delay"] = advancedSettingsItem.DpdDelay
		advancedSettingsData["dpd_timeout"] = advancedSettingsItem.DpdTimeout
		phase1 := advancedSettingsItem.Phase1
		phase2 := advancedSettingsItem.Phase2
		advancedSettingsData["phase1"] = flattenPhasesData(&phase1)
		advancedSettingsData["phase2"] = flattenPhasesData(&phase2)
		advancedSettings[0] = advancedSettingsData
		return advancedSettings
	}

	return make([]interface{}, 0)
}

/*
flattenSharedSettingsData flatten Shared Settings date
  - @param sharedSettingsItem *IpSecSharedSettings - the Ip-Sec Shared settings that need to be flattened

@return []interface{} - the flattened Ip-Sec Shared settings data
*/
func flattenSharedSettingsData(sharedSettingsItem *perimeter81Sdk.IPSecSharedSettings) []interface{} {
	if sharedSettingsItem != nil {
		sharedSettings := make([]interface{}, 1)
		sharedSettingsData := make(map[string]interface{})
		sharedSettingsData["p81_gateway_subnets"] = sharedSettingsItem.P81GatewaySubnets
		sharedSettingsData["remote_gateway_subnets"] = sharedSettingsItem.RemoteGatewaySubnets
		if sharedSettingsItem.PeakBandwidth != nil {
			sharedSettingsData["peak_bandwidth"] = int(*sharedSettingsItem.PeakBandwidth)
		}
		sharedSettings[0] = sharedSettingsData
		return sharedSettings
	}

	return make([]interface{}, 0)
}

/*
flattenTunnelData flatten Tunnel date
  - @param tunnelItem *IPSecRedundantTunnel - the tunnel that need to be flattened

@return []interface{} - the flattened tunnel data
*/
func flattenTunnelData(tunnelItem *perimeter81Sdk.IPSecRedundantTunnel) []interface{} {
	if tunnelItem != nil {
		tunnel := make([]interface{}, 1)
		tunnelData := make(map[string]interface{})
		tunnelData["passphrase"] = tunnelItem.Passphrase
		tunnelData["gateway_id"] = tunnelItem.GatewayID
		// RemoteID is a union type wrapping *string
		if tunnelItem.RemoteID.String != nil {
			tunnelData["remote_id"] = *tunnelItem.RemoteID.String
		} else {
			tunnelData["remote_id"] = ""
		}
		tunnelData["p81_gwinternal_ip"] = tunnelItem.P81GWInternalIP
		tunnelData["remote_gwinternal_ip"] = tunnelItem.RemoteGWInternalIP
		tunnelData["remote_public_ip"] = tunnelItem.RemotePublicIP
		tunnelData["remote_asn"] = fmt.Sprintf("%d", int32(tunnelItem.RemoteASN))
		if tunnelItem.TunnelID != nil {
			tunnelData["tunnel_id"] = *tunnelItem.TunnelID
		} else {
			tunnelData["tunnel_id"] = ""
		}
		tunnel[0] = tunnelData
		return tunnel
	}

	return make([]interface{}, 0)
}

/*
getTunnelId get the tunnel id
  - @param ctx context.Context - the context
  - @param networkId string - the network id
  - @param tunnelBody perimeter81Sdk.BaseTunnelValues - the tunnel body
  - @param client perimeter81Sdk.APIClient - the client
  - @param diags diag.Diagnostics - the diagnostics

@return string - the tunnel id, diag.Diagnostics - the diagnostics
*/
func getTunnelId(ctx context.Context, networkId string, tunnelBody perimeter81Sdk.BaseTunnelValues, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (string, diag.Diagnostics) {
	network, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, networkId).Execute()
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to fetch network", err)
		return "", diags
	}
	// find the tunnel id based on that tunnel name is unique
	for _, region := range network.Regions {
		if region.Id == tunnelBody.RegionID {
			for _, gateway := range region.Instances {
				if gateway.Id == tunnelBody.GatewayID {
					for _, tunnel := range gateway.Tunnels {
						ifName := getNetworkTunnelInterfaceName(tunnel)
						if ifName == tunnelBody.TunnelName {
							return getNetworkTunnelId(tunnel), diags
						}
					}
				}

			}
		}
	}
	diags = appendErrorDiags(diags, "Unable to find tunnel", fmt.Errorf("check tunnel fields there might be overlap error"))
	return "", diags
}

/*
getNetworkTunnelInterfaceName extract the interface name from a NetworkTunnel union type.
*/
func getNetworkTunnelInterfaceName(tunnel perimeter81Sdk.NetworkTunnel) string {
	if tunnel.NetworkTunnelWireguard != nil {
		return tunnel.NetworkTunnelWireguard.InterfaceName
	}
	if tunnel.NetworkTunnelIpsecSingle != nil {
		return tunnel.NetworkTunnelIpsecSingle.InterfaceName
	}
	if tunnel.NetworkTunnelIpsecRedundant != nil {
		return tunnel.NetworkTunnelIpsecRedundant.InterfaceName
	}
	if tunnel.NetworkTunnelOpenvpn != nil {
		return tunnel.NetworkTunnelOpenvpn.InterfaceName
	}
	return ""
}

/*
getNetworkTunnelId extract the id from a NetworkTunnel union type.
*/
func getNetworkTunnelId(tunnel perimeter81Sdk.NetworkTunnel) string {
	if tunnel.NetworkTunnelWireguard != nil {
		return tunnel.NetworkTunnelWireguard.Id
	}
	if tunnel.NetworkTunnelIpsecSingle != nil {
		return tunnel.NetworkTunnelIpsecSingle.Id
	}
	if tunnel.NetworkTunnelIpsecRedundant != nil {
		return tunnel.NetworkTunnelIpsecRedundant.Id
	}
	if tunnel.NetworkTunnelOpenvpn != nil {
		return tunnel.NetworkTunnelOpenvpn.Id
	}
	return ""
}

/*
getNetworkTunnelHaTunnelId extract the HaTunnelID from a NetworkTunnel union type (for redundant tunnels).
*/
func getNetworkTunnelHaTunnelId(tunnel perimeter81Sdk.NetworkTunnel) string {
	if tunnel.NetworkTunnelIpsecRedundant != nil {
		return tunnel.NetworkTunnelIpsecRedundant.HaTunnelID.Id
	}
	return ""
}

/*
getGatewayInfo get the gateway info
  - @param ctx context.Context - the context
  - @param networkId string - the network id
  - @param regionId string - the region id
  - @param client perimeter81Sdk.APIClient - the client
  - @param diags diag.Diagnostics - the diagnostics

@return string - the gateway id, the gateway dns, the gateway ip,  diag.Diagnostics - the diagnostics
*/
func getGatewayInfo(ctx context.Context, networkId string, regionId string, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (string, string, string, diag.Diagnostics) {
	network, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, networkId).Execute()
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to fetch network", err)
		return "", "", "", diags
	}
	// find the gateway id based on that least recently created gateway
	var gatewayId string
	var gatewayDns string
	var gatewayIp string
	for _, region := range network.Regions {
		if region.Id == regionId {
			latest := region.Instances[0].CreatedAt
			for _, gateway := range region.Instances {
				currentTime := gateway.CreatedAt
				gatewayId = gateway.Id
				if currentTime.After(latest) {
					latest = currentTime
					gatewayId = gateway.Id
					gatewayDns = gateway.Dns
					gatewayIp = gateway.Ip
				}
			}
		}
	}
	return gatewayId, gatewayDns, gatewayIp, diags
}

/*
getRedundantTunnelId get the redundant tunnel id
  - @param ctx context.Context - the context
  - @param networkId string - the network id
  - @param tunnelBody perimeter81Sdk.BaseTunnelValues - the tunnel body
  - @param client perimeter81Sdk.APIClient - the client
  - @param diags diag.Diagnostics - the diagnostics

@return string - the redundant tunnel id, diag.Diagnostics - the diagnostics
*/
func getRedundantTunnelId(ctx context.Context, networkId string, tunnelBody perimeter81Sdk.BaseTunnelValues, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (string, diag.Diagnostics) {
	network, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2NetworkFind(ctx, networkId).Execute()
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to fetch network", err)
		return "", diags
	}
	// Find the redundant tunnel by walking ALL gateways in the target region.
	// The wire response for the network-find endpoint does NOT include
	// haTunnelID per-tunnel — instead, redundant tunnel members are returned
	// as type="ipsec" with isHA=true and a per-tunnel id. The SDK's NetworkTunnel
	// union dispatcher falls back to NetworkTunnelBase for these (because the
	// NetworkTunnelIpsecRedundant schema requires haTunnelID which is absent
	// from this endpoint's response). We use the base tunnel's Id as the
	// haTunnelId — the API's GET /tunnels/ipsec/redundant/{id} accepts either
	// member of the pair and returns the full redundant tunnel pair.
	for _, region := range network.Regions {
		if region.Id != tunnelBody.RegionID {
			continue
		}
		for _, gateway := range region.Instances {
			for _, tunnel := range gateway.Tunnels {
				ifName := getNetworkTunnelInterfaceName(tunnel)
				if ifName != tunnelBody.TunnelName+"01" && ifName != tunnelBody.TunnelName+"02" {
					continue
				}
				// Prefer haTunnelID from the redundant-specific variant if present;
				// otherwise fall back to the base/single-routed tunnel id —
				// the API's redundant GET endpoint accepts either pair member's
				// id. The wire structure for redundant tunnel members is
				// identical to a single ipsec tunnel (type:"ipsec" + isHA:true),
				// so the SDK union dispatcher routes redundant pair members
				// into NetworkTunnelIpsecSingle.
				if id := getNetworkTunnelHaTunnelId(tunnel); id != "" {
					return id, diags
				}
				if tunnel.NetworkTunnelIpsecSingle != nil && tunnel.NetworkTunnelIpsecSingle.Id != "" {
					return tunnel.NetworkTunnelIpsecSingle.Id, diags
				}
				if tunnel.NetworkTunnelBase != nil && tunnel.NetworkTunnelBase.Id != "" {
					return tunnel.NetworkTunnelBase.Id, diags
				}
			}
		}
	}
	diags = appendErrorDiags(diags, "Unable to find tunnel",
		fmt.Errorf("no tunnel matched name=%s in region=%s; check tunnel fields or naming convention", tunnelBody.TunnelName, tunnelBody.RegionID))
	return "", diags
}

/*
setNetworkRegionInfos set the network region infos
  - @param regionsData []perimeter81Sdk.Region - the regions data
  - @param networkData *perimeter81Sdk.Network - the network data
  - @param regions []StandardNetworkRegionConfig - the regions

@return void
*/
func setNetworkRegionInfos(regionsData []perimeter81Sdk.Region, networkData *perimeter81Sdk.Network, regions []StandardNetworkRegionConfig) {
	newRegionsData := make([]StandardNetworkRegionConfig, 0)
	for _, networkRegions := range networkData.Regions {
		for _, regionData := range regionsData {
			if networkRegions.Name == regionData.GetDisplayName() {
				newRegionsData = append(newRegionsData, StandardNetworkRegionConfig{RegionID: networkRegions.Id, CpRegionId: regionData.GetId(), Dns: networkRegions.Dns, Name: networkRegions.Name})
			}
		}
	}
	for index, regionData := range regions {
		for _, networkRegions := range newRegionsData {
			if regionData.CpRegionId == networkRegions.CpRegionId {
				regions[index].RegionID = networkRegions.RegionID
				regions[index].Dns = networkRegions.Dns
				regions[index].Name = networkRegions.Name
			}
		}
	}
}

/*
checkNetworkStatus check the network status
  - @param ctx context.Context - the context
  - @param statusId string - the status id
  - @param client perimeter81Sdk.APIClient - the client
  - @param diags diag.Diagnostics - the diagnostics

@return perimeter81Sdk.AsyncOperationStatus, diag.Diagnostics, error - the network status, the diagnostics, the error
*/
func checkNetworkStatus(ctx context.Context, statusId string, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (perimeter81Sdk.AsyncOperationStatus, diag.Diagnostics, error) {
	networkStatus, _, err := client.StandardNetworksAPI.StandardNetworksControllerV2Status(ctx, statusId).Execute()
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to get Network Status", err)
		return perimeter81Sdk.AsyncOperationStatus{}, diags, err
	}
	if networkStatus.Result != nil && networkStatus.Result.StatusCode != nil && *networkStatus.Result.StatusCode == 500 {
		diags = appendErrorDiags(diags, "Unable to get Network Status", fmt.Errorf("%s", strings.Join(networkStatus.Result.Reason, " | ")))
		return *networkStatus, diags, fmt.Errorf("network status error")
	}
	return *networkStatus, diags, err
}

/*
addGatewayToRegion add the gateway to region
  - @param ctx context.Context - the context
  - @param client *perimeter81Sdk.APIClient - the client
  - @param gateways []GatewayConfig - the gateways
  - @param network_id string - the network id
  - @param region_id string - the region id
  - @param diags diag.Diagnostics - the diagnostics

@return diag.Diagnostics, error - the diagnostics, the error
*/
func addGatewayToRegion(ctx context.Context, client *perimeter81Sdk.APIClient, gateways []GatewayConfig, network_id string, region_id string, diags diag.Diagnostics) (diag.Diagnostics, error) {
	if len(gateways) == 0 {
		return diags, nil
	}
	for index, gateway := range gateways {
		gatewayPayload := perimeter81Sdk.CreateInstancesInNetworkPayload{
			RegionId: region_id,
			Idle:     gateway.Idle,
		}
		status, _, err := client.GatewaysAPI.StandardNetworksControllerV2AddNetworkInstance(ctx, network_id).CreateInstancesInNetworkPayload(gatewayPayload).Execute()
		if err != nil {
			diags = appendErrorDiags(diags, "Unable to create gateway", err)
			return diags, err
		}
		statusId := getIdFromUrl(status.GetStatusUrl())
		var gatewayId string
		var gatewayDns string
		var gatewayIp string
		for {
			var networkStatus perimeter81Sdk.AsyncOperationStatus
			networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
			if err != nil {
				return diags, err
			}
			if networkStatus.GetCompleted() {
				gatewayId, gatewayDns, gatewayIp, diags = getGatewayInfo(ctx, network_id, region_id, *client, diags)
				break
			}
			time.Sleep(60 * time.Second)
		}
		gateways[index].Id = gatewayId
		gateways[index].Dns = gatewayDns
		gateways[index].Ip = gatewayIp
	}
	return diags, nil
}

/*
deleteGatewayFromRegion delete the gateway from region
  - @param ctx context.Context - the context
  - @param client *perimeter81Sdk.APIClient - the client
  - @param gateways []GatewayConfig - the gateways
  - @param network_id string - the network id
  - @param region_id string - the region id
  - @param diags diag.Diagnostics - the diagnostics

@return diag.Diagnostics, error - the diagnostics, the error
*/
func deleteGatewayFromRegion(ctx context.Context, client *perimeter81Sdk.APIClient, gateways []GatewayConfig, network_id string, region_id string, diags diag.Diagnostics) (diag.Diagnostics, error) {
	if len(gateways) == 0 {
		return diags, nil
	}
	gatewaysForDelete := perimeter81Sdk.RemoveRegionInstance{
		Regions: []perimeter81Sdk.RemoveRegionPayload{
			{
				RegionId:  &region_id,
				Instances: []perimeter81Sdk.RemoveInstancePayload{},
			},
		},
	}

	for _, gateway := range gateways {
		id := gateway.Id
		gatewaysForDelete.Regions[0].Instances = append(gatewaysForDelete.Regions[0].Instances, perimeter81Sdk.RemoveInstancePayload{
			Id: &id,
		})
	}
	// DeleteNetworkInstance is synchronous — returns AsyncOperationResult (no status URL to poll)
	_, _, err := client.GatewaysAPI.StandardNetworksControllerV2DeleteNetworkInstance(ctx, network_id).RemoveRegionInstance(gatewaysForDelete).Execute()
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to delete gateways", err)
		return diags, err
	}
	return diags, nil
}

/*
getNewGateway get the new gateway
  - @param oldGateways []GatewayConfig - the old gateways
  - @param newGateways []GatewayConfig - the new gateways

@return []GatewayConfig - the new gateways
*/
func getNewGateway(oldGateways []GatewayConfig, newGateways []GatewayConfig) []GatewayConfig {
	var gateways []GatewayConfig
	for _, newGateway := range newGateways {
		if !gatewayExistsInArray(newGateway.Name, oldGateways) {
			gateways = append(gateways, newGateway)
		}
	}
	return gateways
}

/*
getGatewayToBeDeleted get the gateway to be deleted
  - @param oldGateways []GatewayConfig - the old gateways
  - @param newGateways []GatewayConfig - the new gateways

@return []GatewayConfig - the gateways
*/
func getGatewayToBeDeleted(oldGateways []GatewayConfig, newGateways []GatewayConfig) []GatewayConfig {
	var gateways []GatewayConfig
	for _, oldGateway := range oldGateways {
		if !gatewayExistsInArray(oldGateway.Name, newGateways) {
			gateways = append(gateways, oldGateway)
		}
	}
	return gateways
}

/*
appendErrorDiags append the error diagnostics
  - @param diags diag.Diagnostics - the diagnostics
  - @param summary string - the summary
  - @param err error - the error

@return diag.Diagnostics - the diagnostics
*/
func appendErrorDiags(diags diag.Diagnostics, summary string, err error) diag.Diagnostics {
	var errMsg string
	if apiErr, ok := err.(*perimeter81Sdk.GenericOpenAPIError); ok {
		errMsg = string(apiErr.Body())
		if errMsg == "" {
			errMsg = apiErr.Error()
		}
	} else {
		errMsg = err.Error()
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   errMsg,
	})
	return diags
}

/*
testComparableArraiesEq test if the arraies are equal
  - @param a []Type - the a
  - @param b []Type - the b

@return bool - the result
*/
func testComparableArraiesEq[Type comparable](a, b []Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

/*
randStringBytesRmndr generate random string

@return string - the random string
*/
func randStringBytesRmndr() string {
	str := make([]byte, 10)
	for i := range str {
		str[i] = letterBytes[seededRand.Intn(len(letterBytes))]
	}
	return string(str)
}

/*
regionExistsInArray check if region exists in array
  - @param regionId string - the region id
  - @param regions []StandardNetworkRegionConfig - the regions

@return bool - the result
*/
func regionExistsInArray(regionId string, regions []StandardNetworkRegionConfig) bool {
	for _, region := range regions {
		if region.CpRegionId == regionId {
			return true
		}
	}
	return false
}

/*
gatewayExistsInArray check if gateway exists in array
  - @param gateway_name string - the gateway name
  - @param gateways []GatewayConfig - the gateways

@return bool - the result
*/
func gatewayExistsInArray(gateway_name string, gateways []GatewayConfig) bool {
	for _, gateway := range gateways {
		if gateway.Name == gateway_name {
			return true
		}
	}
	return false
}

/*
checkGatewayDuplicatesInArray check if gateway duplicates in array
  - @param gateways []GatewayConfig - the gateways

@return bool - the result, string - the gateway name
*/
func checkGatewayDuplicatesInArray(gateways []GatewayConfig) (bool, string) {
	for _, gatewayToCheck := range gateways {

		var count int
		for _, currentGateway := range gateways {
			if gatewayToCheck.Name == currentGateway.Name {
				count++
			}
		}
		if count > 1 {
			return true, gatewayToCheck.Name
		}
	}

	return false, ""
}

/*
regionClonsInArray get the region clons in array
  - @param regionId string - the region id
  - @param regions []StandardNetworkRegionConfig - the regions

@return []StandardNetworkRegionConfig - the result
*/
func regionClonsInArray(regionId string, regions []StandardNetworkRegionConfig) []StandardNetworkRegionConfig {
	clons := make([]StandardNetworkRegionConfig, 0)
	for _, region := range regions {
		if region.CpRegionId == regionId {
			clons = append(clons, region)
		}
	}
	return clons
}

/*
importRegions import the manually added regions
  - @param networkData *perimeter81Sdk.Network - the network data
  - @param regionsData []perimeter81Sdk.Region - the regions date list
  - @param regions []StandardNetworkRegionConfig - the regions inside the configuration file if exists

@return []StandardNetworkRegionConfig - the result
*/
func importRegions(networkData *perimeter81Sdk.Network, regionsData []perimeter81Sdk.Region, regions []StandardNetworkRegionConfig) []StandardNetworkRegionConfig {
	if len(regions) == 0 {
		regions = make([]StandardNetworkRegionConfig, len(networkData.Regions))
		for i, regionItem := range networkData.Regions {
			region := StandardNetworkRegionConfig{}
			region.Idle = networkData.IsDefault
			region.RegionID = regionItem.Id
			region.Name = regionItem.Name
			region.Dns = regionItem.Dns
			for _, regionInfo := range regionsData {
				if regionInfo.GetDisplayName() == regionItem.Name {
					region.CpRegionId = regionInfo.GetId()
					break
				}
			}
			if region.CpRegionId == "" {
				for _, regionInfo := range regionsData {
					if regionInfo.GetName() == regionItem.Name {
						region.CpRegionId = regionInfo.GetId()
						break
					}
				}
			}
			regions[i] = region
		}
	}
	return regions
}

/*
getGatewaysInArray get the manually added gateways inside a specific region inside a given network
  - @param regionId string - the region id
  - @param network *perimeter81Sdk.Network - the network that has the gateways

@return []perimeter81Sdk.NetworkInstance - the result
*/
func getGatewaysInArray(regionId string, network *perimeter81Sdk.Network) []perimeter81Sdk.NetworkInstance {
	clons := make([]perimeter81Sdk.NetworkInstance, 0)

	for _, region := range network.Regions {
		if region.Id == regionId {
			clons = append(clons, region.Instances...)
			break
		}
	}
	return clons
}

/*
getCurrentObjectServicesInArray get the current object services from all the services by name.
The list API does not return IDs, so matching is done by name.
  - @param objectsServices *perimeter81Sdk.ObjectsServicesResponse - the objects services in the system
  - @param objectServicesName string - the object services name

@return *perimeter81Sdk.ObjectsServicesResponseObj - the result
*/
func getCurrentObjectServicesInArray(objectsServices *perimeter81Sdk.ObjectsServicesResponse, objectServicesName string) *perimeter81Sdk.ObjectsServicesResponseObj {
	for i, service := range objectsServices.Data {
		if service.Name == objectServicesName {
			return &objectsServices.Data[i]
		}
	}
	return nil
}

/*
getTunnelFromNetwork get the wireguard tunnel configs
  - @param tunnelId string - the tunnel id
  - @param network perimeter81Sdk.NetworkInstance - the network instance that has the configs

@return string,string - the result
*/
func getWireguardConfigsFromNetwork(tunnelId string, instances perimeter81Sdk.NetworkInstance) (string, string) {

	for _, tunnel := range instances.Tunnels {
		if tunnel.NetworkTunnelWireguard != nil && tunnel.NetworkTunnelWireguard.Id == tunnelId {
			return tunnel.NetworkTunnelWireguard.RequestConfigToken, tunnel.NetworkTunnelWireguard.Vault
		}
	}
	return "", ""
}

/*
getInstanceFromInstances get the gateway of gateways array
  - @param tunnelId string - the tunnel id
  - @param network []perimeter81Sdk.NetworkInstance - the network that has the gateways

@return perimeter81Sdk.NetworkInstance - the result
*/
func getInstanceFromInstances(gatewayId string, instances []perimeter81Sdk.NetworkInstance) *perimeter81Sdk.NetworkInstance {

	for _, instance := range instances {
		if instance.Id == gatewayId {
			return &instance
		}
	}
	return nil
}

/*
setDefaultGatewayIpForRegions set the default gateway ip for regions
  - @param regions []StandardNetworkRegionConfig - the region list
  - @param networkData *perimeter81Sdk.Network - the network data

@return []StandardNetworkRegionConfig - the result
*/
func setDefaultGatewayIpForRegions(regions []StandardNetworkRegionConfig, networkData *perimeter81Sdk.Network) []StandardNetworkRegionConfig {

	for index, region := range regions {
		gateways := getGatewaysInArray(region.RegionID, networkData)
		if len(gateways) > 0 {
			regions[index].DefaultGatewayIp = gateways[0].Ip
		}
	}
	return regions
}

/*
flattenObjectServicesDataSource flatten object Services data
  - @param objectServicesItems []perimeter81Sdk.ObjectsServicesResponseObj - the object services that need to be flattened

@return []interface{} - the flattened object services data
*/
func flattenObjectServicesData(objectServicesItems []perimeter81Sdk.ObjectsServicesResponseObj) []interface{} {
	if objectServicesItems != nil {
		objectServices := make([]interface{}, len(objectServicesItems))
		for i, objectServicesItem := range objectServicesItems {
			objectService := make(map[string]interface{})
			if objectServicesItem.Id != nil {
				objectService["id"] = *objectServicesItem.Id
			}
			objectService["name"] = objectServicesItem.Name
			if objectServicesItem.Description != nil {
				objectService["description"] = *objectServicesItem.Description
			}
			objectService["protocols"] = flattenProtocolsDataSourceData(objectServicesItem.Protocols)
			objectServices[i] = objectService
		}
		return objectServices
	}
	return make([]interface{}, 0)
}

/*
flattenProtocolsDataSourceData flatten protocols data
  - @param objectServicesItems []perimeter81Sdk.ObjectsServicesProtocolResponseObj - the object services that need to be flattened

@return []interface{} - the flattened object services data
*/
func flattenProtocolsDataSourceData(protocolItems []perimeter81Sdk.ObjectsServicesProtocolResponseObj) []interface{} {
	if protocolItems == nil {
		return make([]interface{}, 0)
	}
	protocols := make([]interface{}, len(protocolItems))
	for i, protocolItem := range protocolItems {
		protocols[i] = map[string]interface{}{
			"protocol":   protocolItem.Protocol,
			"value_type": protocolItem.ValueType,
			"value":      protocolItem.Value,
		}
	}
	return protocols
}

/*
getCurrentObjectAddressesInArray get the current object addresses from all the addresses
  - @param objectsAddresses perimeter81Sdk.ObjectsAddressesResponse - the objects addresses in the system
  - @param objectAddressesId string - the object addresses id

@return *perimeter81Sdk.ObjectsAddressObj - the result
*/
func getCurrentObjectAddressesInArray(objectsAddresses *perimeter81Sdk.ObjectsAddressesResponse, objectAddressesId string) *perimeter81Sdk.ObjectsAddressObj {
	for i, address := range objectsAddresses.Data {
		if address.GetId() == objectAddressesId {
			return &objectsAddresses.Data[i]
		}
	}
	return nil
}

/*
flattenObjectAddressesData flatten ObjectAddresses data
  - @param objectAddressesItems []perimeter81Sdk.ObjectsAddressObj - the object services that need to be flattened

@return []interface{} - the flattened object addressess data
*/
func flattenObjectAddressesData(objectAddressesItems []perimeter81Sdk.ObjectsAddressObj) []interface{} {
	if objectAddressesItems != nil {
		objectAddresses := make([]interface{}, len(objectAddressesItems))
		for i, objectAddressesItem := range objectAddressesItems {
			objectAddress := make(map[string]interface{})
			if objectAddressesItem.Id != nil {
				objectAddress["id"] = *objectAddressesItem.Id
			}
			objectAddress["name"] = objectAddressesItem.Name
			if objectAddressesItem.Description != nil {
				objectAddress["description"] = *objectAddressesItem.Description
			}
			objectAddress["value_type"] = objectAddressesItem.ValueType
			objectAddress["value"] = objectAddressesItem.Value
			objectAddresses[i] = objectAddress
		}
		return objectAddresses
	}
	return make([]interface{}, 0)
}
