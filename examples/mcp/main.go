package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jonathanhecl/gollama"
)

func main() {
	// Configuration for the MCP server (e.g., using Tiny implementation or any other)
	// Here we use a hypothetical 'npx' command to run a server.
	// You can replace this with any valid MCP server configuration.
	// For example, to use the official filesystem server:
	// npx -y @modelcontextprotocol/server-filesystem /path/to/files

	// Check if user provided an argument for the server, else default to a simple echo/test if possible,
	// or warn the user they need to configure it.

	fmt.Println("Starting MCP Client Example...")
	fmt.Println("Note: This example requires 'npx' and an internet connection to download the MCP server if not cached.")

	// Example 1: Filesystem Server (Simple)
	// args: ["-y", "@modelcontextprotocol/server-filesystem", cwd]
	cwd, _ := os.Getwd()
	fsConfig := gollama.McpConfig{
		Command: "npx",
		Args: []string{
			"-y",
			"@modelcontextprotocol/server-filesystem",
			cwd,
		},
	}

	// Example 2: Supabase MCP Server (Complex configuration with Env)
	// This matches the user's request structure:
	// "supabase-mcp-server": {
	//    "args": ["-y", "@supabase/mcp-server-supabase@latest", "--access-token", "..."],
	//    "env": { "SUPABASE_URL": "..." }
	// }
	_ = gollama.McpConfig{
		Command: "npx",
		Args: []string{
			"-y",
			"@supabase/mcp-server-supabase@latest",
			"--access-token", "sbp_your_token_here",
		},
		Env: map[string]string{
			// "SUPABASE_Reference_ID": "...", // If needed by the specific server
		},
	}

	// Example 3: Playwright (Browser automation)
	_ = gollama.McpConfig{
		Command: "npx",
		Args: []string{
			"-y",
			"@playwright/mcp@latest",
		},
	}

	// Select which config to run
	activeConfig := fsConfig

	fmt.Printf("Starting MCP Client with command: %s %v\n", activeConfig.Command, activeConfig.Args)

	client := gollama.NewMcpClient(activeConfig)

	// Optional: Enable logging to see MCP protocol traffic
	// client.SetLogger(func(format string, args ...any) {
	// 	log.Printf(format, args...)
	// })

	ctx := context.Background()

	fmt.Println("Connecting to MCP server...")
	if err := client.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP client: %v", err)
	}
	defer client.Close()
	fmt.Println("Connected!")

	fmt.Println("Listing tools...")
	tools, err := client.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	for _, t := range tools {
		fmt.Printf("- %s: %s\n", t.Function.Name, t.Function.Description)
	}

	if len(tools) == 0 {
		fmt.Println("No tools found on the MCP server.")
		return
	}

	// Initialize Gollama
	g := gollama.New("llama3.2") // Ensure you have this model or change it

	prompt := "List the files in the current directory using the available tools."
	fmt.Printf("\nSending prompt to LLM: %q\n", prompt)

	// Pass the MCP tools to Gollama
	response, err := g.Chat(ctx, prompt, tools)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	if len(response.ToolCalls) > 0 {
		fmt.Println("\nModel requested tool execution:")
		for _, toolCall := range response.ToolCalls {
			fmt.Printf("Tool: %s\nArgs: %v\n", toolCall.Function.Name, toolCall.Function.Arguments)

			// Execute the tool using the MCP client
			result, err := client.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				log.Printf("Tool execution failed: %v", err)
				continue
			}
			fmt.Printf("Result: %s\n", result)
		}
	} else {
		fmt.Println("\nModel response:", response.Content)
	}
}
