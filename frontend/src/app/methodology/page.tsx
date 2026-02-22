import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function MethodologyPage() {
  const locale = await getLocaleFromCookies();
  return (
    <section className="space-y-4">
      <article className="rounded-lg bg-card p-5 shadow-sm">
        <h2 className="mb-3 text-lg font-semibold">{t(locale, "method.title")}</h2>
        <p className="text-sm text-slate-700">{t(locale, "method.body")}</p>
      </article>
      <article className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-xs text-amber-900">
        <p>{t(locale, "method.d1")}</p>
        <p>{t(locale, "method.d2")}</p>
        <p>{t(locale, "method.d3")}</p>
      </article>
    </section>
  );
}
