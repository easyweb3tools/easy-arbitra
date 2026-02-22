"use client";

import { useMemo, useState } from "react";
import type { Locale } from "@/lib/i18n";
import type { WalletShareCard } from "@/lib/types";

type Labels = {
  title: string;
  preview: string;
  copyLink: string;
  copied: string;
  copyFailed: string;
  trades: string;
  realizedPnl: string;
  score: string;
  updatedAt: string;
};

export function ShareCardPanel({ card, locale, labels }: { card: WalletShareCard; locale: Locale; labels: Labels }) {
  const [status, setStatus] = useState<"idle" | "copied" | "error">("idle");

  const shareURL = useMemo(() => {
    if (typeof window === "undefined") return "";
    const url = new URL(`/wallets/${card.wallet.id}`, window.location.origin);
    url.searchParams.set("utm_source", "share_card");
    url.searchParams.set("utm_medium", "copy_link");
    url.searchParams.set("utm_campaign", "wallet_insight");
    url.searchParams.set("locale", locale);
    return url.toString();
  }, [card.wallet.id, locale]);

  async function onCopy() {
    try {
      await navigator.clipboard.writeText(shareURL);
      setStatus("copied");
    } catch {
      setStatus("error");
    }
  }

  return (
    <article className="rounded-lg border border-slate-200 bg-slate-50 p-4">
      <h4 className="text-sm font-semibold">{labels.title}</h4>
      <div className="mt-2 rounded-md border border-slate-200 bg-white p-3">
        <p className="font-medium">{card.wallet.pseudonym || card.wallet.address}</p>
        <p className="mt-1 text-xs text-muted">
          {labels.trades} {card.total_trades} · {labels.realizedPnl} {card.realized_pnl.toFixed(2)} · {labels.score} {card.smart_score}
        </p>
        <p className="mt-1 text-xs text-muted">{card.strategy_type} / {card.info_edge_level}</p>
        <p className="mt-2 text-xs text-slate-700">{card.nl_summary || "-"}</p>
        <p className="mt-2 text-[11px] text-muted">{labels.updatedAt}: {card.updated_at}</p>
      </div>
      <div className="mt-3 flex items-center gap-2">
        <button type="button" onClick={onCopy} className="rounded-md bg-accent px-3 py-2 text-xs font-medium text-white">
          {labels.copyLink}
        </button>
        {status === "copied" && <span className="text-xs text-emerald-700">{labels.copied}</span>}
        {status === "error" && <span className="text-xs text-rose-700">{labels.copyFailed}</span>}
      </div>
      <p className="mt-2 break-all text-[11px] text-muted">{labels.preview}: {shareURL}</p>
    </article>
  );
}
