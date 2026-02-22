"use client";

const KEY = "easy_arbitra_fingerprint";

function randomID() {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `fp_${Math.random().toString(36).slice(2)}_${Date.now()}`;
}

export function ensureFingerprint(): string {
  if (typeof window === "undefined") return "server";
  const current = window.localStorage.getItem(KEY);
  if (current && current.trim()) return current;
  const created = randomID();
  window.localStorage.setItem(KEY, created);
  return created;
}
