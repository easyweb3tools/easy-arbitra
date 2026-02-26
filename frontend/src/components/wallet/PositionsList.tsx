import type { WalletPosition } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { CategoryTag } from "@/components/ui/Badge";

export function PositionsList({
  positions,
  locale,
}: {
  positions: WalletPosition[];
  locale: Locale;
}) {
  if (positions.length === 0) {
    return (
      <Card padding={false}>
        <EmptyState preset="no-activity" locale={locale} />
      </Card>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between pb-3">
        <h2 className="text-title-3 text-label-primary">{t(locale, "positions.title")}</h2>
      </div>
      <Card padding={false}>
        {/* Header */}
        <div className="hidden sm:grid sm:grid-cols-[1fr_5rem_5rem_5rem_5rem_6.5rem] items-center gap-2 border-b border-separator/60 px-5 py-2.5 text-caption-1 font-semibold uppercase tracking-wider text-label-tertiary">
          <span>{t(locale, "positions.market")}</span>
          <span className="text-right">{t(locale, "positions.netSize")}</span>
          <span className="text-right">{t(locale, "positions.avgPrice")}</span>
          <span className="text-right">{t(locale, "positions.totalVolume")}</span>
          <span className="text-right">{t(locale, "positions.tradeCount")}</span>
          <span className="text-right">{t(locale, "positions.lastActive")}</span>
        </div>

        {positions.map((pos) => {
          const sizeColor = pos.net_size > 0 ? "text-tint-green" : pos.net_size < 0 ? "text-tint-red" : "text-label-secondary";
          return (
            <div
              key={pos.market_id}
              className="flex flex-col gap-2 border-b border-separator/60 px-5 py-4 last:border-b-0 sm:grid sm:grid-cols-[1fr_5rem_5rem_5rem_5rem_6.5rem] sm:items-center sm:gap-2"
            >
              <div className="min-w-0">
                <p className="truncate text-subheadline font-medium text-label-primary">{pos.market_title}</p>
                <CategoryTag>{pos.category || "—"}</CategoryTag>
              </div>
              <span className={`text-right text-subheadline font-semibold tabular-nums ${sizeColor}`}>
                {pos.net_size >= 0 ? "+" : ""}{pos.net_size.toFixed(2)}
              </span>
              <span className="text-right text-subheadline tabular-nums text-label-secondary">
                {pos.avg_price.toFixed(3)}
              </span>
              <span className="text-right text-subheadline tabular-nums text-label-secondary">
                {pos.total_volume.toFixed(2)}
              </span>
              <span className="text-right text-subheadline tabular-nums text-label-secondary">
                {pos.trade_count}
              </span>
              <span className="text-right text-caption-1 text-label-quaternary">
                {new Date(pos.last_trade_at).toLocaleDateString()}
              </span>
            </div>
          );
        })}
      </Card>
    </div>
  );
}
