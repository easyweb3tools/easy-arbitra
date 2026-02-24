"use client";

import { useRouter } from "next/navigation";
import { LOCALE_COOKIE, type Locale } from "@/lib/i18n";

export function LanguageSwitcher({
  locale,
  enLabel,
  zhLabel,
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
    <div
      className="flex items-center rounded-full p-0.5"
      style={{ background: "var(--surface-tertiary)" }}
    >
      <button
        type="button"
        onClick={() => setLocale("en")}
        className={[
          "rounded-full px-3 py-1.5 text-caption-1 font-medium",
          "transition-all duration-200 ease-apple",
          locale === "en"
            ? "bg-surface-secondary text-label-primary shadow-elevation-1"
            : "text-label-tertiary hover:text-label-secondary",
        ].join(" ")}
      >
        {enLabel}
      </button>
      <button
        type="button"
        onClick={() => setLocale("zh")}
        className={[
          "rounded-full px-3 py-1.5 text-caption-1 font-medium",
          "transition-all duration-200 ease-apple",
          locale === "zh"
            ? "bg-surface-secondary text-label-primary shadow-elevation-1"
            : "text-label-tertiary hover:text-label-secondary",
        ].join(" ")}
      >
        {zhLabel}
      </button>
    </div>
  );
}
