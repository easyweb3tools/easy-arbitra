import Link from "next/link";
import { Card } from "@/components/ui/Card";
import { TierBadge, CategoryTag, StatusBadge } from "@/components/ui/Badge";
import type { CopyTradingConfig } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

function fmt$(v: number) {
  return v >= 0 ? `+$${v.toFixed(2)}` : `-$${Math.abs(v).toFixed(2)}`;
}

export function ConfigList({
  configs,
  locale,
}: {
  configs: CopyTradingConfig[];
  locale: Locale;
}) {
  if (configs.length === 0) {
    return (
      <p className="py-8 text-center text-subheadline text-label-tertiary">
        {t(locale, "copyTrade.noConfigs")}
      </p>
    );
  }

  return (
    <Card padding={false}>
      {configs.map((cfg) => (
        <Link
          key={cfg.id}
          href={`/wallets/${cfg.wallet_id}`}
          className="flex items-center justify-between gap-3 px-5 py-3.5 border-b border-separator/60 last:border-b-0 hover:bg-surface-tertiary/70 cursor-pointer transition-colors duration-200"
        >
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <span className="text-subheadline font-semibold text-label-primary truncate">
                {cfg.wallet_pseudonym || `0x${cfg.wallet_address.slice(2, 8)}...${cfg.wallet_address.slice(-4)}`}
              </span>
              {cfg.enabled ? (
                <StatusBadge color="green">{t(locale, "copyTrade.simulated")}</StatusBadge>
              ) : (
                <StatusBadge color="gray">{t(locale, "copyTrade.disable")}</StatusBadge>
              )}
            </div>
            <div className="flex items-center gap-3 mt-0.5 text-caption-1 text-label-tertiary">
              <span>
                {t(locale, "copyTrade.maxPosition")}: ${cfg.max_position_usdc.toLocaleString()}
              </span>
              <CategoryTag color="blue">{cfg.risk_preference}</CategoryTag>
            </div>
          </div>
          <div className="text-right">
            <p className={`text-headline font-semibold tabular-nums ${cfg.total_pnl >= 0 ? "text-tint-green" : "text-tint-red"}`}>
              {fmt$(cfg.total_pnl)}
            </p>
            <p className="text-caption-1 text-label-tertiary">
              {cfg.open_positions} {t(locale, "copyTrade.openPositions")}
            </p>
          </div>
        </Link>
      ))}
    </Card>
  );
}
