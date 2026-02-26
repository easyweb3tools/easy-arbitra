"use client";

import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { enableCopyTrading, disableCopyTrading } from "@/lib/api";
import { ensureFingerprint } from "@/lib/fingerprint";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

export function CopyTradeToggle({
  walletID,
  enabled,
  locale,
  onToggle,
}: {
  walletID: number;
  enabled: boolean;
  locale: Locale;
  onToggle?: () => void;
}) {
  const [showSettings, setShowSettings] = useState(false);
  const [loading, setLoading] = useState(false);
  const [maxPosition, setMaxPosition] = useState("1000");
  const [riskPref, setRiskPref] = useState<"conservative" | "moderate" | "aggressive">("moderate");

  async function handleEnable() {
    setLoading(true);
    try {
      const fp = ensureFingerprint();
      await enableCopyTrading(walletID, parseFloat(maxPosition) || 1000, riskPref, fp);
      onToggle?.();
      setShowSettings(false);
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }

  async function handleDisable() {
    setLoading(true);
    try {
      const fp = ensureFingerprint();
      await disableCopyTrading(walletID, fp);
      onToggle?.();
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }

  if (enabled) {
    return (
      <Button variant="destructive" size="small" loading={loading} onClick={handleDisable}>
        {t(locale, "copyTrade.disable")}
      </Button>
    );
  }

  if (showSettings) {
    return (
      <div className="rounded-xl bg-surface-secondary border border-separator/60 p-4 space-y-4 animate-fade-in">
        <h3 className="text-headline font-semibold text-label-primary">{t(locale, "copyTrade.enable")}</h3>

        <div className="space-y-2">
          <label className="text-caption-1 font-medium text-label-tertiary uppercase tracking-wide">
            {t(locale, "copyTrade.maxPosition")}
          </label>
          <input
            type="number"
            value={maxPosition}
            onChange={(e) => setMaxPosition(e.target.value)}
            className="w-full h-10 rounded-lg bg-surface-tertiary border border-separator/60 px-3 text-body text-label-primary tabular-nums focus:outline-none focus:ring-2 focus:ring-tint-blue/40"
            min={1}
            max={100000}
          />
        </div>

        <div className="space-y-2">
          <label className="text-caption-1 font-medium text-label-tertiary uppercase tracking-wide">
            {t(locale, "copyTrade.riskPreference")}
          </label>
          <div className="flex gap-2">
            {(["conservative", "moderate", "aggressive"] as const).map((pref) => (
              <button
                key={pref}
                onClick={() => setRiskPref(pref)}
                className={[
                  "h-9 px-4 rounded-lg text-subheadline font-medium transition-all duration-200",
                  riskPref === pref
                    ? "bg-tint-blue text-white"
                    : "bg-surface-tertiary text-label-secondary hover:bg-surface-tertiary/80",
                ].join(" ")}
              >
                {t(locale, `copyTrade.${pref}`)}
              </button>
            ))}
          </div>
        </div>

        <div className="text-footnote text-label-tertiary space-y-1">
          <p className="font-medium">{t(locale, "copyTrade.agentWill")}</p>
          <ul className="list-disc list-inside space-y-0.5">
            <li>{t(locale, "copyTrade.agentAnalyze")}</li>
            <li>{t(locale, "copyTrade.agentDecide")}</li>
            <li>{t(locale, "copyTrade.agentExplain")}</li>
          </ul>
        </div>

        <div className="flex gap-2">
          <Button variant="filled" size="small" loading={loading} onClick={handleEnable}>
            {t(locale, "copyTrade.confirm")}
          </Button>
          <Button variant="gray" size="small" onClick={() => setShowSettings(false)}>
            Cancel
          </Button>
        </div>
      </div>
    );
  }

  return (
    <Button variant="tinted" size="small" onClick={() => setShowSettings(true)}>
      {t(locale, "copyTrade.enable")}
    </Button>
  );
}
