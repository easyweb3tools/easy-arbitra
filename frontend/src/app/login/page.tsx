import { Suspense } from "react";
import { LoginPageClient } from "@/components/auth/LoginPageClient";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function LoginPage() {
  const locale = await getLocaleFromCookies();
  return (
    <Suspense>
      <LoginPageClient locale={locale} />
    </Suspense>
  );
}
