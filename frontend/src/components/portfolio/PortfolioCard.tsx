"use client";

import Link from "next/link";
import { useMemo, useState } from "react";
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
    <article className="rounded-lg border border-separator bg-surface-secondary p-4 shadow-elevation-1">
      <h3 className="text-headline text-label-primary">{locale === "zh" ? item.name_zh : item.name}</h3>
      <p className="mt-1 text-footnote text-label-tertiary">{item.description}</p>
      <div className="mt-2 flex gap-3 text-footnote text-label-tertiary">
        <span>Risk: <span className="font-medium text-label-secondary">{item.risk_level}</span></span>
        <span>Return: <span className="font-medium text-label-secondary">{item.expected_return}</span></span>
        <span>DD: <span className="font-medium text-label-secondary">{item.max_drawdown}</span></span>
      </div>
      <div className="mt-3 flex flex-wrap gap-2">
        {item.wallets.map((wallet) => (
          <Link
            key={wallet.id}
            href={`/wallets/${wallet.id}`}
            className="rounded-full bg-surface-tertiary px-2.5 py-1 text-caption-1 text-label-secondary"
          >
            {wallet.pseudonym || wallet.address.slice(0, 6)}
          </Link>
        ))}
      </div>
      <button
        type="button"
        disabled={loading}
        onClick={onBatchFollow}
        className="mt-4 inline-flex h-9 items-center rounded-md bg-tint-blue px-4 text-subheadline font-semibold text-white disabled:opacity-40"
      >
        {done ? (locale === "zh" ? "已关注组合" : "Followed") : loading ? (locale === "zh" ? "处理中..." : "Saving...") : (locale === "zh" ? "一键关注全部" : "Follow all")}
      </button>
    </article>
  );
}
