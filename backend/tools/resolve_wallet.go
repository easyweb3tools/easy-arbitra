package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/mark3labs/mcp-go/mcp"
)

var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
var profileURLRegex = regexp.MustCompile(`polymarket\.com/profile/(0x[a-fA-F0-9]{40})`)

type ResolveResult struct {
	WalletAddress string `json:"wallet_address"`
	DisplayName   string `json:"display_name"`
	InputType     string `json:"input_type"`
	ProfileImage  string `json:"profile_image"`
}

func ResolveWalletTarget(client *polymarket.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		input, ok := args["input"].(string)
		if !ok || input == "" {
			return mcp.NewToolResultError("input parameter is required"), nil
		}

		var address string
		var inputType string

		input = strings.TrimSpace(input)

		if ethAddressRegex.MatchString(input) {
			address = input
			inputType = "wallet_address"
		} else if matches := profileURLRegex.FindStringSubmatch(input); len(matches) > 1 {
			address = matches[1]
			inputType = "polymarket_url"
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("cannot parse input: must be an Ethereum address (0x...) or a Polymarket profile URL, got: %s", input)), nil
		}

		displayName := address[:6] + "..." + address[len(address)-4:]
		profileImage := ""

		profile, err := client.GetPublicProfile(address)
		if err == nil && profile != nil {
			if profile.Pseudonym != "" {
				displayName = profile.Pseudonym
			} else if profile.Name != "" {
				displayName = profile.Name
			}
			profileImage = profile.ProfileImage
		}

		result := ResolveResult{
			WalletAddress: address,
			DisplayName:   displayName,
			InputType:     inputType,
			ProfileImage:  profileImage,
		}

		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
