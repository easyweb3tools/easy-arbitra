"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { WalletInput } from "@/components/wallet-input";
import { SettingsPanel } from "@/components/settings-panel";

interface StyleWallet {
  wallet_address: string;
  display_name: string;
  source_rank: number;
  win_rate: number;
  pnl_usd: number;
  nba_trades: number;
  style_label: string;
  style_summary: string;
  explanation_source: "ai" | "fallback";
}

interface StyleGroup {
  label: string;
  wallets: StyleWallet[];
}

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [activeView, setActiveView] = useState<"analyze" | "settings">("analyze");
  const [styleGroups, setStyleGroups] = useState<StyleGroup[]>([]);
  const [isLoadingGroups, setIsLoadingGroups] = useState(true);
  const [groupsError, setGroupsError] = useState("");

  useEffect(() => {
    let cancelled = false;

    const loadStyleGroups = async () => {
      setIsLoadingGroups(true);
      setGroupsError("");
      try {
        const response = await fetch("/api/style-wallets?limit_per_group=6", {
          cache: "no-store",
        });
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.error || "Failed to load style wallets");
        }
        if (!cancelled) {
          setStyleGroups(data.groups || []);
        }
      } catch (error) {
        if (!cancelled) {
          setGroupsError(
            error instanceof Error ? error.message : "Failed to load style wallets"
          );
        }
      } finally {
        if (!cancelled) {
          setIsLoadingGroups(false);
        }
      }
    };

    void loadStyleGroups();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleAnalyze = (input: string) => {
    setIsLoading(true);
    sessionStorage.setItem("walletInput", input);
    sessionStorage.removeItem("analyzeResult");
    router.push("/dashboard");
  };

  return (
    <main className="min-h-screen px-4 bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950">
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 -left-32 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 -right-32 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 mx-auto flex min-h-screen max-w-6xl flex-col py-6">
        <div className="flex justify-end">
          <button
            onClick={() =>
              setActiveView((current) =>
                current === "analyze" ? "settings" : "analyze"
              )
            }
            className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white/70 transition-colors hover:border-white/20 hover:text-white"
          >
            {activeView === "analyze" ? "Settings" : "Back to Analyze"}
          </button>
        </div>

        <div className="flex flex-1 flex-col items-center justify-center gap-8">
          <div className="space-y-4 text-center">
            <h1 className="text-5xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
              SportStyle AI Explainer
            </h1>
            <p className="mx-auto max-w-md text-lg text-white/60">
              Analyze any Polymarket wallet&apos;s NBA trading style with an
              OpenAI-compatible AI backend
            </p>
          </div>

          <div className="flex w-full max-w-5xl justify-center">
            {activeView === "analyze" ? (
              <div className="w-full space-y-10">
                <WalletInput onSubmit={handleAnalyze} isLoading={isLoading} />

                <section className="space-y-5">
                  <div className="text-center">
                    <p className="text-xs uppercase tracking-[0.24em] text-white/35">
                      Style Leaderboard
                    </p>
                    <h2 className="mt-2 text-2xl font-semibold text-white">
                      Top NBA wallets grouped by trading style
                    </h2>
                  </div>

                  {isLoadingGroups ? (
                    <div className="flex justify-center py-8">
                      <span className="h-6 w-6 rounded-full border-2 border-white/20 border-t-white/70 animate-spin" />
                    </div>
                  ) : groupsError ? (
                    <div className="rounded-2xl border border-red-400/20 bg-red-500/10 px-4 py-3 text-sm text-red-200">
                      {groupsError}
                    </div>
                  ) : styleGroups.length === 0 ? (
                    <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-8 text-center text-sm text-white/40">
                      No tagged wallets yet. Run the backend sync job first.
                    </div>
                  ) : (
                    <div className="space-y-8">
                      {styleGroups.map((group) => (
                        <section key={group.label} className="space-y-4">
                          <div className="flex items-center justify-between gap-4">
                            <div>
                              <h3 className="text-xl font-semibold text-white">
                                {group.label}
                              </h3>
                              <p className="text-sm text-white/45">
                                Click a wallet to open a full analysis.
                              </p>
                            </div>
                            <div className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs uppercase tracking-wide text-white/45">
                              {group.wallets.length} wallets
                            </div>
                          </div>

                          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                            {group.wallets.map((wallet) => (
                              <button
                                key={wallet.wallet_address}
                                onClick={() => handleAnalyze(wallet.wallet_address)}
                                className="rounded-2xl border border-white/10 bg-white/5 p-4 text-left transition-colors hover:border-white/20 hover:bg-white/[0.07]"
                              >
                                <div className="flex items-start justify-between gap-3">
                                  <div>
                                    <p className="text-sm font-semibold text-white">
                                      {wallet.display_name}
                                    </p>
                                    <p className="mt-1 font-mono text-xs text-white/40">
                                      {wallet.wallet_address}
                                    </p>
                                  </div>
                                  <span className="rounded-full border border-cyan-400/25 bg-cyan-400/10 px-2 py-1 text-[10px] uppercase tracking-wide text-cyan-200">
                                    #{wallet.source_rank}
                                  </span>
                                </div>

                                <div className="mt-4 grid grid-cols-2 gap-2 text-xs text-white/55">
                                  <p>NBA trades: {wallet.nba_trades}</p>
                                  <p>Win rate: {wallet.win_rate.toFixed(1)}%</p>
                                  <p>PnL: ${wallet.pnl_usd.toLocaleString()}</p>
                                  <p>
                                    Source:{" "}
                                    {wallet.explanation_source === "ai"
                                      ? "AI"
                                      : "Fallback"}
                                  </p>
                                </div>

                                <p className="mt-4 text-sm text-white/65">
                                  {wallet.style_summary}
                                </p>
                              </button>
                            ))}
                          </div>
                        </section>
                      ))}
                    </div>
                  )}
                </section>
              </div>
            ) : (
              <SettingsPanel onAnalyzeWallet={handleAnalyze} />
            )}
          </div>

          <p className="mt-8 text-xs text-white/30">
            Powered by an OpenAI-compatible LLM
          </p>
        </div>
      </div>
    </main>
  );
}
