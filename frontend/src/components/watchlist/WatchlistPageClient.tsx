"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { getWatchlist, getWatchlistFeed } from "@/lib/api";
import { ensureFingerprint } from "@/lib/fingerprint";
import type { WatchlistFeedItem, WatchlistItem } from "@/lib/types";

type Labels = {
  title: string;
  feedTitle: string;
  empty: string;
  trades: string;
  realizedPnl: string;
  score: string;
  loading: string;
  eventType: string;
  eventTime: string;
};

export function WatchlistPageClient({ labels }: { labels: Labels }) {
  const fingerprint = useMemo(() => ensureFingerprint(), []);
  const [loading, setLoading] = useState(true);
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [feed, setFeed] = useState<WatchlistFeedItem[]>([]);

  useEffect(() => {
    let cancelled = false;
    async function run() {
      setLoading(true);
      try {
        const [wl, fd] = await Promise.all([
          getWatchlist(new URLSearchParams({ page: "1", page_size: "50" }), fingerprint),
          getWatchlistFeed(new URLSearchParams({ page: "1", page_size: "50" }), fingerprint)
        ]);
        if (!cancelled) {
          setItems(wl.items);
          setFeed(fd.items);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    run();
    return () => {
      cancelled = true;
    };
  }, [fingerprint]);

  if (loading) {
    return <p className="text-sm text-muted">{labels.loading}</p>;
  }

  return (
    <section className="space-y-4">
      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-3 text-lg font-semibold">{labels.title}</h2>
        {items.length === 0 ? (
          <p className="text-sm text-muted">{labels.empty}</p>
        ) : (
          <div className="space-y-2">
            {items.map((item) => (
              <Link key={item.watchlist_id} href={`/wallets/${item.wallet.id}`} className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50">
                <p className="font-medium">{item.wallet.pseudonym || item.wallet.address}</p>
                <p className="text-xs text-muted">
                  {labels.trades} {item.total_trades} · {labels.realizedPnl} {item.realized_pnl.toFixed(2)} · {labels.score} {item.smart_score}
                </p>
              </Link>
            ))}
          </div>
        )}
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-3 text-lg font-semibold">{labels.feedTitle}</h2>
        {feed.length === 0 ? (
          <p className="text-sm text-muted">{labels.empty}</p>
        ) : (
          <div className="space-y-2">
            {feed.map((item) => (
              <Link key={`${item.event_type}-${item.event_id}`} href={`/wallets/${item.wallet.id}`} className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50">
                <p className="font-medium">{item.wallet.pseudonym || item.wallet.address}</p>
                <p className="text-xs text-muted">
                  {labels.eventType} {item.event_type} · {labels.eventTime} {item.event_time}
                </p>
              </Link>
            ))}
          </div>
        )}
      </article>
    </section>
  );
}
