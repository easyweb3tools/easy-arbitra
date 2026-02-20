import { getLeaderboard, getOverviewStats } from "@/lib/api";

export const dynamic = "force-dynamic";

export default async function HomePage() {
  const [stats, leaderboard] = await Promise.all([
    getOverviewStats(),
    getLeaderboard(new URLSearchParams({ page: "1", page_size: "5", sort_by: "smart_score", order: "desc" }))
  ]);

  return (
    <section className="space-y-4">
      <div className="grid gap-4 md:grid-cols-3">
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <p className="text-sm text-muted">Tracked Wallets</p>
          <p className="mt-2 text-3xl font-semibold">{stats.tracked_wallets}</p>
        </article>
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <p className="text-sm text-muted">Markets Indexed</p>
          <p className="mt-2 text-3xl font-semibold">{stats.indexed_markets}</p>
        </article>
        <article className="rounded-lg bg-card p-5 shadow-sm">
          <p className="text-sm text-muted">Leaderboard Sample</p>
          <p className="mt-2 text-xl font-semibold text-accent">Top {leaderboard.items.length} loaded</p>
        </article>
      </div>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold">Top Smart Wallets</h2>
        <div className="space-y-2">
          {leaderboard.items.map((item) => (
            <div key={item.wallet_id} className="rounded-md border border-slate-200 p-3">
              <p className="font-medium">{item.pseudonym || item.address}</p>
              <p className="text-xs text-muted">
                score {item.smart_score} · {item.strategy_type}
              </p>
            </div>
          ))}
        </div>
      </article>
    </section>
  );
}
