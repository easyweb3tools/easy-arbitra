import { getLeaderboard } from "@/lib/api";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function LeaderboardPage({
  searchParams
}: {
  searchParams: { page?: string; page_size?: string; sort_by?: string; order?: string };
}) {
  const locale = await getLocaleFromCookies();
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
    sort_by: searchParams.sort_by || "smart_score",
    order: searchParams.order || "desc"
  });
  const board = await getLeaderboard(params);

  return (
    <section className="rounded-lg bg-card p-5 shadow-sm">
      <h2 className="mb-4 text-lg font-semibold">{t(locale, "leaderboard.title")}</h2>
      <div className="space-y-2">
        {board.items.map((item, idx) => (
          <div key={`${item.wallet_id}-${idx}`} className="rounded-md border border-slate-200 p-3">
            <p className="font-medium">
              #{idx + 1} {item.pseudonym || item.address}
            </p>
            <p className="text-xs text-muted">
              {t(locale, "leaderboard.score")} {item.smart_score} · {item.strategy_type} · {item.info_edge_level}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}
