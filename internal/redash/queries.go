package redash

import (
	"context"
	"fmt"
	"net/url"
)

func (client *Client) ListQueries(ctx context.Context) ([]map[string]any, error) {
	return client.getList(ctx, "/api/queries")
}

func (client *Client) GetQuery(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/queries/%s", url.PathEscape(id)))
}

func (client *Client) RunQuery(ctx context.Context, id string) (map[string]any, error) {
	return client.postObject(ctx, fmt.Sprintf("/api/queries/%s/refresh", url.PathEscape(id)), map[string]any{})
}

