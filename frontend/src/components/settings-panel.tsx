"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

const STORAGE_KEY = "candidateWallets";
const ETH_ADDRESS = /^0x[a-fA-F0-9]{40}$/;

interface CandidateResult {
  wallet: string;
  display_name: string;
  recent_trades: number;
  recent_markets: number;
  nba_trades: number;
  entry_timing_hours: number;
  size_ratio_pct: number;
  conviction: number;
  style_label: string;
  presentation_score: number;
  reason: string;
}

interface SettingsPanelProps {
  onAnalyzeWallet: (wallet: string) => void;
}

export function SettingsPanel({ onAnalyzeWallet }: SettingsPanelProps) {
  const [walletText, setWalletText] = useState("");
  const [savedWallets, setSavedWallets] = useState<string[]>([]);
  const [discoverResults, setDiscoverResults] = useState<CandidateResult[]>([]);
  const [recentLimit, setRecentLimit] = useState("400");
  const [recentPages, setRecentPages] = useState("4");
  const [isSaving, setIsSaving] = useState(false);
  const [isScoring, setIsScoring] = useState(false);
  const [isDiscovering, setIsDiscovering] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return;
    try {
      const parsed = JSON.parse(stored) as string[];
      setSavedWallets(parsed);
      setWalletText(parsed.join("\n"));
    } catch {
      localStorage.removeItem(STORAGE_KEY);
    }
  }, []);

  const parseWallets = (input: string) =>
    Array.from(
      new Set(
        input
          .split(/\n|,/)
          .map((line) => line.trim())
          .filter(Boolean)
      )
    );

  const validateWallets = (wallets: string[]) => {
    const invalid = wallets.filter((wallet) => !ETH_ADDRESS.test(wallet));
    if (invalid.length > 0) {
      setError(`Invalid wallet address: ${invalid[0]}`);
      return false;
    }
    setError("");
    return true;
  };

  const saveWallets = () => {
    const wallets = parseWallets(walletText);
    if (!validateWallets(wallets)) return;

    setIsSaving(true);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(wallets));
    setSavedWallets(wallets);
    setIsSaving(false);
  };

  const removeWallet = (wallet: string) => {
    const next = savedWallets.filter((item) => item !== wallet);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
    setSavedWallets(next);
    setWalletText(next.join("\n"));
  };

  const scoreSavedWallets = async () => {
    if (savedWallets.length === 0) {
      setError("Add at least one wallet first");
      return;
    }

    setIsScoring(true);
    setError("");
    try {
      const response = await fetch("/api/discover-wallets", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          mode: "wallets",
          wallets: savedWallets,
          output_limit: 12,
        }),
      });
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || "Failed to score wallets");
      }
      setDiscoverResults(data.results || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to score wallets");
    } finally {
      setIsScoring(false);
    }
  };

  const autoDiscover = async () => {
    setIsDiscovering(true);
    setError("");
    try {
      const response = await fetch("/api/discover-wallets", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          mode: "auto",
          recent_limit: Number(recentLimit) || 400,
          recent_pages: Number(recentPages) || 4,
          candidate_limit: 20,
          output_limit: 12,
        }),
      });
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || "Failed to auto-discover wallets");
      }
      setDiscoverResults(data.results || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to auto-discover wallets"
      );
    } finally {
      setIsDiscovering(false);
    }
  };

  return (
    <div className="w-full max-w-6xl grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
      <Card className="border-white/10 bg-white/5">
        <CardHeader>
          <CardTitle className="text-white">Wallet Library</CardTitle>
          <p className="text-sm text-white/50">
            Paste candidate wallet addresses for demo prep. They are stored in
            your browser.
          </p>
        </CardHeader>
        <CardContent className="space-y-4">
          <textarea
            value={walletText}
            onChange={(e) => setWalletText(e.target.value)}
            placeholder={"0x...\n0x...\n0x..."}
            className="min-h-52 w-full rounded-xl border border-white/15 bg-black/20 px-4 py-3 font-mono text-sm text-white outline-none placeholder:text-white/25"
          />
          <div className="flex flex-wrap gap-3">
            <Button
              onClick={saveWallets}
              disabled={isSaving}
              className="bg-white text-gray-950 hover:bg-white/90"
            >
              {isSaving ? "Saving..." : "Save Wallets"}
            </Button>
            <Button
              onClick={scoreSavedWallets}
              disabled={isScoring || savedWallets.length === 0}
              className="bg-gradient-to-r from-blue-500 to-cyan-500 text-white hover:from-blue-600 hover:to-cyan-600"
            >
              {isScoring ? "Scoring..." : "Score Saved Wallets"}
            </Button>
          </div>
          {savedWallets.length > 0 && (
            <div className="space-y-2">
              <p className="text-xs uppercase tracking-wide text-white/40">
                Saved wallets ({savedWallets.length})
              </p>
              <div className="max-h-52 space-y-2 overflow-auto">
                {savedWallets.map((wallet) => (
                  <div
                    key={wallet}
                    className="flex items-center gap-3 rounded-xl border border-white/10 bg-black/20 px-3 py-2"
                  >
                    <code className="flex-1 text-xs text-white/70">{wallet}</code>
                    <Button
                      size="sm"
                      onClick={() => onAnalyzeWallet(wallet)}
                      className="bg-white/10 text-white hover:bg-white/15"
                    >
                      Analyze
                    </Button>
                    <Button
                      size="sm"
                      onClick={() => removeWallet(wallet)}
                      className="bg-red-500/15 text-red-200 hover:bg-red-500/25"
                    >
                      Remove
                    </Button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <Card className="border-white/10 bg-white/5">
        <CardHeader>
          <CardTitle className="text-white">Auto Discover</CardTitle>
          <p className="text-sm text-white/50">
            Scan recent trades for demo-friendly NBA wallets, then save the
            ones you like.
          </p>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-2">
              <label className="text-xs uppercase tracking-wide text-white/40">
                Recent limit
              </label>
              <Input
                value={recentLimit}
                onChange={(e) => setRecentLimit(e.target.value)}
                className="bg-white/10 border-white/15 text-white"
              />
            </div>
            <div className="space-y-2">
              <label className="text-xs uppercase tracking-wide text-white/40">
                Pages
              </label>
              <Input
                value={recentPages}
                onChange={(e) => setRecentPages(e.target.value)}
                className="bg-white/10 border-white/15 text-white"
              />
            </div>
          </div>
          <Button
            onClick={autoDiscover}
            disabled={isDiscovering}
            className="w-full bg-gradient-to-r from-emerald-500 to-teal-500 text-white hover:from-emerald-600 hover:to-teal-600"
          >
            {isDiscovering ? "Discovering..." : "Auto Discover Candidates"}
          </Button>
          {error && <p className="text-sm text-red-300">{error}</p>}

          <div className="space-y-3">
            <p className="text-xs uppercase tracking-wide text-white/40">
              Ranked candidates
            </p>
            <div className="max-h-[28rem] space-y-3 overflow-auto">
              {discoverResults.length === 0 && (
                <p className="text-sm text-white/35">
                  No results yet. Run discovery or score your saved wallets.
                </p>
              )}
              {discoverResults.map((result) => (
                <div
                  key={result.wallet}
                  className="rounded-2xl border border-white/10 bg-black/20 p-4"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-white">
                        {result.display_name}
                      </p>
                      <code className="text-xs text-white/45">{result.wallet}</code>
                    </div>
                    <div className="rounded-full border border-cyan-400/30 bg-cyan-400/10 px-3 py-1 text-xs text-cyan-200">
                      Score {result.presentation_score.toFixed(1)}
                    </div>
                  </div>
                  <div className="mt-3 grid gap-2 text-xs text-white/60 sm:grid-cols-2">
                    <p>NBA trades: {result.nba_trades}</p>
                    <p>Recent sample hits: {result.recent_trades}</p>
                    <p>Style: {result.style_label}</p>
                    <p>Conviction: {result.conviction.toFixed(2)}</p>
                  </div>
                  <p className="mt-3 text-sm text-white/55">{result.reason}</p>
                  <div className="mt-4 flex flex-wrap gap-2">
                    <Button
                      size="sm"
                      onClick={() => onAnalyzeWallet(result.wallet)}
                      className="bg-white text-gray-950 hover:bg-white/90"
                    >
                      Analyze
                    </Button>
                    <Button
                      size="sm"
                      onClick={() => {
                        if (savedWallets.includes(result.wallet)) return;
                        const next = [...savedWallets, result.wallet];
                        localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
                        setSavedWallets(next);
                        setWalletText(next.join("\n"));
                      }}
                      className="bg-white/10 text-white hover:bg-white/15"
                    >
                      Save to Library
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
