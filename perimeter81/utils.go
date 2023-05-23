package perimeter81

import (
	"context"
	"fmt"
	"strings"
	perimeter81Sdk "terraform-provider-perimeter81/perimeter81sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func flattenStringsArrayData(strs []interface{}) []string {
	strsData := make([]string, len(strs))
	for index, str := range strs {
		strsData[index] = fmt.Sprint(str)
	}
	return strsData
}

func getIdFromUrl(url string) string {
	urlSplited := strings.Split(url, "/")
	return urlSplited[len(urlSplited)-1]
}

func flattenRegionsData(regionItems []interface{}) []perimeter81Sdk.CreateRegionInNetworkload {
	if regionItems != nil {
		regions := make([]perimeter81Sdk.CreateRegionInNetworkload, len(regionItems))

		for i, regionItem := range regionItems {
			region := perimeter81Sdk.CreateRegionInNetworkload{}

			region.CpRegionId = regionItem.(map[string]interface{})["cpregionid"].(string)
			region.InstanceCount = int32(regionItem.(map[string]interface{})["instancecount"].(int))
			region.Idle = regionItem.(map[string]interface{})["idle"].(bool)
			regions[i] = region
		}

		return regions
	}

	return make([]perimeter81Sdk.CreateRegionInNetworkload, 0)
}

func flattenNetworkData(networkItems []perimeter81Sdk.CreateNetworkPayload) []interface{} {
	if networkItems != nil {
		networks := make([]interface{}, len(networkItems))

		for i, networkItem := range networkItems {
			network := make(map[string]interface{})

			network["name"] = networkItem.Name
			network["tags"] = networkItem.Tags
			network["subnet"] = networkItem.Subnet
			networks[i] = network
		}

		return networks
	}

	return make([]interface{}, 0)
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

func checkNetworkStatus(ctx context.Context, statusId string, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (perimeter81Sdk.AsyncOperationStatus, diag.Diagnostics, error) {
	networkStatus, _, err := client.NetworksApi.NetworksControllerV2Status(ctx, statusId)
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to get Network Status", err)
	}
	return networkStatus, diags, err
}

func appendErrorDiags(diags diag.Diagnostics, summary string, err error) diag.Diagnostics {
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
	})
	return diags
}
