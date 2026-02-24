import Link from "next/link";
import { TrendingUp, Brain, Activity } from "lucide-react";
import { getAnomalies, getOpsHighlights, getOverviewStats, getPortfolios } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import { CategoryTag, TierBadge } from "@/components/ui/Badge";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { appendUTM, pickUTM } from "@/lib/utm";
import { PortfolioCard } from "@/components/portfolio/PortfolioCard";
import { fallbackSummary } from "@/lib/fallback-summary";

export const dynamic = "force-dynamic";

export default async function HomePage({
  searchParams,
}: {
  searchParams: { utm_source?: string; utm_medium?: string; utm_campaign?: string; utm_term?: string; utm_content?: string };
}) {
  const locale = await getLocaleFromCookies();
  const utm = pickUTM(searchParams);
  const [stats, highlights, portfolios, anomalies] = await Promise.all([
    getOverviewStats(),
    getOpsHighlights(new URLSearchParams({ limit: "5" })),
    getPortfolios(),
    getAnomalies(new URLSearchParams({ page: "1", page_size: "5" })),
  ]);

  return (
    <section className="space-y-8">
      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <StatCell label={t(locale, "home.marketsIndexed")} value={String(stats.indexed_markets)} size="large" />
        </Card>
        <Card>
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-tint-green/10">
              <TrendingUp className="h-5 w-5 text-tint-green" />
            </div>
            <StatCell
              label={t(locale, "home.newPotential24h")}
              value={String(highlights.new_potential_wallets_24h)}
              numericValue={highlights.new_potential_wallets_24h}
              color="positive"
            />
          </div>
        </Card>
        <Card>
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-tint-blue/10">
              <Brain className="h-5 w-5 text-tint-blue" />
            </div>
            <StatCell
              label={t(locale, "home.potentialWallets")}
              value={String(highlights.top_ai_confidence.length)}
            />
          </div>
        </Card>
      </div>

      <div>
        <SectionHeader
          title={locale === "zh" ? "今日精选" : "Featured Star Wallets"}
          action={
            <Link href={appendUTM("/wallets?pool_tier=star&sort_by=smart_score", utm)} className="text-subheadline font-medium text-tint-blue">
              {t(locale, "home.viewAll")} →
            </Link>
          }
        />
        <Card padding={false}>
          {highlights.top_ai_confidence.map((item, idx) => (
            <Link
              key={item.wallet.id}
              href={appendUTM(`/wallets/${item.wallet.id}`, utm)}
              className="flex items-center gap-3 border-b border-separator px-4 py-3 transition-colors duration-150 last:border-b-0 hover:bg-surface-tertiary"
            >
              <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-tint-gold/15 text-caption-1 font-bold text-tint-gold">
                {idx + 1}
              </span>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <p className="text-headline text-label-primary">{item.wallet.pseudonym || item.wallet.address}</p>
                  <TierBadge tier="star" locale={locale} />
                  <CategoryTag>{item.strategy_type || "—"}</CategoryTag>
                </div>
                <p className="mt-1 line-clamp-1 text-footnote text-label-tertiary">
                  {item.nl_summary || fallbackSummary({ strategyType: item.strategy_type, smartScore: item.smart_score, locale })}
                </p>
                <p className="text-caption-1 text-label-tertiary">{item.recommend_reason}</p>
              </div>
              <span className="shrink-0 text-label-quaternary">›</span>
            </Link>
          ))}
        </Card>
      </div>

      <div>
        <SectionHeader title={locale === "zh" ? "新手组合包" : "Starter Portfolios"} />
        <div className="grid gap-4 md:grid-cols-2">
          {portfolios.map((portfolio) => (
            <PortfolioCard key={portfolio.id} item={portfolio} locale={locale} />
          ))}
        </div>
      </div>

      <div>
        <SectionHeader title={locale === "zh" ? "实时信号" : "Realtime Signals"} />
        <Card padding={false}>
          {anomalies.items.map((alert) => (
            <Link
              key={alert.id}
              href={`/anomalies/${alert.id}`}
              className="flex items-center gap-3 border-b border-separator px-4 py-3 last:border-b-0 hover:bg-surface-tertiary"
            >
              <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-tint-orange/10">
                <Activity className="h-4 w-4 text-tint-orange" />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-headline text-label-primary">{alert.alert_type}</p>
                <p className="line-clamp-1 text-footnote text-label-tertiary">{alert.description}</p>
              </div>
              <span className="text-caption-1 text-label-tertiary">S{alert.severity}</span>
            </Link>
          ))}
        </Card>
      </div>
    </section>
  );
}
