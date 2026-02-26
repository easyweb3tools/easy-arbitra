import { cookies } from "next/headers";
import { LOCALE_COOKIE, normalizeLocale, t } from "@/lib/i18n";
import { getCopyTradeDashboard, getCopyTradePositions } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { DashboardStats } from "@/components/copytrade/DashboardStats";
import { ConfigList } from "@/components/copytrade/ConfigList";
import { DecisionFeed } from "@/components/copytrade/DecisionFeed";
import { OpenPositions } from "@/components/copytrade/OpenPositions";

export const dynamic = "force-dynamic";

export default async function CopyTradingPage() {
  const store = await cookies();
  const locale = normalizeLocale(store.get(LOCALE_COOKIE)?.value);

  // Use a server-side fingerprint from cookie for SSR
  const fp = store.get("user_fingerprint")?.value || "server";

  let dashboard = null;
  let positions: Awaited<ReturnType<typeof getCopyTradePositions>> = [];

  try {
    dashboard = await getCopyTradeDashboard(fp);
    positions = await getCopyTradePositions(fp);
  } catch {
    // will show empty state
  }

  return (
    <section className="space-y-8 animate-fade-in">
      <h1 className="text-large-title font-bold text-label-primary">
        {t(locale, "copyTrade.dashboard")}
      </h1>

      {dashboard ? (
        <>
          {/* Stats strip */}
          <DashboardStats data={dashboard} locale={locale} />

          {/* Active configs */}
          <div>
            <SectionHeader title={t(locale, "copyTrade.activeConfigs")} />
            <ConfigList configs={dashboard.configs} locale={locale} />
          </div>

          {/* Open positions */}
          <div>
            <SectionHeader title={t(locale, "copyTrade.positions")} />
            <OpenPositions positions={positions} locale={locale} />
          </div>

          {/* Recent decisions */}
          <div>
            <SectionHeader title={t(locale, "copyTrade.recentDecisions")} />
            <DecisionFeed decisions={dashboard.recent_decisions} locale={locale} />
          </div>
        </>
      ) : (
        <Card>
          <p className="py-12 text-center text-subheadline text-label-tertiary">
            {t(locale, "copyTrade.noConfigs")}
          </p>
        </Card>
      )}
    </section>
  );
}
