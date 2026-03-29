package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/decoch/dashcli/internal/exitcode"
)

func TestQueryList_DefaultFlags(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var requestPath string
	var requestQuery map[string]string
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		requestPath = request.URL.Path
		requestQuery = map[string]string{
			"page":      request.URL.Query().Get("page"),
			"page_size": request.URL.Query().Get("page_size"),
			"order":     request.URL.Query().Get("order"),
			"search":    request.URL.Query().Get("search"),
		}
		return appJSONResponse(http.StatusOK, `{"results":[{"id":1,"name":"Q1"}]}`), nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"query", "list"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if requestPath != "/api/queries" {
		t.Fatalf("path = %q, want %q", requestPath, "/api/queries")
	}
	if requestQuery["page"] != "1" {
		t.Fatalf("page = %q, want %q", requestQuery["page"], "1")
	}
	if requestQuery["page_size"] != "20" {
		t.Fatalf("page_size = %q, want %q", requestQuery["page_size"], "20")
	}
	if requestQuery["order"] != "-updated_at" {
		t.Fatalf("order = %q, want %q", requestQuery["order"], "-updated_at")
	}
	if requestQuery["search"] != "" {
		t.Fatalf("search = %q, want empty", requestQuery["search"])
	}
	if got, want := stdout.String(), "1\tQ1\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestDashboardList_CustomFlags(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var requestPath string
	var requestQuery map[string]string
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		requestPath = request.URL.Path
		requestQuery = map[string]string{
			"page":      request.URL.Query().Get("page"),
			"page_size": request.URL.Query().Get("page_size"),
			"order":     request.URL.Query().Get("order"),
			"search":    request.URL.Query().Get("search"),
		}
		return appJSONResponse(http.StatusOK, `{"results":[{"id":2,"slug":"sales","name":"Sales"}]}`), nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"dashboard", "list", "--page", "3", "--page-size", "50", "--order", "name", "--search", "sales"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if requestPath != "/api/dashboards" {
		t.Fatalf("path = %q, want %q", requestPath, "/api/dashboards")
	}
	if requestQuery["page"] != "3" {
		t.Fatalf("page = %q, want %q", requestQuery["page"], "3")
	}
	if requestQuery["page_size"] != "50" {
		t.Fatalf("page_size = %q, want %q", requestQuery["page_size"], "50")
	}
	if requestQuery["order"] != "name" {
		t.Fatalf("order = %q, want %q", requestQuery["order"], "name")
	}
	if requestQuery["search"] != "sales" {
		t.Fatalf("search = %q, want %q", requestQuery["search"], "sales")
	}
	if got, want := stdout.String(), "2\tsales\tSales\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestQueryCreate_Success(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var requestMethod string
	var requestPath string
	var requestBody map[string]any
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		requestMethod = request.Method
		requestPath = request.URL.Path
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		return appJSONResponse(http.StatusOK, `{"id":101,"name":"Sample"}`), nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"query", "create", "--name", "Sample", "--sql", "select 1", "--datasource", "7", "--description", "desc"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if requestMethod != http.MethodPost {
		t.Fatalf("method = %q, want %q", requestMethod, http.MethodPost)
	}
	if requestPath != "/api/queries" {
		t.Fatalf("path = %q, want %q", requestPath, "/api/queries")
	}
	if requestBody["name"] != "Sample" {
		t.Fatalf("name = %v, want %v", requestBody["name"], "Sample")
	}
	if requestBody["query"] != "select 1" {
		t.Fatalf("query = %v, want %v", requestBody["query"], "select 1")
	}
	if requestBody["data_source_id"] != float64(7) {
		t.Fatalf("data_source_id = %v, want %v", requestBody["data_source_id"], float64(7))
	}
	if requestBody["description"] != "desc" {
		t.Fatalf("description = %v, want %v", requestBody["description"], "desc")
	}
	if got, want := stdout.String(), "Query created: id=101 name=Sample\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestQueryUpdate_ChangedFieldsOnly(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var requestPath string
	var requestBody map[string]any
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		requestPath = request.URL.Path
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		return appJSONResponse(http.StatusOK, `{"id":42}`), nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"query", "update", "42", "--name", "Renamed"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if requestPath != "/api/queries/42" {
		t.Fatalf("path = %q, want %q", requestPath, "/api/queries/42")
	}
	if len(requestBody) != 1 {
		t.Fatalf("len(body) = %d, want %d", len(requestBody), 1)
	}
	if requestBody["name"] != "Renamed" {
		t.Fatalf("name = %v, want %v", requestBody["name"], "Renamed")
	}
	if got, want := stdout.String(), "Query updated: id=42\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestQueryArchive_Success(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var requestMethod string
	var requestPath string
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		requestMethod = request.Method
		requestPath = request.URL.Path
		return appJSONResponse(http.StatusNoContent, ``), nil
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	code := Run(context.Background(), []string{"query", "archive", "42"}, stdout, stderr)

	if code != exitcode.CodeSuccess {
		t.Fatalf("Run() code = %d, want %d", code, exitcode.CodeSuccess)
	}
	if requestMethod != http.MethodDelete {
		t.Fatalf("method = %q, want %q", requestMethod, http.MethodDelete)
	}
	if requestPath != "/api/queries/42" {
		t.Fatalf("path = %q, want %q", requestPath, "/api/queries/42")
	}
	if got, want := stdout.String(), "Query archived: id=42\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestQueryResultsSQLRunAndSchema_JSON(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	var lastBody map[string]any
	stubDefaultTransport(t, func(request *http.Request) (*http.Response, error) {
		switch request.URL.Path {
		case "/api/queries/9/results.json":
			return appJSONResponse(http.StatusOK, `{"query_result":{"data":{"rows":[{"value":1}]}}}`), nil
		case "/api/query_results":
			if err := json.NewDecoder(request.Body).Decode(&lastBody); err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			return appJSONResponse(http.StatusOK, `{"query_result":{"id":77}}`), nil
		case "/api/data_sources/3/schema":
			return appJSONResponse(http.StatusOK, `{"schema":[{"name":"id"}]}`), nil
		default:
			return appJSONResponse(http.StatusNotFound, `{}`), nil
		}
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	resultsCode := Run(context.Background(), []string{"--json", "query", "results", "9"}, stdout, stderr)
	if resultsCode != exitcode.CodeSuccess {
		t.Fatalf("query results code = %d, want %d", resultsCode, exitcode.CodeSuccess)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("\"query_result\"")) {
		t.Fatalf("query results stdout = %q, want JSON response", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	sqlCode := Run(context.Background(), []string{"--json", "sql", "run", "--datasource", "11", "--query", "select 1", "--max-age", "60"}, stdout, stderr)
	if sqlCode != exitcode.CodeSuccess {
		t.Fatalf("sql run code = %d, want %d", sqlCode, exitcode.CodeSuccess)
	}
	if lastBody["data_source_id"] != float64(11) {
		t.Fatalf("data_source_id = %v, want %v", lastBody["data_source_id"], float64(11))
	}
	if lastBody["query"] != "select 1" {
		t.Fatalf("query = %v, want %v", lastBody["query"], "select 1")
	}
	if lastBody["max_age"] != float64(60) {
		t.Fatalf("max_age = %v, want %v", lastBody["max_age"], float64(60))
	}

	stdout.Reset()
	stderr.Reset()
	schemaCode := Run(context.Background(), []string{"--json", "datasource", "schema", "3"}, stdout, stderr)
	if schemaCode != exitcode.CodeSuccess {
		t.Fatalf("datasource schema code = %d, want %d", schemaCode, exitcode.CodeSuccess)
	}
	if !bytes.Contains(stdout.Bytes(), []byte("\"schema\"")) {
		t.Fatalf("datasource schema stdout = %q, want JSON response", stdout.String())
	}
}

func TestCreateAndSQLRun_RequiredFlags(t *testing.T) {
	stubCredentials(t, "https://redash.example.com")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCode := Run(context.Background(), []string{"query", "create", "--sql", "select 1", "--datasource", "1"}, stdout, stderr)
	if createCode == exitcode.CodeSuccess {
		t.Fatalf("query create code = %d, want non-success", createCode)
	}

	stdout.Reset()
	stderr.Reset()
	sqlCode := Run(context.Background(), []string{"sql", "run", "--datasource", "1"}, stdout, stderr)
	if sqlCode == exitcode.CodeSuccess {
		t.Fatalf("sql run code = %d, want non-success", sqlCode)
	}
}

func stubCredentials(t *testing.T, baseURL string) {
	t.Helper()

	previousGetBaseURL := getBaseURL
	previousGetAPIKey := getAPIKey
	previousLookupEnv := lookupEnv
	getBaseURL = func() (string, error) {
		return baseURL, nil
	}
	getAPIKey = func() (string, error) {
		return "test-api-key", nil
	}
	lookupEnv = func(string) (string, bool) {
		return "", false
	}

	t.Cleanup(func() {
		getBaseURL = previousGetBaseURL
		getAPIKey = previousGetAPIKey
		lookupEnv = previousLookupEnv
	})
}

func stubDefaultTransport(t *testing.T, fn appRoundTripFunc) {
	t.Helper()

	previousTransport := http.DefaultTransport
	http.DefaultTransport = fn
	t.Cleanup(func() {
		http.DefaultTransport = previousTransport
	})
}

type appRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn appRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func appJSONResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
}
