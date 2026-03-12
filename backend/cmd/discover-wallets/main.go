package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brucexwang/easy-arbitra/backend/discovery"
	"github.com/brucexwang/easy-arbitra/backend/polymarket"
)

func main() {
	var (
		recentLimit     = flag.Int("recent-limit", 400, "number of recent global trades to scan")
		recentPages     = flag.Int("recent-pages", 4, "number of recent trade pages to scan")
		perWalletLimit  = flag.Int("wallet-limit", 500, "trade history limit when scoring each wallet")
		candidateLimit  = flag.Int("candidates", 20, "number of candidate wallets to inspect deeply")
		outputLimit     = flag.Int("output", 10, "number of top wallets to print")
		minRecentTrades = flag.Int("min-recent-trades", 2, "minimum recent NBA trades to consider a wallet")
		sport           = flag.String("sport", "nba", "sport to search for")
		jsonOutput      = flag.Bool("json", false, "print machine-readable JSON")
		walletsFile     = flag.String("wallets-file", "", "newline-delimited wallet list to score directly")
	)
	flag.Parse()

	client := polymarket.NewClient()
	opts := discovery.Options{
		Sport:           *sport,
		RecentLimit:     *recentLimit,
		RecentPages:     *recentPages,
		CandidateLimit:  *candidateLimit,
		OutputLimit:     *outputLimit,
		MinRecentTrades: *minRecentTrades,
		WalletLimit:     *perWalletLimit,
	}

	var (
		results []discovery.Candidate
		err     error
	)
	if *walletsFile != "" {
		wallets, readErr := readWallets(*walletsFile)
		if readErr != nil {
			log.Fatal(readErr)
		}
		results, err = discovery.ScoreWallets(context.Background(), client, wallets, opts)
	} else {
		results, err = discovery.DiscoverFromRecent(context.Background(), client, opts)
	}
	if err != nil {
		log.Fatal(err)
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
		return
	}

	fmt.Printf("Top %d %s wallets for demo\n\n", len(results), strings.ToUpper(*sport))
	for idx, result := range results {
		fmt.Printf("%d. %s (%s)\n", idx+1, result.DisplayName, result.Wallet)
		fmt.Printf("   NBA trades: %d | recent sample hits: %d | recent markets: %d\n", result.NbaTrades, result.RecentTrades, result.RecentMarkets)
		fmt.Printf("   Style: %s | conviction: %.2f | size ratio: %.4f%% | entry timing: %.1fh\n", result.StyleLabel, result.Conviction, result.SizeRatioPct, result.EntryTimingHours)
		fmt.Printf("   Demo reason: %s\n\n", result.Reason)
	}
}

func readWallets(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wallets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wallet := strings.TrimSpace(scanner.Text())
		if wallet == "" || strings.HasPrefix(wallet, "#") {
			continue
		}
		wallets = append(wallets, wallet)
	}
	return wallets, scanner.Err()
}
