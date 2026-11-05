package mcpwhiteboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpsdk "github.com/mark3labs/mcp-go/server"
)

type Server struct {
	httpServer *mcpsdk.StreamableHTTPServer
}

type visitNode struct {
	ID          string              `json:"Id"`
	Workspace   string              `json:"Workspace" jsonschema_description:"所属工作区的实际名称，保留大小写"`
	Title       string              `json:"Title" jsonschema_description:"纯文本标题；不做 Markdown 转义，创建时去除首尾空白"`
	Progress    string              `json:"Progress"`
	Session     []string            `json:"Session"`
	Description map[string][]string `json:"Description" jsonschema_description:"各分区字符串值为规范化后的 Markdown 源文本；比对应以规范化后文本或渲染后文本为准"`
}

type visitEdge struct {
	Source string `json:"Source"`
	Target string `json:"Target"`
}

type visitNodeCreateInput struct {
	Workspace   string              `json:"Workspace"`
	Title       string              `json:"Title"`
	Session     []string            `json:"Session"`
	Description map[string][]string `json:"Description"`
}

type visitNodeDeleteInput struct {
	ID string `json:"Id"`
}

type visitEdgesInput struct {
	Source string `json:"Source"`
	Target string `json:"Target"`
}

type visitNodesOutput struct {
	VisitNodes []visitNode `json:"visit_nodes"`
}

type visitEdgesOutput struct {
	VisitEdges []visitEdge `json:"visit_edges"`
}

type deleteOutput struct {
	Deleted bool `json:"deleted"`
}

