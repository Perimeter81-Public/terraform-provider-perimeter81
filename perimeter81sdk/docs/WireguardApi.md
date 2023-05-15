# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateWireguardTunnel**](WireguardApi.md#CreateWireguardTunnel) | **Post** /networks/{networkId}/tunnels/wireguard | Create a new Wireguard tunnel
[**DeleteWireguardTunnel**](WireguardApi.md#DeleteWireguardTunnel) | **Delete** /networks/{networkId}/tunnels/wireguard/{tunnelId} | Delete Wireguard tunnel
[**GetWireguardTunnel**](WireguardApi.md#GetWireguardTunnel) | **Get** /networks/{networkId}/tunnels/wireguard/{tunnelId} | Get a Wireguard tunnel
[**UpdateWireguardTunnel**](WireguardApi.md#UpdateWireguardTunnel) | **Put** /networks/{networkId}/tunnels/wireguard/{tunnelId} | Update a Wireguard tunnel

# **CreateWireguardTunnel**
> AsyncOperationResponse CreateWireguardTunnel(ctx, body, networkId)
Create a new Wireguard tunnel

Required permissions: `[\"network:write\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateWireguardTunnelPayload**](CreateWireguardTunnelPayload.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteWireguardTunnel**
> AsyncOperationResponse DeleteWireguardTunnel(ctx, networkId, tunnelId)
Delete Wireguard tunnel

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

# **GetWireguardTunnel**
> WireguardTunnel GetWireguardTunnel(ctx, networkId, tunnelId)
Get a Wireguard tunnel

Required permissions: `[\"network:read\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**WireguardTunnel**](WireguardTunnel.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateWireguardTunnel**
> AsyncOperationResponse UpdateWireguardTunnel(ctx, body, networkId, tunnelId)
Update a Wireguard tunnel

Required permissions: `[\"network:write\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**WireGuradDetails**](WireGuradDetails.md)|  | 
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

