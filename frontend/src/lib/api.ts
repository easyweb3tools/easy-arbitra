import {
  CandidateScore,
  DailyPick,
  DecisionExplanation,
  LeaderboardItem,
  Market,
  NovaMemory,
  NovaSession,
  NovaStatus,
  OverviewStats,
  Paged,
  PnLHistoryPoint,
  PotentialWallet,
  ThinkingRound,
  TradeHistoryItem,
  Wallet,
  WalletPosition,
  WalletProfile,
} from "@/lib/types";

const API_BASE =
  typeof window === "undefined"
    ? process.env.API_SERVER_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || "http://backend:8080/api/v1"
    : process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1";

interface ApiEnvelope<T> {
  data: T;
}

async function getJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
    ...init,
  });
  if (!res.ok) {
    throw new Error(`API ${res.status}`);
  }
  const body = (await res.json()) as ApiEnvelope<T>;
  return body.data;
}

// ── Wallets ──

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

export function getWalletTrades(id: string, params?: URLSearchParams) {
  const q = params?.toString() || "";
  return getJSON<Paged<TradeHistoryItem>>(`/wallets/${id}/trades${q ? `?${q}` : ""}`);
}

export function getWalletPositions(id: string) {
  return getJSON<WalletPosition[]>(`/wallets/${id}/positions`);
}

// ── Markets ──

export function getMarkets(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<Market>>(`/markets${q ? `?${q}` : ""}`);
}

// ── Stats & Leaderboard ──

export function getOverviewStats() {
  return getJSON<OverviewStats>("/stats/overview");
}

export function getLeaderboard(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<LeaderboardItem>>(`/leaderboard${q ? `?${q}` : ""}`);
}

// ── Daily Pick ──

export function getDailyPick() {
  return getJSON<{ pick: DailyPick; wallet: Wallet }>("/daily-pick");
}

export function getDailyPickHistory(limit = 14) {
  return getJSON<DailyPick[]>(`/daily-pick/history?limit=${limit}`);
}

// ── Nova Sessions ──

export function getNovaSessions(date?: string) {
  const q = date ? `?date=${date}` : "";
  return getJSON<NovaSession[]>(`/nova/sessions${q}`);
}


// ── Nova Insight (AI Brain Visualization) ──

export function getNovaStatus() {
  return getJSON<NovaStatus>("/nova/status");
}

export function getNovaTimeline(date: string) {
  return getJSON<ThinkingRound[]>(`/nova/timeline/${date}`);
}

export function getNovaCandidates(date: string) {
  return getJSON<CandidateScore[]>(`/nova/candidates/${date}`);
}

export function getNovaDecisionExplanation(pickId: number) {
  return getJSON<DecisionExplanation>(`/nova/decision-explain/${pickId}`);
}


export function getNovaMemory(limit = 30) {
  return getJSON<NovaMemory>(`/nova/memory?limit=${limit}`);
}
