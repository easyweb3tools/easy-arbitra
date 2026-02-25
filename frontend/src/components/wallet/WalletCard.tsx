import Link from "next/link";
import type { PotentialWallet } from "@/lib/types";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { CategoryTag, TierBadge } from "@/components/ui/Badge";
import { fallbackSummary } from "@/lib/fallback-summary";
import { appendUTM } from "@/lib/utm";
import { ChevronRight } from "lucide-react";

export function WalletCard({
  item,
  locale,
  labels,
  utm,
}: {
  item: PotentialWallet;
  locale: "en" | "zh";
  labels: {
    unnamed: string;
    follow: string;
    unfollow: string;
    following: string;
    failed: string;
    strategyUnknown: string;
  };
  utm: URLSearchParams;
}) {
  const summary = item.summary || item.nl_summary || fallbackSummary({
    strategyType: item.strategy_type,
    smartScore: item.smart_score,
    tradeCount: item.total_trades,
    realizedPnl: item.realized_pnl,
    locale,
  });

  return (
    <div className="group border-b border-separator/60 px-5 py-4 transition-all duration-200 ease-apple last:border-b-0 hover:bg-surface-tertiary/60">
      <div className="flex items-start justify-between gap-3">
        <Link href={appendUTM(`/wallets/${item.wallet.id}`, utm)} className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <p className="text-headline text-label-primary">
              {item.wallet.pseudonym || labels.unnamed}
            </p>
            <TierBadge tier={item.pool_tier} locale={locale} />
            <CategoryTag color={item.has_ai_report ? "blue" : "gray"}>
              {item.has_ai_report ? "AI" : "No AI"}
            </CategoryTag>
          </div>
          <p className="mt-0.5 truncate font-mono text-caption-1 text-label-quaternary">{item.wallet.address}</p>
        </Link>
        <div className="flex items-center gap-2">
          <WatchlistToggleButton
            walletID={item.wallet.id}
            labels={{
              follow: labels.follow,
              unfollow: labels.unfollow,
              following: labels.following,
              failed: labels.failed,
            }}
          />
        </div>
      </div>

      <Link href={appendUTM(`/wallets/${item.wallet.id}`, utm)} className="mt-2.5 block">
        <div className="flex flex-wrap items-center gap-x-5 gap-y-1">
          <span className="text-footnote text-label-tertiary">
            trades <span className="font-semibold tabular-nums text-label-secondary">{item.total_trades}</span>
          </span>
          <span className="text-footnote text-label-tertiary">
            PnL{" "}
            <span className={`font-semibold tabular-nums ${item.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
              {item.realized_pnl >= 0 ? "+" : ""}{item.realized_pnl.toFixed(2)}
            </span>
          </span>
          <span className="text-footnote text-label-tertiary">
            score <span className="font-semibold tabular-nums text-label-secondary">{item.smart_score}</span>
          </span>
          <span className="text-footnote text-label-tertiary">{item.strategy_type || labels.strategyUnknown}</span>
        </div>
        <p className="mt-1.5 line-clamp-2 text-footnote text-label-tertiary leading-relaxed">{summary}</p>
      </Link>
    </div>
  );
}
