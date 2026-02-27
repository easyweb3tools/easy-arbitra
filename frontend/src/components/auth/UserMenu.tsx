"use client";

import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/components/auth/AuthProvider";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

export function UserMenu({ locale }: { locale: Locale }) {
  const { user, loading, logout } = useAuth();
  const router = useRouter();

  if (loading) return null;

  if (!user) {
    return (
      <Link
        href="/login"
        className="inline-flex h-8 items-center rounded-lg bg-tint-blue/[0.12] px-3 text-caption-1 font-semibold text-tint-blue transition-all hover:bg-tint-blue/[0.15]"
      >
        {t(locale, "auth.login")}
      </Link>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <span className="text-caption-1 text-label-secondary truncate max-w-[120px]">
        {user.name || user.email}
      </span>
      <button
        type="button"
        onClick={async () => {
          await logout();
          router.push("/");
        }}
        className="inline-flex h-8 items-center rounded-lg bg-surface-tertiary px-3 text-caption-1 font-medium text-label-secondary transition-all hover:bg-surface-tertiary/80"
      >
        {t(locale, "auth.logout")}
      </button>
    </div>
  );
}
