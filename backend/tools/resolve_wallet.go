package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/mark3labs/mcp-go/mcp"
)

var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
var profileURLRegex = regexp.MustCompile(`polymarket\.com/profile/(0x[a-fA-F0-9]{40})`)
var profileSlugURLRegex = regexp.MustCompile(`(?:https?://)?(?:www\.)?polymarket\.com/(@[^/?#]+)`)
var slugInputRegex = regexp.MustCompile(`^@[^/?#]+$`)
var proxyAddressRegex = regexp.MustCompile(`"(?:proxyAddress|baseAddress)":"(0x[a-f0-9]{40})"`)

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
		} else if matches := profileSlugURLRegex.FindStringSubmatch(input); len(matches) > 1 {
			resolved, err := resolveProfileSlugAddress(ctx, client, matches[1])
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to resolve Polymarket profile slug: %v", err)), nil
			}
			address = resolved
			inputType = "polymarket_slug"
		} else if slugInputRegex.MatchString(input) {
			resolved, err := resolveProfileSlugAddress(ctx, client, input)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to resolve Polymarket profile slug: %v", err)), nil
			}
			address = resolved
			inputType = "polymarket_slug"
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

func resolveProfileSlugAddress(ctx context.Context, client *polymarket.Client, slug string) (string, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", fmt.Errorf("empty profile slug")
	}
	if !strings.HasPrefix(slug, "@") {
		slug = "@" + slug
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://polymarket.com/"+slug, nil)
	if err != nil {
		return "", fmt.Errorf("build profile request: %w", err)
	}

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("profile page request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("profile page returned %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read profile page: %w", err)
	}

	matches := proxyAddressRegex.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("no wallet address found in profile page")
	}
	return string(matches[1]), nil
}
