package redash

import "context"

func (client *Client) ListDataSources(ctx context.Context) ([]map[string]any, error) {
	return client.getList(ctx, "/api/data_sources")
}

