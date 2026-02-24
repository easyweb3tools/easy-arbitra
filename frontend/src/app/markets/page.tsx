import { Search } from "lucide-react";
import { getMarkets } from "@/lib/api";
import { Card, SectionHeader } from "@/components/ui/Card";
import { CategoryTag } from "@/components/ui/Badge";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function MarketsPage({
  searchParams,
}: {
  searchParams: { page?: string; page_size?: string; category?: string; sort_by?: string; order?: string };
}) {
  const locale = await getLocaleFromCookies();
  const params = new URLSearchParams({
    page: searchParams.page || "1",
    page_size: searchParams.page_size || "20",
    sort_by: searchParams.sort_by || "updated_at",
    order: searchParams.order || "desc",
  });
  if (searchParams.category) params.set("category", searchParams.category);
  const markets = await getMarkets(params);

  return (
    <section className="space-y-6">
      {/* Search / Filter */}
      <Card variant="flat">
        <form className="flex flex-wrap items-center gap-3 p-4" method="get">
          <Search className="h-4 w-4 text-label-tertiary" />
          <input
            name="category"
            defaultValue={searchParams.category}
            placeholder={t(locale, "markets.categoryPlaceholder")}
            className="h-9 w-full rounded-md border-0 bg-surface-secondary px-3 text-subheadline text-label-primary shadow-elevation-1 placeholder:text-label-quaternary focus:ring-2 focus:ring-tint-blue sm:w-56"
          />
          <button
            className="inline-flex h-9 items-center rounded-md bg-tint-blue px-5 text-subheadline font-semibold text-white transition-all duration-200 ease-apple hover:brightness-110 active:scale-[0.98]"
            type="submit"
          >
            {t(locale, "markets.apply")}
          </button>
        </form>
      </Card>

      {/* Market List */}
      <div>
        <SectionHeader title={t(locale, "markets.title")} />
        <Card padding={false}>
          {markets.items.map((market) => (
            <div
              key={market.id}
              className="border-b border-separator px-4 py-3 last:border-b-0"
            >
              <p className="text-headline text-label-primary">{market.title}</p>
              <div className="mt-1 flex flex-wrap items-center gap-2">
                <CategoryTag>{market.category || "—"}</CategoryTag>
                <span className="text-footnote text-label-tertiary">
                  {t(locale, "markets.vol")}{" "}
                  <span className="font-medium text-label-secondary">{market.volume}</span>
                </span>
                <span className="text-footnote text-label-tertiary">
                  {t(locale, "markets.liq")}{" "}
                  <span className="font-medium text-label-secondary">{market.liquidity}</span>
                </span>
              </div>
            </div>
          ))}
        </Card>
      </div>
    </section>
  );
}
