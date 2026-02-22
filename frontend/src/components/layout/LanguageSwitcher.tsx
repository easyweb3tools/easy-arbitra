"use client";

import { useRouter } from "next/navigation";
import { LOCALE_COOKIE, type Locale } from "@/lib/i18n";

export function LanguageSwitcher({
  locale,
  enLabel,
  zhLabel
}: {
  locale: Locale;
  enLabel: string;
  zhLabel: string;
}) {
  const router = useRouter();

  function setLocale(next: Locale) {
    document.cookie = `${LOCALE_COOKIE}=${next}; path=/; max-age=31536000; samesite=lax`;
    router.refresh();
  }

  return (
    <div className="flex items-center gap-1 rounded-md border border-slate-200 bg-white p-1">
      <button
        type="button"
        onClick={() => setLocale("en")}
        className={`rounded px-2 py-1 text-xs font-medium ${locale === "en" ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"}`}
      >
        {enLabel}
      </button>
      <button
        type="button"
        onClick={() => setLocale("zh")}
        className={`rounded px-2 py-1 text-xs font-medium ${locale === "zh" ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"}`}
      >
        {zhLabel}
      </button>
    </div>
  );
}
