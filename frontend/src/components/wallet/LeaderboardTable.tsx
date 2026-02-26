import Link from "next/link";
import type { PotentialWallet } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";
import { Card } from "@/components/ui/Card";
import { TierBadge, CategoryTag } from "@/components/ui/Badge";

export function LeaderboardTable({
  items,
  locale,
  page = 1,
  pageSize = 20,
}: {
  items: PotentialWallet[];
  locale: Locale;
  page?: number;
  pageSize?: number;
}) {
  return (
    <Card padding={false}>
      {/* Header row */}
      <div className="hidden sm:grid sm:grid-cols-[3rem_1fr_6rem_5.5rem_6rem_5rem_5rem] items-center gap-2 border-b border-separator/60 px-5 py-2.5 text-caption-1 font-semibold uppercase tracking-wider text-label-tertiary">
        <span>{t(locale, "leaderboard.rank")}</span>
        <span>{t(locale, "leaderboard.name")}</span>
        <span>{t(locale, "leaderboard.strategy")}</span>
        <span className="text-right">{t(locale, "leaderboard.smartScore")}</span>
        <span className="text-right">{t(locale, "leaderboard.pnl")}</span>
        <span className="text-right">{t(locale, "leaderboard.trades")}</span>
        <span className="text-center">{t(locale, "leaderboard.tier")}</span>
      </div>

      {items.map((item, idx) => {
        const rank = (page - 1) * pageSize + idx + 1;
        const pnlColor = item.realized_pnl > 0 ? "text-tint-green" : item.realized_pnl < 0 ? "text-tint-red" : "text-label-secondary";

        return (
          <Link
            key={item.wallet.id}
            href={`/wallets/${item.wallet.id}`}
            className="group flex flex-col gap-2 border-b border-separator/60 px-5 py-4 transition-all duration-200 ease-apple last:border-b-0 hover:bg-surface-tertiary/70 sm:grid sm:grid-cols-[3rem_1fr_6rem_5.5rem_6rem_5rem_5rem] sm:items-center sm:gap-2"
          >
            {/* Rank */}
            <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-surface-tertiary/80 text-caption-1 font-bold tabular-nums text-label-secondary sm:h-7 sm:w-7">
              {rank}
            </span>

            {/* Name / Address */}
            <div className="min-w-0 flex-1">
              <p className="truncate text-subheadline font-medium text-label-primary">
                {item.wallet.pseudonym || `${item.wallet.address.slice(0, 6)}...${item.wallet.address.slice(-4)}`}
              </p>
              <p className="mt-0.5 line-clamp-1 text-caption-1 text-label-quaternary sm:hidden">
                {item.summary}
              </p>
            </div>

            {/* Strategy */}
            <div>
              <CategoryTag>{item.strategy_type || "—"}</CategoryTag>
            </div>

            {/* Smart Score */}
            <span className="text-right text-subheadline font-semibold tabular-nums text-label-primary">
              {item.smart_score}
            </span>

            {/* PnL */}
            <span className={`text-right text-subheadline font-semibold tabular-nums ${pnlColor}`}>
              {item.realized_pnl >= 0 ? "+" : ""}{item.realized_pnl.toFixed(2)}
            </span>

            {/* Trades */}
            <span className="text-right text-subheadline tabular-nums text-label-secondary">
              {item.total_trades.toLocaleString()}
            </span>

            {/* Tier */}
            <div className="flex justify-center">
              <TierBadge tier={item.pool_tier} locale={locale} />
            </div>
          </Link>
        );
      })}
    </Card>
  );
}
