package checkpointsase

// BUG-17 workaround: bypass the SDK's broken protocols deserialization for
// object_services by issuing a raw GET against /v2.3/objects/services and
// parsing the flat wire shape directly. See TEST-PLAN.md BUG-17 for the
// root cause (swagger's nested oneOf/anyOf vs the flat wire payload).
//
// When the swagger is fixed and the SDK is regenerated, this file becomes
// dead code and the callers in resource_object_services.go +
// data_source_object_services.go can revert to SDK-only Reads.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"
)

// rawObjectServiceProtocol is the flat per-protocol shape the public-api
// actually returns (see perimeter81-public-api ServiceGetResponse interceptor).
type rawObjectServiceProtocol struct {
	Protocol  string `json:"protocol"`
	ValueType string `json:"valueType,omitempty"`
	Value     []int  `json:"value,omitempty"`
}

// rawObjectService is the per-entry shape returned by GET /v2.3/objects/services.
type rawObjectService struct {
	Id          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description,omitempty"`
	Protocols   []rawObjectServiceProtocol `json:"protocols"`
}

type rawObjectServicesResponse struct {
	Data []rawObjectService `json:"data"`
}

// fetchRawObjectServices does an authenticated GET against
// /v2.3/objects/services and returns the parsed flat response.
//
// Authentication: re-uses the SDK's exported GetBearerTokenFromApiKey helper
// against the api_key + base_url that the provider was configured with. We
// could cache the bearer to avoid one auth round-trip per Read, but the SDK
// already caches it internally for its own requests; doing one extra auth
// per object_service Read is acceptable for a workaround. If this ever
// becomes a hotspot, swap in a local cache mirroring SDK behavior.
func fetchRawObjectServices(ctx context.Context, client *perimeter81Sdk.APIClient) ([]rawObjectService, error) {
	providerAuthInfo.RLock()
	apiKey := providerAuthInfo.apiKey
	baseURL := providerAuthInfo.baseURL
	providerAuthInfo.RUnlock()

	if apiKey == "" || baseURL == "" {
		return nil, fmt.Errorf("provider auth info not initialized; providerConfigure must run before raw fetch")
	}

	token, err := client.GetBearerTokenFromApiKey(apiKey, baseURL)
	if err != nil {
		return nil, fmt.Errorf("BUG-17 workaround: bearer token exchange failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v2.3/objects/services", nil)
	if err != nil {
		return nil, fmt.Errorf("BUG-17 workaround: building GET request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BUG-17 workaround: GET /v2.3/objects/services failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("BUG-17 workaround: reading response body failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("BUG-17 workaround: GET /v2.3/objects/services returned %d: %s", resp.StatusCode, string(body))
	}

	var parsed rawObjectServicesResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("BUG-17 workaround: unmarshalling response failed: %w", err)
	}
	return parsed.Data, nil
}

// rawProtocolsToTerraform converts the flat raw protocols slice into the
// []interface{} shape terraform's d.Set expects for the "protocols" field.
func rawProtocolsToTerraform(protocols []rawObjectServiceProtocol) []interface{} {
	out := make([]interface{}, 0, len(protocols))
	for _, p := range protocols {
		entry := map[string]interface{}{
			"protocol":   p.Protocol,
			"value_type": p.ValueType,
		}
		value := make([]interface{}, 0, len(p.Value))
		for _, v := range p.Value {
			value = append(value, v)
		}
		entry["value"] = value
		out = append(out, entry)
	}
	return out
}

// hclProtocolsToRaw converts the []interface{} that terraform stores under
// the "protocols" attribute back into the flat wire shape the public-api
// expects. Mirror image of rawProtocolsToTerraform.
func hclProtocolsToRaw(protocolItems []interface{}) []rawObjectServiceProtocol {
	out := make([]rawObjectServiceProtocol, 0, len(protocolItems))
	for _, item := range protocolItems {
		m := item.(map[string]interface{})
		entry := rawObjectServiceProtocol{
			Protocol:  asString(m["protocol"]),
			ValueType: asString(m["value_type"]),
		}
		if vals, ok := m["value"].([]interface{}); ok {
			for _, v := range vals {
				entry.Value = append(entry.Value, asInt(v))
			}
		}
		out = append(out, entry)
	}
	return out
}

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func asInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

// rawObjectServiceCreateRequest is the body shape the public-api expects
// for POST /v2.3/objects/services (mirrors the actual flat wire payload).
type rawObjectServiceCreateRequest struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description,omitempty"`
	Protocols   []rawObjectServiceProtocol `json:"protocols"`
}

