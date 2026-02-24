import Link from "next/link";
import { AlertTriangle, ChevronRight, Filter } from "lucide-react";
import { getAnomalies } from "@/lib/api";
import { AcknowledgeButton } from "@/components/anomaly/AcknowledgeButton";
import { Card, SectionHeader } from "@/components/ui/Card";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

const severityConfig: Record<number, { color: string; bg: string }> = {
  1: { color: "text-tint-green", bg: "bg-tint-green/10" },
  2: { color: "text-tint-orange", bg: "bg-tint-orange/10" },
  3: { color: "text-tint-red", bg: "bg-tint-red/10" },
};

export default async function AnomaliesPage({
  searchParams,
}: {
  searchParams: { page?: string; page_size?: string; severity?: string; type?: string; acknowledged?: string };
}) {
  const locale = await getLocaleFromCookies();
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
  });
  if (searchParams.severity) params.set("severity", searchParams.severity);
  if (searchParams.type) params.set("type", searchParams.type);
  if (searchParams.acknowledged) params.set("acknowledged", searchParams.acknowledged);
  const feed = await getAnomalies(params);

  return (
    <section className="space-y-6">
      {/* Filter */}
      <Card variant="flat">
        <form className="flex flex-wrap items-center gap-3 p-4" method="get">
          <Filter className="h-4 w-4 text-label-tertiary" />
          <select
            name="severity"
            defaultValue={searchParams.severity || ""}
            className="h-9 rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1 focus:ring-2 focus:ring-tint-blue"
          >
            <option value="">{t(locale, "anomalies.allSeverity")}</option>
            <option value="1">1 — {t(locale, "anomalies.low")}</option>
            <option value="2">2 — {t(locale, "anomalies.medium")}</option>
            <option value="3">3 — {t(locale, "anomalies.high")}</option>
          </select>
          <input
            name="type"
            defaultValue={searchParams.type}
            placeholder={t(locale, "anomalies.typePlaceholder")}
            className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1 placeholder:text-label-quaternary focus:ring-2 focus:ring-tint-blue sm:w-48"
          />
          <button
            className="inline-flex h-9 items-center rounded-md bg-tint-blue px-5 text-subheadline font-semibold text-white transition-all duration-200 ease-apple hover:brightness-110 active:scale-[0.98]"
            type="submit"
          >
            {t(locale, "anomalies.apply")}
          </button>
        </form>
      </Card>

      {/* Feed */}
      <div>
        <SectionHeader title={t(locale, "anomalies.feedTitle")} />
        <div className="space-y-3">
          {feed.items.map((alert) => {
            const sev = severityConfig[alert.severity] || severityConfig[1];
            return (
              <Card key={alert.id}>
                <div className="flex items-start gap-3">
                  <div className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full ${sev.bg}`}>
                    <AlertTriangle className={`h-4 w-4 ${sev.color}`} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <Link
                      href={`/anomalies/${alert.id}`}
                      className="text-headline text-label-primary transition-colors hover:text-tint-blue"
                    >
                      {alert.alert_type}
                    </Link>
                    <div className="mt-0.5 flex gap-3">
                      <span className="text-footnote text-label-tertiary">
                        {t(locale, "anomalies.wallet")} #{alert.wallet_id}
                      </span>
                      <span className={`text-footnote font-medium ${sev.color}`}>
                        {t(locale, "anomalies.severity")} {alert.severity}
                      </span>
                    </div>
                    <p className="mt-2 text-callout text-label-secondary">{alert.description}</p>

                    {/* Collapsible evidence */}
                    <details className="group mt-3 rounded-md border border-separator">
                      <summary className="flex cursor-pointer items-center gap-1.5 px-3 py-2 text-caption-1 font-medium text-label-tertiary transition-colors hover:bg-surface-tertiary">
                        <ChevronRight className="h-3 w-3 transition-transform duration-200 group-open:rotate-90" />
                        {t(locale, "anomalies.evidence")}
                      </summary>
                      <pre className="overflow-auto border-t border-separator p-3 font-mono text-caption-1 text-label-tertiary">
                        {JSON.stringify(alert.evidence, null, 2)}
                      </pre>
                    </details>

                    {!alert.acknowledged && (
                      <div className="mt-3">
                        <AcknowledgeButton
                          alertID={alert.id}
                          labels={{ ack: t(locale, "anomalies.ack"), acking: t(locale, "anomalies.acking") }}
                        />
                      </div>
                    )}
                  </div>
                </div>
              </Card>
            );
          })}
        </div>
      </div>
    </section>
  );
}
