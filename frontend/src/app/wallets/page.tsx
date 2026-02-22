import Link from "next/link";
import { getPotentialWallets } from "@/lib/api";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { appendUTM, pickUTM } from "@/lib/utm";

export const dynamic = "force-dynamic";

export default async function WalletsPage({
  searchParams
}: {
  searchParams: {
    page?: string;
    page_size?: string;
    min_trades?: string;
    min_realized_pnl?: string;
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
    min_realized_pnl: searchParams.min_realized_pnl || "0"
  });

  const wallets = await getPotentialWallets(params);
  const analyzedCount = wallets.items.filter((w) => w.has_ai_report).length;

  return (
    <section className="space-y-4">
      <form className="flex flex-wrap gap-2 rounded-lg bg-card p-4 shadow-sm" method="get">
        <input
          name="min_trades"
          type="number"
          min="1"
          defaultValue={searchParams.min_trades || "100"}
          placeholder={t(locale, "wallets.minTrades")}
          className="w-40 rounded-md border border-slate-300 px-3 py-2 text-sm"
        />
        <input
          name="min_realized_pnl"
          type="number"
          step="0.01"
          defaultValue={searchParams.min_realized_pnl || "0"}
          placeholder={t(locale, "wallets.minPnl")}
          className="w-48 rounded-md border border-slate-300 px-3 py-2 text-sm"
        />
        <button className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-white" type="submit">
          {t(locale, "wallets.apply")}
        </button>
        {utm.get("utm_source") ? <input type="hidden" name="utm_source" value={utm.get("utm_source") || ""} /> : null}
        {utm.get("utm_medium") ? <input type="hidden" name="utm_medium" value={utm.get("utm_medium") || ""} /> : null}
        {utm.get("utm_campaign") ? <input type="hidden" name="utm_campaign" value={utm.get("utm_campaign") || ""} /> : null}
        {utm.get("utm_term") ? <input type="hidden" name="utm_term" value={utm.get("utm_term") || ""} /> : null}
        {utm.get("utm_content") ? <input type="hidden" name="utm_content" value={utm.get("utm_content") || ""} /> : null}
      </form>

      <article className="grid gap-3 sm:grid-cols-3">
        <div className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
          <p className="text-xs uppercase tracking-wide text-muted">{t(locale, "wallets.potentialCount")}</p>
          <p className="mt-1 text-2xl font-semibold">{wallets.pagination.total}</p>
        </div>
        <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 shadow-sm">
          <p className="text-xs uppercase tracking-wide text-blue-700">{t(locale, "wallets.analyzedPage")}</p>
          <p className="mt-1 text-2xl font-semibold text-blue-900">{analyzedCount}</p>
        </div>
        <div className="rounded-lg border border-emerald-200 bg-emerald-50 p-4 shadow-sm">
          <p className="text-xs uppercase tracking-wide text-emerald-700">{t(locale, "wallets.filterRule")}</p>
          <p className="mt-1 text-sm font-semibold text-emerald-900">{t(locale, "wallets.ruleValue")}</p>
        </div>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-2 text-lg font-semibold">{t(locale, "wallets.title")}</h2>
        <p className="mb-4 text-xs text-muted">{t(locale, "wallets.sortHint")}</p>
        <div className="space-y-2">
          {wallets.items.map((item) => (
            <div key={item.wallet.id} className="rounded-md border border-slate-200 p-3 hover:bg-slate-50">
              <div className="flex items-start justify-between gap-2">
                <Link href={appendUTM(`/wallets/${item.wallet.id}`, utm)} className="min-w-0 flex-1">
                  <p className="font-medium">{item.wallet.pseudonym || t(locale, "wallets.unnamed")}</p>
                  <p className="truncate text-xs text-muted">{item.wallet.address}</p>
                </Link>
                <div className="flex flex-col items-end gap-1">
                  <WatchlistToggleButton
                    walletID={item.wallet.id}
                    labels={{
                      follow: t(locale, "watchlist.follow"),
                      unfollow: t(locale, "watchlist.unfollow"),
                      following: t(locale, "watchlist.following"),
                      failed: t(locale, "watchlist.failed")
                    }}
                  />
                  <div className="flex flex-wrap justify-end gap-1">
                    <span className="rounded bg-emerald-100 px-2 py-0.5 text-xs font-medium text-emerald-700">{t(locale, "wallets.tagPotential")}</span>
                    <span className={`rounded px-2 py-0.5 text-xs font-medium ${item.has_ai_report ? "bg-blue-100 text-blue-700" : "bg-slate-100 text-slate-600"}`}>
                      {item.has_ai_report ? t(locale, "wallets.tagAnalyzed") : t(locale, "wallets.tagNotAnalyzed")}
                    </span>
                  </div>
                </div>
              </div>
              <Link href={appendUTM(`/wallets/${item.wallet.id}`, utm)} className="block">
                <p className="mt-2 text-xs text-muted">
                  {t(locale, "home.trades")} {item.total_trades} · {t(locale, "home.realizedPnl")} {item.realized_pnl.toFixed(2)} · {t(locale, "home.score")}{" "}
                  {item.smart_score} · {item.strategy_type || t(locale, "wallets.strategyUnknown")} / {item.info_edge_level || t(locale, "wallets.strategyUnknown")}
                </p>
              </Link>
            </div>
          ))}
        </div>
      </article>
    </section>
  );
}
