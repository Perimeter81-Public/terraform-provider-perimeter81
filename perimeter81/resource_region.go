package perimeter81

import (
	"context"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
resourceRegionCreate Create a Region inside a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param networkId string - the network id
  - @param oldRegions []StandardNetworkRegionConfig - the old regions
  - @param newRegions []StandardNetworkRegionConfig - the new regions
  - @param d *schema.ResourceData - the terraform resource data
  - @param client *perimeter81Sdk.APIClient - the api client

@return string, string, error
*/
func resourceRegionCreate(ctx context.Context, networkId string, oldRegions []StandardNetworkRegionConfig, newRegions []StandardNetworkRegionConfig, d *schema.ResourceData, client *perimeter81Sdk.APIClient) (string, string, error) {
	// Get the new regions that need to be created inside a network
	for _, newRegion := range newRegions {
		// If the region does not exist in the old regions, create it
		if !regionExistsInArray(newRegion.CpRegionId, oldRegions) {
			// Create the region inside the network using the standard region payload
			regionPayload := perimeter81Sdk.CreateRegionInNetworkPayload{
				HarmonySaseRegionId: newRegion.CpRegionId,
				Idle:                newRegion.Idle,
			}
			_, _, err := client.RegionsAPI.StandardNetworksControllerV2AddNetworkRegion(ctx, networkId).CreateRegionInNetworkPayload(regionPayload).Execute()
			if err != nil {
				d.Partial(true)
				return "", newRegion.CpRegionId, err
			}
		}
	}
	// Region add/delete operations are synchronous in the new SDK — no status polling needed.
	return "", "", nil
}

/*
resourceRegionDelete Delete a Region from a Network
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param networkId string - the network id
  - @param oldRegions []StandardNetworkRegionConfig - the old regions
  - @param newRegions []StandardNetworkRegionConfig - the new regions
  - @param d *schema.ResourceData - the terraform resource data
  - @param client *perimeter81Sdk.APIClient - the api client
  - @param oldStatusId string - the old status id (kept for API compatibility)

@return string, string, error
*/
func resourceRegionDelete(ctx context.Context, networkId string, oldRegions []StandardNetworkRegionConfig, newRegions []StandardNetworkRegionConfig, d *schema.ResourceData, client *perimeter81Sdk.APIClient, oldStatusId string) (string, string, error) {
	// Get the old regions that need to be deleted from the network
	for _, oldRegion := range oldRegions {
		// If the region does not exist in the new regions, delete it
		if !regionExistsInArray(oldRegion.CpRegionId, newRegions) {
			// Delete the region from the network
			_, _, err := client.RegionsAPI.StandardNetworksControllerV2DeleteNetworkRegion(ctx, networkId).RemoveRegionDTO(perimeter81Sdk.RemoveRegionDTO{RegionId: oldRegion.RegionID}).Execute()
			if err != nil {
				d.Partial(true)
				return "", oldRegion.RegionID, err
			}
		}
	}
	// Region add/delete operations are synchronous in the new SDK — no status polling needed.
	return oldStatusId, "", nil
}
