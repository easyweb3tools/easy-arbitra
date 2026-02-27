import type { Metadata } from "next";
import "./globals.css";
import { Header } from "@/components/layout/Header";
import { TabBar } from "@/components/layout/TabBar";
import { AuthProvider } from "@/components/auth/AuthProvider";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const metadata: Metadata = {
  title: "Easy Arbitra",
  description: "Polymarket smart wallet analyzer",
};

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const locale = await getLocaleFromCookies();
  return (
    <html lang={locale} suppressHydrationWarning>
      <head>
        <meta name="color-scheme" content="light dark" />
        <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover" />
      </head>
      <body>
        <AuthProvider>
          <Header locale={locale} />
          <main className="mx-auto max-w-5xl px-4 pb-28 pt-4 sm:px-6 sm:pt-6 md:pb-10">
            {children}
          </main>
          <TabBar locale={locale} />
        </AuthProvider>
      </body>
    </html>
  );
}
