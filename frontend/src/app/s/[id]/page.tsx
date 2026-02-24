import Link from "next/link";
import { getWalletShareLanding } from "@/lib/api";
import { TierBadge, CategoryTag } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { t } from "@/lib/i18n";

export const dynamic = "force-dynamic";

export default async function ShareLandingPage({ params }: { params: { id: string } }) {
  const locale = await getLocaleFromCookies();
  const data = await getWalletShareLanding(params.id);

  return (
    <section className="mx-auto max-w-2xl space-y-6 py-4">
      <Card variant="prominent">
        <div className="flex items-center gap-2">
          <h1 className="text-title-1 text-label-primary">{data.wallet.pseudonym || data.wallet.address}</h1>
          <TierBadge tier={data.pool_tier} locale={locale} />
        </div>
        <p className="mt-1 font-mono text-footnote text-label-tertiary">{data.wallet.address}</p>

        <div className="mt-4 grid gap-3 sm:grid-cols-3">
          <Metric label="7D" value={`${data.pnl_7d >= 0 ? "+" : ""}${data.pnl_7d.toFixed(2)}`} positive={data.pnl_7d >= 0} />
          <Metric label="Drawdown" value={`-${data.max_drawdown_7d.toFixed(2)}`} positive={false} />
          <Metric label="Stability" value={String(data.stability_score)} positive={true} />
        </div>

        <div className="mt-3 flex flex-wrap items-center gap-2">
          <CategoryTag>{data.strategy_type || "—"}</CategoryTag>
          <span className="text-footnote text-label-tertiary">Followers <span className="font-medium text-label-secondary">{data.follower_count}</span></span>
          <span className="text-footnote text-label-tertiary">7D +<span className="font-medium text-label-secondary">{data.new_followers_7d}</span></span>
        </div>

        <p className="mt-4 rounded-lg bg-surface-tertiary p-4 text-callout text-label-secondary">{data.nl_summary}</p>

        <div className="mt-4 flex flex-wrap items-center gap-3">
          <WatchlistToggleButton
            walletID={data.wallet.id}
            labels={{
              follow: t(locale, "watchlist.follow"),
              unfollow: t(locale, "watchlist.unfollow"),
              following: t(locale, "watchlist.following"),
              failed: t(locale, "watchlist.failed"),
            }}
          />
          <Link href={`/wallets/${data.wallet.id}`} className="inline-flex h-9 items-center rounded-md bg-surface-tertiary px-4 text-subheadline font-semibold text-label-primary">
            {locale === "zh" ? "查看完整分析" : "View full analysis"}
          </Link>
        </div>
      </Card>

      <Card>
        <h2 className="text-title-3 text-label-primary">{locale === "zh" ? "跟单建议" : "Decision"}</h2>
        <div className="mt-2 grid gap-2 sm:grid-cols-2">
          <p className="text-footnote text-label-tertiary">{locale === "zh" ? "适合人群" : "Suitable for"}: <span className="font-medium text-label-secondary">{data.decision_card.suitable_for}</span></p>
          <p className="text-footnote text-label-tertiary">{locale === "zh" ? "风险等级" : "Risk"}: <span className="font-medium text-label-secondary">{data.decision_card.risk_level}</span></p>
          <p className="text-footnote text-label-tertiary">{locale === "zh" ? "建议仓位" : "Suggested position"}: <span className="font-medium text-label-secondary">{data.decision_card.suggested_position}</span></p>
          <p className="text-footnote text-label-tertiary">{locale === "zh" ? "近7天状态" : "7D momentum"}: <span className="font-medium text-label-secondary">{data.decision_card.momentum}</span></p>
        </div>
        <p className="mt-3 text-callout text-label-secondary">{locale === "zh" ? data.decision_card.recommendation_zh : data.decision_card.recommendation}</p>
        <p className="mt-2 text-caption-1 text-tint-orange">{locale === "zh" ? data.decision_card.disclaimer_zh : data.decision_card.disclaimer}</p>
      </Card>
    </section>
  );
}

function Metric({ label, value, positive }: { label: string; value: string; positive: boolean }) {
  return (
    <div className="rounded-md bg-surface-tertiary p-3">
      <p className="text-caption-1 text-label-tertiary">{label}</p>
      <p className={`text-headline ${positive ? "text-tint-green" : "text-label-secondary"}`}>{value}</p>
    </div>
  );
}
