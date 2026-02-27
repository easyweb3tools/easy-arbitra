"use client";

import { useEffect, useState } from "react";
import { Star, Loader2 } from "lucide-react";
import { addToWatchlist, removeFromWatchlist } from "@/lib/api";

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

  useEffect(() => {
    const ids = readIDs();
    setWatching(ids.includes(walletID));
  }, [walletID]);

  async function onToggle() {
    setLoading(true);
    setError("");
    try {
      if (watching) {
        await removeFromWatchlist(walletID);
        const ids = readIDs().filter((id) => id !== walletID);
        writeIDs(ids);
        setWatching(false);
      } else {
        await addToWatchlist(walletID);
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
    <div className="flex flex-col items-end gap-1">
      <button
        type="button"
        onClick={onToggle}
        disabled={loading}
        aria-label={watching ? labels.unfollow : labels.follow}
        className={[
          "inline-flex h-9 items-center gap-1.5 rounded-md px-4 text-subheadline font-semibold",
          "transition-all duration-200 ease-apple",
          "disabled:opacity-35 disabled:pointer-events-none",
          "active:scale-[0.97]",
          watching
            ? "bg-surface-tertiary text-label-primary"
            : "bg-tint-blue/[0.12] text-tint-blue hover:bg-tint-blue/[0.15]",
        ].join(" ")}
      >
        {loading ? (
          <Loader2 className="h-4 w-4 animate-spin" />
        ) : (
          <Star className={`h-4 w-4 ${watching ? "fill-current" : ""}`} />
        )}
        {loading ? labels.following : watching ? labels.unfollow : labels.follow}
      </button>
      {error && <p className="text-caption-1 text-tint-red">{error}</p>}
    </div>
  );
}
