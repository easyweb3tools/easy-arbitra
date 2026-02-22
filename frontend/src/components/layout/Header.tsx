import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { t, type Locale } from "@/lib/i18n";

export function Header({ locale }: { locale: Locale }) {
  return (
    <header className="flex items-start justify-between gap-3 border-b border-slate-200 bg-white/90 px-6 py-4 backdrop-blur">
      <div>
        <h1 className="text-lg font-semibold text-ink">Easy Arbitra</h1>
        <p className="text-sm text-muted">{t(locale, "app.subtitle")}</p>
      </div>
      <LanguageSwitcher locale={locale} enLabel={t(locale, "lang.en")} zhLabel={t(locale, "lang.zh")} />
    </header>
  );
}
