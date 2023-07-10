package perimeter81

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

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

			region["country_code"] = regionData.CountryCode
			region["continent_code"] = regionData.ContinentCode
			region["display_name"] = regionData.DisplayName
			region["name"] = regionData.Name
			region["class_name"] = regionData.ClassName
			region["object_id"] = regionData.ObjectId
			region["id"] = regionData.Id
			regions[i] = region
		}

		return regions
	}

	return make([]interface{}, 0)
}

/*
flattenRegionsData flatten regions data
  - @param regionItems []interface{} - the regions that need to be flattened

@return []perimeter81Sdk.CreateRegionInNetworkload - the flattened regions
*/
func flattenRegionsData(regionItems []interface{}) []perimeter81Sdk.CreateRegionInNetworkload {
	if regionItems != nil {
		regions := make([]perimeter81Sdk.CreateRegionInNetworkload, len(regionItems))

		for i, regionItem := range regionItems {
			region := perimeter81Sdk.CreateRegionInNetworkload{}

			region.CpRegionId = regionItem.(map[string]interface{})["cpregion_id"].(string)
			region.InstanceCount = int32(regionItem.(map[string]interface{})["instance_count"].(int))
			region.Idle = regionItem.(map[string]interface{})["idle"].(bool)
			region_id := regionItem.(map[string]interface{})["region_id"]
			if region_id != nil {
				region.RegionID = region_id.(string)
			}
			regions[i] = region
		}

		return regions
	}

	return make([]perimeter81Sdk.CreateRegionInNetworkload, 0)
}

