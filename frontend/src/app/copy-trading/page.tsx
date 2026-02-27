import { LOCALE_COOKIE, normalizeLocale, t } from "@/lib/i18n";
import { cookies } from "next/headers";
import { CopyTradingPageClient } from "@/components/copytrade/CopyTradingPageClient";

export const dynamic = "force-dynamic";

export default async function CopyTradingPage() {
  const store = await cookies();
  const locale = normalizeLocale(store.get(LOCALE_COOKIE)?.value);

  return (
    <CopyTradingPageClient
      locale={locale}
      labels={{
        title: t(locale, "copyTrade.dashboard"),
        activeConfigs: t(locale, "copyTrade.activeConfigs"),
        positions: t(locale, "copyTrade.positions"),
        recentDecisions: t(locale, "copyTrade.recentDecisions"),
        noConfigs: t(locale, "copyTrade.noConfigs"),
      }}
    />
  );
}
