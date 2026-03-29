package redash

import (
	"context"
	"fmt"
	"net/url"
)

func (client *Client) ListDataSources(ctx context.Context) ([]map[string]any, error) {
	return client.getList(ctx, "/api/data_sources")
}

func (client *Client) GetDataSourceSchema(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/data_sources/%s/schema", url.PathEscape(id)))
}
