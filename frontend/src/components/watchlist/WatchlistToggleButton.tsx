"use client";

import { useEffect, useMemo, useState } from "react";
import { addToWatchlist, removeFromWatchlist } from "@/lib/api";
import { ensureFingerprint } from "@/lib/fingerprint";

type Labels = {
  follow: string;
  unfollow: string;
  following: string;
  failed: string;
};

const IDS_KEY = "easy_arbitra_watchlist_ids";

function readIDs(): number[] {
  if (typeof window === "undefined") return [];
  const raw = window.localStorage.getItem(IDS_KEY);
  if (!raw) return [];
  try {
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed.map((v) => Number(v)).filter((v) => Number.isFinite(v) && v > 0);
  } catch {
    return [];
  }
}

function writeIDs(ids: number[]) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(IDS_KEY, JSON.stringify(Array.from(new Set(ids)).slice(0, 500)));
}

export function WatchlistToggleButton({ walletID, labels }: { walletID: number; labels: Labels }) {
  const [watching, setWatching] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const fingerprint = useMemo(() => ensureFingerprint(), []);

  useEffect(() => {
    const ids = readIDs();
    setWatching(ids.includes(walletID));
  }, [walletID]);

  async function onToggle() {
    setLoading(true);
    setError("");
    try {
      if (watching) {
        await removeFromWatchlist(walletID, fingerprint);
        const ids = readIDs().filter((id) => id !== walletID);
        writeIDs(ids);
        setWatching(false);
      } else {
        await addToWatchlist(walletID, fingerprint);
        const ids = readIDs();
        ids.push(walletID);
        writeIDs(ids);
        setWatching(true);
      }
    } catch {
      setError(labels.failed);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-1">
      <button
        type="button"
        onClick={onToggle}
        disabled={loading}
        className={`rounded-md px-3 py-1 text-xs font-medium text-white disabled:opacity-60 ${watching ? "bg-slate-600" : "bg-emerald-600"}`}
      >
        {loading ? labels.following : watching ? labels.unfollow : labels.follow}
      </button>
      {error ? <p className="text-xs text-rose-700">{error}</p> : null}
    </div>
  );
}
