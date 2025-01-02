package gollama

import "context"

func (c *Gollama) Version(ctx context.Context) (string, error) {
	var resp versionResponse
	c.apiGet(ctx, "/api/version", &resp)

	return resp.Version, nil
}
