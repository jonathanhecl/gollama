package gollama

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
)

// McpConfig defines how to start an MCP server.
type McpConfig struct {
	Command string            // Executable to run (e.g., "npx", "python")
	Args    []string          // Arguments for the command
	Env     map[string]string // Environment variables to set
}

// McpClient manages the connection to an MCP server.
type McpClient struct {
	config McpConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	seq       int64
	pending   map[string]chan *JsonRpcResponse
	pendingMu sync.Mutex

	logger func(string, ...any) // Simple logger
}

// NewMcpClient creates a new client instance.
func NewMcpClient(config McpConfig) *McpClient {
	return &McpClient{
		config:  config,
		pending: make(map[string]chan *JsonRpcResponse),
		logger: func(format string, args ...any) {
			// Default no-op logger, user can override if needed
		},
	}
}

// SetLogger sets a logging function for debugging MCP traffic.
func (c *McpClient) SetLogger(l func(string, ...any)) {
	c.logger = l
}

// Start launches the MCP server process and initializes the connection.
func (c *McpClient) Start(ctx context.Context) error {
	c.cmd = exec.CommandContext(ctx, c.config.Command, c.config.Args...)

	// Setup environment
	c.cmd.Env = os.Environ()
	for k, v := range c.config.Env {
		c.cmd.Env = append(c.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var err error
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	c.stdout, err = c.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	c.stderr, err = c.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start reading stderr in background for debugging
	go func() {
		scanner := bufio.NewScanner(c.stderr)
		for scanner.Scan() {
			c.logger("MCP STDERR: %s", scanner.Text())
		}
	}()

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Start reading stdout loop
	go c.readLoop()

	// Perform MCP Handshake
	return c.initialize(ctx)
}

// Close terminates the MCP server process.
func (c *McpClient) Close() error {
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}

func (c *McpClient) readLoop() {
	reader := bufio.NewReader(c.stdout)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				c.logger("Error reading stdout: %v", err)
			}
			break
		}

		c.logger("MCP < %s", string(line))

		var response JsonRpcResponse
		if err := json.Unmarshal(line, &response); err != nil {
			c.logger("Error unmarshaling JSON-RPC response: %v", err)
			continue
		}

		if response.ID != nil {
			id := string(response.ID)
			c.pendingMu.Lock()
			ch, ok := c.pending[id]
			if ok {
				delete(c.pending, id)
			}
			c.pendingMu.Unlock()

			if ok {
				ch <- &response
			}
		} else {
			// Notification or request from server (not handled in this simple client yet)
		}
	}
}

func (c *McpClient) sendRequest(method string, params interface{}) (*JsonRpcResponse, error) {
	id := atomic.AddInt64(&c.seq, 1)
	idStr := fmt.Sprintf("%d", id)
	idRaw := json.RawMessage(idStr)

	req := JsonRpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      idRaw,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *JsonRpcResponse, 1)
	c.pendingMu.Lock()
	c.pending[idStr] = ch
	c.pendingMu.Unlock()

	c.logger("MCP > %s", string(reqBytes))
	if _, err := c.stdin.Write(append(reqBytes, '\n')); err != nil {
		return nil, err
	}

	// Wait for response
	resp := <-ch
	if resp.Error != nil {
		return nil, fmt.Errorf("MCP error (%d): %s", resp.Error.Code, resp.Error.Message)
	}

	return resp, nil
}

func (c *McpClient) sendNotification(method string, params interface{}) error {
	req := JsonRpcRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	c.logger("MCP > %s", string(reqBytes))
	_, err = c.stdin.Write(append(reqBytes, '\n'))
	return err
}

func (c *McpClient) initialize(ctx context.Context) error {
	// 1. Send initialize request
	params := McpInitializeParams{
		ProtocolVersion: "2024-11-05", // Example version
		Capabilities:    McpCapabilities{},
		ClientInfo: McpClientInfo{
			Name:    "gollama-mcp-client",
			Version: "1.0.0",
		},
	}

	resp, err := c.sendRequest("initialize", params)
	if err != nil {
		return err
	}

	var result McpInitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return err
	}

	// 2. Send initialized notification
	return c.sendNotification("notifications/initialized", map[string]interface{}{})
}

// ListTools returns the list of tools available on the MCP server converted to Gollama Tools.
func (c *McpClient) ListTools() ([]Tool, error) {
	resp, err := c.sendRequest("tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result McpToolListResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	var tools []Tool
	for _, mcpTool := range result.Tools {
		tools = append(tools, Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        mcpTool.Name,
				Description: mcpTool.Description,
				Parameters:  mcpTool.InputSchema,
			},
		})
	}

	return tools, nil
}

// CallTool executes a tool on the MCP server.
func (c *McpClient) CallTool(name string, arguments map[string]any) (string, error) {
	params := McpCallToolParams{
		Name:      name,
		Arguments: arguments,
	}

	resp, err := c.sendRequest("tools/call", params)
	if err != nil {
		return "", err
	}

	var result McpCallToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", err
	}

	if result.IsError {
		return "", fmt.Errorf("tool execution failed")
	}

	// Combine all text content
	var output string
	for _, content := range result.Content {
		if content.Type == "text" {
			output += content.Text + "\n"
		}
	}

	return output, nil
}
