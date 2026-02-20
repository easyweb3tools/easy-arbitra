"use client";

import { useState } from "react";

export function AcknowledgeButton({ alertID }: { alertID: number }) {
  const [loading, setLoading] = useState(false);

  async function onAck() {
    if (loading) return;
    setLoading(true);
    try {
      const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";
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
      className="rounded-md border border-slate-300 px-2 py-1 text-xs text-slate-700 disabled:opacity-60"
    >
      {loading ? "..." : "Acknowledge"}
    </button>
  );
}
