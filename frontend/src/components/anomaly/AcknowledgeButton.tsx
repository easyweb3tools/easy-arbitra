"use client";

import { useState } from "react";
import { Check, Loader2 } from "lucide-react";

type Labels = {
  ack: string;
  acking: string;
};

const defaultLabels: Labels = {
  ack: "Acknowledge",
  acking: "...",
};

export function AcknowledgeButton({ alertID, labels }: { alertID: number; labels?: Labels }) {
  const text = labels || defaultLabels;
  const [loading, setLoading] = useState(false);

  async function onAck() {
    if (loading) return;
    setLoading(true);
    try {
      const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "/api/v1";
      await fetch(`${apiBase}/anomalies/${alertID}/acknowledge`, { method: "PATCH" });
      window.location.reload();
    } finally {
      setLoading(false);
    }
  }

  return (
    <button
      type="button"
      onClick={onAck}
      disabled={loading}
      className={[
        "inline-flex h-9 items-center gap-1.5 rounded-md px-4",
        "bg-surface-tertiary text-subheadline font-semibold text-label-primary",
        "transition-all duration-200 ease-apple",
        "hover:brightness-95 active:scale-[0.97]",
        "disabled:opacity-35 disabled:pointer-events-none",
      ].join(" ")}
    >
      {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}
      {loading ? text.acking : text.ack}
    </button>
  );
}
