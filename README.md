# gollama
Easy Ollama package for Golang

### Example use

> go get -u github.com/jonathanhecl/gollama


```go
package main

import (
	"context"
	"fmt"

	"github.com/jonathanhecl/gollama"
)

func main() {
	ctx := context.Background()
	g := gollama.New("llama3.2")
	g.Verbose = true
	if err := g.PullIfMissing(ctx); err != nil {
		fmt.Println("Error:", err)
		return
	}

	prompt := "what is the capital of Argentina?"

	type Capital struct {
		Capital string `required:"true"`
	}

	option, _ := gollama.StructToStructuredFormat(Capital{})

	fmt.Printf("Option: %+v\n", option)

	output, err := g.Chat(ctx, prompt, option)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\n%s\n", output.Content)
}
```

### Features

- Support Vision models
- Support Tools models
- Support Structured format
- Downloads model if missing
- Chat with model
- Generates embeddings with model
- Get model details

### Functions

- `New(model string) *Gollama` - Create a new Gollama
- `NewWithConfig(config Gollama) *Gollama` - Create a new Gollama with a pre-populated config
- `Chat(prompt string, ...ChatOption) (*gollama.ChatOutput, error)` - Generate a response
- `Embedding(prompt string) ([]float64, error)` - Generate embeddings
- `ListModels() ([]ModelInfo, error)` - List models available on ollama
- `HasModel(model string) (bool, error)` - Check if model is available
- `ModelSize(model string) (int, error)` - Get model size from ollama
- `PullModel(model string) error` - Pull model
- `PullIfMissing(model ...string) error` - Pull model if missing
- `GetDetails(model ...string) ([]ModelDetails, error)` - Get model details from ollama
- `Version() (string, error)` - Get ollama version
- `StructToStructuredFormat(interface{}) (StructuredFormat, error)` - Converts a Go struct to a Gollama structured format
- `CosenoSimilarity(vector1, vector2 []float64) float64` - Calculates the cosine similarity between two vectors
- output.`DecodeContent(output interface{}) error` - Decodes the content of a Gollama response
