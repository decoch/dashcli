package redash

import (
	"context"
	"fmt"
	"net/url"
)

func (client *Client) GetQueryResult(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/query_results/%s", url.PathEscape(id)))
}

func (client *Client) CreateQueryResult(ctx context.Context, dataSourceID int, query string, maxAge int) (map[string]any, error) {
	body := map[string]any{
		"query":          query,
		"data_source_id": dataSourceID,
		"max_age":        maxAge,
	}
	return client.postObject(ctx, "/api/query_results", body)
}
