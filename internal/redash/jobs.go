package redash

import (
	"context"
	"fmt"
	"net/url"
)

func (client *Client) GetJob(ctx context.Context, id string) (map[string]any, error) {
	return client.getObject(ctx, fmt.Sprintf("/api/jobs/%s", url.PathEscape(id)))
}
