# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetNetworks**](NetworksApi.md#GetNetworks) | **Get** /v2/networks | Get all Networks
[**NetworksControllerV2NetworkCreate**](NetworksApi.md#NetworksControllerV2NetworkCreate) | **Post** /v2/networks | Create network
[**NetworksControllerV2NetworkDelete**](NetworksApi.md#NetworksControllerV2NetworkDelete) | **Delete** /v2/networks/{networkId} | Delete network
[**NetworksControllerV2NetworkFind**](NetworksApi.md#NetworksControllerV2NetworkFind) | **Get** /v2/networks/{networkId} | Get network by Id
[**NetworksControllerV2NetworkUpdate**](NetworksApi.md#NetworksControllerV2NetworkUpdate) | **Put** /v2/networks/{networkId} | Update network
[**NetworksControllerV2Status**](NetworksApi.md#NetworksControllerV2Status) | **Get** /v2/networks/status/{statusId} | Get status of asynchronous operations.

# **GetNetworks**
> Network GetNetworks(ctx, )
Get all Networks

List all networks

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**Network**](Network.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2NetworkCreate**
> AsyncOperationResponse NetworksControllerV2NetworkCreate(ctx, body)
Create network

Required permissions: `[\"network:create\"]`<br><br>Create networks.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**DeployNetworkPayload**](DeployNetworkPayload.md)|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2NetworkDelete**
> AsyncOperationResult NetworksControllerV2NetworkDelete(ctx, networkId)
Delete network

Required permissions: `[\"network:delete\"]`<br><br>Delete network.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResult**](AsyncOperationResult.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2NetworkFind**
> Network NetworksControllerV2NetworkFind(ctx, networkId)
Get network by Id

Required permissions: `[\"network:read\"]`<br><br>List network.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 

### Return type

[**Network**](Network.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2NetworkUpdate**
> Network NetworksControllerV2NetworkUpdate(ctx, body, networkId)
Update network

Required permissions: `[\"network:update\"]`<br><br>Update network.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**UpdateNetworkDto**](UpdateNetworkDto.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**Network**](Network.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **NetworksControllerV2Status**
> AsyncOperationStatus NetworksControllerV2Status(ctx, statusId)
Get status of asynchronous operations.

Required permissions: `[\"network:read\"]`<br><br> status of asynchronous operations.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **statusId** | **string**|  | 

### Return type

[**AsyncOperationStatus**](AsyncOperationStatus.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

