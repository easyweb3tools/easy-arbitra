import { getAnomalies } from "@/lib/api";
import { AcknowledgeButton } from "@/components/anomaly/AcknowledgeButton";
import Link from "next/link";

export const dynamic = "force-dynamic";

export default async function AnomaliesPage({
  searchParams
}: {
  searchParams: { page?: string; page_size?: string; severity?: string; type?: string; acknowledged?: string };
}) {
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20"
  });
  if (searchParams.severity) params.set("severity", searchParams.severity);
  if (searchParams.type) params.set("type", searchParams.type);
  if (searchParams.acknowledged) params.set("acknowledged", searchParams.acknowledged);

  const feed = await getAnomalies(params);

  return (
    <section className="space-y-4">
      <form className="flex flex-wrap gap-2 rounded-lg bg-card p-4 shadow-sm" method="get">
        <select name="severity" defaultValue={searchParams.severity || ""} className="rounded-md border border-slate-300 px-3 py-2 text-sm">
          <option value="">All severity</option>
          <option value="1">1 - Low</option>
          <option value="2">2 - Medium</option>
          <option value="3">3 - High</option>
        </select>
        <input name="type" defaultValue={searchParams.type} placeholder="Type" className="rounded-md border border-slate-300 px-3 py-2 text-sm" />
        <button className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-white" type="submit">
          Apply
        </button>
      </form>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold">Anomaly Feed</h2>
        <div className="space-y-2">
          {feed.items.map((alert) => (
            <div key={alert.id} className="rounded-md border border-slate-200 p-3">
              <Link href={`/anomalies/${alert.id}`} className="font-medium text-slate-900 underline-offset-2 hover:underline">
                {alert.alert_type}
              </Link>
              <p className="text-xs text-muted">wallet #{alert.wallet_id} · severity {alert.severity}</p>
              <p className="text-sm">{alert.description}</p>
              <pre className="mt-2 overflow-auto rounded bg-slate-100 p-2 text-xs text-slate-700">
                {JSON.stringify(alert.evidence, null, 2)}
              </pre>
              {!alert.acknowledged && (
                <div className="mt-2">
                  <AcknowledgeButton alertID={alert.id} />
                </div>
              )}
            </div>
          ))}
        </div>
      </article>
    </section>
  );
}
