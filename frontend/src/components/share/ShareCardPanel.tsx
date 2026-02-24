"use client";

import { useMemo, useState } from "react";
import { Copy, Check, XCircle, Share2 } from "lucide-react";
import type { Locale } from "@/lib/i18n";
import type { WalletShareCard } from "@/lib/types";
import { TierBadge } from "@/components/ui/Badge";

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
    const url = new URL(`/s/${card.wallet.id}`, window.location.origin);
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
      setTimeout(() => setStatus("idle"), 2000);
    } catch {
      setStatus("error");
    }
  }

  return (
    <article className="rounded-lg border border-separator bg-surface-tertiary p-5">
      <div className="flex items-center gap-2">
        <Share2 className="h-4 w-4 text-label-tertiary" />
        <h4 className="text-headline text-label-primary">{labels.title}</h4>
      </div>

      {/* Preview Card */}
      <div className="mt-3 rounded-lg border border-separator bg-surface-secondary p-4 shadow-elevation-1">
        <p className="text-headline text-label-primary">
          {card.wallet.pseudonym || card.wallet.address}
        </p>
        <div className="mt-1">
          <TierBadge tier={card.pool_tier} locale={locale} />
        </div>
        <div className="mt-2 flex flex-wrap gap-x-4 gap-y-0.5">
          <span className="text-footnote text-label-tertiary">
            {labels.trades}{" "}
            <span className="font-medium text-label-secondary">{card.total_trades}</span>
          </span>
          <span className="text-footnote text-label-tertiary">
            {labels.realizedPnl}{" "}
            <span className={`font-medium ${card.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
              {card.realized_pnl >= 0 ? "+" : ""}{card.realized_pnl.toFixed(2)}
            </span>
          </span>
          <span className="text-footnote text-label-tertiary">
            {labels.score}{" "}
            <span className="font-medium text-label-secondary">{card.smart_score}</span>
          </span>
        </div>
        <p className="mt-1 text-footnote text-label-tertiary">
          {card.strategy_type} / {card.info_edge_level}
        </p>
        <p className="mt-1 text-footnote text-label-tertiary">
          Followers {card.follower_count} · 7D +{card.new_followers_7d}
        </p>
        {card.nl_summary && card.nl_summary !== "-" && (
          <p className="mt-2 text-callout text-label-secondary">{card.nl_summary}</p>
        )}
        <p className="mt-2 text-caption-2 text-label-quaternary">
          {labels.updatedAt}: {card.updated_at}
        </p>
      </div>

      {/* Actions */}
      <div className="mt-4 flex items-center gap-3">
        <button
          type="button"
          onClick={onCopy}
          className={[
            "inline-flex h-9 items-center gap-1.5 rounded-md px-4 text-subheadline font-semibold",
            "transition-all duration-200 ease-apple active:scale-[0.97]",
            status === "copied"
              ? "bg-tint-green/[0.12] text-tint-green"
              : "bg-tint-blue/[0.12] text-tint-blue hover:bg-tint-blue/[0.15]",
          ].join(" ")}
        >
          {status === "copied" ? (
            <Check className="h-4 w-4" />
          ) : (
            <Copy className="h-4 w-4" />
          )}
          {status === "copied" ? labels.copied : labels.copyLink}
        </button>
        {status === "error" && (
          <span className="flex items-center gap-1 text-caption-1 text-tint-red">
            <XCircle className="h-3.5 w-3.5" />
            {labels.copyFailed}
          </span>
        )}
      </div>
      <p className="mt-2 break-all font-mono text-caption-2 text-label-quaternary">
        {labels.preview}: {shareURL}
      </p>
    </article>
  );
}
