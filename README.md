# gollama

**The easiest way to integrate Ollama into your Go applications.**

`gollama` provides a simple, idiomatic Go wrapper for the [Ollama](https://ollama.com/) API, enabling you to build powerful AI applications with local LLMs. It supports advanced features like Structured Outputs, Vision, Function Calling, and the new **Model Context Protocol (MCP)**.

## 🚀 Features

- **Model Context Protocol (MCP)**: First-class support for connecting to external MCP servers (Filesystem, PostgreSQL, Supabase, Brave Search, etc.).
- **Structured Outputs**: Automatically convert Go structs into JSON schemas for type-safe LLM responses.
- **Vision Support**: Easily pass images to multimodal models like `llama3.2-vision`.
- **Function Calling**: Define tools and let the model decide when to execute them.
- **Auto-Management**: Automatically checks for and pulls models if they are missing.
- **Embeddings**: Generate vector embeddings for RAG applications.

## 📦 Installation

```bash
go get -u github.com/jonathanhecl/gollama
```

## 💡 Usage Examples

### 1. Basic Chat
The simplest way to interact with a model.

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

	// Automatically pull the model if not present
	if err := g.PullIfMissing(ctx); err != nil {
		panic(err)
	}

	response, err := g.Chat(ctx, "Why is the sky blue?")
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Content)
}
```

### 2. Model Context Protocol (MCP) 🌟
Connect your LLM to the outside world using the standard Model Context Protocol. This allows you to use pre-built MCP servers without writing custom tool logic.

**Example: Connecting to Supabase or Filesystem**

```go
// Configure the MCP Client
// Example: Using the Supabase MCP Server
config := gollama.McpConfig{
    Command: "npx",
    Args: []string{
        "-y",
        "@supabase/mcp-server-supabase@latest",
        "--access-token", "sbp_your_token_here",
    },
    Env: map[string]string{
        "SUPABASE_URL": "https://your-project.supabase.co",
    },
}

// Or use the Filesystem server
// config := gollama.McpConfig{
//     Command: "npx",
//     Args: []string{"-y", "@modelcontextprotocol/server-filesystem", "."},
// }

client := gollama.NewMcpClient(config)
defer client.Close()

// Start connection
if err := client.Start(ctx); err != nil {
    panic(err)
}

// Get tools from the MCP server
tools, _ := client.ListTools()

// Pass them to Gollama
output, _ := g.Chat(ctx, "List the users in the database", tools)

// Execute tool calls requested by the model
for _, call := range output.ToolCalls {
    result, _ := client.CallTool(call.Function.Name, call.Function.Arguments)
    fmt.Println("Tool Result:", result)
}
```

### 3. Structured Outputs (JSON)
Force the model to return data matching your Go struct definition.

```go
type Capital struct {
    Country string `json:"country" description:"The name of the country"`
    City    string `json:"city" description:"The capital city"`
    Population int `json:"population" description:"Approximate population"`
}

// Convert struct to schema
schema := gollama.StructToStructuredFormat(Capital{})

resp, err := g.Chat(ctx, "Tell me about France", schema)

// Decode directly into your struct
var result Capital
resp.DecodeContent(&result)

fmt.Printf("%+v\n", result)
```

### 4. Vision
Analyze images with multimodal models.

```go
g := gollama.New("llama3.2-vision")

image := gollama.PromptImage{Filename: "./photo.png"}

resp, err := g.Chat(ctx, "Describe this image", image)
```

## 📚 API Reference

### Core Functions
- `New(model string) *Gollama`: Initialize a new client.
- `g.Chat(ctx, prompt, options...)`: Main entry point for interaction. Options can be `Tool`, `PromptImage`, or `StructuredFormat`.
- `g.PullIfMissing(ctx)`: Ensures the model exists locally before running.

### Utilities
- `StructToStructuredFormat(v interface{})`: Generates a JSON schema from a Go struct.
- `DecodeContent(v interface{})`: Unmarshals the JSON response into a struct.
- `CosenoSimilarity(v1, v2 []float64)`: Helper for RAG/Embedding comparisons.

### MCP (Model Context Protocol)
- `NewMcpClient(config McpConfig)`: Creates a client to talk to any MCP-compliant server.
- `client.ListTools()`: Discovers available tools on the server.
- `client.CallTool(name, args)`: Executes a tool on the server.

## 📄 License

MIT
