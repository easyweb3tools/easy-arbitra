"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
import { Check, Users } from "lucide-react";
import { addBatchToWatchlist } from "@/lib/api";
import { ensureFingerprint } from "@/lib/fingerprint";
import type { PortfolioItem } from "@/lib/types";

export function PortfolioCard({
  item,
  locale,
}: {
  item: PortfolioItem;
  locale: "en" | "zh";
}) {
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);
  const fingerprint = useMemo(() => ensureFingerprint(), []);

  async function onBatchFollow() {
    if (loading) return;
    setLoading(true);
    try {
      await addBatchToWatchlist(item.wallet_ids, fingerprint);
      setDone(true);
    } finally {
      setLoading(false);
    }
  }

  return (
    <article className="rounded-2xl border border-separator/50 bg-surface-secondary p-5 shadow-elevation-1 transition-all duration-300 ease-apple hover:shadow-[var(--card-hover-shadow)]">
      <div className="flex items-start justify-between gap-3">
        <h3 className="text-headline text-label-primary">{locale === "zh" ? item.name_zh : item.name}</h3>
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-tint-blue/10">
          <Users className="h-4 w-4 text-tint-blue" />
        </div>
      </div>
      <p className="mt-1.5 text-footnote text-label-tertiary leading-relaxed">{item.description}</p>
      <div className="mt-3 flex gap-4 text-footnote">
        <span className="text-label-tertiary">Risk <span className="font-semibold tabular-nums text-label-secondary">{item.risk_level}</span></span>
        <span className="text-label-tertiary">Return <span className="font-semibold tabular-nums text-label-secondary">{item.expected_return}</span></span>
        <span className="text-label-tertiary">DD <span className="font-semibold tabular-nums text-label-secondary">{item.max_drawdown}</span></span>
      </div>
      <div className="mt-3 flex flex-wrap gap-2">
        {item.wallets.map((wallet) => (
          <Link
            key={wallet.id}
            href={`/wallets/${wallet.id}`}
            className="rounded-full bg-surface-tertiary/70 px-2.5 py-1 text-caption-1 text-label-secondary transition-colors duration-150 hover:bg-surface-tertiary hover:text-label-primary"
          >
            {wallet.pseudonym || wallet.address.slice(0, 6)}
          </Link>
        ))}
      </div>
      <button
        type="button"
        disabled={loading || done}
        onClick={onBatchFollow}
        className={[
          "mt-4 inline-flex h-9 items-center gap-2 rounded-lg px-4 text-subheadline font-semibold transition-all duration-250 ease-apple",
          done
            ? "bg-tint-green/12 text-tint-green"
            : "bg-tint-blue text-white shadow-[0_1px_3px_rgba(0,122,255,0.25)] hover:shadow-[0_4px_12px_rgba(0,122,255,0.3)] active:scale-[0.97]",
          "disabled:opacity-40 disabled:shadow-none",
        ].join(" ")}
      >
        {done ? (
          <>
            <Check className="h-4 w-4" />
            {locale === "zh" ? "已关注组合" : "Followed"}
          </>
        ) : loading ? (
          locale === "zh" ? "处理中..." : "Saving..."
        ) : (
          locale === "zh" ? "一键关注全部" : "Follow all"
        )}
      </button>
    </article>
  );
}
