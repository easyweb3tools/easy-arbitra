package leaderboard

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const nbaLeaderboardURL = "https://r.jina.ai/https://polymarketanalytics.com/traders?overallCategory=NBA&sortBy=rank&sortDesc=false"

var rowPattern = regexp.MustCompile(`^\|\s*\|\s*(\d+)\s*\|\s*\[([^\]]+)\]\(https://polymarketanalytics\.com/traders/(0x[a-f0-9]{40})\)\s*\|\s*([0-9,]+)\s*\|\s*([0-9,]+)\s*\|\s*\$([0-9,.\-]+)\s*\|\s*\$([0-9,.\-]+)\s*\|\s*([0-9.]+)%\s*\|\s*\$([0-9,.\-]+)\s*\|\s*\$([0-9,.\-]+)\s*\|$`)

type Entry struct {
	Rank             int
	DisplayName      string
	WalletAddress    string
	Predictions      int
	Wins             int
	VolumeUSD        float64
	LossUSD          float64
	WinRate          float64
	OpenPositionsUSD float64
	PnlUSD           float64
	FetchedAt        time.Time
}

func FetchNBALeaderboard(ctx context.Context, limit int) ([]Entry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, nbaLeaderboardURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build leaderboard request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch leaderboard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("leaderboard returned %d", resp.StatusCode)
	}

	fetchedAt := time.Now().UTC()
	entries := make([]Entry, 0, limit)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := rowPattern.FindStringSubmatch(line)
		if len(matches) != 11 {
			continue
		}

		entry, err := parseRow(matches, fetchedAt)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
		if limit > 0 && len(entries) >= limit {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan leaderboard: %w", err)
	}
	return entries, nil
}

func parseRow(matches []string, fetchedAt time.Time) (Entry, error) {
	rank, err := parseInt(matches[1])
	if err != nil {
		return Entry{}, err
	}
	predictions, err := parseInt(matches[4])
	if err != nil {
		return Entry{}, err
	}
	wins, err := parseInt(matches[5])
	if err != nil {
		return Entry{}, err
	}
	volumeUSD, err := parseFloat(matches[6])
	if err != nil {
		return Entry{}, err
	}
	lossUSD, err := parseFloat(matches[7])
	if err != nil {
		return Entry{}, err
	}
	winRate, err := parseFloat(matches[8])
	if err != nil {
		return Entry{}, err
	}
	openPositionsUSD, err := parseFloat(matches[9])
	if err != nil {
		return Entry{}, err
	}
	pnlUSD, err := parseFloat(matches[10])
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Rank:             rank,
		DisplayName:      matches[2],
		WalletAddress:    strings.ToLower(matches[3]),
		Predictions:      predictions,
		Wins:             wins,
		VolumeUSD:        volumeUSD,
		LossUSD:          lossUSD,
		WinRate:          winRate,
		OpenPositionsUSD: openPositionsUSD,
		PnlUSD:           pnlUSD,
		FetchedAt:        fetchedAt,
	}, nil
}

func parseInt(value string) (int, error) {
	normalized := strings.ReplaceAll(value, ",", "")
	return strconv.Atoi(normalized)
}

func parseFloat(value string) (float64, error) {
	normalized := strings.ReplaceAll(value, ",", "")
	return strconv.ParseFloat(normalized, 64)
}
