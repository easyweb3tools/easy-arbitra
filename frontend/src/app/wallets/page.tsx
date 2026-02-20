import Link from "next/link";
import { getWallets } from "@/lib/api";

export const dynamic = "force-dynamic";

export default async function WalletsPage({
  searchParams
}: {
  searchParams: { page?: string; page_size?: string; q?: string; tracked?: string; sort_by?: string; order?: string };
}) {
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
    sort_by: searchParams.sort_by || "updated_at",
    order: searchParams.order || "desc"
  });
  if (searchParams.q) params.set("q", searchParams.q);
  if (searchParams.tracked) params.set("tracked", searchParams.tracked);

  const wallets = await getWallets(params);

  return (
    <section className="space-y-4">
      <form className="flex flex-wrap gap-2 rounded-lg bg-card p-4 shadow-sm" method="get">
        <input name="q" defaultValue={searchParams.q} placeholder="Search pseudonym" className="rounded-md border border-slate-300 px-3 py-2 text-sm" />
        <select name="tracked" defaultValue={searchParams.tracked || ""} className="rounded-md border border-slate-300 px-3 py-2 text-sm">
          <option value="">All</option>
          <option value="true">Tracked</option>
          <option value="false">Untracked</option>
        </select>
        <button className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-white" type="submit">
          Apply
        </button>
      </form>

      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold">Wallets</h2>
        <div className="space-y-2">
          {wallets.items.map((wallet) => (
            <Link key={wallet.id} href={`/wallets/${wallet.id}`} className="block rounded-md border border-slate-200 p-3 hover:bg-slate-50">
              <p className="font-medium">{wallet.pseudonym || "Unnamed"}</p>
              <p className="text-xs text-muted">{wallet.address}</p>
            </Link>
          ))}
        </div>
      </article>
    </section>
  );
}
