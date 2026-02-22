import type { Metadata } from "next";
import "./globals.css";
import { Header } from "@/components/layout/Header";
import { Sidebar } from "@/components/layout/Sidebar";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const metadata: Metadata = {
  title: "Easy Arbitra",
  description: "Polymarket smart wallet analyzer"
};

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const locale = await getLocaleFromCookies();
  return (
    <html lang={locale}>
      <body>
        <Header locale={locale} />
        <div className="flex min-h-[calc(100vh-73px)] flex-col md:flex-row">
          <Sidebar locale={locale} />
          <main className="flex-1 px-6 py-6">{children}</main>
        </div>
      </body>
    </html>
  );
}
