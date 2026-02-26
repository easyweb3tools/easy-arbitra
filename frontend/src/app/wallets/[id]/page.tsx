import { ChevronRight, AlertTriangle } from "lucide-react";
import {
  getWalletAIReport,
  getWalletAIReportHistory,
  getWalletDecisionCard,
  getWalletPnLHistory,
  getWalletPositions,
  getWalletProfile,
  getWalletTrades,
  getCopyTradeConfig,
} from "@/lib/api";
import { TriggerAnalysisButton } from "@/components/ai/TriggerAnalysisButton";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { CopyTradeToggle } from "@/components/copytrade/CopyTradeToggle";
import { Card, SectionHeader } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import { EmptyState } from "@/components/ui/EmptyState";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { DecisionCard } from "@/components/wallet/DecisionCard";
import { AccountInfoCard } from "@/components/wallet/AccountInfoCard";
import { PnLChart } from "@/components/wallet/PnLChart";
import { PositionsList } from "@/components/wallet/PositionsList";
import { TradeHistory } from "@/components/wallet/TradeHistory";

export const dynamic = "force-dynamic";

function toRiskWarnings(report: { risk_warnings?: string[]; report: unknown }): string[] {
  if (Array.isArray(report.risk_warnings)) {
    return report.risk_warnings.filter((v): v is string => typeof v === "string" && v.trim() !== "");
  }
  const payload = report.report;
  if (!payload || typeof payload !== "object") return [];
  const fromReport = (payload as Record<string, unknown>).risk_warnings;
  if (!Array.isArray(fromReport)) return [];
  return fromReport.filter((v): v is string => typeof v === "string" && v.trim() !== "");
}

