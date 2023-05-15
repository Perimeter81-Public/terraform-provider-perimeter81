# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateIPSecRedundantTunnel**](IPSecRedundantApi.md#CreateIPSecRedundantTunnel) | **Post** /networks/{networkId}/tunnels/ipsec/redundant | Create a new IPSec Redundant tunnel
[**DeleteIPSecRedundantTunnel**](IPSecRedundantApi.md#DeleteIPSecRedundantTunnel) | **Delete** /networks/{networkId}/tunnels/ipsec/redundant/{tunnelId} | Delete IPSec Redundant tunnel
[**GetIPSecRedundantTunnel**](IPSecRedundantApi.md#GetIPSecRedundantTunnel) | **Get** /networks/{networkId}/tunnels/ipsec/redundant/{tunnelId} | Get one IPSec Redundant tunnel
[**UpdateIPSecRedundantTunnel**](IPSecRedundantApi.md#UpdateIPSecRedundantTunnel) | **Put** /networks/{networkId}/tunnels/ipsec/redundant/{tunnelId} | Update IPSec Redundant Tunnel

# **CreateIPSecRedundantTunnel**
> AsyncOperationResponse CreateIPSecRedundantTunnel(ctx, body, networkId)
Create a new IPSec Redundant tunnel

Required permissions: `[\"network:write]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateIpSecRedundantPayload**](CreateIpSecRedundantPayload.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteIPSecRedundantTunnel**
> AsyncOperationResponse DeleteIPSecRedundantTunnel(ctx, networkId, tunnelId)
Delete IPSec Redundant tunnel

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

# **GetIPSecRedundantTunnel**
> IpSecRedundantTunnels GetIPSecRedundantTunnel(ctx, networkId, tunnelId)
Get one IPSec Redundant tunnel

Required permissions: `[\"network:read\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**IpSecRedundantTunnels**](IPSecRedundantTunnels.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateIPSecRedundantTunnel**
> AsyncOperationResponse UpdateIPSecRedundantTunnel(ctx, body, networkId, tunnelId)
Update IPSec Redundant Tunnel

Required permissions: `[\"network:write\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**UpdateIpSecRedundantPayload**](UpdateIpSecRedundantPayload.md)|  | 
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

