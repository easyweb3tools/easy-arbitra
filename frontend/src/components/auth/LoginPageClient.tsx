"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { login } from "@/lib/auth";
import { googleLoginUrl } from "@/lib/auth";
import { useAuth } from "@/components/auth/AuthProvider";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

export function LoginPageClient({ locale }: { locale: Locale }) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { refresh } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(email, password);
      await refresh();
      const redirect = searchParams.get("redirect") || "/";
      router.push(redirect);
    } catch (err: any) {
      setError(err.message || t(locale, "auth.loginFailed"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-sm space-y-6">
      <h1 className="text-large-title font-bold text-label-primary">
        {t(locale, "auth.login")}
      </h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="text-caption-1 font-medium text-label-tertiary uppercase tracking-wide">
            {t(locale, "auth.email")}
          </label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="mt-1 w-full h-10 rounded-lg bg-surface-tertiary border border-separator/60 px-3 text-body text-label-primary focus:outline-none focus:ring-2 focus:ring-tint-blue/40"
          />
        </div>
        <div>
          <label className="text-caption-1 font-medium text-label-tertiary uppercase tracking-wide">
            {t(locale, "auth.password")}
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="mt-1 w-full h-10 rounded-lg bg-surface-tertiary border border-separator/60 px-3 text-body text-label-primary focus:outline-none focus:ring-2 focus:ring-tint-blue/40"
          />
        </div>

        {error && (
          <p className="text-footnote text-tint-red">{error}</p>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full h-10 rounded-lg bg-tint-blue text-white font-semibold text-subheadline transition-all duration-200 hover:opacity-90 disabled:opacity-40"
        >
          {loading ? "..." : t(locale, "auth.login")}
        </button>
      </form>

      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-separator/60" />
        </div>
        <div className="relative flex justify-center text-caption-1">
          <span className="bg-surface-primary px-2 text-label-tertiary">or</span>
        </div>
      </div>

      <button
        type="button"
        onClick={() => { window.location.href = googleLoginUrl(); }}
        className="w-full h-10 rounded-lg border border-separator/60 bg-surface-secondary text-label-primary font-semibold text-subheadline transition-all duration-200 hover:bg-surface-tertiary"
      >
        {t(locale, "auth.loginWithGoogle")}
      </button>

      <p className="text-center text-footnote text-label-tertiary">
        {t(locale, "auth.noAccount")}{" "}
        <Link href="/register" className="text-tint-blue hover:underline">
          {t(locale, "auth.register")}
        </Link>
      </p>
    </div>
  );
}
