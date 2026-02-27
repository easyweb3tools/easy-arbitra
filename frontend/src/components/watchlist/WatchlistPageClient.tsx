"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Star, Activity, AlertTriangle } from "lucide-react";
import { getWatchlist, getWatchlistFeed, getWatchlistSummary } from "@/lib/api";
import { EmptyState } from "@/components/ui/EmptyState";
import { Card, SectionHeader } from "@/components/ui/Card";
import { SkeletonRow } from "@/components/ui/Skeleton";
import type { Locale } from "@/lib/i18n";
import type { WatchlistFeedItem, WatchlistItem, WatchlistSummary } from "@/lib/types";

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

export function WatchlistPageClient({ labels, locale }: { labels: Labels; locale?: Locale }) {
  const [loading, setLoading] = useState(true);
  const [items, setItems] = useState<WatchlistItem[]>([]);
  const [feed, setFeed] = useState<WatchlistFeedItem[]>([]);
  const [summary, setSummary] = useState<WatchlistSummary | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function run() {
      setLoading(true);
      try {
        const [wl, fd, sm] = await Promise.all([
          getWatchlist(new URLSearchParams({ page: "1", page_size: "50" })),
          getWatchlistFeed(new URLSearchParams({ page: "1", page_size: "50" })),
          getWatchlistSummary(),
        ]);
        if (!cancelled) {
          setItems(wl.items);
          setFeed(fd.items);
          setSummary(sm);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    run();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return (
      <section className="space-y-6">
        <Card padding={false}>
          <SkeletonRow />
          <SkeletonRow />
          <SkeletonRow />
        </Card>
      </section>
    );
  }

  const actionItems = feed.filter((item) => item.action_required);
  const normalItems = feed.filter((item) => !item.action_required);

  return (
    <section className="space-y-8 animate-fade-in">
      <div>
        <SectionHeader title={locale === "zh" ? "组合概览" : "Portfolio Summary"} />
        <Card>
          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <p className="text-caption-1 text-label-tertiary">{locale === "zh" ? "关注钱包" : "Followed wallets"}</p>
              <p className="text-title-3 text-label-primary">{summary?.followed_wallets ?? 0}</p>
            </div>
            <div>
              <p className="text-caption-1 text-label-tertiary">{locale === "zh" ? "需行动" : "Action required"}</p>
              <p className="text-title-3 text-tint-red">{summary?.action_required ?? 0}</p>
            </div>
            <div>
              <p className="text-caption-1 text-label-tertiary">{locale === "zh" ? "健康状态" : "Healthy"}</p>
              <p className="text-title-3 text-tint-green">{summary?.healthy_wallets ?? 0}</p>
            </div>
          </div>
          {summary && Object.keys(summary.style_distribution || {}).length > 0 && (
            <div className="mt-3 flex flex-wrap gap-2">
              {Object.entries(summary.style_distribution).map(([style, count]) => (
                <span key={style} className="rounded-full bg-surface-tertiary px-2.5 py-1 text-caption-1 text-label-secondary">
                  {style}: {count}
                </span>
              ))}
            </div>
          )}
        </Card>
      </div>

      <div>
        <SectionHeader title={labels.title} />
        {items.length === 0 ? (
          <Card>
            <EmptyState
              preset="watchlist-empty"
              locale={locale}
              action={
                <Link
                  href="/wallets"
                  className="inline-flex h-9 items-center gap-1.5 rounded-md bg-tint-blue/[0.12] px-4 text-subheadline font-semibold text-tint-blue transition-all hover:bg-tint-blue/[0.15]"
                >
                  <Star className="h-4 w-4" />
                  {locale === "zh" ? "去探索钱包" : "Explore Wallets"}
                </Link>
              }
            />
          </Card>
        ) : (
          <Card padding={false}>
            {items.map((item) => (
              <Link
                key={item.watchlist_id}
                href={`/wallets/${item.wallet.id}`}
                className="flex items-center gap-3 border-b border-separator px-4 py-3 transition-colors duration-150 ease-apple last:border-b-0 hover:bg-surface-tertiary"
              >
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-tint-blue/10">
                  <Star className="h-4 w-4 fill-tint-blue text-tint-blue" />
                </span>
                <div className="min-w-0 flex-1">
                  <p className="text-headline text-label-primary">{item.wallet.pseudonym || item.wallet.address}</p>
                  <div className="mt-0.5 flex flex-wrap gap-x-3 gap-y-0.5">
                    <span className="text-footnote text-label-tertiary">{labels.trades} <span className="font-medium text-label-secondary">{item.total_trades}</span></span>
                    <span className="text-footnote text-label-tertiary">PnL <span className={`font-medium ${item.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>{item.realized_pnl >= 0 ? "+" : ""}{item.realized_pnl.toFixed(2)}</span></span>
                    <span className="text-footnote text-label-tertiary">{labels.score} <span className="font-medium text-label-secondary">{item.smart_score}</span></span>
                  </div>
                </div>
                <span className="shrink-0 text-label-quaternary">›</span>
              </Link>
            ))}
          </Card>
        )}
      </div>

      <div>
        <SectionHeader title={locale === "zh" ? "需要行动" : "Action Required"} />
        {actionItems.length === 0 ? (
          <Card>
            <EmptyState preset="all-clear" locale={locale} />
          </Card>
        ) : (
          <Card padding={false}>
            {actionItems.map((item) => (
              <Link
                key={`${item.event_type}-${item.event_id}`}
                href={`/wallets/${item.wallet.id}`}
                className="flex items-center gap-3 border-b border-separator px-4 py-3 transition-colors duration-150 ease-apple last:border-b-0 hover:bg-surface-tertiary"
              >
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-tint-red/10">
                  <AlertTriangle className="h-4 w-4 text-tint-red" />
                </span>
                <div className="min-w-0 flex-1">
                  <p className="text-headline text-label-primary">{item.wallet.pseudonym || item.wallet.address}</p>
                  <p className="mt-0.5 text-footnote text-label-tertiary">{locale === "zh" ? item.suggestion_zh : item.suggestion}</p>
                  <p className="mt-1 text-caption-2 text-label-quaternary">{item.event_time}</p>
                </div>
              </Link>
            ))}
          </Card>
        )}
      </div>

      <div>
        <SectionHeader title={labels.feedTitle} />
        {normalItems.length === 0 ? (
          <Card>
            <EmptyState preset="feed-empty" locale={locale} />
          </Card>
        ) : (
          <Card padding={false}>
            {normalItems.map((item) => (
              <Link
                key={`${item.event_type}-${item.event_id}`}
                href={`/wallets/${item.wallet.id}`}
                className="flex items-center gap-3 border-b border-separator px-4 py-3 transition-colors duration-150 ease-apple last:border-b-0 hover:bg-surface-tertiary"
              >
                <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-tint-orange/10">
                  <Activity className="h-4 w-4 text-tint-orange" />
                </span>
                <div className="min-w-0 flex-1">
                  <p className="text-headline text-label-primary">{item.wallet.pseudonym || item.wallet.address}</p>
                  <div className="mt-0.5 flex gap-3">
                    <span className="text-footnote text-label-tertiary">{labels.eventType} <span className="font-medium text-label-secondary">{item.event_type}</span></span>
                    <span className="text-footnote text-label-tertiary">{labels.eventTime} {item.event_time}</span>
                  </div>
                </div>
                <span className="shrink-0 text-label-quaternary">›</span>
              </Link>
            ))}
          </Card>
        )}
      </div>
    </section>
  );
}
