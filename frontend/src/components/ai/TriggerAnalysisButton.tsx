"use client";

import { useState } from "react";

export function TriggerAnalysisButton({ walletID }: { walletID: string }) {
  const [state, setState] = useState<"idle" | "loading" | "done" | "error">("idle");

  async function onTrigger() {
    setState("loading");
    const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";
    try {
      const res = await fetch(`${apiBase}/ai/analyze/${walletID}`, { method: "POST" });
      if (!res.ok) {
        throw new Error(`status ${res.status}`);
      }
      setState("done");
      window.location.reload();
    } catch {
      setState("error");
    }
  }

  return (
    <button
      type="button"
      onClick={onTrigger}
      disabled={state === "loading"}
      className="rounded-md bg-accent px-3 py-2 text-xs font-medium text-white disabled:opacity-60"
    >
      {state === "loading" ? "Analyzing..." : "Trigger AI Analysis"}
    </button>
  );
}