func New(version, vikunjaURL, vikunjaToken string) *Server {
	server := mcpsdk.NewMCPServer(
		"backstage-whiteboard",
		version,
		mcpsdk.WithInstructions("Whiteboard 工具仅用于处理 Whiteboard 相关问题。Description 按分区组织，每个字符串值为 Markdown。写入时，孤立的强调符号（_、*）会被转义，成对标记（如 **bold**）保留；规范化具有幂等性，重复写入不叠加转义。读取返回规范化后的 Markdown 源文本，上述规范化保留用户预期的渲染效果。自动化校验应以规范化后文本或渲染后文本为基准，不得对写入前后的源文本做字节级 diff。Title 为纯文本，不做 Markdown 转义，创建时去除首尾空白。Workspace 为工作区名称，返回所属工作区的实际名称并保留大小写。"),
		mcpsdk.WithToolCapabilities(false),
	)
	client := &http.Client{Timeout: 30 * time.Second}
	request := func(ctx context.Context, method, path string, input, output any) error {
		if strings.TrimSpace(vikunjaURL) == "" || strings.TrimSpace(vikunjaToken) == "" {
			return fmt.Errorf("whiteboard service is not configured")
		}

		var body *bytes.Reader
		if input == nil {
			body = bytes.NewReader(nil)
		} else {
			payload, err := json.Marshal(input)
			if err != nil {
				return fmt.Errorf("whiteboard request is invalid")
			}
			body = bytes.NewReader(payload)
		}

		endpointURL := strings.TrimRight(vikunjaURL, "/") + "/api/v1/plugins/whiteboard/" + path
		httpRequest, err := http.NewRequestWithContext(ctx, method, endpointURL, body)
		if err != nil {
			return fmt.Errorf("whiteboard service request failed")
		}
		httpRequest.Header.Set("Accept", "application/json")
		httpRequest.Header.Set("Authorization", "Bearer "+vikunjaToken)
		if input != nil {
			httpRequest.Header.Set("Content-Type", "application/json")
		}

		response, err := client.Do(httpRequest)
		if err != nil {
			return fmt.Errorf("whiteboard service request failed")
		}
		defer response.Body.Close()
		if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
			return fmt.Errorf("whiteboard service returned %s", response.Status)
		}
		if output != nil && response.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(response.Body).Decode(output); err != nil {
				return fmt.Errorf("whiteboard service returned an invalid response")
			}
		}
		return nil
	}

	server.AddTool(mcp.NewTool(
		"whiteboard_list_visit_nodes",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。列出 Whiteboard 中的全部节点，包含所属工作区名称 Workspace。Description 各字符串值为规范化后的 Markdown 源文本；比对应以规范化后文本或渲染后文本为准，不得对写入前后的源文本做字节级 diff。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)),
		mcp.WithOutputSchema[visitNodesOutput](),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
	), func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var output visitNodesOutput
		if err := request(ctx, http.MethodGet, "visit-nodes", nil, &output.VisitNodes); err != nil {
			return nil, err
		}
		if output.VisitNodes == nil {
			output.VisitNodes = []visitNode{}
		}
		return mcp.NewToolResultJSON(output)
	})

	server.AddTool(mcp.NewTool(
		"whiteboard_create_visit_node",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。创建节点；需要指定工作区名称、标题，并可提供 Session 与 Description。Description 各字符串值为 Markdown：写入时转义孤立的 _、*，保留成对标记（如 **bold**），处理幂等且保留预期渲染效果。返回规范化后的 Markdown 源文本；比对应以规范化后文本或渲染后文本为准，不得对写入前后的源文本做字节级 diff。Title 为纯文本，不做 Markdown 转义，创建时去除首尾空白。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{"Workspace":{"type":"string","description":"目标工作区名称；返回所属工作区的实际名称并保留大小写"},"Title":{"type":"string","description":"纯文本节点标题；不做 Markdown 转义，创建时去除首尾空白"},"Session":{"type":"array","items":{"type":"string"},"description":"节点所属会话"},"Description":{"type":"object","additionalProperties":{"type":"array","items":{"type":"string","description":"Markdown 源文本；写入时转义孤立的 _、*，保留成对标记，规范化幂等；读取返回规范化后的文本"}},"description":"按分区组织的 Markdown 描述；比对应以规范化后文本或渲染后文本为准，不得对写入前后的源文本做字节级 diff"}},"required":["Workspace","Title"],"additionalProperties":false}`)),
		mcp.WithOutputSchema[visitNode](),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
	), func(ctx context.Context, call mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input visitNodeCreateInput
		if err := call.BindArguments(&input); err != nil {
			return nil, fmt.Errorf("invalid visit node input")
		}
		var output visitNode
		if err := request(ctx, http.MethodPost, "visit-nodes.create", input, &output); err != nil {
			return nil, err
		}
		return mcp.NewToolResultJSON(output)
	})

	server.AddTool(mcp.NewTool(
		"whiteboard_delete_visit_node",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。按节点 ID 删除 Whiteboard 节点。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{"Id":{"type":"string","description":"要删除的节点 ID"}},"required":["Id"],"additionalProperties":false}`)),
		mcp.WithOutputSchema[deleteOutput](),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
	), func(ctx context.Context, call mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input visitNodeDeleteInput
		if err := call.BindArguments(&input); err != nil || strings.TrimSpace(input.ID) == "" {
			return nil, fmt.Errorf("invalid visit node ID")
		}
		path := "visit-nodes.delete/" + url.PathEscape(input.ID)
		if err := request(ctx, http.MethodDelete, path, nil, nil); err != nil {
			return nil, err
		}
		return mcp.NewToolResultJSON(deleteOutput{Deleted: true})
	})

	server.AddTool(mcp.NewTool(
		"whiteboard_list_visit_edges",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。列出 Whiteboard 中的全部边。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{},"additionalProperties":false}`)),
		mcp.WithOutputSchema[visitEdgesOutput](),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
	), func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var output visitEdgesOutput
		if err := request(ctx, http.MethodGet, "visit-edges", nil, &output.VisitEdges); err != nil {
			return nil, err
		}
		if output.VisitEdges == nil {
			output.VisitEdges = []visitEdge{}
		}
		return mcp.NewToolResultJSON(output)
	})

	server.AddTool(mcp.NewTool(
		"whiteboard_create_visit_edge",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。创建一条从 Source 指向 Target 的 Whiteboard 边。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{"Source":{"type":"string","description":"起点节点 ID"},"Target":{"type":"string","description":"终点节点 ID"}},"required":["Source","Target"],"additionalProperties":false}`)),
		mcp.WithOutputSchema[visitEdge](),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(false),
	), func(ctx context.Context, call mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input visitEdgesInput
		if err := call.BindArguments(&input); err != nil {
			return nil, fmt.Errorf("invalid visit edge input")
		}
		var output visitEdge
		if err := request(ctx, http.MethodPost, "visit-edge.create", input, &output); err != nil {
			return nil, err
		}
		return mcp.NewToolResultJSON(output)
	})

	server.AddTool(mcp.NewTool(
		"whiteboard_delete_visit_edge",
		mcp.WithDescription("仅在处理 Whiteboard 相关问题时调用。删除 Source 到 Target 的 Whiteboard 边。"),
		mcp.WithRawInputSchema(json.RawMessage(`{"type":"object","properties":{"Source":{"type":"string","description":"起点节点 ID"},"Target":{"type":"string","description":"终点节点 ID"}},"required":["Source","Target"],"additionalProperties":false}`)),
		mcp.WithOutputSchema[deleteOutput](),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(false),
	), func(ctx context.Context, call mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input visitEdgesInput
		if err := call.BindArguments(&input); err != nil {
			return nil, fmt.Errorf("invalid visit edge input")
		}
		if err := request(ctx, http.MethodPost, "visit-edges.delete", input, nil); err != nil {
			return nil, err
		}
		return mcp.NewToolResultJSON(deleteOutput{Deleted: true})
	})

	return &Server{
		httpServer: mcpsdk.NewStreamableHTTPServer(server, mcpsdk.WithStateful(true)),
	}
}

func (s *Server) Handler() http.Handler {
	return s.httpServer
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
