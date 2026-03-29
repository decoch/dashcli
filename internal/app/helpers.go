package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/redash"
)

func (state *appState) apiClient() (*redash.Client, error) {
	if strings.TrimSpace(state.resolved.BaseURL) == "" {
		return nil, exitcode.Usagef("base URL is required; set --base-url, keyring base URL, or REDASH_BASE_URL")
	}
	if strings.TrimSpace(state.resolved.APIKey) == "" {
		return nil, exitcode.Usagef("API key is required; set --api-key, keyring API key, or REDASH_API_KEY")
	}

	client, err := redash.NewClient(state.resolved.BaseURL, state.resolved.APIKey, state.resolved.UserAgent, state.resolved.Timeout)
	if err != nil {
		return nil, exitcode.WrapUsage(err)
	}
	return client, nil
}

func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case float64:
		if typed == math.Trunc(typed) && !math.IsInf(typed, 0) && !math.IsNaN(typed) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		if typed == nil {
			return ""
		}
		return fmt.Sprintf("%v", typed)
	}
}

func asInt(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func extractJobObject(response map[string]any) map[string]any {
	if response == nil {
		return nil
	}
	if job, ok := response["job"].(map[string]any); ok {
		return job
	}
	return response
}
