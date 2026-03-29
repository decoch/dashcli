package redash

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (client *Client) ListDashboards(ctx context.Context, page int, pageSize int, order string, search string) ([]map[string]any, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("page_size", strconv.Itoa(pageSize))
	params.Set("order", order)
	if strings.TrimSpace(search) != "" {
		params.Set("search", search)
	}
	return client.getListWithParams(ctx, "/api/dashboards", params)
}

func (client *Client) GetDashboard(ctx context.Context, slugOrID string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/dashboards/%s", url.PathEscape(slugOrID)))
}
