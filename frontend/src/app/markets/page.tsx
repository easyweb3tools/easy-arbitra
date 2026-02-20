import { getMarkets } from "@/lib/api";

export const dynamic = "force-dynamic";

export default async function MarketsPage({
  searchParams
}: {
  searchParams: { page?: string; page_size?: string; category?: string; sort_by?: string; order?: string };
}) {
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
    sort_by: searchParams.sort_by || "updated_at",
    order: searchParams.order || "desc"
  });
  if (searchParams.category) params.set("category", searchParams.category);

  const markets = await getMarkets(params);

  return (
    <section className="space-y-4">
      <form className="flex flex-wrap gap-2 rounded-lg bg-card p-4 shadow-sm" method="get">
        <input name="category" defaultValue={searchParams.category} placeholder="Category e.g. Politics" className="rounded-md border border-slate-300 px-3 py-2 text-sm" />
        <button className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-white" type="submit">
          Apply
        </button>
      </form>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold">Markets</h2>
        <div className="space-y-2">
          {markets.items.map((market) => (
            <div key={market.id} className="rounded-md border border-slate-200 p-3">
              <p className="font-medium">{market.title}</p>
              <p className="text-xs text-muted">
                {market.category} · vol {market.volume} · liq {market.liquidity}
              </p>
            </div>
          ))}
        </div>
      </article>
    </section>
  );
}
