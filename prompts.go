package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterAllPrompts(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(
		mcp.NewPrompt("daily_digest",
			mcp.WithPromptDescription("Summarize today's unread entries across all feeds"),
		),
		func(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return mcp.NewGetPromptResult(
				"Summarize today's unread entries",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(
						mcp.RoleUser,
						mcp.NewTextContent(
							"Please create a digest of today's unread feed entries using the following steps:\n\n"+
								"1. Call get_entries with status=\"unread\" to fetch all unread entries.\n"+
								"2. Filter for entries published today (use the published_at field).\n"+
								"3. Group them by feed or category.\n"+
								"4. For each entry, provide the title and a 2-3 sentence summary.\n"+
								"5. Highlight any particularly noteworthy articles.\n\n"+
								"Present the result in a clear, readable format.",
						),
					),
				},
			), nil
		},
	)

	mcpServer.AddPrompt(
		mcp.NewPrompt("unread_summary",
			mcp.WithPromptDescription("Show unread counts grouped by category"),
		),
		func(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return mcp.NewGetPromptResult(
				"Unread entries summary by category",
				[]mcp.PromptMessage{
					mcp.NewPromptMessage(
						mcp.RoleUser,
						mcp.NewTextContent(
							"Please summarize the current unread status in Miniflux using the following steps:\n\n"+
								"1. Call fetch_counters to get the overall unread count.\n"+
								"2. Call get_categories to retrieve all categories.\n"+
								"3. For each category, call get_category_entries with status=\"unread\" to count unread entries.\n"+
								"4. Rank categories by unread count (highest first).\n"+
								"5. Include the total unread count and an estimated reading time.\n\n"+
								"Present the result as a concise ranked list.",
						),
					),
				},
			), nil
		},
	)
}
