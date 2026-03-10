"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { WalletInput } from "@/components/wallet-input";

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleAnalyze = async (input: string) => {
    setIsLoading(true);
    setError("");

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

      const result = await response.json();
      sessionStorage.setItem("analyzeResult", JSON.stringify(result));
      router.push("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
      setIsLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4 bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950">
      {/* Background decorations */}
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
            Analyze any Polymarket wallet&apos;s NBA trading style using Amazon Nova AI
          </p>
        </div>

        <WalletInput onSubmit={handleAnalyze} isLoading={isLoading} />

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-3 text-red-400 text-sm max-w-xl w-full">
            {error}
          </div>
        )}

        {isLoading && (
          <div className="text-center space-y-3 animate-pulse">
            <p className="text-white/50 text-sm">
              Nova AI is analyzing the wallet...
            </p>
            <div className="flex items-center justify-center gap-2">
              {["Resolving wallet", "Fetching trades", "Computing metrics", "Building report"].map(
                (step, i) => (
                  <span
                    key={step}
                    className="text-xs px-2 py-1 rounded bg-white/5 text-white/30"
                    style={{ animationDelay: `${i * 0.2}s` }}
                  >
                    {step}
                  </span>
                )
              )}
            </div>
          </div>
        )}

        <p className="text-xs text-white/30 mt-8">
          Powered by Amazon Nova
        </p>
      </div>
    </main>
  );
}
