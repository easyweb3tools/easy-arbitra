import Link from "next/link";

const nav = [
  { href: "/", label: "Dashboard" },
  { href: "/wallets", label: "Wallets" },
  { href: "/markets", label: "Markets" },
  { href: "/leaderboard", label: "Leaderboard" },
  { href: "/anomalies", label: "Anomalies" },
  { href: "/methodology", label: "Methodology" }
];

export function Sidebar() {
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
