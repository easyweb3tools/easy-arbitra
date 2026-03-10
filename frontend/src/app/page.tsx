"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { WalletInput } from "@/components/wallet-input";

export default function Home() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const handleAnalyze = (input: string) => {
    setIsLoading(true);
    sessionStorage.setItem("walletInput", input);
    sessionStorage.removeItem("analyzeResult");
    router.push("/dashboard");
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

        <p className="text-xs text-white/30 mt-8">Powered by Amazon Nova</p>
      </div>
    </main>
  );
}