export default async function WalletProfilePage({
  params,
  searchParams,
}: {
  params: { id: string };
  searchParams: { trades_page?: string };
}) {
  const locale = await getLocaleFromCookies();
  const tradesPage = Math.max(1, parseInt(searchParams.trades_page || "1", 10) || 1);

  const [profile, decisionCard, pnlHistory, positions, trades] = await Promise.all([
    getWalletProfile(params.id),
    getWalletDecisionCard(params.id).catch(() => null),
    getWalletPnLHistory(params.id, 90).catch(() => []),
    getWalletPositions(params.id).catch(() => []),
    getWalletTrades(params.id, new URLSearchParams({ page: String(tradesPage), page_size: "20" })).catch(() => ({
      items: [],
      pagination: { page: tradesPage, page_size: 20, total: 0 },
    })),
  ]);

  let aiReport = null;
  let aiHistory: Awaited<ReturnType<typeof getWalletAIReportHistory>> = [];
  try {
    aiReport = await getWalletAIReport(params.id);
    aiHistory = await getWalletAIReportHistory(params.id);
  } catch {
    aiReport = null;
  }

  const warnings = aiReport ? toRiskWarnings(aiReport) : [];

  return (
    <section className="space-y-8 animate-fade-in">
      {/* ── 1. Account Info Card ── */}
      <AccountInfoCard profile={profile} decisionCard={decisionCard} locale={locale} />

      {/* ── AI Copy Trading Toggle ── */}
      <div className="flex items-center gap-3">
        <CopyTradeToggle
          walletID={parseInt(params.id, 10)}
          enabled={false}
          locale={locale}
        />
      </div>

      {/* ── 2. P&L Chart ── */}
      {pnlHistory.length > 0 && (
        <PnLChart
          data={pnlHistory}
          labels={{
            title: t(locale, "pnlChart.title"),
            d7: t(locale, "pnlChart.7d"),
            d30: t(locale, "pnlChart.30d"),
            d90: t(locale, "pnlChart.90d"),
          }}
        />
      )}

      {/* ── 3. Positions ── */}
      <PositionsList positions={positions} locale={locale} />

      {/* ── 4. Trade History ── */}
      <TradeHistory
        items={trades.items}
        locale={locale}
        page={trades.pagination.page}
        pageSize={trades.pagination.page_size}
        total={trades.pagination.total}
        walletId={params.id}
      />

      {/* ── 5. Decision Card ── */}
      {decisionCard && <DecisionCard card={decisionCard} locale={locale} />}

      {/* ── 6. AI Analysis ── */}
      <Card variant="prominent">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-title-3 text-label-primary">{t(locale, "profile.aiAnalysis")}</h2>
          <TriggerAnalysisButton
            walletID={params.id}
            labels={{
              trigger: t(locale, "ai.trigger"),
              loading: t(locale, "ai.loading"),
              updatedAt: t(locale, "ai.updatedAt"),
              failedPrefix: t(locale, "ai.failedPrefix"),
              requestFailed: t(locale, "ai.requestFailed"),
            }}
          />
        </div>
        {aiReport ? (
          <div className="space-y-4">
            <div className="rounded-xl bg-tint-green/[0.06] p-5">
              <p className="text-callout font-medium leading-relaxed text-label-primary">{aiReport.nl_summary || t(locale, "profile.noSummary")}</p>
            </div>

            <div className="flex flex-wrap gap-x-4 gap-y-1 text-caption-1 text-label-quaternary">
              <span>{t(locale, "profile.model")} {aiReport.model_id}</span>
              <span>{aiReport.input_tokens}/{aiReport.output_tokens} {t(locale, "profile.tokens")}</span>
              <span>{aiReport.latency_ms}ms</span>
              <span>{t(locale, "profile.aiGeneratedAt")} {aiReport.created_at}</span>
            </div>

            {warnings.length > 0 && (
              <div className="rounded-xl bg-tint-orange/[0.06] p-5">
                <h4 className="mb-2 flex items-center gap-2 text-caption-1 font-bold uppercase tracking-widest text-tint-orange">
                  <AlertTriangle className="h-3.5 w-3.5" />
                  {t(locale, "profile.aiWarnings")}
                </h4>
                <ul className="space-y-1.5 pl-4">
                  {warnings.map((w, i) => (
                    <li key={`${w}-${i}`} className="list-disc text-footnote leading-relaxed text-label-secondary">{w}</li>
                  ))}
                </ul>
              </div>
            )}

            <details className="group rounded-xl border border-separator/50 overflow-hidden">
              <summary className="flex cursor-pointer items-center gap-2 px-5 py-3.5 text-subheadline font-medium text-label-secondary transition-colors duration-200 hover:bg-surface-tertiary/70">
                <ChevronRight className="h-4 w-4 transition-transform duration-200 group-open:rotate-90" />
                {t(locale, "profile.aiRawJson")}
              </summary>
              <div className="border-t border-separator/50 p-5">
                <pre className="overflow-auto rounded-lg bg-surface-tertiary/70 p-4 font-mono text-caption-1 text-label-secondary">
                  {JSON.stringify(aiReport.report, null, 2)}
                </pre>
              </div>
            </details>
          </div>
        ) : (
          <EmptyState preset="no-ai-report" locale={locale} />
        )}
        {aiHistory.length > 0 && (
          <div className="mt-5 border-t border-separator/50 pt-3 text-caption-1 text-label-quaternary">
            {aiHistory.map((h) => (
              <p key={h.id} className="py-0.5">{t(locale, "profile.history")}: {h.created_at}</p>
            ))}
          </div>
        )}
      </Card>

      {/* ── 7. Layer 3 Info-Edge (compact) ── */}
      <Card>
        <h3 className="mb-3 text-subheadline font-semibold text-label-secondary">{t(locale, "profile.layer3")}</h3>
        <div className="grid gap-4 sm:grid-cols-4">
          <StatCell label={t(locale, "profile.label")} value={profile.layer3_info_edge.label} />
          <StatCell label={t(locale, "profile.meanDt")} value={`${profile.layer3_info_edge.mean_delta_minutes.toFixed(2)} min`} />
          <StatCell label={t(locale, "profile.samples")} value={String(profile.layer3_info_edge.samples)} />
          <StatCell label={t(locale, "profile.pvalue")} value={profile.layer3_info_edge.p_value.toFixed(4)} />
        </div>
      </Card>

      {/* ── 8. Disclosures ── */}
      <div className="rounded-2xl bg-tint-orange/[0.05] p-5">
        <div className="space-y-1.5">
          {profile.meta.disclosures.map((d) => (
            <p key={d} className="text-caption-1 leading-relaxed text-tint-orange">{d}</p>
          ))}
        </div>
      </div>
    </section>
  );
}
