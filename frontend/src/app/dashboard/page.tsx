"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DecisionLog } from "@/components/decision-log";
import { WalletCard } from "@/components/wallet-card";
import { RadarChart } from "@/components/radar-chart";
import { ReportSummary } from "@/components/report-summary";
import type { AnalyzeResponse } from "@/lib/types";

export default function Dashboard() {
  const router = useRouter();
  const [data, setData] = useState<AnalyzeResponse | null>(null);

  useEffect(() => {
    const stored = sessionStorage.getItem("analyzeResult");
    if (!stored) {
      router.push("/");
      return;
    }
    try {
      setData(JSON.parse(stored));
    } catch {
      router.push("/");
    }
  }, [router]);

  if (!data) {
    return (
      <main className="min-h-screen flex items-center justify-center bg-gray-950">
        <div className="h-8 w-8 border-2 border-white/30 border-t-white rounded-full animate-spin" />
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950 p-6">
      {/* Header */}
      <div className="max-w-7xl mx-auto mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              SportStyle AI Explainer
            </h1>
            <p className="text-sm text-white/40 mt-1">NBA Trading Style Dashboard</p>
          </div>
          <button
            onClick={() => {
              sessionStorage.removeItem("analyzeResult");
              router.push("/");
            }}
            className="text-sm text-white/50 hover:text-white/80 px-4 py-2 rounded-lg border border-white/10 hover:border-white/20 transition-colors"
          >
            New Analysis
          </button>
        </div>
      </div>

      {/* Dashboard Grid */}
      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column: Decision Log */}
        <div className="lg:col-span-1">
          <Card className="bg-white/5 border-white/10 h-full">
            <CardHeader>
              <CardTitle className="text-white/80 text-sm font-medium">
                Nova Decision Log
              </CardTitle>
            </CardHeader>
            <CardContent>
              <DecisionLog steps={data.decisionLog} />
            </CardContent>
          </Card>
        </div>

        {/* Right Column: Cards */}
        <div className="lg:col-span-2 space-y-6">
          {/* Wallet Card */}
          <WalletCard data={data.reportPayload.wallet_card} />

          {/* Radar Chart */}
          <Card className="bg-white/5 border-white/10">
            <CardHeader>
              <CardTitle className="text-white/80 text-sm font-medium">
                Trading Style Radar
              </CardTitle>
            </CardHeader>
            <CardContent>
              <RadarChart data={data.reportPayload.radar_chart} />
            </CardContent>
          </Card>

          {/* Report Summary */}
          <ReportSummary
            report={data.reportPayload.report}
            explanation={data.explanation}
          />
        </div>
      </div>

      {/* Footer */}
      <div className="max-w-7xl mx-auto mt-8 text-center">
        <p className="text-xs text-white/30">
          Powered by Amazon Nova
        </p>
      </div>
    </main>
  );
}
