package gollama

import "context"

func (c *Gollama) Version(ctx context.Context) (string, error) {
	var resp versionResponse
	err := c.apiGet(ctx, "/api/version", &resp)
	if err != nil {
		return "", err
	}

	return resp.Version, nil
}
