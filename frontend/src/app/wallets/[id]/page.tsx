import { getWalletAIReport, getWalletAIReportHistory, getWalletExplanation, getWalletProfile } from "@/lib/api";
import { TriggerAnalysisButton } from "@/components/ai/TriggerAnalysisButton";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function WalletProfilePage({ params }: { params: { id: string } }) {
  const locale = await getLocaleFromCookies();
  const profile = await getWalletProfile(params.id);
  const explanation = await getWalletExplanation(params.id);

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
        <h2 className="text-lg font-semibold">{profile.wallet.pseudonym || profile.wallet.address}</h2>
        <p className="text-xs text-muted">{profile.wallet.address}</p>
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
            <p className="text-sm">{aiReport.nl_summary || t(locale, "profile.noSummary")}</p>
            <p className="mt-2 text-xs text-muted">
              {t(locale, "profile.model")} {aiReport.model_id} · {aiReport.input_tokens}/{aiReport.output_tokens} {t(locale, "profile.tokens")} · {aiReport.latency_ms} ms
            </p>
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

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "profile.evidence")}</h3>
        <pre className="overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">{JSON.stringify(explanation.layer2, null, 2)}</pre>
      </article>
    </section>
  );
}
