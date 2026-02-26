"use client";

const KEY = "easy_arbitra_fingerprint";
const COOKIE_NAME = "user_fingerprint";

function randomID() {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `fp_${Math.random().toString(36).slice(2)}_${Date.now()}`;
}

function syncToCookie(value: string) {
  document.cookie = `${COOKIE_NAME}=${encodeURIComponent(value)};path=/;max-age=${60 * 60 * 24 * 365};samesite=lax`;
}

export function ensureFingerprint(): string {
  if (typeof window === "undefined") return "server";
  const current = window.localStorage.getItem(KEY);
  if (current && current.trim()) {
    syncToCookie(current);
    return current;
  }
  const created = randomID();
  window.localStorage.setItem(KEY, created);
  syncToCookie(created);
  return created;
}
