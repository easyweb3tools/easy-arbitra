import {
  DailyPick,
  LeaderboardItem,
  Market,
  OverviewStats,
  Paged,
  PnLHistoryPoint,
  PotentialWallet,
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
