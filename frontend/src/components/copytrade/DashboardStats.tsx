import { Card } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import type { CopyTradeDashboard } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

function fmt$(v: number) {
  const sign = v >= 0 ? "+" : "-";
  return `${sign}$${Math.abs(v).toFixed(2)}`;
}

export function DashboardStats({
  data,
  locale,
}: {
  data: CopyTradeDashboard;
  locale: Locale;
}) {
  return (
    <Card variant="glass">
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
        <StatCell
          label={t(locale, "copyTrade.totalPnl")}
          value={fmt$(data.total_pnl)}
          numericValue={data.total_pnl}
          size="large"
        />
        <StatCell
          label={t(locale, "copyTrade.winRate")}
          value={`${(data.win_rate * 100).toFixed(0)}%`}
        />
        <StatCell
          label={t(locale, "copyTrade.totalCopies")}
          value={String(data.total_copies)}
        />
        <StatCell
          label={t(locale, "copyTrade.totalSkipped")}
          value={String(data.total_skipped)}
        />
        <StatCell
          label={t(locale, "copyTrade.openPositions")}
          value={String(data.open_positions)}
        />
        <StatCell
          label={t(locale, "copyTrade.activeConfigs")}
          value={String(data.active_configs)}
        />
      </div>
    </Card>
  );
}
