# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateIPSecSingleTunnel**](IPSecSingleApi.md#CreateIPSecSingleTunnel) | **Post** /networks/{networkId}/tunnels/ipsec/single | Create a new IPSec Single tunnel
[**DeleteIPSecSingleTunnel**](IPSecSingleApi.md#DeleteIPSecSingleTunnel) | **Delete** /networks/{networkId}/tunnels/ipsec/single/{tunnelId} | Delete IPSec Single tunnel
[**GetIPSecSingleTunnel**](IPSecSingleApi.md#GetIPSecSingleTunnel) | **Get** /networks/{networkId}/tunnels/ipsec/single/{tunnelId} | Get one IPSec Single tunnel
[**UpdateIPSecSingleTunnel**](IPSecSingleApi.md#UpdateIPSecSingleTunnel) | **Put** /networks/{networkId}/tunnels/ipsec/single/{tunnelId} | Update IPSec Single Tunnel

# **CreateIPSecSingleTunnel**
> AsyncOperationResponse CreateIPSecSingleTunnel(ctx, body, networkId)
Create a new IPSec Single tunnel

Required permissions: `[\"network:write]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateIpSecSinglePayload**](CreateIpSecSinglePayload.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteIPSecSingleTunnel**
> AsyncOperationResponse DeleteIPSecSingleTunnel(ctx, networkId, tunnelId)
Delete IPSec Single tunnel

Required permissions: `[\"network:delete\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetIPSecSingleTunnel**
> IpSecSingleTunnel GetIPSecSingleTunnel(ctx, networkId, tunnelId)
Get one IPSec Single tunnel

Required permissions: `[\"network:read\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**IpSecSingleTunnel**](IPSecSingleTunnel.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateIPSecSingleTunnel**
> AsyncOperationResponse UpdateIPSecSingleTunnel(ctx, body, networkId, tunnelId)
Update IPSec Single Tunnel

Required permissions: `[\"network:write\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**IpSecSingleDetails**](IpSecSingleDetails.md)|  | 
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

