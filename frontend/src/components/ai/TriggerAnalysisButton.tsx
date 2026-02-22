"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

type Labels = {
  trigger: string;
  loading: string;
  updatedAt: string;
  failedPrefix: string;
  requestFailed: string;
};

const defaultLabels: Labels = {
  trigger: "Trigger AI Analysis",
  loading: "Analyzing...",
  updatedAt: "Updated at",
  failedPrefix: "Analyze failed",
  requestFailed: "request failed"
};

export function TriggerAnalysisButton({ walletID, labels }: { walletID: string; labels?: Labels }) {
  const text = labels || defaultLabels;
  const [state, setState] = useState<"idle" | "loading" | "done" | "error">("idle");
  const [message, setMessage] = useState<string>("");
  const router = useRouter();

  async function onTrigger() {
    setState("loading");
    setMessage("");
    const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1";
    try {
      const res = await fetch(`${apiBase}/ai/analyze/${walletID}?force=true`, { method: "POST" });
      if (!res.ok) {
        const errBody = await res.json().catch(() => null);
        const detail = errBody?.error || `status ${res.status}`;
        throw new Error(detail);
      }
      const body = await res.json();
      setState("done");
      setMessage(`${text.updatedAt} ${body?.data?.created_at || "now"}`);
      router.refresh();
    } catch (err) {
      setState("error");
      setMessage(err instanceof Error ? err.message : text.requestFailed);
    }
  }

  return (
    <div className="space-y-2">
      <button
        type="button"
        onClick={onTrigger}
        disabled={state === "loading"}
        className="rounded-md bg-accent px-3 py-2 text-xs font-medium text-white disabled:opacity-60"
      >
        {state === "loading" ? text.loading : text.trigger}
      </button>
      {state === "done" && <p className="text-xs text-emerald-700">{message}</p>}
      {state === "error" && <p className="text-xs text-rose-700">{text.failedPrefix}: {message}</p>}
    </div>
  );
}
