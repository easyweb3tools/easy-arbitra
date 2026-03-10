// Bedrock Converse API tool definitions mirroring the 4 MCP tools

export const toolConfig = {
  tools: [
    {
      toolSpec: {
        name: "resolve_wallet_target",
        description:
          "Resolve a wallet address or Polymarket profile URL to a standardized wallet target with display name and profile image.",
        inputSchema: {
          json: {
            type: "object",
            properties: {
              input: {
                type: "string",
                description:
                  "Wallet address (0x...) or Polymarket profile URL",
              },
            },
            required: ["input"],
          },
        },
      },
    },
    {
      toolSpec: {
        name: "fetch_sports_trades",
        description:
          "Fetch sports trades for a given wallet address from Polymarket. Returns enriched trade data with market metadata.",
        inputSchema: {
          json: {
            type: "object",
            properties: {
              wallet: {
                type: "string",
                description: "Standardized wallet address (0x...)",
              },
              sport: {
                type: "string",
                description: "Sport to filter trades by (e.g., 'nba')",
              },
              limit: {
                type: "number",
                description: "Maximum number of trades to fetch (default 500)",
              },
            },
            required: ["wallet"],
          },
        },
      },
    },
    {
      toolSpec: {
        name: "calculate_style_metrics",
        description:
          "Calculate deterministic trading style metrics (entry timing, position size ratio, conviction) from enriched trade data.",
        inputSchema: {
          json: {
            type: "object",
            properties: {
              wallet: {
                type: "string",
                description: "Wallet address",
              },
              trades_json: {
                type: "string",
                description:
                  "JSON array of enriched trades from fetch_sports_trades",
              },
            },
            required: ["wallet", "trades_json"],
          },
        },
      },
    },
    {
      toolSpec: {
        name: "build_report_payload",
        description:
          "Build the final report payload including wallet card, radar chart data, and style summary for frontend rendering.",
        inputSchema: {
          json: {
            type: "object",
            properties: {
              wallet_info: {
                type: "string",
                description:
                  "JSON string of wallet info from resolve_wallet_target",
              },
              metrics_json: {
                type: "string",
                description:
                  "JSON string of metrics result from calculate_style_metrics",
              },
              trades_summary: {
                type: "string",
                description:
                  "Optional JSON string of trades summary for additional context",
              },
            },
            required: ["wallet_info", "metrics_json"],
          },
        },
      },
    },
  ],
};
