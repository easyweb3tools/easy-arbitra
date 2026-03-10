package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/brucexwang/easy-arbitra/backend/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	client := polymarket.NewClient()

	// Create MCP Server
	mcpServer := server.NewMCPServer(
		"SportStyle",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register tools
	registerTools(mcpServer, client)

	// Start SSE server on :8081
	sseServer := server.NewSSEServer(mcpServer)
	go func() {
		log.Println("MCP SSE Server starting on :8081")
		if err := sseServer.Start(":8081"); err != nil {
			log.Printf("SSE server error: %v", err)
			os.Exit(1)
		}
	}()

	// Start REST bridge on :8082
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tools/call", corsMiddleware(restBridge(client)))
	mux.HandleFunc("/api/health", corsMiddleware(healthHandler))

	log.Println("REST bridge starting on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatalf("REST server error: %v", err)
	}
}

func registerTools(s *server.MCPServer, client *polymarket.Client) {
	// 1. resolve_wallet_target
	s.AddTool(mcp.NewTool("resolve_wallet_target",
		mcp.WithDescription("Resolve a wallet address or Polymarket profile URL to a standardized wallet target with display name and profile image."),
		mcp.WithString("input",
			mcp.Description("Wallet address (0x...) or Polymarket profile URL"),
			mcp.Required(),
		),
	), tools.ResolveWalletTarget(client))

	// 2. fetch_sports_trades
	s.AddTool(mcp.NewTool("fetch_sports_trades",
		mcp.WithDescription("Fetch NBA sports trades for a given wallet address from Polymarket. Returns enriched trade data with market metadata."),
		mcp.WithString("wallet",
			mcp.Description("Standardized wallet address (0x...)"),
			mcp.Required(),
		),
	), tools.FetchSportsTrades(client))

	// 3. calculate_style_metrics
	s.AddTool(mcp.NewTool("calculate_style_metrics",
		mcp.WithDescription("Calculate deterministic trading style metrics (entry timing, position size ratio, conviction) from enriched trade data."),
		mcp.WithString("wallet",
			mcp.Description("Wallet address"),
			mcp.Required(),
		),
		mcp.WithString("trades_json",
			mcp.Description("JSON array of enriched trades from fetch_sports_trades"),
			mcp.Required(),
		),
	), tools.CalculateStyleMetrics())

	// 4. build_report_payload
	s.AddTool(mcp.NewTool("build_report_payload",
		mcp.WithDescription("Build the final report payload including wallet card, radar chart data, and style summary for frontend rendering."),
		mcp.WithString("wallet_info",
			mcp.Description("JSON string of wallet info from resolve_wallet_target"),
			mcp.Required(),
		),
		mcp.WithString("metrics_json",
			mcp.Description("JSON string of metrics result from calculate_style_metrics"),
			mcp.Required(),
		),
		mcp.WithString("trades_summary",
			mcp.Description("Optional JSON string of trades summary for additional context"),
		),
	), tools.BuildReportPayload())
}

// REST bridge handler
type ToolCallRequest struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

func restBridge(client *polymarket.Client) http.HandlerFunc {
	handlers := map[string]func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error){
		"resolve_wallet_target":  tools.ResolveWalletTarget(client),
		"fetch_sports_trades":    tools.FetchSportsTrades(client),
		"calculate_style_metrics": tools.CalculateStyleMetrics(),
		"build_report_payload":   tools.BuildReportPayload(),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ToolCallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			return
		}

		handler, ok := handlers[req.Tool]
		if !ok {
			http.Error(w, fmt.Sprintf("unknown tool: %s", req.Tool), http.StatusBadRequest)
			return
		}

		mcpReq := mcp.CallToolRequest{}
		mcpReq.Params.Name = req.Tool
		mcpReq.Params.Arguments = req.Args

		result, err := handler(r.Context(), mcpReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("tool error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
