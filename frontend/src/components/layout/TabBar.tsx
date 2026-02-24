"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, Search, Star, MoreHorizontal } from "lucide-react";
import { t, type Locale } from "@/lib/i18n";

const tabs = [
  { href: "/", icon: Home, labelKey: "nav.discover" as const },
  { href: "/wallets", icon: Search, labelKey: "nav.explore" as const },
  { href: "/watchlist", icon: Star, labelKey: "nav.tracking" as const },
  { href: "/markets", icon: MoreHorizontal, labelKey: "nav.more" as const },
] as const;

export function TabBar({ locale }: { locale: Locale }) {
  const pathname = usePathname();

  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-50 border-t md:hidden"
      style={{
        background: "var(--header-blur-bg)",
        backdropFilter: "saturate(180%) blur(20px)",
        WebkitBackdropFilter: "saturate(180%) blur(20px)",
        borderColor: "var(--separator)",
        paddingBottom: "env(safe-area-inset-bottom)",
      }}
    >
      <div className="flex items-stretch">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive =
            tab.href === "/"
              ? pathname === "/"
              : pathname.startsWith(tab.href);
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={[
                "flex flex-1 flex-col items-center gap-0.5 py-2",
                "transition-colors duration-150 ease-apple",
                isActive ? "text-tint-blue" : "text-label-tertiary",
              ].join(" ")}
            >
              <Icon className="h-6 w-6" strokeWidth={isActive ? 2.2 : 1.5} />
              <span className="text-caption-2">{t(locale, tab.labelKey)}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
