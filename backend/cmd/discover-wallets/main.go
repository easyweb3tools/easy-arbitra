package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/brucexwang/easy-arbitra/backend/metrics"
	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/brucexwang/easy-arbitra/backend/tools"
)

type walletSeed struct {
	Wallet          string
	RecentTrades    int
	UniqueMarkets   map[string]bool
	RecentBuyVolume float64
}

type candidate struct {
	Wallet           string  `json:"wallet"`
	DisplayName      string  `json:"display_name"`
	RecentTrades     int     `json:"recent_trades"`
	RecentMarkets    int     `json:"recent_markets"`
	NbaTrades        int     `json:"nba_trades"`
	EntryTimingHours float64 `json:"entry_timing_hours"`
	SizeRatioPct     float64 `json:"size_ratio_pct"`
	Conviction       float64 `json:"conviction"`
	StyleLabel       string  `json:"style_label"`
	PresentationScore float64 `json:"presentation_score"`
	Reason           string  `json:"reason"`
}

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
	ctx := context.Background()

	var seeds []walletSeed
	var err error
	if *walletsFile != "" {
		seeds, err = readWalletSeeds(*walletsFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		seeds, err = discoverSeeds(client, *sport, *recentLimit, *recentPages, *minRecentTrades)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(seeds) == 0 {
		log.Fatalf("no %s candidate wallets found; try a larger recent sample or provide -wallets-file", strings.ToUpper(*sport))
	}

	sort.Slice(seeds, func(i, j int) bool {
		if seeds[i].RecentTrades != seeds[j].RecentTrades {
			return seeds[i].RecentTrades > seeds[j].RecentTrades
		}
		return len(seeds[i].UniqueMarkets) > len(seeds[j].UniqueMarkets)
	})

	if len(seeds) > *candidateLimit {
		seeds = seeds[:*candidateLimit]
	}

	results := make([]candidate, 0, len(seeds))
	for _, seed := range seeds {
		profile, _ := client.GetPublicProfile(seed.Wallet)
		displayName := shortWallet(seed.Wallet)
		if profile != nil {
			if profile.Pseudonym != "" {
				displayName = profile.Pseudonym
			} else if profile.Name != "" {
				displayName = profile.Name
			}
		}

		fetchResult, err := tools.FetchSportsTradesData(ctx, client, seed.Wallet, *sport, *perWalletLimit)
		if err != nil {
			continue
		}
		if fetchResult.TotalTrades == 0 {
			continue
		}

		entryTiming := metrics.EntryTimingHours(fetchResult.Trades)
		sizeRatio := metrics.SizeRatioPct(fetchResult.Trades)
		conviction := metrics.Conviction(fetchResult.Trades)
		uniqueMarkets := countUniqueMarkets(fetchResult.Trades)
		styleLabel := tools.DetermineStyleLabel(
			max(0, 1-entryTiming/24),
			minFloat(1, sizeRatio/10),
			conviction,
		)

		score := presentationScore(fetchResult.TotalTrades, uniqueMarkets, conviction, sizeRatio)
		results = append(results, candidate{
			Wallet:            seed.Wallet,
			DisplayName:       displayName,
			RecentTrades:      seed.RecentTrades,
			RecentMarkets:     len(seed.UniqueMarkets),
			NbaTrades:         fetchResult.TotalTrades,
			EntryTimingHours:  entryTiming,
			SizeRatioPct:      sizeRatio,
			Conviction:        conviction,
			StyleLabel:        styleLabel,
			PresentationScore: score,
			Reason:            buildReason(fetchResult.TotalTrades, uniqueMarkets, conviction, sizeRatio),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].PresentationScore != results[j].PresentationScore {
			return results[i].PresentationScore > results[j].PresentationScore
		}
		return results[i].NbaTrades > results[j].NbaTrades
	})

	if len(results) > *outputLimit {
		results = results[:*outputLimit]
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

func discoverSeeds(client *polymarket.Client, sport string, limit, pages, minRecentTrades int) ([]walletSeed, error) {
	seeds := map[string]*walletSeed{}

	for page := 0; page < pages; page++ {
		offset := page * limit
		trades, err := client.GetRecentTrades(limit, offset)
		if err != nil {
			return nil, fmt.Errorf("fetch recent trades page %d: %w", page, err)
		}
		if len(trades) == 0 {
			break
		}

		for _, trade := range trades {
			if !isSportText(trade.Title, sport) {
				continue
			}
			seed, ok := seeds[trade.ProxyWallet]
			if !ok {
				seed = &walletSeed{
					Wallet:        trade.ProxyWallet,
					UniqueMarkets: map[string]bool{},
				}
				seeds[trade.ProxyWallet] = seed
			}

			seed.RecentTrades++
			seed.UniqueMarkets[trade.ConditionID] = true
			if trade.Side == "BUY" {
				seed.RecentBuyVolume += trade.Size * trade.Price
			}
		}
	}

	list := make([]walletSeed, 0, len(seeds))
	for _, seed := range seeds {
		if seed.RecentTrades < minRecentTrades {
			continue
		}
		list = append(list, *seed)
	}
	return list, nil
}

func readWalletSeeds(path string) ([]walletSeed, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var seeds []walletSeed
	seen := map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wallet := strings.TrimSpace(scanner.Text())
		if wallet == "" || strings.HasPrefix(wallet, "#") || seen[wallet] {
			continue
		}
		seen[wallet] = true
		seeds = append(seeds, walletSeed{
			Wallet:        wallet,
			UniqueMarkets: map[string]bool{},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return seeds, nil
}

func presentationScore(nbaTrades, uniqueMarkets int, conviction, sizeRatio float64) float64 {
	score := float64(min(nbaTrades, 40))*2.5 + float64(min(uniqueMarkets, 12))*3
	if conviction >= 0.35 && conviction <= 0.8 {
		score += 8
	}
	if sizeRatio > 0 {
		score += minFloat(sizeRatio*3, 12)
	}
	return score
}

func buildReason(nbaTrades, uniqueMarkets int, conviction, sizeRatio float64) string {
	reasons := []string{}
	if nbaTrades >= 20 {
		reasons = append(reasons, "large NBA sample")
	}
	if uniqueMarkets >= 5 {
		reasons = append(reasons, "diverse market coverage")
	}
	if conviction >= 0.35 && conviction <= 0.8 {
		reasons = append(reasons, "balanced conviction profile")
	}
	if sizeRatio > 0.05 {
		reasons = append(reasons, "visible position sizing")
	}
	if len(reasons) == 0 {
		return "worth checking manually; enough NBA activity to show on the dashboard"
	}
	return strings.Join(reasons, ", ")
}

func countUniqueMarkets(trades []polymarket.EnrichedTrade) int {
	seen := map[string]bool{}
	for _, trade := range trades {
		if trade.ConditionID != "" {
			seen[trade.ConditionID] = true
		}
	}
	return len(seen)
}

func shortWallet(wallet string) string {
	if len(wallet) < 10 {
		return wallet
	}
	return wallet[:6] + "..." + wallet[len(wallet)-4:]
}

func isSportText(text, sport string) bool {
	value := strings.ToLower(text)
	if strings.Contains(value, sport) {
		return true
	}
	return sport == "nba" && strings.Contains(value, "basketball")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
