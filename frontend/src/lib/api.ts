import {
  AIReport,
  AnomalyAlert,
  ApiEnvelope,
  CopyTradeDashboard,
  CopyTradeDecision,
  CopyTradeMonitor,
  CopyTradePerformance,
  CopyTradingConfig,
  LeaderboardItem,
  Market,
  OpsHighlights,
  OverviewStats,
  Paged,
  PnLHistoryPoint,
  PortfolioItem,
  PotentialWallet,
  TradeHistoryItem,
  Wallet,
  WalletDecisionCard,
  WalletPosition,
  WatchlistFeedItem,
  WatchlistItem,
  WatchlistSummary,
  WalletExplanation,
  WalletProfile,
  WalletShareLanding,
  WalletShareCard
} from "@/lib/types";

const API_BASE =
  typeof window === "undefined"
    ? process.env.API_SERVER_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://backend:8080/api/v1"
    : process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1";

async function getJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { cache: "no-store", ...init });
  if (!res.ok) {
    throw new Error(`API ${res.status}`);
  }
  const body = (await res.json()) as ApiEnvelope<T>;
  return body.data;
}

export function getWallets(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<Wallet>>(`/wallets${q ? `?${q}` : ""}`);
}

export function getPotentialWallets(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<PotentialWallet>>(`/wallets/potential${q ? `?${q}` : ""}`);
}

export function getWalletProfile(id: string) {
  return getJSON<WalletProfile>(`/wallets/${id}/profile`);
}

export function getWalletShareCard(id: string) {
  return getJSON<WalletShareCard>(`/wallets/${id}/share-card`);
}

export function getWalletDecisionCard(id: string) {
  return getJSON<WalletDecisionCard>(`/wallets/${id}/decision-card`);
}

export function getWalletShareLanding(id: string) {
  return getJSON<WalletShareLanding>(`/wallets/${id}/share-landing`);
}

export function getWalletExplanation(id: string) {
  return getJSON<WalletExplanation>(`/wallets/${id}/explanations`);
}

export function getWalletAIReport(id: string) {
  return getJSON<AIReport>(`/ai/report/${id}`);
}

export function getWalletAIReportHistory(id: string) {
  return getJSON<AIReport[]>(`/ai/report/${id}/history?limit=5`);
}

export function getMarkets(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<Market>>(`/markets${q ? `?${q}` : ""}`);
}

export function getOverviewStats() {
  return getJSON<OverviewStats>("/stats/overview");
}

export function getOpsHighlights(params?: URLSearchParams) {
  const q = params?.toString() || "";
  return getJSON<OpsHighlights>(`/ops/highlights${q ? `?${q}` : ""}`);
}

export function getLeaderboard(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<LeaderboardItem>>(`/leaderboard${q ? `?${q}` : ""}`);
}

export function getPortfolios() {
  return getJSON<PortfolioItem[]>("/portfolios");
}

export function getAnomalies(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<AnomalyAlert>>(`/anomalies${q ? `?${q}` : ""}`);
}

export function getAnomaly(id: string) {
  return getJSON<AnomalyAlert>(`/anomalies/${id}`);
}

export function getWatchlist(params: URLSearchParams, fingerprint: string) {
  const q = params.toString();
  return getJSON<Paged<WatchlistItem>>(`/watchlist${q ? `?${q}` : ""}`, {
    headers: { "X-User-Fingerprint": fingerprint }
  });
}

export function getWatchlistSummary(fingerprint: string) {
  return getJSON<WatchlistSummary>("/watchlist/summary", {
    headers: { "X-User-Fingerprint": fingerprint }
  });
}

export function getWatchlistFeed(params: URLSearchParams, fingerprint: string) {
  const q = params.toString();
  return getJSON<Paged<WatchlistFeedItem>>(`/watchlist/feed${q ? `?${q}` : ""}`, {
    headers: { "X-User-Fingerprint": fingerprint }
  });
}

