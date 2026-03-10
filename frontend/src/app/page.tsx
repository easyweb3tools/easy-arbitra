"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { WalletInput } from "@/components/wallet-input";
import type { DecisionStep } from "@/lib/types";

const STEP_LABELS = [
  "Resolving wallet",
  "Fetching trades",
  "Computing metrics",
  "Building report",
];

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [completedSteps, setCompletedSteps] = useState<DecisionStep[]>([]);
  const [activeStep, setActiveStep] = useState(-1);

  const handleAnalyze = async (input: string) => {
    setIsLoading(true);
    setError("");
    setCompletedSteps([]);
    setActiveStep(0);

    try {
      const response = await fetch("/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ walletInput: input }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || "Analysis failed");
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
            case "error":
              throw new Error(event.data);
            case "done":
              sessionStorage.setItem(
                "analyzeResult",
                JSON.stringify(event.data)
              );
              router.push("/dashboard");
              return;
          }
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
      setIsLoading(false);
      setActiveStep(-1);
    }
  };

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4 bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950">
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 -left-32 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl" />
        <div className="absolute bottom-1/4 -right-32 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 flex flex-col items-center gap-8 max-w-2xl w-full">
        <div className="text-center space-y-4">
          <h1 className="text-5xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
            SportStyle AI Explainer
          </h1>
          <p className="text-lg text-white/60 max-w-md mx-auto">
            Analyze any Polymarket wallet&apos;s NBA trading style using Amazon
            Nova AI
          </p>
        </div>

        <WalletInput onSubmit={handleAnalyze} isLoading={isLoading} />

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 text-red-400 text-sm max-w-xl w-full">
            {error}
          </div>
        )}

        {isLoading && (
          <div className="w-full max-w-md space-y-3">
            {STEP_LABELS.map((label, i) => {
              const completed = completedSteps.find((s) => s.step === i + 1);
              const isActive = activeStep === i || (activeStep > i && !completed);
              const isPending = activeStep < i;
              return (
                <div
                  key={label}
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
              );
            })}
          </div>
        )}

        <p className="text-xs text-white/30 mt-8">Powered by Amazon Nova</p>
      </div>
    </main>
  );
}
