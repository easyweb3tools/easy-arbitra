import Link from "next/link";
import { BarChart3, TrendingUp, Sparkles } from "lucide-react";
import { getPotentialWallets, getOverviewStats, getDailyPick } from "@/lib/api";
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

  const [wallets, stats] = await Promise.all([
    getPotentialWallets(params),
    getOverviewStats(),
  ]);

  let dailyPick: Awaited<ReturnType<typeof getDailyPick>> | null = null;
  try {
    dailyPick = await getDailyPick();
  } catch {
    // no daily pick yet
  }

  const totalPages = Math.ceil(wallets.pagination.total / pageSize);

  return (
    <section className="space-y-6">
      {/* ── Daily Pick Banner ── */}
      {dailyPick && (
        <Link
          href="/daily-picks"
          className="group block rounded-2xl bg-gradient-to-r from-tint-blue/10 via-tint-purple/10 to-tint-green/10 p-5 shadow-elevation-1 transition-all hover:shadow-elevation-2 opacity-0 animate-slide-up stagger-1"
        >
          <div className="flex items-center gap-3 mb-2">
            <Sparkles className="h-5 w-5 text-tint-blue" />
            <span className="text-subheadline font-bold text-label-primary">
              {locale === "zh" ? "🏆 今日推荐交易者" : "🏆 Today's Recommended Trader"}
            </span>
          </div>
          <div className="flex items-center gap-4 flex-wrap">
            <span className="text-caption-1 font-mono text-label-secondary">
              {dailyPick.wallet?.address
                ? `${dailyPick.wallet.address.slice(0, 6)}…${dailyPick.wallet.address.slice(-4)}`
                : `Wallet #${dailyPick.pick.wallet_id}`}
            </span>
            <span className="text-caption-1 text-tint-green font-semibold">
              Score: {dailyPick.pick.smart_score}
            </span>
            <span className="text-caption-1 text-tint-purple font-semibold">
              PnL: {dailyPick.pick.realized_pnl >= 0 ? "+" : ""}{dailyPick.pick.realized_pnl.toFixed(2)} USDC
            </span>
          </div>
          {dailyPick.pick.reason_summary && (
            <p className="mt-2 text-caption-1 text-label-tertiary line-clamp-2">
              {locale === "zh" && dailyPick.pick.reason_summary_zh
                ? dailyPick.pick.reason_summary_zh
                : dailyPick.pick.reason_summary}
            </p>
          )}
        </Link>
      )}

      {/* ── Compact Stats Strip ── */}
      <div className="flex flex-wrap items-center gap-x-6 gap-y-2 rounded-xl bg-surface-secondary px-5 py-3 shadow-elevation-1 opacity-0 animate-slide-up stagger-2">
        <div className="flex items-center gap-2">
          <BarChart3 className="h-4 w-4 text-tint-purple" />
          <span className="text-caption-1 text-label-tertiary">{t(locale, "home.marketsIndexed")}</span>
          <span className="text-subheadline font-bold tabular-nums text-label-primary">{stats.indexed_markets}</span>
        </div>
        <span className="h-4 w-px bg-separator hidden sm:block" />
        <div className="flex items-center gap-2">
          <TrendingUp className="h-4 w-4 text-tint-green" />
          <span className="text-caption-1 text-label-tertiary">{t(locale, "home.potentialWallets")}</span>
          <span className="text-subheadline font-bold tabular-nums text-label-primary">{wallets.pagination.total}</span>
        </div>
      </div>

      {/* ── Sort Toggle Bar ── */}
      <div className="flex items-center justify-between opacity-0 animate-slide-up stagger-3">
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
      <div className="opacity-0 animate-slide-up stagger-4">
        <LeaderboardTable
          items={wallets.items}
          locale={locale}
          page={page}
          pageSize={pageSize}
        />
      </div>

      {/* ── Pagination ── */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3 opacity-0 animate-slide-up stagger-5">
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
