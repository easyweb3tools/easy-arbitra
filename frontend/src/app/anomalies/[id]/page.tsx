import { getAnomaly } from "@/lib/api";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function AnomalyDetailPage({ params }: { params: { id: string } }) {
  const locale = await getLocaleFromCookies();
  const anomaly = await getAnomaly(params.id);

  return (
    <section className="space-y-4">
      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="text-lg font-semibold">{t(locale, "anomalies.detailTitle")} #{anomaly.id}</h2>
        <p className="text-sm text-muted">
          {t(locale, "anomalies.type")} {anomaly.alert_type} · {t(locale, "anomalies.severity")} {anomaly.severity} · {t(locale, "anomalies.wallet")} #{anomaly.wallet_id}
        </p>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "anomalies.description")}</h3>
        <p className="text-sm">{anomaly.description}</p>
      </article>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h3 className="mb-2 text-base font-semibold">{t(locale, "anomalies.evidence")}</h3>
        <pre className="overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">{JSON.stringify(anomaly.evidence, null, 2)}</pre>
      </article>
    </section>
  );
}
