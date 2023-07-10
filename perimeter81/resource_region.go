package perimeter81

import (
	"context"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceRegionCreate Create a Region inside a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param networkId string - the network id
  - @param oldRegions []perimeter81Sdk.CreateRegionInNetworkload - the old regions
  - @param newRegions []perimeter81Sdk.CreateRegionInNetworkload - the new regions
  - @param d *schema.ResourceData - the terraform resource data
  - @param client *perimeter81Sdk.APIClient - the api client

@return string, string, error
*/

func resourceRegionCreate(ctx context.Context, networkId string, oldRegions []perimeter81Sdk.CreateRegionInNetworkload, newRegions []perimeter81Sdk.CreateRegionInNetworkload, d *schema.ResourceData, client *perimeter81Sdk.APIClient) (string, string, error) {
	var status perimeter81Sdk.AsyncOperationResponse
	// Get the new regions that need to be created inside a network
	for _, newRegion := range newRegions {
		// If the region does not exist in the old regions, create it
		if !regionExistsInArray(newRegion.CpRegionId, oldRegions) {
			// Create the region inside the network
			statusPut, _, err := client.RegionsApi.NetworksControllerV2AddNetworkRegion(ctx, newRegion, networkId)
			// Get the status of the async operation and check for errors
			status = statusPut
			if err != nil {
				d.Partial(true)
				return "", newRegion.CpRegionId, err
			}
		}
	}
	// Get the status id from the status url
	statusId := getIdFromUrl(status.StatusUrl)
	return statusId, "", nil
}

/*
resourceRegionDelete Delete a Region from a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param networkId string - the network id
  - @param oldRegions []perimeter81Sdk.CreateRegionInNetworkload - the old regions
  - @param newRegions []perimeter81Sdk.CreateRegionInNetworkload - the new regions
  - @param d *schema.ResourceData - the terraform resource data
  - @param client *perimeter81Sdk.APIClient - the api client
  - @param oldStatusId string - the old status id

@return string, string, error
*/
func resourceRegionDelete(ctx context.Context, networkId string, oldRegions []perimeter81Sdk.CreateRegionInNetworkload, newRegions []perimeter81Sdk.CreateRegionInNetworkload, d *schema.ResourceData, client *perimeter81Sdk.APIClient, oldStatusId string) (string, string, error) {
	var status perimeter81Sdk.AsyncOperationResponse
	// Get the old regions that need to be deleted from the network
	for _, oldRegion := range oldRegions {
		// If the region does not exist in the new regions, delete it
		if !regionExistsInArray(oldRegion.CpRegionId, newRegions) {
			// Delete the region from the network
			statusDelete, _, err := client.RegionsApi.NetworksControllerV2DeleteNetworkRegion(ctx, perimeter81Sdk.RemoveRegionDto{RegionId: oldRegion.RegionID}, networkId)
			// Get the status of the async operation and check for errors
			status = statusDelete
			if err != nil {
				d.Partial(true)
				return "", oldRegion.RegionID, err
			}
		}
	}
	// Get the status id from the status url
	statusId := getIdFromUrl(status.StatusUrl)
	// If the status id is empty, return the old status id
	if statusId == "" {
		return oldStatusId, "", nil
	}
	return statusId, "", nil
}