// createRawObjectService issues a POST against /v2.3/objects/services with
// the flat protocols shape (BUG-17 workaround for the CREATE direction).
// Returns the id assigned by the server.
func createRawObjectService(ctx context.Context, client *perimeter81Sdk.APIClient, name, description string, protocols []rawObjectServiceProtocol) (string, error) {
	id, err := mutatingObjectServiceRequest(ctx, client, http.MethodPost, "/v2.3/objects/services", rawObjectServiceCreateRequest{
		Name:        name,
		Description: description,
		Protocols:   protocols,
	})
	return id, err
}

// updateRawObjectService issues a PUT against /v2.3/objects/services/{id}
// (BUG-17 workaround for the UPDATE direction).
func updateRawObjectService(ctx context.Context, client *perimeter81Sdk.APIClient, id, name, description string, protocols []rawObjectServiceProtocol) error {
	_, err := mutatingObjectServiceRequest(ctx, client, http.MethodPut, "/v2.3/objects/services/"+id, rawObjectServiceCreateRequest{
		Name:        name,
		Description: description,
		Protocols:   protocols,
	})
	return err
}

// postRawAsync issues an authenticated POST against an async public-api
// endpoint (one that returns `{ statusUrl, ... }`). Used by BUG-23/24
// workarounds where the SDK-generated Create types don't match the actual
// wire shape (stale swagger). Returns the `statusUrl` so the caller can
// poll via the existing `checkNetworkStatus` helper.
func postRawAsync(ctx context.Context, client *perimeter81Sdk.APIClient, path string, body interface{}) (string, error) {
	providerAuthInfo.RLock()
	apiKey := providerAuthInfo.apiKey
	baseURL := providerAuthInfo.baseURL
	providerAuthInfo.RUnlock()

	if apiKey == "" || baseURL == "" {
		return "", fmt.Errorf("provider auth info not initialized")
	}

	token, err := client.GetBearerTokenFromApiKey(apiKey, baseURL)
	if err != nil {
		return "", fmt.Errorf("postRawAsync: bearer token exchange failed: %w", err)
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("postRawAsync: marshalling request body failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("postRawAsync: building POST %s request failed: %w", path, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("postRawAsync: POST %s failed: %w", path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("postRawAsync: POST %s returned %d: %s", path, resp.StatusCode, string(respBody))
	}

	var parsed struct {
		StatusUrl string `json:"statusUrl"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", fmt.Errorf("postRawAsync: parsing async response failed: %w (body=%s)", err, string(respBody))
	}
	if parsed.StatusUrl == "" {
		return "", fmt.Errorf("postRawAsync: async response missing statusUrl (body=%s)", string(respBody))
	}
	return parsed.StatusUrl, nil
}

func mutatingObjectServiceRequest(ctx context.Context, client *perimeter81Sdk.APIClient, method, path string, body interface{}) (string, error) {
	providerAuthInfo.RLock()
	apiKey := providerAuthInfo.apiKey
	baseURL := providerAuthInfo.baseURL
	providerAuthInfo.RUnlock()

	if apiKey == "" || baseURL == "" {
		return "", fmt.Errorf("provider auth info not initialized")
	}

	token, err := client.GetBearerTokenFromApiKey(apiKey, baseURL)
	if err != nil {
		return "", fmt.Errorf("BUG-17 workaround: bearer token exchange failed: %w", err)
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("BUG-17 workaround: marshalling request body failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("BUG-17 workaround: building %s %s request failed: %w", method, path, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("BUG-17 workaround: %s %s failed: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("BUG-17 workaround: %s %s returned %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	// POST returns the created entity (we need its id); PUT may return 204 or the entity.
	if len(respBody) == 0 {
		return "", nil
	}
	var parsed struct {
		Id string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		// Some responses wrap data; ignore parse failure here since callers can
		// fall back to Read.
		return "", nil
	}
	return parsed.Id, nil
}
