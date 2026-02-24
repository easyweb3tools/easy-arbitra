import { AlertTriangle, ChevronRight } from "lucide-react";
import { getAnomaly } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

const severityConfig: Record<number, { label: string; color: string; bg: string }> = {
  1: { label: "Low", color: "text-tint-green", bg: "bg-tint-green/10" },
  2: { label: "Medium", color: "text-tint-orange", bg: "bg-tint-orange/10" },
  3: { label: "High", color: "text-tint-red", bg: "bg-tint-red/10" },
};

export default async function AnomalyDetailPage({ params }: { params: { id: string } }) {
  const locale = await getLocaleFromCookies();
  const anomaly = await getAnomaly(params.id);
  const sev = severityConfig[anomaly.severity] || severityConfig[1];

  return (
    <section className="space-y-6 animate-fade-in">
      {/* Header */}
      <div>
        <div className="flex items-center gap-3">
          <div className={`flex h-10 w-10 items-center justify-center rounded-full ${sev.bg}`}>
            <AlertTriangle className={`h-5 w-5 ${sev.color}`} />
          </div>
          <div>
            <h1 className="text-title-1 text-label-primary">
              {t(locale, "anomalies.detailTitle")} #{anomaly.id}
            </h1>
            <div className="mt-0.5 flex gap-3">
              <span className="text-footnote text-label-tertiary">
                {t(locale, "anomalies.type")} {anomaly.alert_type}
              </span>
              <span className={`text-footnote font-medium ${sev.color}`}>
                {t(locale, "anomalies.severity")} {anomaly.severity}
              </span>
              <span className="text-footnote text-label-tertiary">
                {t(locale, "anomalies.wallet")} #{anomaly.wallet_id}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Description */}
      <div>
        <SectionHeader title={t(locale, "anomalies.description")} />
        <Card>
          <p className="text-body text-label-secondary">{anomaly.description}</p>
        </Card>
      </div>

      {/* Evidence */}
      <details className="group rounded-lg border border-separator bg-surface-secondary shadow-elevation-1">
        <summary className="flex cursor-pointer items-center gap-2 px-5 py-4 text-headline text-label-primary transition-colors hover:bg-surface-tertiary">
          <ChevronRight className="h-4 w-4 transition-transform duration-200 group-open:rotate-90" />
          {t(locale, "anomalies.evidence")}
        </summary>
        <div className="border-t border-separator p-5">
          <pre className="overflow-auto rounded-md bg-surface-tertiary p-3 font-mono text-caption-1 text-label-secondary">
            {JSON.stringify(anomaly.evidence, null, 2)}
          </pre>
        </div>
      </details>
    </section>
  );
}
