package gollama

import (
	"context"
	"errors"
	"fmt"
)

// SetModel sets the model to use for the Gollama object.
func (c *Gollama) SetModel(model string) {
	c.ModelName = model
}

// ListModels lists all available models on the server.
//
// The function will return an error if the request fails.
func (c *Gollama) ListModels(ctx context.Context) ([]ModelInfo, error) {
	type tagsResponse struct {
		Models []ModelInfo `json:"models"`
	}

	var r tagsResponse
	err := c.apiGet(ctx, "/api/tags", &r)
	if err != nil {
		return []ModelInfo{}, err
	}

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
	err := c.apiPost(ctx, "/api/pull", &resp, req)
	if err != nil {
		return err
	}

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

// GetModels retrieves a list of available models from the server.
//
// It returns a slice of strings containing model names, or an error if the request fails.
func (c *Gollama) GetModels(ctx context.Context) ([]string, error) {
	type detailsResponse struct {
		Family string `json:"family"`
	}

	type modelResponse struct {
		Name       string          `json:"model"`
		ModifiedAt string          `json:"modified_at"`
		Details    detailsResponse `json:"details"`
	}

	type tagsResponse struct {
		Models []modelResponse `json:"models"`
	}

	var resp tagsResponse
	err := c.apiGet(ctx, "/api/tags", &resp)
	if err != nil {
		return nil, err
	}

	res := []string{}

	for _, m := range resp.Models {
		res = append(res, m.Name)
	}

	return res, nil
}
