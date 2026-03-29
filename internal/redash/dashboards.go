package redash

import (
	"context"
	"fmt"
	"net/url"
)

func (client *Client) ListDashboards(ctx context.Context) ([]map[string]any, error) {
	return client.getList(ctx, "/api/dashboards")
}

func (client *Client) GetDashboard(ctx context.Context, slugOrID string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/dashboards/%s", url.PathEscape(slugOrID)))
}

