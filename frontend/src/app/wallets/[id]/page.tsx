import { getWalletAIReport, getWalletAIReportHistory, getWalletExplanation, getWalletProfile } from "@/lib/api";
import { TriggerAnalysisButton } from "@/components/ai/TriggerAnalysisButton";

export const dynamic = "force-dynamic";

export default async function WalletProfilePage({ params }: { params: { id: string } }) {
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
        <h3 className="mb-3 text-base font-semibold">Layer1 PnL</h3>
        <div className="grid gap-3 text-sm md:grid-cols-3">
          <div>Realized PnL: {profile.layer1_facts.realized_pnl.toFixed(2)}</div>
          <div>Trading PnL: {profile.layer1_facts.trading_pnl.toFixed(2)}</div>
          <div>Maker Rebates: {profile.layer1_facts.maker_rebates.toFixed(2)}</div>
          <div>Fees Paid: {profile.layer1_facts.fees_paid.toFixed(2)}</div>
          <div>Total Trades: {profile.layer1_facts.total_trades}</div>
          <div>Volume 30D: {profile.layer1_facts.volume_30d.toFixed(2)}</div>
        </div>
      </article>

      {profile.strategy && (
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <h3 className="mb-2 text-base font-semibold">Strategy Snapshot</h3>
          <p className="text-sm">
            {profile.strategy.strategy_type} · score {profile.strategy.smart_score} · {profile.strategy.info_edge_level}
          </p>
        </article>
      )}

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">Layer3 Timing Edge</h3>
        <p className="text-sm">
          label {profile.layer3_info_edge.label} · mean Δt {profile.layer3_info_edge.mean_delta_minutes.toFixed(2)} min · samples{" "}
          {profile.layer3_info_edge.samples}
        </p>
        <p className="mt-1 text-xs text-muted">
          stddev {profile.layer3_info_edge.stddev_minutes.toFixed(2)} · p-value {profile.layer3_info_edge.p_value.toFixed(4)}
        </p>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">AI Analysis</h3>
        <div className="mb-3">
          <TriggerAnalysisButton walletID={params.id} />
        </div>
        {aiReport ? (
          <>
            <p className="text-sm">{aiReport.nl_summary || "No summary from model."}</p>
            <p className="mt-2 text-xs text-muted">
              model {aiReport.model_id} · {aiReport.input_tokens}/{aiReport.output_tokens} tokens · {aiReport.latency_ms} ms
            </p>
          </>
        ) : (
          <p className="text-sm text-muted">No AI report yet. Trigger analysis via backend API.</p>
        )}
        {aiHistory.length > 0 && (
          <div className="mt-3 space-y-1 text-xs text-muted">
            {aiHistory.map((h) => (
              <p key={h.id}>history: {h.created_at}</p>
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
        <h3 className="mb-2 text-base font-semibold">Evidence Snapshot</h3>
        <pre className="overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">{JSON.stringify(explanation.layer2, null, 2)}</pre>
      </article>
    </section>
  );
}
