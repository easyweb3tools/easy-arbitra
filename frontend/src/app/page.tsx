"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { WalletInput } from "@/components/wallet-input";
import { SettingsPanel } from "@/components/settings-panel";

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [activeView, setActiveView] = useState<"analyze" | "settings">("analyze");

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
        <div className="text-center space-y-4">
          <h1 className="text-5xl font-bold bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
            SportStyle AI Explainer
          </h1>
          <p className="text-lg text-white/60 max-w-md mx-auto">
            Analyze any Polymarket wallet&apos;s NBA trading style with an
            OpenAI-compatible AI backend
          </p>
        </div>

        <div className="w-full max-w-5xl">
          {activeView === "analyze" ? (
            <WalletInput onSubmit={handleAnalyze} isLoading={isLoading} />
          ) : (
            <SettingsPanel onAnalyzeWallet={handleAnalyze} />
          )}
        </div>

        <p className="text-xs text-white/30 mt-8">
          Powered by an OpenAI-compatible LLM
        </p>
        </div>
      </div>
    </main>
  );
}