export async function addBatchToWatchlist(walletIDs: number[], fingerprint: string) {
  return getJSON<{ wallet_ids: number[]; watching: boolean }>(`/watchlist/batch`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User-Fingerprint": fingerprint
    },
    body: JSON.stringify({ wallet_ids: walletIDs })
  });
}

export async function addToWatchlist(walletID: number, fingerprint: string) {
  return getJSON<{ wallet_id: number; watching: boolean }>(`/watchlist`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-User-Fingerprint": fingerprint
    },
    body: JSON.stringify({ wallet_id: walletID })
  });
}

export async function removeFromWatchlist(walletID: number, fingerprint: string) {
  return getJSON<{ wallet_id: number; watching: boolean }>(`/watchlist/${walletID}`, {
    method: "DELETE",
    headers: { "X-User-Fingerprint": fingerprint }
  });
}

export function getWalletPnLHistory(id: string, limit = 90) {
  return getJSON<PnLHistoryPoint[]>(`/wallets/${id}/pnl-history?limit=${limit}`);
}

export function getWalletTrades(id: string, params?: URLSearchParams) {
  const q = params?.toString() || "";
  return getJSON<Paged<TradeHistoryItem>>(`/wallets/${id}/trades${q ? `?${q}` : ""}`);
}

export function getWalletPositions(id: string) {
  return getJSON<WalletPosition[]>(`/wallets/${id}/positions`);
}

// ── Copy Trading API ──

export function enableCopyTrading(walletID: number, maxPositionUSDC: number, riskPreference: string, fingerprint: string) {
  return getJSON<CopyTradingConfig>("/copy-trading/enable", {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-User-Fingerprint": fingerprint },
    body: JSON.stringify({ wallet_id: walletID, max_position_usdc: maxPositionUSDC, risk_preference: riskPreference }),
  });
}

export function disableCopyTrading(walletID: number, fingerprint: string) {
  return getJSON<{ disabled: boolean }>("/copy-trading/disable", {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-User-Fingerprint": fingerprint },
    body: JSON.stringify({ wallet_id: walletID }),
  });
}

export function updateCopyTradeSettings(walletID: number, maxPositionUSDC: number, riskPreference: string, fingerprint: string) {
  return getJSON<CopyTradingConfig>("/copy-trading/settings", {
    method: "PUT",
    headers: { "Content-Type": "application/json", "X-User-Fingerprint": fingerprint },
    body: JSON.stringify({ wallet_id: walletID, max_position_usdc: maxPositionUSDC, risk_preference: riskPreference }),
  });
}

export function getCopyTradeConfigs(fingerprint: string) {
  return getJSON<CopyTradingConfig[]>("/copy-trading/configs", {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradeConfig(walletID: number, fingerprint: string) {
  return getJSON<CopyTradingConfig>(`/copy-trading/${walletID}`, {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradeDashboard(fingerprint: string) {
  return getJSON<CopyTradeDashboard>("/copy-trading/dashboard", {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradeDecisions(walletID: number, params: URLSearchParams, fingerprint: string) {
  const q = params.toString();
  return getJSON<Paged<CopyTradeDecision>>(`/copy-trading/${walletID}/decisions${q ? `?${q}` : ""}`, {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradePerformance(walletID: number, fingerprint: string) {
  return getJSON<CopyTradePerformance>(`/copy-trading/${walletID}/performance`, {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradePositions(fingerprint: string) {
  return getJSON<CopyTradeDecision[]>("/copy-trading/positions", {
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function closeCopyTradePosition(decisionID: number, fingerprint: string) {
  return getJSON<CopyTradeDecision>(`/copy-trading/decisions/${decisionID}/close`, {
    method: "POST",
    headers: { "X-User-Fingerprint": fingerprint },
  });
}

export function getCopyTradeMonitor() {
  return getJSON<CopyTradeMonitor>("/copy-trading/monitor");
}
