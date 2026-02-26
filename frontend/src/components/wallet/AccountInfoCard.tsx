import type { WalletProfile, WalletDecisionCard } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";
import { TierBadge, CategoryTag } from "@/components/ui/Badge";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";

export function AccountInfoCard({
  profile,
  decisionCard,
  locale,
}: {
  profile: WalletProfile;
  decisionCard?: WalletDecisionCard | null;
  locale: Locale;
}) {
  const pnlColor =
    profile.layer1_facts.realized_pnl > 0
      ? "text-tint-green"
      : profile.layer1_facts.realized_pnl < 0
      ? "text-tint-red"
      : "text-label-primary";

  return (
    <article className="rounded-2xl border border-separator/50 bg-surface-secondary p-5 shadow-elevation-2 transition-all duration-300 ease-apple sm:p-6">
      {/* Top row: identity + follow */}
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <h1 className="text-title-1 tracking-tight text-label-primary">
            {profile.wallet.pseudonym || profile.wallet.address}
          </h1>
          <p className="mt-1 truncate font-mono text-footnote text-label-quaternary">
            {profile.wallet.address}
          </p>
          {profile.strategy && (
            <div className="mt-3 flex flex-wrap items-center gap-2.5">
              <CategoryTag>{profile.strategy.strategy_type}</CategoryTag>
              <TierBadge tier={profile.strategy.pool_tier} locale={locale} />
              <span className="h-3.5 w-px bg-separator" />
              <span className="text-footnote text-label-tertiary">
                {profile.strategy.info_edge_level}
              </span>
            </div>
          )}
        </div>
        <WatchlistToggleButton
          walletID={profile.wallet.id}
          labels={{
            follow: t(locale, "watchlist.follow"),
            unfollow: t(locale, "watchlist.unfollow"),
            following: t(locale, "watchlist.following"),
            failed: t(locale, "watchlist.failed"),
          }}
        />
      </div>

      {/* Stats grid */}
      <div className="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.realizedPnl")}
          </span>
          <p className={`mt-0.5 text-subheadline font-bold tabular-nums ${pnlColor}`}>
            {profile.layer1_facts.realized_pnl >= 0 ? "+" : ""}
            {profile.layer1_facts.realized_pnl.toFixed(2)}
          </p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.totalTrades")}
          </span>
          <p className="mt-0.5 text-subheadline font-bold tabular-nums text-label-primary">
            {profile.layer1_facts.total_trades.toLocaleString()}
          </p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.volume30d")}
          </span>
          <p className="mt-0.5 text-subheadline font-bold tabular-nums text-label-primary">
            {profile.layer1_facts.volume_30d.toFixed(2)}
          </p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.smartScore")}
          </span>
          <p className="mt-0.5 text-subheadline font-bold tabular-nums text-label-primary">
            {profile.strategy?.smart_score ?? "—"}
          </p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.risk")}
          </span>
          <p className={[
            "mt-0.5 text-subheadline font-bold",
            decisionCard?.risk_level === "low" ? "text-tint-green" :
            decisionCard?.risk_level === "high" ? "text-tint-red" : "text-tint-orange"
          ].join(" ")}>
            {decisionCard?.risk_level || "—"}
          </p>
        </div>
        <div className="rounded-xl bg-surface-tertiary/60 px-4 py-3">
          <span className="text-caption-1 font-medium uppercase tracking-wide text-label-tertiary">
            {t(locale, "accountInfo.momentum")}
          </span>
          <p className={[
            "mt-0.5 text-subheadline font-bold",
            decisionCard?.momentum === "heating" ? "text-tint-green" :
            decisionCard?.momentum === "cooling" ? "text-tint-red" : "text-label-secondary"
          ].join(" ")}>
            {decisionCard?.momentum || "—"}
          </p>
        </div>
      </div>
    </article>
  );
}
