package gollama

import (
	"context"
	"fmt"
	"strings"
)

type ChatOption interface{}

// Chat generates a response to a prompt using the Ollama API.
//
// The first argument is the prompt to generate a response to.
//
// The function takes a variable number of options as arguments. The options are:
//   - A slice of strings representing the paths to images that should be passed as vision input.
//   - A slice of Tool objects representing the tools that should be available to the model.
//
// The function returns a pointer to a ChatOuput object, which contains the response to the prompt,
// as well as some additional information about the response. If an error occurs, the function
// returns nil and an error.
func (c *Gollama) Chat(ctx context.Context, prompt string, options ...ChatOption) (*ChatOuput, error) {
	var (
		temperature   float64
		seed          = c.SeedOrNegative
		contextLength = c.ContextLength
		promptImages  = []PromptImage{}
		tools         = []Tool{}
		format        = StructuredFormat{}
	)

	for _, option := range options {
		switch opt := option.(type) {
		case PromptImage:
			promptImages = append(promptImages, opt)
		case []PromptImage:
			promptImages = opt
		case Tool:
			tools = append(tools, opt)
		case []Tool:
			tools = opt
		case StructuredFormat:
			format = opt
		default:
			continue
		}
	}

	if seed < 0 {
		temperature = c.TemperatureIfNegativeSeed
	}

	messages := []chatMessage{}
	if c.SystemPrompt != "" {
		messages = append(messages, chatMessage{
			Role:    "system",
			Content: c.SystemPrompt,
		})
	}

	userMessage := chatMessage{
		Role:    "user",
		Content: prompt,
	}

	base64VisionImages := make([]string, 0)
	for _, image := range promptImages {
		base64image, err := base64EncodeFile(image.Filename)
		if err != nil {
			return nil, err
		}
		base64VisionImages = append(base64VisionImages, base64image)
	}

	if len(base64VisionImages) > 0 {
		userMessage.Images = base64VisionImages
	}

	messages = append(messages, userMessage)

	req := chatRequest{
		Stream:   false,
		Model:    c.ModelName,
		Messages: messages,
		Options: chatOptionsRequest{
			Seed:          seed,
			Temperature:   temperature,
			ContextLength: contextLength,
		},
	}

	if len(tools) > 0 {
		req.Tools = &tools
	}

	if len(format.Properties) > 0 {
		req.Format = &format
	}

	if c.ContextLength != 0 {
		req.Options.ContextLength = c.ContextLength
	}

	var resp chatResponse
	err := c.apiPost(ctx, "/api/chat", &resp, req)
	if err != nil {
		return nil, err
	}

	if resp.Model != c.ModelName {
		return nil, fmt.Errorf("model don't found")
	}

	out := &ChatOuput{
		Role:           resp.Message.Role,
		Content:        resp.Message.Content,
		ToolCalls:      resp.Message.ToolCalls,
		PromptTokens:   resp.PromptEvalCount,
		ResponseTokens: resp.EvalCount,
	}

	if c.TrimSpace {
		out.Content = strings.TrimSpace(out.Content)
	}

	return out, nil
}
