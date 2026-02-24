import Link from "next/link";
import type { PotentialWallet } from "@/lib/types";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { CategoryTag, TierBadge } from "@/components/ui/Badge";
import { fallbackSummary } from "@/lib/fallback-summary";
import { appendUTM } from "@/lib/utm";

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
    <div className="border-b border-separator px-4 py-3 transition-colors duration-150 ease-apple last:border-b-0 hover:bg-surface-tertiary">
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
          <p className="mt-0.5 truncate font-mono text-caption-1 text-label-tertiary">{item.wallet.address}</p>
        </Link>
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

      <Link href={appendUTM(`/wallets/${item.wallet.id}`, utm)} className="mt-2 block">
        <div className="flex flex-wrap gap-x-4 gap-y-0.5">
          <span className="text-footnote text-label-tertiary">
            trades <span className="font-medium text-label-secondary">{item.total_trades}</span>
          </span>
          <span className="text-footnote text-label-tertiary">
            PnL{" "}
            <span className={`font-medium ${item.realized_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
              {item.realized_pnl >= 0 ? "+" : ""}{item.realized_pnl.toFixed(2)}
            </span>
          </span>
          <span className="text-footnote text-label-tertiary">
            score <span className="font-medium text-label-secondary">{item.smart_score}</span>
          </span>
          <span className="text-footnote text-label-tertiary">{item.strategy_type || labels.strategyUnknown}</span>
        </div>
        <p className="mt-1 line-clamp-2 text-footnote text-label-tertiary">{summary}</p>
      </Link>
    </div>
  );
}
