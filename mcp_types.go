package gollama

import "encoding/json"

// JSON-RPC 2.0 Types

type JsonRpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  interface{}     `json:"params,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"` // Can be string or number
}

type JsonRpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JsonRpcError   `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// MCP Protocol Types

type McpInitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    McpCapabilities `json:"capabilities"`
	ClientInfo      McpClientInfo   `json:"clientInfo"`
}

type McpCapabilities struct {
	Roots    map[string]interface{} `json:"roots,omitempty"`
	Sampling map[string]interface{} `json:"sampling,omitempty"`
}

type McpClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type McpInitializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ServerInfo      McpClientInfo  `json:"serverInfo"`
}

type McpToolListResult struct {
	Tools []McpTool `json:"tools"`
}

type McpTool struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	InputSchema StructuredFormat `json:"inputSchema"`
}

type McpCallToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type McpCallToolResult struct {
	Content []McpContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type McpContent struct {
	Type string `json:"type"` // "text" or "image" or "resource"
	Text string `json:"text,omitempty"`
	// Add other fields for images/resources as needed later
}
