# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**NetworksControllerV2AddNetworkRegion**](RegionsApi.md#NetworksControllerV2AddNetworkRegion) | **Put** /v2/networks/{networkId}/regions | Add regions to a network
[**NetworksControllerV2DeleteNetworkRegion**](RegionsApi.md#NetworksControllerV2DeleteNetworkRegion) | **Delete** /v2/networks/{networkId}/regions | Remove regions from network
[**NetworksControllerV2GetRegions**](RegionsApi.md#NetworksControllerV2GetRegions) | **Get** /v2/regions | List of available regions

# **NetworksControllerV2AddNetworkRegion**
> AsyncOperationResult NetworksControllerV2AddNetworkRegion(ctx, body, networkId)
Add regions to a network

Required permissions: `[\"Network:Update\"]`<br><br>Network Update.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateRegionPayload**](CreateRegionPayload.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResult**](AsyncOperationResult.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2DeleteNetworkRegion**
> CreateRegionPayload NetworksControllerV2DeleteNetworkRegion(ctx, body, networkId)
Remove regions from network

Required permissions: `[\"Network:Delete\"]`<br><br>Remove Region. Gateways will still be avaidble through the remaining regions, in case you removed the last region, the gateways will be removed as well.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RemoveRegionDto**](RemoveRegionDto.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**CreateRegionPayload**](CreateRegionPayload.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2GetRegions**
> RegionsList NetworksControllerV2GetRegions(ctx, )
List of available regions

Required permissions: `[\"addon:read\"]`<br><br>List of regions.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**RegionsList**](RegionsList.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

