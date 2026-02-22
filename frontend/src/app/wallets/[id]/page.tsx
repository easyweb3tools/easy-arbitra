import { getWalletAIReport, getWalletAIReportHistory, getWalletExplanation, getWalletProfile, getWalletShareCard } from "@/lib/api";
import { TriggerAnalysisButton } from "@/components/ai/TriggerAnalysisButton";
import { ShareCardPanel } from "@/components/share/ShareCardPanel";
import { WatchlistToggleButton } from "@/components/watchlist/WatchlistToggleButton";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

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
  const profile = await getWalletProfile(params.id);
  const explanation = await getWalletExplanation(params.id);
  const shareCard = await getWalletShareCard(params.id);

  let aiReport = null;
  let aiHistory: Awaited<ReturnType<typeof getWalletAIReportHistory>> = [];
  try {
    aiReport = await getWalletAIReport(params.id);
    aiHistory = await getWalletAIReportHistory(params.id);
  } catch {
    aiReport = null;
  }

  return (
    <section className="space-y-4">
      <article className="rounded-lg bg-card p-5 shadow-sm">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h2 className="text-lg font-semibold">{profile.wallet.pseudonym || profile.wallet.address}</h2>
            <p className="text-xs text-muted">{profile.wallet.address}</p>
          </div>
          <WatchlistToggleButton
            walletID={profile.wallet.id}
            labels={{
              follow: t(locale, "watchlist.follow"),
              unfollow: t(locale, "watchlist.unfollow"),
              following: t(locale, "watchlist.following"),
              failed: t(locale, "watchlist.failed")
            }}
          />
        </div>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-3 text-base font-semibold">{t(locale, "profile.layer1")}</h3>
        <div className="grid gap-3 text-sm md:grid-cols-3">
          <div>{t(locale, "profile.realizedPnl")}: {profile.layer1_facts.realized_pnl.toFixed(2)}</div>
          <div>{t(locale, "profile.tradingPnl")}: {profile.layer1_facts.trading_pnl.toFixed(2)}</div>
          <div>{t(locale, "profile.makerRebates")}: {profile.layer1_facts.maker_rebates.toFixed(2)}</div>
          <div>{t(locale, "profile.feesPaid")}: {profile.layer1_facts.fees_paid.toFixed(2)}</div>
          <div>{t(locale, "profile.totalTrades")}: {profile.layer1_facts.total_trades}</div>
          <div>{t(locale, "profile.volume30d")}: {profile.layer1_facts.volume_30d.toFixed(2)}</div>
        </div>
      </article>

      {profile.strategy && (
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <h3 className="mb-2 text-base font-semibold">{t(locale, "profile.strategySnapshot")}</h3>
          <p className="text-sm">
            {profile.strategy.strategy_type} · {t(locale, "profile.score")} {profile.strategy.smart_score} · {profile.strategy.info_edge_level}
          </p>
        </article>
      )}

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "profile.layer3")}</h3>
        <p className="text-sm">
          {t(locale, "profile.label")} {profile.layer3_info_edge.label} · {t(locale, "profile.meanDt")}{" "}
          {profile.layer3_info_edge.mean_delta_minutes.toFixed(2)} min · {t(locale, "profile.samples")}{" "}
          {profile.layer3_info_edge.samples}
        </p>
        <p className="mt-1 text-xs text-muted">
          {t(locale, "profile.stddev")} {profile.layer3_info_edge.stddev_minutes.toFixed(2)} · {t(locale, "profile.pvalue")}{" "}
          {profile.layer3_info_edge.p_value.toFixed(4)}
        </p>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "profile.aiAnalysis")}</h3>
        <div className="mb-3">
          <TriggerAnalysisButton
            walletID={params.id}
            labels={{
              trigger: t(locale, "ai.trigger"),
              loading: t(locale, "ai.loading"),
              updatedAt: t(locale, "ai.updatedAt"),
              failedPrefix: t(locale, "ai.failedPrefix"),
              requestFailed: t(locale, "ai.requestFailed")
            }}
          />
        </div>
        {aiReport ? (
          <>
            <div className="rounded-md border border-emerald-200 bg-emerald-50 p-3">
              <p className="text-sm font-medium">{aiReport.nl_summary || t(locale, "profile.noSummary")}</p>
            </div>

            <div className="mt-3 rounded-md border border-slate-200 bg-slate-50 p-3">
              <h4 className="text-xs font-semibold uppercase tracking-wide text-slate-700">{t(locale, "profile.aiKeySignals")}</h4>
              <p className="mt-2 text-xs text-muted">
                {t(locale, "profile.model")} {aiReport.model_id} · {aiReport.input_tokens}/{aiReport.output_tokens} {t(locale, "profile.tokens")} ·{" "}
                {aiReport.latency_ms} ms · {t(locale, "profile.aiGeneratedAt")} {aiReport.created_at}
              </p>
            </div>

            <div className="mt-3 rounded-md border border-amber-200 bg-amber-50 p-3">
              <h4 className="text-xs font-semibold uppercase tracking-wide text-amber-800">{t(locale, "profile.aiWarnings")}</h4>
              {toRiskWarnings(aiReport).length > 0 ? (
                <ul className="mt-2 list-disc space-y-1 pl-5 text-xs text-amber-900">
                  {toRiskWarnings(aiReport).map((warning, idx) => (
                    <li key={`${warning}-${idx}`}>{warning}</li>
                  ))}
                </ul>
              ) : (
                <p className="mt-2 text-xs text-amber-900">{t(locale, "profile.aiNoWarnings")}</p>
              )}
            </div>

            <details className="mt-3 rounded-md border border-slate-200 bg-white p-3">
              <summary className="cursor-pointer text-xs font-medium text-slate-700">{t(locale, "profile.aiOpenJson")} · {t(locale, "profile.aiRawJson")}</summary>
              <pre className="mt-2 overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">{JSON.stringify(aiReport.report, null, 2)}</pre>
            </details>
          </>
        ) : (
          <p className="text-sm text-muted">{t(locale, "profile.noReport")}</p>
        )}
        {aiHistory.length > 0 && (
          <div className="mt-3 space-y-1 text-xs text-muted">
            {aiHistory.map((h) => (
              <p key={h.id}>{t(locale, "profile.history")}: {h.created_at}</p>
            ))}
          </div>
        )}
      </article>

      <article className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-xs text-amber-900">
        {profile.meta.disclosures.map((d) => (
          <p key={d}>{d}</p>
        ))}
      </article>

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
          updatedAt: t(locale, "share.updatedAt")
        }}
      />

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "profile.evidence")}</h3>
        <pre className="overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">{JSON.stringify(explanation.layer2, null, 2)}</pre>
      </article>
    </section>
  );
}
