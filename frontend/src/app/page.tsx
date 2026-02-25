import Link from "next/link";
import { TrendingUp, Brain, Activity, BarChart3, ChevronRight } from "lucide-react";
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
    <section className="space-y-10">
      {/* ── Stats Overview ── */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Card className="opacity-0 animate-slide-up stagger-1">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-tint-purple/10">
              <BarChart3 className="h-5 w-5 text-tint-purple" />
            </div>
            <StatCell label={t(locale, "home.marketsIndexed")} value={String(stats.indexed_markets)} size="large" />
          </div>
        </Card>
        <Card className="opacity-0 animate-slide-up stagger-2">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-tint-green/10">
              <TrendingUp className="h-5 w-5 text-tint-green" />
            </div>
            <StatCell
              label={t(locale, "home.newPotential24h")}
              value={String(highlights.new_potential_wallets_24h)}
              numericValue={highlights.new_potential_wallets_24h}
              color="positive"
              size="large"
            />
          </div>
        </Card>
        <Card className="opacity-0 animate-slide-up stagger-3">
          <div className="flex items-start gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-tint-blue/10">
              <Brain className="h-5 w-5 text-tint-blue" />
            </div>
            <StatCell
              label={t(locale, "home.potentialWallets")}
              value={String(highlights.top_ai_confidence.length)}
              size="large"
            />
          </div>
        </Card>
      </div>

      {/* ── Featured Star Wallets ── */}
      <div className="opacity-0 animate-slide-up stagger-4">
        <SectionHeader
          title={locale === "zh" ? "今日精选" : "Featured Star Wallets"}
          action={
            <Link href={appendUTM("/wallets?pool_tier=star&sort_by=smart_score", utm)} className="inline-flex items-center gap-1 text-subheadline font-medium text-tint-blue transition-opacity hover:opacity-70">
              {t(locale, "home.viewAll")}
              <ChevronRight className="h-4 w-4" />
            </Link>
          }
        />
        <Card padding={false}>
          {highlights.top_ai_confidence.map((item, idx) => (
            <Link
              key={item.wallet.id}
              href={appendUTM(`/wallets/${item.wallet.id}`, utm)}
              className="group flex items-center gap-4 border-b border-separator/60 px-5 py-4 transition-all duration-200 ease-apple last:border-b-0 hover:bg-surface-tertiary/70"
            >
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-amber-400/20 to-orange-400/20 text-caption-1 font-bold tabular-nums text-tint-gold">
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
                {item.recommend_reason && (
                  <p className="mt-0.5 text-caption-1 text-label-quaternary">{item.recommend_reason}</p>
                )}
              </div>
              <ChevronRight className="h-4 w-4 shrink-0 text-label-quaternary transition-transform duration-200 group-hover:translate-x-0.5" />
            </Link>
          ))}
        </Card>
      </div>

      {/* ── Starter Portfolios ── */}
      <div className="opacity-0 animate-slide-up stagger-5">
        <SectionHeader title={locale === "zh" ? "新手组合包" : "Starter Portfolios"} />
        <div className="grid gap-4 md:grid-cols-2">
          {portfolios.map((portfolio) => (
            <PortfolioCard key={portfolio.id} item={portfolio} locale={locale} />
          ))}
        </div>
      </div>

      {/* ── Realtime Signals ── */}
      <div className="opacity-0 animate-slide-up stagger-6">
        <SectionHeader title={locale === "zh" ? "实时信号" : "Realtime Signals"} />
        <Card padding={false}>
          {anomalies.items.map((alert) => (
            <Link
              key={alert.id}
              href={`/anomalies/${alert.id}`}
              className="group flex items-center gap-4 border-b border-separator/60 px-5 py-4 last:border-b-0 transition-all duration-200 ease-apple hover:bg-surface-tertiary/70"
            >
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-tint-orange/10">
                <Activity className="h-4.5 w-4.5 text-tint-orange" />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-headline text-label-primary">{alert.alert_type}</p>
                <p className="line-clamp-1 text-footnote text-label-tertiary">{alert.description}</p>
              </div>
              <div className="flex shrink-0 items-center gap-2">
                <span className={[
                  "rounded-full px-2 py-0.5 text-caption-2 font-semibold",
                  alert.severity >= 4 ? "bg-tint-red/10 text-tint-red" :
                  alert.severity >= 2 ? "bg-tint-orange/10 text-tint-orange" :
                  "bg-surface-tertiary text-label-tertiary"
                ].join(" ")}>
                  S{alert.severity}
                </span>
                <ChevronRight className="h-4 w-4 text-label-quaternary transition-transform duration-200 group-hover:translate-x-0.5" />
              </div>
            </Link>
          ))}
        </Card>
      </div>
    </section>
  );
}
