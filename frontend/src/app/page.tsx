import Link from "next/link";
import { getOpsHighlights, getOverviewStats } from "@/lib/api";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { appendUTM, pickUTM } from "@/lib/utm";

export const dynamic = "force-dynamic";

export default async function HomePage({
  searchParams
}: {
  searchParams: { utm_source?: string; utm_medium?: string; utm_campaign?: string; utm_term?: string; utm_content?: string };
}) {
  const locale = await getLocaleFromCookies();
  const utm = pickUTM(searchParams);
  const [stats, highlights] = await Promise.all([
    getOverviewStats(),
    getOpsHighlights(new URLSearchParams({ limit: "5" }))
  ]);

  return (
    <section className="space-y-4">
      <div className="grid gap-4 md:grid-cols-3">
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <p className="text-sm text-muted">{t(locale, "home.marketsIndexed")}</p>
          <p className="mt-2 text-3xl font-semibold">{stats.indexed_markets}</p>
        </article>
        <article className="rounded-lg border border-emerald-200 bg-emerald-50 p-5 shadow-sm">
          <p className="text-sm text-emerald-700">{t(locale, "home.newPotential24h")}</p>
          <p className="mt-2 text-3xl font-semibold text-emerald-900">{highlights.new_potential_wallets_24h}</p>
        </article>
        <article className="rounded-lg border border-blue-200 bg-blue-50 p-5 shadow-sm">
          <p className="text-sm text-blue-700">{t(locale, "home.potentialWallets")}</p>
          <p className="mt-2 text-3xl font-semibold text-blue-900">{highlights.top_realized_pnl_24h.length}</p>
        </article>
      </div>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <div className="mb-4 flex items-center justify-between gap-2">
          <h2 className="text-lg font-semibold">{t(locale, "home.topRealized24h")}</h2>
          <Link href={appendUTM("/wallets", utm)} className="text-sm font-medium text-accent hover:underline">
            {t(locale, "home.viewAll")}
          </Link>
        </div>
        <div className="space-y-2">
          {highlights.top_realized_pnl_24h.map((item, idx) => (
            <Link
              key={item.wallet.id}
              href={appendUTM(`/wallets/${item.wallet.id}`, utm)}
              className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50"
            >
              <p className="font-medium">
                #{idx + 1} {item.wallet.pseudonym || item.wallet.address}
              </p>
              <p className="text-xs text-muted">
                {t(locale, "home.trades")} {item.trade_count} · {t(locale, "home.realizedPnl")} {item.realized_pnl.toFixed(2)} · 24h {item.realized_pnl_24h.toFixed(2)}
              </p>
              <p className="mt-1 text-xs text-muted">
                {item.nl_summary || t(locale, "home.summaryFallback")}
              </p>
            </Link>
          ))}
        </div>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold">{t(locale, "home.topAIConfidence")}</h2>
        <div className="space-y-2">
          {highlights.top_ai_confidence.map((item, idx) => (
            <Link
              key={item.wallet.id}
              href={appendUTM(`/wallets/${item.wallet.id}`, utm)}
              className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50"
            >
              <p className="font-medium">
                #{idx + 1} {item.wallet.pseudonym || item.wallet.address}
              </p>
              <p className="text-xs text-muted">
                {t(locale, "home.score")} {item.smart_score} · {item.strategy_type} · {item.info_edge_level}
              </p>
              <p className="mt-1 text-xs text-muted">{item.nl_summary || t(locale, "home.summaryFallback")}</p>
            </Link>
          ))}
        </div>
      </article>
    </section>
  );
}
