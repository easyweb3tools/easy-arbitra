import { RegisterPageClient } from "@/components/auth/RegisterPageClient";
import { getLocaleFromCookies } from "@/lib/i18n-server";

export const dynamic = "force-dynamic";

export default async function RegisterPage() {
  const locale = await getLocaleFromCookies();
  return <RegisterPageClient locale={locale} />;
}
