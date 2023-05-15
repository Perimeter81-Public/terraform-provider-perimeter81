# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOpenVPNTunnel**](OpenVPNApi.md#CreateOpenVPNTunnel) | **Post** /networks/{networkId}/tunnels/openvpn | Create a new OpenVPN tunnel
[**DeleteOpenVPNTunnel**](OpenVPNApi.md#DeleteOpenVPNTunnel) | **Delete** /networks/{networkId}/tunnels/openvpn/{tunnelId} | Delete OpenVPN tunnel
[**GetOpenVPNTunnel**](OpenVPNApi.md#GetOpenVPNTunnel) | **Get** /networks/{networkId}/tunnels/openvpn/{tunnelId} | Get one openVPN tunnel
[**UpdateOpenVPNTunnel**](OpenVPNApi.md#UpdateOpenVPNTunnel) | **Put** /networks/{networkId}/tunnels/openvpn/{tunnelId} | Update openVPN Tunnel

# **CreateOpenVPNTunnel**
> AsyncOperationResponse CreateOpenVPNTunnel(ctx, body, networkId)
Create a new OpenVPN tunnel

Required permissions: `[\"network:write]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BaseTunnelValues**](BaseTunnelValues.md)|  | 
  **networkId** | **string**|  | 

### Return type

[**AsyncOperationResponse**](AsyncOperationResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteOpenVPNTunnel**
> AsyncOperationResponse DeleteOpenVPNTunnel(ctx, networkId, tunnelId)
Delete OpenVPN tunnel

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

# **GetOpenVPNTunnel**
> OpenVpnTunnel GetOpenVPNTunnel(ctx, networkId, tunnelId)
Get one openVPN tunnel

Required permissions: `[\"network:read\"]`

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **networkId** | **string**|  | 
  **tunnelId** | **string**|  | 

### Return type

[**OpenVpnTunnel**](OpenVPNTunnel.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOpenVPNTunnel**
> AsyncOperationResponse UpdateOpenVPNTunnel(ctx, networkId, tunnelId)
Update openVPN Tunnel

Required permissions: `[\"network:write\"]`

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

