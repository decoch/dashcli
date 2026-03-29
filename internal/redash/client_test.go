package redash

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNewClient_InvalidBaseURL(t *testing.T) {
	t.Parallel()

	_, err := NewClient("", "test-key", "", time.Second)
	if err == nil {
		t.Fatal("NewClient() error = nil, want error")
	}
}

func TestListQueries_AuthorizationHeaderAndResults(t *testing.T) {
	t.Parallel()

	var gotAuth string
	var gotPath string
	var gotQuery string

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotAuth = request.Header.Get("Authorization")
		gotPath = request.URL.Path
		gotQuery = request.URL.RawQuery
		return jsonResponse(http.StatusOK, `{"results":[{"id":1,"name":"Q1"}]}`), nil
	})}

	queries, err := client.ListQueries(context.Background(), 1, 20, "-updated_at", "name")
	if err != nil {
		t.Fatalf("ListQueries() error = %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("len(ListQueries()) = %d, want %d", len(queries), 1)
	}
	if gotAuth != "Key test-key" {
		t.Fatalf("Authorization = %q, want %q", gotAuth, "Key test-key")
	}
	if gotPath != "/api/queries" {
		t.Fatalf("Path = %q, want %q", gotPath, "/api/queries")
	}
	if gotQuery != "order=-updated_at&page=1&page_size=20&search=name" {
		t.Fatalf("Query = %q, want %q", gotQuery, "order=-updated_at&page=1&page_size=20&search=name")
	}
}

func TestNewClient_UserAgentHeader(t *testing.T) {
	t.Parallel()

	var gotUserAgent string

	client, err := NewClient("https://redash.example.com", "test-key", "myapp/1.0", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotUserAgent = request.Header.Get("User-Agent")
		return jsonResponse(http.StatusOK, `{"results":[]}`), nil
	})}

	_, _ = client.ListQueries(context.Background(), 1, 20, "-updated_at", "")
	if gotUserAgent != "myapp/1.0" {
		t.Fatalf("User-Agent = %q, want %q", gotUserAgent, "myapp/1.0")
	}
}

func TestNewClient_EmptyUserAgent_NotSet(t *testing.T) {
	t.Parallel()

	var gotUserAgent string

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotUserAgent = request.Header.Get("User-Agent")
		return jsonResponse(http.StatusOK, `{"results":[]}`), nil
	})}

	_, _ = client.ListQueries(context.Background(), 1, 20, "-updated_at", "")
	if gotUserAgent != "" {
		t.Fatalf("User-Agent = %q, want empty", gotUserAgent)
	}
}

func TestArchiveQuery_UsesDelete(t *testing.T) {
	t.Parallel()

	var gotMethod string
	var gotPath string

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotMethod = request.Method
		gotPath = request.URL.Path
		return jsonResponse(http.StatusNoContent, ``), nil
	})}

	if err := client.ArchiveQuery(context.Background(), "123"); err != nil {
		t.Fatalf("ArchiveQuery() error = %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Fatalf("Method = %q, want %q", gotMethod, http.MethodDelete)
	}
	if gotPath != "/api/queries/123" {
		t.Fatalf("Path = %q, want %q", gotPath, "/api/queries/123")
	}
}

func TestGetQueryResult_Path(t *testing.T) {
	t.Parallel()

	var gotMethod string
	var gotPath string

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotMethod = request.Method
		gotPath = request.URL.Path
		return jsonResponse(http.StatusOK, `{"query_result":{"id":42}}`), nil
	})}

	if _, err := client.GetQueryResult(context.Background(), "42"); err != nil {
		t.Fatalf("GetQueryResult() error = %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("Method = %q, want %q", gotMethod, http.MethodGet)
	}
	if gotPath != "/api/query_results/42" {
		t.Fatalf("Path = %q, want %q", gotPath, "/api/query_results/42")
	}
}

func TestCreateQueryResult_RequestBody(t *testing.T) {
	t.Parallel()

	var gotPath string
	var requestBody map[string]any

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		gotPath = request.URL.Path
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		return jsonResponse(http.StatusOK, `{"query_result":{"id":1}}`), nil
	})}

	response, err := client.CreateQueryResult(context.Background(), 9, "select 1", 10)
	if err != nil {
		t.Fatalf("CreateQueryResult() error = %v", err)
	}
	if gotPath != "/api/query_results" {
		t.Fatalf("Path = %q, want %q", gotPath, "/api/query_results")
	}
	if requestBody["query"] != "select 1" {
		t.Fatalf("query = %v, want %v", requestBody["query"], "select 1")
	}
	if requestBody["data_source_id"] != float64(9) {
		t.Fatalf("data_source_id = %v, want %v", requestBody["data_source_id"], float64(9))
	}
	if requestBody["max_age"] != float64(10) {
		t.Fatalf("max_age = %v, want %v", requestBody["max_age"], float64(10))
	}
	if _, ok := response["query_result"]; !ok {
		t.Fatalf("response = %v, want query_result field", response)
	}
}

func TestAPIError_MessageParsing(t *testing.T) {
	t.Parallel()

	client, err := NewClient("https://redash.example.com", "test-key", "", time.Second)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.httpClient = &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"message":"invalid api key"}`), nil
	})}

	_, err = client.GetQuery(context.Background(), "1")
	if err == nil {
		t.Fatal("GetQuery() error = nil, want error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("GetQuery() error type = %T, want %T", err, &APIError{})
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusUnauthorized)
	}
	if apiErr.Message != "invalid api key" {
		t.Fatalf("Message = %q, want %q", apiErr.Message, "invalid api key")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
}
