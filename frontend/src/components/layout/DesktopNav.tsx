"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { t, type Locale } from "@/lib/i18n";

const links = [
  { href: "/", labelKey: "nav.discover" as const },
  { href: "/wallets", labelKey: "nav.explore" as const },
  { href: "/watchlist", labelKey: "nav.tracking" as const },
  { href: "/copy-trading", labelKey: "copyTrade.title" as const },
  { href: "/monitor", labelKey: "monitor.title" as const },
  { href: "/markets", labelKey: "nav.more" as const },
] as const;

export function DesktopNav({ locale }: { locale: Locale }) {
  const pathname = usePathname();

  return (
    <nav className="mx-auto hidden max-w-5xl gap-0.5 overflow-x-auto px-4 sm:px-6 md:flex">
      {links.map((link) => {
        const isActive =
          link.href === "/"
            ? pathname === "/"
            : pathname.startsWith(link.href);
        return (
          <Link
            key={link.href}
            href={link.href}
            className={[
              "relative whitespace-nowrap px-4 py-2.5 text-subheadline",
              "transition-colors duration-200 ease-apple",
              isActive
                ? "text-tint-blue font-semibold"
                : "text-label-tertiary hover:text-label-primary",
            ].join(" ")}
          >
            {t(locale, link.labelKey)}
            {isActive && (
              <span className="absolute bottom-0 left-4 right-4 h-[2.5px] rounded-full bg-tint-blue" />
            )}
          </Link>
        );
      })}
    </nav>
  );
}
