import Link from "next/link";
import type { TradeHistoryItem } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";

export function TradeHistory({
  items,
  locale,
  page,
  pageSize,
  total,
  walletId,
}: {
  items: TradeHistoryItem[];
  locale: Locale;
  page: number;
  pageSize: number;
  total: number;
  walletId: string;
}) {
  const totalPages = Math.ceil(total / pageSize);

  return (
    <div>
      <div className="flex items-center justify-between pb-3">
        <h2 className="text-title-3 text-label-primary">{t(locale, "tradeHistory.title")}</h2>
        <span className="text-caption-1 text-label-quaternary">
          {total.toLocaleString()} {locale === "zh" ? "笔" : "total"}
        </span>
      </div>
      <Card padding={false}>
        {items.length === 0 ? (
          <EmptyState preset="no-activity" locale={locale} />
        ) : (
          <>
            {/* Header */}
            <div className="hidden sm:grid sm:grid-cols-[7rem_1fr_3.5rem_3.5rem_4rem_4rem_4rem] items-center gap-2 border-b border-separator/60 px-5 py-2.5 text-caption-1 font-semibold uppercase tracking-wider text-label-tertiary">
              <span>{t(locale, "tradeHistory.time")}</span>
              <span>{t(locale, "tradeHistory.market")}</span>
              <span>{t(locale, "tradeHistory.outcome")}</span>
              <span>{t(locale, "tradeHistory.action")}</span>
              <span className="text-right">{t(locale, "tradeHistory.price")}</span>
              <span className="text-right">{t(locale, "tradeHistory.size")}</span>
              <span className="text-right">{t(locale, "tradeHistory.fee")}</span>
            </div>

            {items.map((trade) => {
              const actionColor = trade.action === "Buy" ? "text-tint-green" : "text-tint-red";
              const outcomeLabel = trade.outcome === "Yes" ? t(locale, "tradeHistory.yes") : t(locale, "tradeHistory.no");
              const actionLabel = trade.action === "Buy" ? t(locale, "tradeHistory.buy") : t(locale, "tradeHistory.sell");

              return (
                <div
                  key={trade.id}
                  className="flex flex-col gap-1.5 border-b border-separator/60 px-5 py-3.5 last:border-b-0 sm:grid sm:grid-cols-[7rem_1fr_3.5rem_3.5rem_4rem_4rem_4rem] sm:items-center sm:gap-2"
                >
                  <span className="text-caption-1 tabular-nums text-label-quaternary">
                    {new Date(trade.block_time).toLocaleString(locale === "zh" ? "zh-CN" : "en-US", {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                  <p className="truncate text-subheadline text-label-primary">{trade.market_title}</p>
                  <span className="text-subheadline font-medium text-label-secondary">{outcomeLabel}</span>
                  <span className={`text-subheadline font-semibold ${actionColor}`}>{actionLabel}</span>
                  <span className="text-right text-subheadline tabular-nums text-label-secondary">
                    {trade.price.toFixed(3)}
                  </span>
                  <span className="text-right text-subheadline tabular-nums text-label-secondary">
                    {trade.size.toFixed(2)}
                  </span>
                  <span className="text-right text-subheadline tabular-nums text-label-quaternary">
                    {trade.fee_paid.toFixed(4)}
                  </span>
                </div>
              );
            })}

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2 border-t border-separator/60 px-5 py-3">
                {page > 1 && (
                  <Link
                    href={`/wallets/${walletId}?trades_page=${page - 1}`}
                    className="rounded-lg bg-surface-tertiary/80 px-3 py-1.5 text-caption-1 font-semibold text-label-secondary transition-colors hover:bg-surface-tertiary"
                  >
                    {locale === "zh" ? "上一页" : "Prev"}
                  </Link>
                )}
                <span className="text-caption-1 tabular-nums text-label-tertiary">
                  {page} / {totalPages}
                </span>
                {page < totalPages && (
                  <Link
                    href={`/wallets/${walletId}?trades_page=${page + 1}`}
                    className="rounded-lg bg-surface-tertiary/80 px-3 py-1.5 text-caption-1 font-semibold text-label-secondary transition-colors hover:bg-surface-tertiary"
                  >
                    {locale === "zh" ? "下一页" : "Next"}
                  </Link>
                )}
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
}
