import {
  AIReport,
  AnomalyAlert,
  ApiEnvelope,
  LeaderboardItem,
  Market,
  OverviewStats,
  Paged,
  Wallet,
  WalletExplanation,
  WalletProfile
} from "@/lib/types";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

async function getJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, { cache: "no-store" });
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

export function getWalletProfile(id: string) {
  return getJSON<WalletProfile>(`/wallets/${id}/profile`);
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

export function getLeaderboard(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<LeaderboardItem>>(`/leaderboard${q ? `?${q}` : ""}`);
}

export function getAnomalies(params: URLSearchParams) {
  const q = params.toString();
  return getJSON<Paged<AnomalyAlert>>(`/anomalies${q ? `?${q}` : ""}`);
}

export function getAnomaly(id: string) {
  return getJSON<AnomalyAlert>(`/anomalies/${id}`);
}
