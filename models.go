package gollama

import (
	"context"
	"errors"
	"fmt"
)

func (c *Gollama) ListModels(ctx context.Context) ([]ModelInfo, error) {
	type tagsResponse struct {
		Models []ModelInfo `json:"models"`
	}

	var r tagsResponse
	c.apiGet(ctx, "/api/tags", &r)

	return r.Models, nil
}

// HasModel checks if a given model is available on the server.
//
// The function will return an error if the request fails.
//
// The function will return false if the model is not found on the server.
func (c *Gollama) HasModel(ctx context.Context, model string) (bool, error) {
	models, err := c.ListModels(ctx)
	if err != nil {
		return false, err
	}

	for _, m := range models {
		if m.Model == model || m.Model == model+":latest" {
			return true, nil
		}
	}

	return false, nil
}

// ModelSize returns the size of a model on the server.
//
// The function will return an error if the model is not found.
//
// The function will return 0 if the model is not found.
func (c *Gollama) ModelSize(ctx context.Context, model string) (int, error) {
	models, err := c.ListModels(ctx)
	if err != nil {
		return 0, err
	}

	for _, m := range models {
		if m.Model == model || m.Model == model+":latest" {
			return m.Size, nil
		}
	}

	return 0, errors.New("model not found")
}

// PullModel pulls a model from the server if it is not available locally.
//
// The function will return an error if the request fails.
//
// The function will return an error if the model is not found on the server.
func (c *Gollama) PullModel(ctx context.Context, model string) error {
	req := pullRequest{
		Model:  model,
		Stream: false,
	}

	var resp pullResponse
	c.apiPost(ctx, "/api/pull", &resp, req)

	if resp.Status != "success" {
		return fmt.Errorf("failed to pull model %s", model)
	}
	return nil
}

// PullIfMissing pulls a model from the server if it is not available locally.
//
// The function will return an error if the request fails.
//
// The function will return an error if the model is not found on the server.
//
// If no model is specified, the model name set in the Gollama object is used.
func (c *Gollama) PullIfMissing(ctx context.Context, model ...string) error {
	if len(model) == 0 {
		model = []string{c.ModelName}
	}

	for _, m := range model {
		hasModel, err := c.HasModel(ctx, m)
		if err != nil {
			return err
		}

		if !hasModel {
			return c.PullModel(ctx, m)
		}
	}

	return nil
}

// GetDetails retrieves the details of specified models from the server.
//
// The function accepts a variadic parameter of model names. If no model names are provided,
// it defaults to using the model name set in the Gollama object.
//
// It returns a slice of ModelDetails for each requested model, or an error if the request fails.
func (c *Gollama) GetDetails(ctx context.Context, model ...string) ([]ModelDetails, error) {
	if len(model) == 0 {
		model = []string{c.ModelName}
	}

	ret := make([]ModelDetails, 0)

	for _, m := range model {
		req := showRequest{
			Model: m,
		}

		var resp ModelDetails
		err := c.apiPost(ctx, "/api/show", &resp, req)
		if err != nil {
			return nil, err
		}

		ret = append(ret, resp)
	}

	return ret, nil
}
