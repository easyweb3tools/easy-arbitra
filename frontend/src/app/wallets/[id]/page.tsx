import { ChevronRight, AlertTriangle } from "lucide-react";
import {
  getWalletAIReport,
  getWalletAIReportHistory,
  getWalletDecisionCard,
  getWalletExplanation,
  getWalletProfile,
  getWalletShareCard,
} from "@/lib/api";
import { TriggerAnalysisButton } from "@/components/ai/TriggerAnalysisButton";
import { ShareCardPanel } from "@/components/share/ShareCardPanel";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { Card, SectionHeader } from "@/components/ui/Card";
import { StatCell } from "@/components/ui/StatCell";
import { CategoryTag } from "@/components/ui/Badge";
import { EmptyState } from "@/components/ui/EmptyState";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";
import { DecisionCard } from "@/components/wallet/DecisionCard";

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

export default async function WalletProfilePage({ params }: { params: { id: string } }) {
  const locale = await getLocaleFromCookies();
  const [profile, explanation, shareCard, decisionCard] = await Promise.all([
    getWalletProfile(params.id),
    getWalletExplanation(params.id),
    getWalletShareCard(params.id),
    getWalletDecisionCard(params.id),
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
      {/* ── Profile Header ── */}
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <h1 className="text-title-1 tracking-tight text-label-primary">{profile.wallet.pseudonym || profile.wallet.address}</h1>
          <p className="mt-1 truncate font-mono text-footnote text-label-quaternary">{profile.wallet.address}</p>
          {profile.strategy && (
            <div className="mt-3 flex flex-wrap items-center gap-2.5">
              <CategoryTag>{profile.strategy.strategy_type}</CategoryTag>
              <span className="text-footnote text-label-tertiary">
                Score <span className="font-bold tabular-nums text-label-primary">{profile.strategy.smart_score}</span>
              </span>
              <span className="h-3.5 w-px bg-separator" />
              <span className="text-footnote text-label-tertiary">{profile.strategy.info_edge_level}</span>
            </div>
          )}
        </div>
        <WatchlistToggleButton
          walletID={profile.wallet.id}
          labels={{
            follow: t(locale, "watchlist.follow"),
            unfollow: t(locale, "watchlist.unfollow"),
            following: t(locale, "watchlist.following"),
            failed: t(locale, "watchlist.failed"),
          }}
        />
      </div>

      {/* ── Decision Card ── */}
      <DecisionCard card={decisionCard} locale={locale} />

      {/* ── Recent Events ── */}
      {profile.recent_events && profile.recent_events.length > 0 && (
        <div>
          <SectionHeader title={locale === "zh" ? "近期动态" : "Recent Events"} />
          <Card padding={false}>
            {profile.recent_events.slice(0, 6).map((event) => (
              <div key={`${event.event_type}-${event.event_id}`} className="border-b border-separator/60 px-5 py-4 last:border-b-0">
                <div className="flex items-center justify-between gap-2">
                  <p className="text-subheadline font-medium text-label-primary">{event.event_type}</p>
                  {event.action_required && (
                    <span className="rounded-full bg-tint-red/10 px-2.5 py-0.5 text-caption-1 font-semibold text-tint-red">{locale === "zh" ? "需行动" : "Action"}</span>
                  )}
                </div>
                <p className="mt-1 text-footnote text-label-tertiary leading-relaxed">{locale === "zh" ? event.suggestion_zh : event.suggestion}</p>
                <p className="mt-1 text-caption-2 text-label-quaternary">{event.event_time}</p>
              </div>
            ))}
          </Card>
        </div>
      )}

      {/* ── AI Analysis ── */}
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

      {/* ── Layer 1 Facts ── */}
      <div>
        <SectionHeader title={t(locale, "profile.layer1")} />
        <Card>
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <StatCell label={t(locale, "profile.realizedPnl")} value={profile.layer1_facts.realized_pnl.toFixed(2)} numericValue={profile.layer1_facts.realized_pnl} />
            <StatCell label={t(locale, "profile.tradingPnl")} value={profile.layer1_facts.trading_pnl.toFixed(2)} numericValue={profile.layer1_facts.trading_pnl} />
            <StatCell label={t(locale, "profile.makerRebates")} value={profile.layer1_facts.maker_rebates.toFixed(2)} numericValue={profile.layer1_facts.maker_rebates} />
            <StatCell label={t(locale, "profile.feesPaid")} value={profile.layer1_facts.fees_paid.toFixed(2)} numericValue={profile.layer1_facts.fees_paid} />
            <StatCell label={t(locale, "profile.totalTrades")} value={String(profile.layer1_facts.total_trades)} />
            <StatCell label={t(locale, "profile.volume30d")} value={profile.layer1_facts.volume_30d.toFixed(2)} />
          </div>
        </Card>
      </div>

      {/* ── Layer 3 Info-Edge ── */}
      <div>
        <SectionHeader title={t(locale, "profile.layer3")} />
        <Card>
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            <StatCell label={t(locale, "profile.label")} value={profile.layer3_info_edge.label} />
            <StatCell label={t(locale, "profile.meanDt")} value={`${profile.layer3_info_edge.mean_delta_minutes.toFixed(2)} min`} />
            <StatCell label={t(locale, "profile.samples")} value={String(profile.layer3_info_edge.samples)} />
            <StatCell label={t(locale, "profile.pvalue")} value={profile.layer3_info_edge.p_value.toFixed(4)} />
          </div>
        </Card>
      </div>

      {/* ── Disclosures ── */}
      <div className="rounded-2xl bg-tint-orange/[0.05] p-5">
        <div className="space-y-1.5">
          {profile.meta.disclosures.map((d) => (
            <p key={d} className="text-caption-1 leading-relaxed text-tint-orange">{d}</p>
          ))}
        </div>
      </div>

      {/* ── Share Card ── */}
      <ShareCardPanel
        card={shareCard}
        locale={locale}
        labels={{
          title: t(locale, "share.title"),
          preview: t(locale, "share.preview"),
          copyLink: t(locale, "share.copyLink"),
          copied: t(locale, "share.copied"),
          copyFailed: t(locale, "share.copyFailed"),
          trades: t(locale, "share.trades"),
          realizedPnl: t(locale, "share.realizedPnl"),
          score: t(locale, "share.score"),
          updatedAt: t(locale, "share.updatedAt"),
        }}
      />

      {/* ── Evidence (raw JSON) ── */}
      <details className="group rounded-2xl border border-separator/50 bg-surface-secondary shadow-elevation-1 overflow-hidden">
        <summary className="flex cursor-pointer items-center gap-2.5 px-6 py-4 text-headline text-label-primary transition-colors duration-200 hover:bg-surface-tertiary/70">
          <ChevronRight className="h-4 w-4 transition-transform duration-200 group-open:rotate-90" />
          {t(locale, "profile.evidence")}
        </summary>
        <div className="border-t border-separator/50 p-6">
          <pre className="overflow-auto rounded-xl bg-surface-tertiary/70 p-4 font-mono text-caption-1 text-label-secondary">
            {JSON.stringify(explanation.layer2, null, 2)}
          </pre>
        </div>
      </details>
    </section>
  );
}
