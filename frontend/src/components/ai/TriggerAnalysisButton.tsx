"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Brain, Loader2, CheckCircle, XCircle } from "lucide-react";

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
  requestFailed: "request failed",
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
    <div className="flex flex-wrap items-center gap-3">
      <button
        type="button"
        onClick={onTrigger}
        disabled={state === "loading"}
        className={[
          "inline-flex h-9 items-center gap-1.5 rounded-md px-4 text-subheadline font-semibold",
          "transition-all duration-200 ease-apple",
          "disabled:opacity-35 disabled:pointer-events-none",
          "active:scale-[0.97]",
          "bg-tint-blue/[0.12] text-tint-blue hover:bg-tint-blue/[0.15]",
        ].join(" ")}
      >
        {state === "loading" ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Brain className="h-4 w-4" />
        )}
        {state === "loading" ? text.loading : text.trigger}
      </button>
      {state === "done" && (
        <span className="flex items-center gap-1 text-caption-1 text-tint-green">
          <CheckCircle className="h-3.5 w-3.5" />
          {message}
        </span>
      )}
      {state === "error" && (
        <span className="flex items-center gap-1 text-caption-1 text-tint-red">
          <XCircle className="h-3.5 w-3.5" />
          {text.failedPrefix}: {message}
        </span>
      )}
    </div>
  );
}
