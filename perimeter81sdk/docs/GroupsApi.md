# {{classname}}

All URIs are relative to *https://virtserver.swaggerhub.com/perimeter81/public-api-yaml/1.0.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddGroupMember**](GroupsApi.md#AddGroupMember) | **Post** /v2/groups/{groupId}/member/{userId} | Add a member to a group.
[**CreateGroup**](GroupsApi.md#CreateGroup) | **Post** /v2/groups | Creates a new Group
[**GetGroups**](GroupsApi.md#GetGroups) | **Get** /v2/groups | Returns paginated list of groups
[**RemoveGroup**](GroupsApi.md#RemoveGroup) | **Delete** /v2/groups/{groupId} | Remove a group by ID
[**RemoveGroupMember**](GroupsApi.md#RemoveGroupMember) | **Delete** /v2/groups/{groupId}/member/{userId} | Remove a member from a group.

# **AddGroupMember**
> Group AddGroupMember(ctx, groupId, userId)
Add a member to a group.

Required permissions: `[\"group.member:create\"]`<br><br>Add a member to a group.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**|  | 
  **userId** | **string**|  | 

### Return type

[**Group**](Group.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateGroup**
> Group CreateGroup(ctx, body)
Creates a new Group

Required permissions: `[\"group:create\"]`<br><br>Creates a new group.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateGroupPayload**](CreateGroupPayload.md)|  | 

### Return type

[**Group**](Group.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetGroups**
> GroupList GetGroups(ctx, optional)
Returns paginated list of groups

Required permissions: `[\"group:read\"]`<br><br>Returns paginated list of groups.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GroupsApiGetGroupsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GroupsApiGetGroupsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int32**| Page to start from | 
 **limit** | **optional.Int32**| Amount of records per page | 

### Return type

[**GroupList**](GroupList.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RemoveGroup**
> Group RemoveGroup(ctx, groupId)
Remove a group by ID

Required permissions: `[\"group:delete\"]`<br><br>Remove a group by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**|  | 

### Return type

[**Group**](Group.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RemoveGroupMember**
> Group RemoveGroupMember(ctx, groupId, userId)
Remove a member from a group.

Required permissions: `[\"group.member:delete\"]`<br><br>Remove a member from a group.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **groupId** | **string**|  | 
  **userId** | **string**|  | 

### Return type

[**Group**](Group.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