/*
flattenGatewaysData flatten gateways data
  - @param gatewaysItems []interface{} - the gateways data that need to be flattened

@return []perimeter81Sdk.Gateway - the flattened gateways
*/
func flattenGatewaysData(gatewaysItems []interface{}) []perimeter81Sdk.Gateway {
	if gatewaysItems != nil {
		gateways := make([]perimeter81Sdk.Gateway, len(gatewaysItems))

		for i, gatewayItem := range gatewaysItems {
			gateway := perimeter81Sdk.Gateway{}

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

	return make([]perimeter81Sdk.Gateway, 0)
}

/*
flattenGatewaysData flatten gateways data
  - @param gatewaysItems []perimeter81Sdk.Gateway - the gateways that need to be flattened

@return []interface{} - the flattened gateways data
*/
func flattenGateways(gatewaysItems []perimeter81Sdk.Gateway) []interface{} {
	if gatewaysItems != nil {
		gateways := make([]interface{}, len(gatewaysItems))

		for i, gatewayItems := range gatewaysItems {
			gateway := make(map[string]interface{})

			gateway["name"] = gatewayItems.Name
			gateway["idle"] = gatewayItems.Idle
			gateway["id"] = gatewayItems.Id
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
			network["subnet"] = networkItem.Subnet
			networks[i] = network
		}

		return networks
	}

	return make([]interface{}, 0)
}

/*
flattenNetworkRegions flatten network regions
  - @param regionItems []perimeter81Sdk.CreateRegionInNetworkload - the network regions that need to be flattened

@return []interface{} - the flattened  network regions
*/
func flattenNetworkRegions(regionItems []perimeter81Sdk.CreateRegionInNetworkload) []interface{} {
	if regionItems != nil {
		regions := make([]interface{}, len(regionItems))

		for i, regionItem := range regionItems {
			region := make(map[string]interface{})

			region["cpregion_id"] = regionItem.CpRegionId
			region["region_id"] = regionItem.RegionID
			region["instance_count"] = regionItem.InstanceCount
			region["idle"] = regionItem.Idle
			regions[i] = region
		}

		return regions
	}

	return make([]interface{}, 0)
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
			network["subnet"] = networkItem.Subnet
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
			network["createdat"] = serverItem.CreatedAt
			network["updatedat"] = serverItem.UpdatedAt
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
			instance["createdat"] = instanceItem.CreatedAt
			instance["updatedat"] = instanceItem.UpdatedAt
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
	network, _, err := client.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkId)
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
						if tunnel.InterfaceName == tunnelBody.TunnelName {
							return tunnel.Id, diags
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
getGatewayId get the gateway id
  - @param ctx context.Context - the context
  - @param networkId string - the network id
  - @param regionId string - the region id
  - @param client perimeter81Sdk.APIClient - the client
  - @param diags diag.Diagnostics - the diagnostics

@return string - the gateway id, diag.Diagnostics - the diagnostics
*/
func getGatewayId(ctx context.Context, networkId string, regionId string, client perimeter81Sdk.APIClient, diags diag.Diagnostics) (string, diag.Diagnostics) {
	network, _, err := client.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkId)
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to fetch network", err)
		return "", diags
	}
	// find the gateway id based on that least recently created gateway
	var gatewayId string
	for _, region := range network.Regions {
		if region.Id == regionId {
			latest, _ := time.Parse("2006-01-02T15:04:05.000Z", region.Instances[0].CreatedAt)
			for _, gateway := range region.Instances {
				currentTime, _ := time.Parse("2006-01-02T15:04:05.000Z", gateway.CreatedAt)
				gatewayId = gateway.Id
				if currentTime.After(latest) {
					latest = currentTime
					gatewayId = gateway.Id
				}
			}
		}
	}
	return gatewayId, diags
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
	network, _, err := client.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkId)
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
						if tunnel.InterfaceName == tunnelBody.TunnelName+"01" || tunnel.InterfaceName == tunnelBody.TunnelName+"02" {
							return tunnel.HaTunnelID.Id, diags
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
	setNetworkRegionIds set the network region ids
	 - @param regionsData perimeter81Sdk.RegionsList - the regions data
	 - @param networkData perimeter81Sdk.Network - the network data
	 - @param regions []perimeter81Sdk.CreateRegionInNetworkload - the regions

@return void
*/
func setNetworkRegionIds(regionsData perimeter81Sdk.RegionsList, networkData perimeter81Sdk.Network, regions []perimeter81Sdk.CreateRegionInNetworkload) {
	newRegionsData := make([]perimeter81Sdk.CreateRegionInNetworkload, 0)
	for _, networkRegions := range networkData.Regions {
		for _, regionData := range regionsData.Regions {
			if networkRegions.Name == regionData.Name {
				newRegionsData = append(newRegionsData, perimeter81Sdk.CreateRegionInNetworkload{RegionID: networkRegions.Id, CpRegionId: regionData.Id})
			}
		}
	}
	for index, regionData := range regions {
		for _, networkRegions := range newRegionsData {
			if regionData.CpRegionId == networkRegions.CpRegionId {
				regions[index].RegionID = networkRegions.RegionID
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
	networkStatus, _, err := client.NetworksApi.NetworksControllerV2Status(ctx, statusId)
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to get Network Status", err)
	} else if networkStatus.Result != nil && networkStatus.Result.StatusCode == 500 {
		diags = appendErrorDiags(diags, "Unable to get Network Status", fmt.Errorf(strings.Join(networkStatus.Result.Reason, " | ")))
	}
	return networkStatus, diags, err
}

/*
addGatewayToRegion add the gateway to region
  - @param ctx context.Context - the context
  - @param client *perimeter81Sdk.APIClient - the client
  - @param gateways []perimeter81Sdk.Gateway - the gateways
  - @param network_id string - the network id
  - @param region_id string - the region id
  - @param diags diag.Diagnostics - the diagnostics

@return diag.Diagnostics, error - the diagnostics, the error
*/
func addGatewayToRegion(ctx context.Context, client *perimeter81Sdk.APIClient, gateways []perimeter81Sdk.Gateway, network_id string, region_id string, diags diag.Diagnostics) (diag.Diagnostics, error) {
	if len(gateways) == 0 {
		return diags, nil
	}
	for index, gateway := range gateways {
		gatewayPayload := perimeter81Sdk.CreateGatewayInRegionload{
			InstanceCount: 1,
			RegionID:      region_id,
			Idle:          gateway.Idle,
		}
		status, _, err := client.GatewaysApi.NetworksControllerV2AddNetworkInstance(ctx, gatewayPayload, network_id)
		if err != nil {
			diags = appendErrorDiags(diags, "Unable to create gateway", err)
			return diags, err
		}
		statusId := getIdFromUrl(status.StatusUrl)
		var gatewayId string
		for {
			var networkStatus perimeter81Sdk.AsyncOperationStatus
			networkStatus, diags, err = checkNetworkStatus(ctx, statusId, *client, diags)
			if err != nil {
				return diags, err
			}
			if networkStatus.Completed {
				gatewayId, diags = getGatewayId(ctx, network_id, region_id, *client, diags)
				break
			}
			time.Sleep(60 * time.Second)
		}
		gateways[index].Id = gatewayId
	}
	return diags, nil
}

/*
deleteGatewayFromRegion delete the gateway from region
  - @param ctx context.Context - the context
  - @param client *perimeter81Sdk.APIClient - the client
  - @param gateways []perimeter81Sdk.Gateway - the gateways
  - @param network_id string - the network id
  - @param region_id string - the region id
  - @param diags diag.Diagnostics - the diagnostics

@return diag.Diagnostics, error - the diagnostics, the error
*/
func deleteGatewayFromRegion(ctx context.Context, client *perimeter81Sdk.APIClient, gateways []perimeter81Sdk.Gateway, network_id string, region_id string, diags diag.Diagnostics) (diag.Diagnostics, error) {
	if len(gateways) == 0 {
		return diags, nil
	}
	gatewaysForDelete := perimeter81Sdk.RemoveRegionInstance{
		Regions: []perimeter81Sdk.RemoveRegionPayload{
			{
				RegionId:  region_id,
				Instances: []perimeter81Sdk.RemoveInstancePayload{},
			},
		},
	}

	for _, gateway := range gateways {
		gatewaysForDelete.Regions[0].Instances = append(gatewaysForDelete.Regions[0].Instances, perimeter81Sdk.RemoveInstancePayload{
			Id: gateway.Id,
		})
	}
	status, _, err := client.GatewaysApi.NetworksControllerV2DeleteNetworkInstance(ctx, gatewaysForDelete, network_id)
	if err != nil {
		diags = appendErrorDiags(diags, "Unable to delete gateways", err)
		return diags, err
	}
	statusId := getIdFromUrl(status.StatusUrl)
	for {
		networkStatus, diags, err := checkNetworkStatus(ctx, statusId, *client, diags)
		if err != nil {
			return diags, err
		}
		if networkStatus.Completed {
			break
		}
		time.Sleep(20 * time.Second)
	}
	return diags, nil
}

/*
getNewGateway get the new gateway
  - @param oldGateways []perimeter81Sdk.Gateway - the old gateways
  - @param newGateways []perimeter81Sdk.Gateway - the new gateways

@return []perimeter81Sdk.Gateway - the new gateways
*/
func getNewGateway(oldGateways []perimeter81Sdk.Gateway, newGateways []perimeter81Sdk.Gateway) []perimeter81Sdk.Gateway {
	var gateways []perimeter81Sdk.Gateway
	for _, newGateway := range newGateways {
		if !gatewayExistsInArray(newGateway.Name, oldGateways) {
			gateways = append(gateways, newGateway)
		}
	}
	return gateways
}

/*
getGatewayToBeDeleted get the gateway to be deleted
  - @param oldGateways []perimeter81Sdk.Gateway - the old gateways
  - @param newGateways []perimeter81Sdk.Gateway - the new gateways

@return []perimeter81Sdk.Gateway - the gateways
*/
func getGatewayToBeDeleted(oldGateways []perimeter81Sdk.Gateway, newGateways []perimeter81Sdk.Gateway) []perimeter81Sdk.Gateway {
	var gateways []perimeter81Sdk.Gateway
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

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
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
  - @param regions []perimeter81Sdk.CreateRegionInNetworkload - the regions

@return bool - the result
*/
func regionExistsInArray(regionId string, regions []perimeter81Sdk.CreateRegionInNetworkload) bool {
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
  - @param gateways []perimeter81Sdk.Gateway - the gateways

@return bool - the result
*/
func gatewayExistsInArray(gateway_name string, gateways []perimeter81Sdk.Gateway) bool {
	for _, gateway := range gateways {
		if gateway.Name == gateway_name {
			return true
		}
	}
	return false
}

/*
checkGatewayDuplicatesInArray check if gateway duplicates in array
  - @param gateways []perimeter81Sdk.Gateway - the gateways

@return bool - the result, string - the gateway name
*/
func checkGatewayDuplicatesInArray(gateways []perimeter81Sdk.Gateway) (bool, string) {
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
  - @param regions []perimeter81Sdk.CreateRegionInNetworkload - the regions

@return []perimeter81Sdk.CreateRegionInNetworkload - the result
*/
func regionClonsInArray(regionId string, regions []perimeter81Sdk.CreateRegionInNetworkload) []perimeter81Sdk.CreateRegionInNetworkload {
	clons := make([]perimeter81Sdk.CreateRegionInNetworkload, 0)
	for _, region := range regions {
		if region.CpRegionId == regionId {
			clons = append(clons, region)
		}
	}
	return clons
}
