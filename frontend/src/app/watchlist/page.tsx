import { WatchlistPageClient } from "@/components/watchlist/WatchlistPageClient";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function WatchlistPage() {
  const locale = await getLocaleFromCookies();
  return (
    <WatchlistPageClient
      locale={locale}
      labels={{
        title: t(locale, "watchlist.title"),
        feedTitle: t(locale, "watchlist.feedTitle"),
        empty: t(locale, "watchlist.empty"),
        trades: t(locale, "home.trades"),
        realizedPnl: t(locale, "home.realizedPnl"),
        score: t(locale, "home.score"),
        loading: t(locale, "watchlist.loading"),
        eventType: t(locale, "watchlist.eventType"),
        eventTime: t(locale, "watchlist.eventTime"),
      }}
    />
  );
}
