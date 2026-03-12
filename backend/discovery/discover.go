package discovery

import (
	"context"
	"strings"

	"github.com/brucexwang/easy-arbitra/backend/metrics"
	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/brucexwang/easy-arbitra/backend/tools"
)

type Options struct {
	Sport           string
	RecentLimit     int
	RecentPages     int
	CandidateLimit  int
	OutputLimit     int
	MinRecentTrades int
	WalletLimit     int
}

type Candidate struct {
	Wallet            string  `json:"wallet"`
	DisplayName       string  `json:"display_name"`
	RecentTrades      int     `json:"recent_trades"`
	RecentMarkets     int     `json:"recent_markets"`
	NbaTrades         int     `json:"nba_trades"`
	EntryTimingHours  float64 `json:"entry_timing_hours"`
	SizeRatioPct      float64 `json:"size_ratio_pct"`
	Conviction        float64 `json:"conviction"`
	StyleLabel        string  `json:"style_label"`
	PresentationScore float64 `json:"presentation_score"`
	Reason            string  `json:"reason"`
}

type walletSeed struct {
	Wallet          string
	RecentTrades    int
	UniqueMarkets   map[string]bool
	RecentBuyVolume float64
}

func DiscoverFromRecent(ctx context.Context, client *polymarket.Client, opts Options) ([]Candidate, error) {
	seeds, err := discoverSeeds(client, opts.Sport, opts.RecentLimit, opts.RecentPages, opts.MinRecentTrades)
	if err != nil {
		return nil, err
	}
	return scoreSeeds(ctx, client, seeds, opts), nil
}

func ScoreWallets(ctx context.Context, client *polymarket.Client, wallets []string, opts Options) ([]Candidate, error) {
	seeds := make([]walletSeed, 0, len(wallets))
	seen := map[string]bool{}
	for _, wallet := range wallets {
		trimmed := strings.TrimSpace(wallet)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		seeds = append(seeds, walletSeed{
			Wallet:        trimmed,
			UniqueMarkets: map[string]bool{},
		})
	}
	return scoreSeeds(ctx, client, seeds, opts), nil
}

func discoverSeeds(client *polymarket.Client, sport string, limit, pages, minRecentTrades int) ([]walletSeed, error) {
	seeds := map[string]*walletSeed{}

	for page := 0; page < pages; page++ {
		offset := page * limit
		trades, err := client.GetRecentTrades(limit, offset)
		if err != nil {
			return nil, err
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

func scoreSeeds(ctx context.Context, client *polymarket.Client, seeds []walletSeed, opts Options) []Candidate {
	if len(seeds) > opts.CandidateLimit && opts.CandidateLimit > 0 {
		seeds = seeds[:opts.CandidateLimit]
	}

	results := make([]Candidate, 0, len(seeds))
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

		fetchResult, err := tools.FetchSportsTradesData(ctx, client, seed.Wallet, opts.Sport, opts.WalletLimit)
		if err != nil || fetchResult.TotalTrades == 0 {
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

		results = append(results, Candidate{
			Wallet:            seed.Wallet,
			DisplayName:       displayName,
			RecentTrades:      seed.RecentTrades,
			RecentMarkets:     len(seed.UniqueMarkets),
			NbaTrades:         fetchResult.TotalTrades,
			EntryTimingHours:  entryTiming,
			SizeRatioPct:      sizeRatio,
			Conviction:        conviction,
			StyleLabel:        styleLabel,
			PresentationScore: presentationScore(fetchResult.TotalTrades, uniqueMarkets, conviction, sizeRatio),
			Reason:            buildReason(fetchResult.TotalTrades, uniqueMarkets, conviction, sizeRatio),
		})
	}

	sortCandidates(results)
	if len(results) > opts.OutputLimit && opts.OutputLimit > 0 {
		results = results[:opts.OutputLimit]
	}
	return results
}

func sortCandidates(results []Candidate) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].PresentationScore > results[i].PresentationScore ||
				(results[j].PresentationScore == results[i].PresentationScore && results[j].NbaTrades > results[i].NbaTrades) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
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
