import { cookies } from "next/headers";
import Link from "next/link";
import { LOCALE_COOKIE, normalizeLocale, t } from "@/lib/i18n";
import { getCopyTradeMonitor } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import type { CopyTradeMonitor } from "@/lib/types";

export const dynamic = "force-dynamic";

function StatusBadge({ status }: { status: string }) {
  const color =
    status === "done"
      ? "bg-tint-green/15 text-tint-green"
      : status === "running"
        ? "bg-tint-blue/15 text-tint-blue"
        : "bg-tint-red/15 text-tint-red";
  return (
    <span className={`inline-block rounded-full px-2 py-0.5 text-caption-1 font-medium ${color}`}>
      {status}
    </span>
  );
}

function formatDuration(start: string, end?: string): string {
  if (!end) return "—";
  const ms = new Date(end).getTime() - new Date(start).getTime();
  if (ms < 1000) return `${ms}ms`;
  return `${(ms / 1000).toFixed(1)}s`;
}

function formatHour(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleString("en-US", { month: "short", day: "numeric", hour: "2-digit", minute: "2-digit", hour12: false });
}

export default async function MonitorPage() {
  const store = await cookies();
  const locale = normalizeLocale(store.get(LOCALE_COOKIE)?.value);

  let data: CopyTradeMonitor | null = null;
  try {
    data = await getCopyTradeMonitor();
  } catch {
    // will show empty state
  }

  const totalRuns24h = data?.hourly_stats?.reduce((sum, h) => sum + h.runs, 0) ?? 0;

  return (
    <section className="space-y-8 animate-fade-in">
      <h1 className="text-large-title font-bold text-label-primary">
        {t(locale, "monitor.title")}
      </h1>

      {/* Overview strip */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <Card>
          <p className="text-caption-1 text-label-tertiary">{t(locale, "monitor.enabledConfigs")}</p>
          <p className="text-title-1 font-bold text-label-primary">{data?.enabled_configs ?? 0}</p>
        </Card>
        <Card>
          <p className="text-caption-1 text-label-tertiary">{t(locale, "monitor.totalRuns24h")}</p>
          <p className="text-title-1 font-bold text-label-primary">{totalRuns24h}</p>
        </Card>
      </div>

      {/* Hourly Stats */}
      <div>
        <SectionHeader title={t(locale, "monitor.hourlyStats")} />
        <Card padding={false}>
          {data?.hourly_stats && data.hourly_stats.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-subheadline">
                <thead>
                  <tr className="border-b border-separator/60 text-left text-caption-1 text-label-tertiary">
                    <th className="px-4 py-2.5">{t(locale, "monitor.hour")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.runs")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.walletsChecked")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.newTrades")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.copies")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.skips")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.errors")}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.hourly_stats.map((h) => (
                    <tr key={h.hour} className="border-b border-separator/40 last:border-b-0">
                      <td className="px-4 py-2.5 text-label-secondary">{formatHour(h.hour)}</td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{h.runs}</td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{h.wallets_checked}</td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{h.new_trades}</td>
                      <td className="px-4 py-2.5 text-right text-tint-green">{h.decisions_copy}</td>
                      <td className="px-4 py-2.5 text-right text-label-tertiary">{h.decisions_skip}</td>
                      <td className="px-4 py-2.5 text-right text-tint-red">{h.errors}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="py-8 text-center text-subheadline text-label-tertiary">
              {t(locale, "monitor.noHourly")}
            </p>
          )}
        </Card>
      </div>

      {/* Recent Syncer Runs */}
      <div>
        <SectionHeader title={t(locale, "monitor.recentRuns")} />
        <Card padding={false}>
          {data?.recent_runs && data.recent_runs.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-subheadline">
                <thead>
                  <tr className="border-b border-separator/60 text-left text-caption-1 text-label-tertiary">
                    <th className="px-4 py-2.5">{t(locale, "monitor.status")}</th>
                    <th className="px-4 py-2.5">{t(locale, "monitor.hour")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.duration")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.walletsChecked")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.newTrades")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.copies")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.errors")}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.recent_runs.slice(0, 10).map((run) => (
                    <tr key={run.id} className="border-b border-separator/40 last:border-b-0">
                      <td className="px-4 py-2.5"><StatusBadge status={run.status} /></td>
                      <td className="px-4 py-2.5 text-label-secondary">{formatHour(run.started_at)}</td>
                      <td className="px-4 py-2.5 text-right text-label-primary">
                        {formatDuration(run.started_at, run.ended_at)}
                      </td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{run.stats?.wallets_checked ?? 0}</td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{run.stats?.new_trades ?? 0}</td>
                      <td className="px-4 py-2.5 text-right text-tint-green">{run.stats?.decisions_copy ?? 0}</td>
                      <td className="px-4 py-2.5 text-right text-tint-red">{run.stats?.errors ?? 0}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="py-8 text-center text-subheadline text-label-tertiary">
              {t(locale, "monitor.noRuns")}
            </p>
          )}
        </Card>
      </div>

      {/* Copyable Wallets */}
      <div>
        <SectionHeader title={t(locale, "monitor.copyableWallets")} />
        <Card padding={false}>
          {data?.copyable_wallets && data.copyable_wallets.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-subheadline">
                <thead>
                  <tr className="border-b border-separator/60 text-left text-caption-1 text-label-tertiary">
                    <th className="px-4 py-2.5">{t(locale, "leaderboard.name")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.score")}</th>
                    <th className="px-4 py-2.5">{t(locale, "monitor.tier")}</th>
                    <th className="px-4 py-2.5">{t(locale, "monitor.strategy")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.pnl30d")}</th>
                    <th className="px-4 py-2.5 text-right">{t(locale, "monitor.trades30d")}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.copyable_wallets.map((w) => (
                    <tr key={w.id} className="border-b border-separator/40 last:border-b-0 hover:bg-surface-tertiary/50">
                      <td className="px-4 py-2.5">
                        <Link href={`/wallets/${w.id}`} className="text-tint-blue hover:underline">
                          {w.pseudonym || `0x${w.address.slice(0, 8)}…`}
                        </Link>
                      </td>
                      <td className="px-4 py-2.5 text-right font-medium text-label-primary">{w.smart_score}</td>
                      <td className="px-4 py-2.5 text-label-secondary">{w.pool_tier}</td>
                      <td className="px-4 py-2.5 text-label-secondary">{w.strategy_type}</td>
                      <td className="px-4 py-2.5 text-right font-medium text-tint-green">
                        ${w.pnl_30d.toFixed(2)}
                      </td>
                      <td className="px-4 py-2.5 text-right text-label-primary">{w.trade_count_30d}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="py-8 text-center text-subheadline text-label-tertiary">
              {t(locale, "monitor.noCopyable")}
            </p>
          )}
        </Card>
      </div>
    </section>
  );
}
