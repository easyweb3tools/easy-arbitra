import { LanguageSwitcher } from "@/components/layout/LanguageSwitcher";
import { DesktopNav } from "@/components/layout/DesktopNav";
import { t, type Locale } from "@/lib/i18n";

export function Header({ locale }: { locale: Locale }) {
  return (
    <header
      className="sticky top-0 z-50 border-b"
      style={{
        background: "var(--header-blur-bg)",
        backdropFilter: "saturate(180%) blur(20px)",
        WebkitBackdropFilter: "saturate(180%) blur(20px)",
        borderColor: "var(--separator)",
      }}
    >
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3.5 sm:px-6">
        <div className="min-w-0">
          <h1 className="text-title-2 tracking-tight text-label-primary">Easy Arbitra</h1>
          <p className="text-caption-1 text-label-tertiary">{t(locale, "app.subtitle")}</p>
        </div>
        <LanguageSwitcher
          locale={locale}
          enLabel={t(locale, "lang.en")}
          zhLabel={t(locale, "lang.zh")}
        />
      </div>
      <DesktopNav locale={locale} />
    </header>
  );
}
