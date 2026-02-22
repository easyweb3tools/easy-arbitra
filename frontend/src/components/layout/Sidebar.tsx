import Link from "next/link";
import { t, type Locale } from "@/lib/i18n";

export function Sidebar({ locale }: { locale: Locale }) {
  const nav = [
    { href: "/", label: t(locale, "nav.dashboard") },
    { href: "/wallets", label: t(locale, "nav.wallets") },
    { href: "/markets", label: t(locale, "nav.markets") },
    { href: "/leaderboard", label: t(locale, "nav.leaderboard") },
    { href: "/anomalies", label: t(locale, "nav.anomalies") },
    { href: "/methodology", label: t(locale, "nav.methodology") }
  ];
  return (
    <aside className="w-full border-b border-slate-200 bg-white px-4 py-3 md:w-64 md:border-b-0 md:border-r md:px-5 md:py-6">
      <nav className="flex gap-3 md:flex-col">
        {nav.map((item) => (
          <Link key={item.href} href={item.href} className="rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100">
            {item.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
