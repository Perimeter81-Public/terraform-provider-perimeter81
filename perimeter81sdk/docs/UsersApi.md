# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateUser**](UsersApi.md#CreateUser) | **Post** /v2/users | Creates a new user
[**DeleteUser**](UsersApi.md#DeleteUser) | **Delete** /v2/users/{userId} | Remove User
[**GetUsers**](UsersApi.md#GetUsers) | **Get** /v2/users | Returns paginated list of users

# **CreateUser**
> User CreateUser(ctx, body)
Creates a new user

Required permissions: `[\"user:create\"]`<br><br>Creates a new user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateUserPayloadModel**](CreateUserPayloadModel.md)|  | 

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteUser**
> User DeleteUser(ctx, userId)
Remove User

Required permissions: `[\"user:delete\"]`<br><br>Marks a User as discontinued - user will remain in the system as inactive. The license of this user will be freed.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **userId** | **string**|  | 

### Return type

[**User**](User.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUsers**
> InlineResponse200 GetUsers(ctx, optional)
Returns paginated list of users

Required permissions: `[\"user:read\"]`<br><br>Returns paginated list of users.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiGetUsersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiGetUsersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int32**| Page number to start from | 
 **limit** | **optional.Int32**| Amount of users per page | 
 **q** | **optional.String**| Search string or URL encoded JSON | 
 **qType** | **optional.String**| Type of search. | [default to full]
 **qOperator** | **optional.String**| Applicable only if &#x60;qType &#x3D;&#x3D; partial&#x60;. | [default to or]

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

