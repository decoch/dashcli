package redash

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (client *Client) ListQueries(ctx context.Context, page int, pageSize int, order string, search string) ([]map[string]any, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	params.Set("order", order)
	if strings.TrimSpace(search) != "" {
		params.Set("search", search)
	}
	return client.getListWithParams(ctx, "/api/queries", params)
}

func (client *Client) GetQuery(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/queries/%s", url.PathEscape(id)))
}

func (client *Client) RunQuery(ctx context.Context, id string) (map[string]any, error) {
	return client.postObject(ctx, fmt.Sprintf("/api/queries/%s/refresh", url.PathEscape(id)), map[string]any{})
}

func (client *Client) CreateQuery(ctx context.Context, name string, sql string, dataSourceID int, description string) (map[string]any, error) {
	body := map[string]any{
		"name":           name,
		"query":          sql,
		"data_source_id": dataSourceID,
		"description":    description,
	}
	return client.postObject(ctx, "/api/queries", body)
}

func (client *Client) UpdateQuery(ctx context.Context, id string, fields map[string]any) (map[string]any, error) {
	if fields == nil {
		fields = map[string]any{}
	}
	return client.postObject(ctx, fmt.Sprintf("/api/queries/%s", url.PathEscape(id)), fields)
}

func (client *Client) ArchiveQuery(ctx context.Context, id string) error {
	return client.deleteObject(ctx, fmt.Sprintf("/api/queries/%s", url.PathEscape(id)))
}

func (client *Client) GetQueryCachedResult(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/queries/%s/results", url.PathEscape(id)))
}
