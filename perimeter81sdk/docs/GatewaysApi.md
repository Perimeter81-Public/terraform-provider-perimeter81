# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**NetworksControllerV2AddNetworkInstance**](GatewaysApi.md#NetworksControllerV2AddNetworkInstance) | **Post** /v2/networks/{networkId}/instances | Add gateway
[**NetworksControllerV2DeleteNetworkInstance**](GatewaysApi.md#NetworksControllerV2DeleteNetworkInstance) | **Delete** /v2/networks/{networkId}/instances | Remove Gateways from Network

# **NetworksControllerV2AddNetworkInstance**
> AsyncOperationResponse NetworksControllerV2AddNetworkInstance(ctx, body, networkId)
Add gateway

Required permissions: `[\"network:create\"]`<br><br>Add gateway.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateRegionInNetworkload**](CreateRegionInNetworkload.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2DeleteNetworkInstance**
> AsyncOperationResult NetworksControllerV2DeleteNetworkInstance(ctx, body, networkId)
Remove Gateways from Network

Required permissions: `[\"network:update\"]`<br><br>Remove Gateways from Network.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**RemoveRegionInstance**](RemoveRegionInstance.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResult**](AsyncOperationResult.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

