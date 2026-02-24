import { Filter } from "lucide-react";
import { getPotentialWallets } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { appendUTM, pickUTM } from "@/lib/utm";
import { StrategyTabs } from "@/components/wallet/StrategyTabs";
import { WalletCard } from "@/components/wallet/WalletCard";

export const dynamic = "force-dynamic";

export default async function WalletsPage({
  searchParams,
}: {
  searchParams: {
    page?: string;
    page_size?: string;
    min_trades?: string;
    min_realized_pnl?: string;
    strategy_type?: string;
    pool_tier?: string;
    has_ai_report?: string;
    sort_by?: string;
    order?: string;
    utm_source?: string;
    utm_medium?: string;
    utm_campaign?: string;
    utm_term?: string;
    utm_content?: string;
  };
}) {
  const locale = await getLocaleFromCookies();
  const utm = pickUTM(searchParams);

  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
    min_trades: searchParams.min_trades || "100",
    min_realized_pnl: searchParams.min_realized_pnl || "0",
    sort_by: searchParams.sort_by || "smart_score",
    order: searchParams.order || "desc",
  });
  if (searchParams.strategy_type) params.set("strategy_type", searchParams.strategy_type);
  if (searchParams.pool_tier) params.set("pool_tier", searchParams.pool_tier);
  if (searchParams.has_ai_report) params.set("has_ai_report", searchParams.has_ai_report);

  const wallets = await getPotentialWallets(params);
  const analyzedCount = wallets.items.filter((w) => w.has_ai_report).length;

  const makeStrategyHref = (strategyType: string) => {
    const q = new URLSearchParams(params);
    if (strategyType) q.set("strategy_type", strategyType);
    else q.delete("strategy_type");
    return appendUTM(`/wallets?${q.toString()}`, utm);
  };

  return (
    <section className="space-y-6">
      <Card variant="flat">
        <form className="grid gap-3 p-4 md:grid-cols-6" method="get">
          <div className="md:col-span-6">
            <StrategyTabs current={searchParams.strategy_type} makeHref={makeStrategyHref} />
          </div>

          <div className="md:col-span-2">
            <label className="mb-1 block text-caption-1 text-label-tertiary">{t(locale, "wallets.minTrades")}</label>
            <input
              name="min_trades"
              type="number"
              min="1"
              defaultValue={searchParams.min_trades || "100"}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            />
          </div>
          <div className="md:col-span-2">
            <label className="mb-1 block text-caption-1 text-label-tertiary">{t(locale, "wallets.minPnl")}</label>
            <input
              name="min_realized_pnl"
              type="number"
              step="0.01"
              defaultValue={searchParams.min_realized_pnl || "0"}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            />
          </div>
          <div className="md:col-span-2">
            <label className="mb-1 block text-caption-1 text-label-tertiary">Pool tier</label>
            <select
              name="pool_tier"
              defaultValue={searchParams.pool_tier || ""}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            >
              <option value="">All</option>
              <option value="star">Star</option>
              <option value="strategy">Strategy</option>
              <option value="observation">Observation</option>
            </select>
          </div>
          <div className="md:col-span-2">
            <label className="mb-1 block text-caption-1 text-label-tertiary">Sort by</label>
            <select
              name="sort_by"
              defaultValue={searchParams.sort_by || "smart_score"}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            >
              <option value="smart_score">Smart score</option>
              <option value="realized_pnl">Realized PnL</option>
              <option value="trade_count">Trade count</option>
              <option value="last_analyzed_at">AI freshness</option>
            </select>
          </div>
          <div className="md:col-span-2">
            <label className="mb-1 block text-caption-1 text-label-tertiary">AI report</label>
            <select
              name="has_ai_report"
              defaultValue={searchParams.has_ai_report || ""}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            >
              <option value="">All</option>
              <option value="true">With AI</option>
              <option value="false">No AI</option>
            </select>
          </div>
          <div className="md:col-span-1">
            <label className="mb-1 block text-caption-1 text-label-tertiary">Order</label>
            <select
              name="order"
              defaultValue={searchParams.order || "desc"}
              className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1"
            >
              <option value="desc">Desc</option>
              <option value="asc">Asc</option>
            </select>
          </div>
          <div className="flex items-end md:col-span-1">
            <button className="inline-flex h-9 w-full items-center justify-center rounded-md bg-tint-blue px-5 text-subheadline font-semibold text-white" type="submit">
              <Filter className="mr-1.5 h-4 w-4" /> {t(locale, "wallets.apply")}
            </button>
          </div>

          {utm.get("utm_source") && <input type="hidden" name="utm_source" value={utm.get("utm_source") || ""} />}
          {utm.get("utm_medium") && <input type="hidden" name="utm_medium" value={utm.get("utm_medium") || ""} />}
          {utm.get("utm_campaign") && <input type="hidden" name="utm_campaign" value={utm.get("utm_campaign") || ""} />}
          {utm.get("utm_term") && <input type="hidden" name="utm_term" value={utm.get("utm_term") || ""} />}
          {utm.get("utm_content") && <input type="hidden" name="utm_content" value={utm.get("utm_content") || ""} />}
        </form>
      </Card>

      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <StatCell label={t(locale, "wallets.potentialCount")} value={String(wallets.pagination.total)} size="large" />
        </Card>
        <Card>
          <StatCell label={t(locale, "wallets.analyzedPage")} value={String(analyzedCount)} numericValue={analyzedCount} color="positive" />
        </Card>
        <Card>
          <div className="flex flex-col gap-0.5">
            <span className="text-caption-1 text-label-tertiary">Pool distribution</span>
            <span className="text-subheadline font-medium text-label-secondary">
              Star {wallets.items.filter((i) => i.pool_tier === "star").length} · Strategy {wallets.items.filter((i) => i.pool_tier === "strategy").length}
            </span>
          </div>
        </Card>
      </div>

      <div>
        <SectionHeader title={t(locale, "wallets.title")} />
        <p className="mb-3 text-footnote text-label-tertiary">Now merged with ranking filters to support direct follow decisions.</p>
        <Card padding={false}>
          {wallets.items.map((item) => (
            <WalletCard
              key={item.wallet.id}
              item={item}
              locale={locale}
              utm={utm}
              labels={{
                unnamed: t(locale, "wallets.unnamed"),
                follow: t(locale, "watchlist.follow"),
                unfollow: t(locale, "watchlist.unfollow"),
                following: t(locale, "watchlist.following"),
                failed: t(locale, "watchlist.failed"),
                strategyUnknown: t(locale, "wallets.strategyUnknown"),
              }}
            />
          ))}
        </Card>
      </div>
    </section>
  );
}
