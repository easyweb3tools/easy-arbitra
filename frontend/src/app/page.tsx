import Link from "next/link";
import { getOverviewStats, getPotentialWallets } from "@/lib/api";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function HomePage() {
  const locale = await getLocaleFromCookies();
  const [stats, potential] = await Promise.all([
    getOverviewStats(),
    getPotentialWallets(new URLSearchParams({ page: "1", page_size: "5", min_trades: "100", min_realized_pnl: "0" }))
  ]);
  const analyzedCount = potential.items.filter((w) => w.has_ai_report).length;

  return (
    <section className="space-y-4">
      <div className="grid gap-4 md:grid-cols-3">
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <p className="text-sm text-muted">{t(locale, "home.marketsIndexed")}</p>
          <p className="mt-2 text-3xl font-semibold">{stats.indexed_markets}</p>
        </article>
        <article className="rounded-lg border border-emerald-200 bg-emerald-50 p-5 shadow-sm">
          <p className="text-sm text-emerald-700">{t(locale, "home.potentialWallets")}</p>
          <p className="mt-2 text-3xl font-semibold text-emerald-900">{potential.pagination.total}</p>
        </article>
        <article className="rounded-lg border border-blue-200 bg-blue-50 p-5 shadow-sm">
          <p className="text-sm text-blue-700">{t(locale, "home.aiAnalyzedTop5")}</p>
          <p className="mt-2 text-3xl font-semibold text-blue-900">{analyzedCount}</p>
        </article>
      </div>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <div className="mb-4 flex items-center justify-between gap-2">
          <h2 className="text-lg font-semibold">{t(locale, "home.topPotential")}</h2>
          <Link href="/wallets" className="text-sm font-medium text-accent hover:underline">
            {t(locale, "home.viewAll")}
          </Link>
        </div>
        <div className="space-y-2">
          {potential.items.map((item, idx) => (
            <Link key={item.wallet.id} href={`/wallets/${item.wallet.id}`} className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50">
              <p className="font-medium">
                #{idx + 1} {item.wallet.pseudonym || item.wallet.address}
              </p>
              <p className="text-xs text-muted">
                {t(locale, "home.trades")} {item.total_trades} · {t(locale, "home.realizedPnl")} {item.realized_pnl.toFixed(2)} · {t(locale, "home.score")}{" "}
                {item.smart_score} · {item.strategy_type}
              </p>
            </Link>
          ))}
        </div>
      </article>
    </section>
  );
}
