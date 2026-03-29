package redash

import "context"

func (client *Client) Me(ctx context.Context) (map[string]any, error) {
	response, err := client.getObject(ctx, "/api/me")
	if err == nil {
		return response, nil
	}
	if !IsStatus(err, 404) {
		return nil, err
	}

	return client.getObject(ctx, "/api/users/me")
}
