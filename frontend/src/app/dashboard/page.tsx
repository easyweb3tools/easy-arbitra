"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DecisionLog } from "@/components/decision-log";
import { WalletCard } from "@/components/wallet-card";
import { RadarChart } from "@/components/radar-chart";
import { ReportSummary } from "@/components/report-summary";
import type { AnalyzeResponse, DecisionStep, ToolLogEntry } from "@/lib/types";

const STEP_LABELS = [
  "Resolving wallet",
  "Fetching trades",
  "Computing metrics",
  "Building report",
];

export default function Dashboard() {
  const router = useRouter();
  const [data, setData] = useState<AnalyzeResponse | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [error, setError] = useState("");
  const [completedSteps, setCompletedSteps] = useState<DecisionStep[]>([]);
  const [activeStep, setActiveStep] = useState(-1);
  const [toolLogs, setToolLogs] = useState<ToolLogEntry[]>([]);
  const hasStarted = useRef(false);

  useEffect(() => {
    // If we already have results, show them
    const stored = sessionStorage.getItem("analyzeResult");
    if (stored) {
      try {
        setData(JSON.parse(stored));
        return;
      } catch {
        // fall through to check walletInput
      }
    }

    // If we have a wallet input, start analysis
    const walletInput = sessionStorage.getItem("walletInput");
    if (!walletInput) {
      router.push("/");
      return;
    }

    // Prevent double-start in strict mode
    if (hasStarted.current) return;
    hasStarted.current = true;

    sessionStorage.removeItem("walletInput");
    startAnalysis(walletInput);
  }, [router]);

  const startAnalysis = async (input: string) => {
    setIsAnalyzing(true);
    setError("");
    setCompletedSteps([]);
    setActiveStep(0);
    setToolLogs([]);

    try {
      const response = await fetch("/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ walletInput: input }),
      });

      if (!response.ok) {
        const respData = await response.json();
        throw new Error(respData.error || "Analysis failed");
      }

      const reader = response.body?.getReader();
      if (!reader) throw new Error("No response stream");

      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          const dataLine = line.replace(/^data: /, "").trim();
          if (!dataLine) continue;

          const event = JSON.parse(dataLine);

          switch (event.type) {
            case "step":
              setCompletedSteps((prev) => [...prev, event.data]);
              setActiveStep(event.data.step);
              break;
            case "tool_log":
              setToolLogs((prev) => [...prev, event.data]);
              break;
            case "error":
              throw new Error(event.data);
            case "done":
              sessionStorage.setItem(
                "analyzeResult",
                JSON.stringify(event.data)
              );
              setData(event.data);
              setIsAnalyzing(false);
              return;
          }
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
      setIsAnalyzing(false);
      setActiveStep(-1);
    }
  };

  // Analyzing state: show progress
  if (isAnalyzing || (!data && !error)) {
    return (
      <main className="min-h-screen bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950 p-6">
        <div className="max-w-7xl mx-auto mb-8">
          <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
            SportStyle AI Explainer
          </h1>
          <p className="text-sm text-white/40 mt-1">Analyzing wallet...</p>
        </div>

        <div className="max-w-md mx-auto space-y-3 mt-16">
          {STEP_LABELS.map((label, i) => {
            const completed = completedSteps.find((s) => s.step === i + 1);
            const isActive = activeStep === i || (activeStep > i && !completed);
            const isPending = activeStep < i;
            const toolName = [
              "resolve_wallet_target",
              "fetch_sports_trades",
              "calculate_style_metrics",
              "build_report_payload",
            ][i];
            const latestLog = [...toolLogs]
              .reverse()
              .find((entry) => entry.tool === toolName)?.message;
            return (
              <div key={label}>
                <div
                  className={`flex items-center gap-3 px-4 py-2 rounded-lg transition-all duration-300 ${
                    completed
                      ? "bg-green-500/10 border border-green-500/20"
                      : isActive
                        ? "bg-blue-500/10 border border-blue-500/20"
                        : "bg-white/5 border border-white/5"
                  }`}
                >
                  <div className="w-6 h-6 flex items-center justify-center">
                    {completed ? (
                      <svg
                        className="w-5 h-5 text-green-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                    ) : isActive ? (
                      <span className="h-4 w-4 border-2 border-blue-400/30 border-t-blue-400 rounded-full animate-spin" />
                    ) : (
                      <span
                        className={`text-sm font-mono ${isPending ? "text-white/20" : "text-white/50"}`}
                      >
                        {i + 1}
                      </span>
                    )}
                  </div>
                  <span
                    className={`text-sm ${
                      completed
                        ? "text-green-300"
                        : isActive
                          ? "text-blue-300"
                          : "text-white/30"
                    }`}
                  >
                    {label}
                  </span>
                  {completed?.result_summary && (
                    <span className="ml-auto text-xs text-white/40 truncate max-w-[200px]">
                      {completed.result_summary}
                    </span>
                  )}
                </div>
                {!completed?.result_summary && latestLog && (
                  <div className="ml-9 mt-1 text-xs text-white/45">
                    {latestLog}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {toolLogs.length > 0 && (
          <div className="max-w-2xl mx-auto mt-8 rounded-lg border border-white/10 bg-white/5 p-4">
            <p className="text-xs font-medium uppercase tracking-wide text-white/50">
              Live backend log
            </p>
            <div className="mt-3 max-h-64 space-y-2 overflow-auto font-mono text-xs text-white/60">
              {toolLogs.map((entry, idx) => (
                <div key={`${entry.timestamp}-${idx}`}>
                  [{new Date(entry.timestamp).toLocaleTimeString()}] {entry.tool}:{" "}
                  {entry.message}
                </div>
              ))}
            </div>
          </div>
        )}

        {error && (
          <div className="max-w-md mx-auto mt-6 bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 text-red-400 text-sm">
            {error}
            <button
              onClick={() => router.push("/")}
              className="block mt-2 text-white/50 hover:text-white/80 underline"
            >
              Try again
            </button>
          </div>
        )}
      </main>
    );
  }

  // Error state without data
  if (error && !data) {
    return (
      <main className="min-h-screen flex flex-col items-center justify-center bg-gray-950 px-4">
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg px-6 py-4 text-red-400 text-center max-w-md">
          <p>{error}</p>
          <button
            onClick={() => router.push("/")}
            className="mt-4 text-white/50 hover:text-white/80 underline text-sm"
          >
            Back to home
          </button>
        </div>
      </main>
    );
  }

  // Results state
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
                Analysis Decision Log
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
          Powered by an OpenAI-compatible LLM
        </p>
      </div>
    </main>
  );
}
