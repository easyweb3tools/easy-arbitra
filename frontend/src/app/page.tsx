import Link from "next/link";
import { BarChart3, TrendingUp, Users } from "lucide-react";
import { getPotentialWallets, getOverviewStats, getOpsHighlights } from "@/lib/api";
import { SortToggle } from "@/components/ui/SortToggle";
import { LeaderboardTable } from "@/components/wallet/LeaderboardTable";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function HomePage({
  searchParams,
}: {
  searchParams: {
    sort_by?: string;
    order?: string;
    page?: string;
    page_size?: string;
  };
}) {
  const locale = await getLocaleFromCookies();
  const sortBy = searchParams.sort_by || "smart_score";
  const order = searchParams.order || "desc";
  const page = Math.max(1, parseInt(searchParams.page || "1", 10) || 1);
  const pageSize = 20;

  const params = new URLSearchParams({
    sort_by: sortBy,
    order,
    page: String(page),
    page_size: String(pageSize),
    min_trades: "100",
    min_realized_pnl: "0",
  });

  const [wallets, stats, highlights] = await Promise.all([
    getPotentialWallets(params),
    getOverviewStats(),
    getOpsHighlights(new URLSearchParams({ limit: "1" })),
  ]);

  const totalPages = Math.ceil(wallets.pagination.total / pageSize);

  return (
    <section className="space-y-6">
      {/* ── Compact Stats Strip ── */}
      <div className="flex flex-wrap items-center gap-x-6 gap-y-2 rounded-xl bg-surface-secondary px-5 py-3 shadow-elevation-1 opacity-0 animate-slide-up stagger-1">
        <div className="flex items-center gap-2">
          <BarChart3 className="h-4 w-4 text-tint-purple" />
          <span className="text-caption-1 text-label-tertiary">{t(locale, "home.marketsIndexed")}</span>
          <span className="text-subheadline font-bold tabular-nums text-label-primary">{stats.indexed_markets}</span>
        </div>
        <span className="h-4 w-px bg-separator hidden sm:block" />
        <div className="flex items-center gap-2">
          <TrendingUp className="h-4 w-4 text-tint-green" />
          <span className="text-caption-1 text-label-tertiary">{t(locale, "home.newPotential24h")}</span>
          <span className="text-subheadline font-bold tabular-nums text-tint-green">{highlights.new_potential_wallets_24h}</span>
        </div>
        <span className="h-4 w-px bg-separator hidden sm:block" />
        <div className="flex items-center gap-2">
          <Users className="h-4 w-4 text-tint-blue" />
          <span className="text-caption-1 text-label-tertiary">{t(locale, "home.potentialWallets")}</span>
          <span className="text-subheadline font-bold tabular-nums text-label-primary">{wallets.pagination.total}</span>
        </div>
      </div>

      {/* ── Sort Toggle Bar ── */}
      <div className="flex items-center justify-between opacity-0 animate-slide-up stagger-2">
        <h1 className="text-title-2 text-label-primary">{t(locale, "leaderboard.title")}</h1>
        <SortToggle
          options={[
            { label: t(locale, "leaderboard.sortByScore"), value: "smart_score" },
            { label: t(locale, "leaderboard.sortByPnl"), value: "realized_pnl" },
            { label: t(locale, "leaderboard.sortByTrades"), value: "trade_count" },
          ]}
          defaultValue={sortBy}
        />
      </div>

      {/* ── Leaderboard Table ── */}
      <div className="opacity-0 animate-slide-up stagger-3">
        <LeaderboardTable
          items={wallets.items}
          locale={locale}
          page={page}
          pageSize={pageSize}
        />
      </div>

      {/* ── Pagination ── */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3 opacity-0 animate-slide-up stagger-4">
          {page > 1 && (
            <Link
              href={`/?sort_by=${sortBy}&order=${order}&page=${page - 1}`}
              className="rounded-lg bg-surface-tertiary/80 px-4 py-2 text-subheadline font-semibold text-label-secondary transition-colors hover:bg-surface-tertiary"
            >
              {locale === "zh" ? "上一页" : "Previous"}
            </Link>
          )}
          <span className="text-subheadline tabular-nums text-label-tertiary">
            {page} / {totalPages}
          </span>
          {page < totalPages && (
            <Link
              href={`/?sort_by=${sortBy}&order=${order}&page=${page + 1}`}
              className="rounded-lg bg-surface-tertiary/80 px-4 py-2 text-subheadline font-semibold text-label-secondary transition-colors hover:bg-surface-tertiary"
            >
              {locale === "zh" ? "下一页" : "Next"}
            </Link>
          )}
        </div>
      )}
    </section>
  );
}
