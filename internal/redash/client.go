package redash

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, timeout time.Duration) (*Client, error) {
	trimmedBaseURL := strings.TrimSpace(baseURL)
	if trimmedBaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	parsed, err := url.Parse(trimmedBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid base URL: %q", baseURL)
	}
	if strings.EqualFold(parsed.Scheme, "http") {
		return nil, fmt.Errorf("base URL uses http:// — API key will be sent unencrypted; use https://")
	}

	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	return &Client{
		baseURL: strings.TrimRight(trimmedBaseURL, "/"),
		apiKey:  strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (client *Client) getObject(ctx context.Context, path string) (map[string]any, error) {
	var response map[string]any
	if err := client.doJSON(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (client *Client) postObject(ctx context.Context, path string, body any) (map[string]any, error) {
	var response map[string]any
	if err := client.doJSON(ctx, http.MethodPost, path, body, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (client *Client) deleteObject(ctx context.Context, path string) error {
	return client.doJSON(ctx, http.MethodDelete, path, nil, nil)
}

func (client *Client) getList(ctx context.Context, path string) ([]map[string]any, error) {
	var raw any
	if err := client.doJSON(ctx, http.MethodGet, path, nil, &raw); err != nil {
		return nil, err
	}
	return normalizeList(raw), nil
}

func (client *Client) getListWithParams(ctx context.Context, path string, params url.Values) ([]map[string]any, error) {
	var raw any
	if err := client.doJSONWithParams(ctx, http.MethodGet, path, params, nil, &raw); err != nil {
		return nil, err
	}
	return normalizeList(raw), nil
}

func (client *Client) doJSON(ctx context.Context, method, path string, requestBody any, responseBody any) error {
	return client.doJSONWithParams(ctx, method, path, nil, requestBody, responseBody)
}

func (client *Client) doJSONWithParams(ctx context.Context, method, path string, params url.Values, requestBody any, responseBody any) error {
	requestURL, err := client.buildRequestURL(path, params)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if requestBody != nil {
		encoded, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if client.apiKey != "" {
		request.Header.Set("Authorization", "Key "+client.apiKey)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return parseAPIError(response.StatusCode, payload)
	}
	if responseBody == nil || len(payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(payload, responseBody); err != nil {
		return fmt.Errorf("parse response body: %w", err)
	}
	return nil
}

func (client *Client) buildRequestURL(path string, params url.Values) (string, error) {
	requestURL, err := url.JoinPath(client.baseURL, strings.TrimLeft(path, "/"))
	if err != nil {
		return "", fmt.Errorf("build request URL: %w", err)
	}
	if len(params) == 0 {
		return requestURL, nil
	}

	parsed, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("parse request URL: %w", err)
	}
	parsed.RawQuery = params.Encode()
	return parsed.String(), nil
}

func parseAPIError(statusCode int, payload []byte) error {
	message := ""
	bodyText := strings.TrimSpace(string(payload))

	if len(payload) > 0 {
		var parsed map[string]any
		if err := json.Unmarshal(payload, &parsed); err == nil {
			for _, key := range []string{"message", "error", "detail"} {
				if value, ok := parsed[key].(string); ok && strings.TrimSpace(value) != "" {
					message = strings.TrimSpace(value)
					break
				}
			}
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Body:       bodyText,
	}
}

func normalizeList(raw any) []map[string]any {
	if raw == nil {
		return nil
	}

	if list, ok := raw.([]any); ok {
		return toObjectList(list)
	}

	if object, ok := raw.(map[string]any); ok {
		for _, key := range []string{"results", "items", "data"} {
			if nested, ok := object[key].([]any); ok {
				return toObjectList(nested)
			}
		}
	}

	return nil
}

func toObjectList(raw []any) []map[string]any {
	objects := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		object, ok := item.(map[string]any)
		if !ok {
			continue
		}
		objects = append(objects, object)
	}
	return objects
}
