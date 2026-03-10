package metrics

import (
	"math"
	"time"

	"github.com/brucexwang/easy-arbitra/backend/polymarket"
)

// EntryTimingHours calculates the average hours between market start and trade time.
func EntryTimingHours(trades []polymarket.EnrichedTrade) float64 {
	if len(trades) == 0 {
		return 0
	}
	var totalHours float64
	validCount := 0
	for _, t := range trades {
		tradeTime, err1 := time.Parse(time.RFC3339, t.TradeTime)
		marketStart, err2 := time.Parse(time.RFC3339, t.MarketStartTime)
		if err1 != nil || err2 != nil {
			continue
		}
		hours := tradeTime.Sub(marketStart).Hours()
		if hours >= 0 {
			totalHours += hours
			validCount++
		}
	}
	if validCount == 0 {
		return 0
	}
	return math.Round(totalHours/float64(validCount)*100) / 100
}

// SizeRatioPct calculates the average trade size relative to market volume (%).
func SizeRatioPct(trades []polymarket.EnrichedTrade) float64 {
	if len(trades) == 0 {
		return 0
	}
	var totalRatio float64
	validCount := 0
	for _, t := range trades {
		if t.MarketVolume <= 0 {
			continue
		}
		ratio := (t.Size * t.Price) / t.MarketVolume * 100
		totalRatio += ratio
		validCount++
	}
	if validCount == 0 {
		return 0
	}
	return math.Round(totalRatio/float64(validCount)*10000) / 10000
}

// Conviction calculates the average BUY price (0-1 scale).
// High conviction (>0.7) = bets on favorites, low conviction (<0.4) = contrarian/underdog bets.
// This is directly measurable from trade data without requiring settlement info.
func Conviction(trades []polymarket.EnrichedTrade) float64 {
	var totalPrice float64
	var buyCount int
	for _, t := range trades {
		if t.Side == "BUY" && t.Price > 0 && t.Price <= 1 {
			totalPrice += t.Price
			buyCount++
		}
	}
	if buyCount == 0 {
		return 0
	}
	return math.Round(totalPrice/float64(buyCount)*100) / 100
}
