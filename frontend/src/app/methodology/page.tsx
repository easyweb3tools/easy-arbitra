import { Info, AlertTriangle } from "lucide-react";
import { Card, SectionHeader } from "@/components/ui/Card";
import { t } from "@/lib/i18n";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function MethodologyPage() {
  const locale = await getLocaleFromCookies();
  return (
    <section className="space-y-6 animate-fade-in">
      <SectionHeader title={t(locale, "method.title")} />
      <Card variant="prominent">
        <div className="flex gap-3">
          <Info className="mt-0.5 h-5 w-5 shrink-0 text-tint-blue" />
          <p className="text-body text-label-secondary">{t(locale, "method.body")}</p>
        </div>
      </Card>

      <div className="rounded-lg bg-tint-orange/[0.06] p-5">
        <div className="flex gap-3">
          <AlertTriangle className="mt-0.5 h-5 w-5 shrink-0 text-tint-orange" />
          <div className="space-y-1.5">
            <p className="text-footnote text-tint-orange">{t(locale, "method.d1")}</p>
            <p className="text-footnote text-tint-orange">{t(locale, "method.d2")}</p>
            <p className="text-footnote text-tint-orange">{t(locale, "method.d3")}</p>
          </div>
        </div>
      </div>
    </section>
  );
}
