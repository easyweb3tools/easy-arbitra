"use client";

import { useState } from "react";
import type { CopyTradePerformance } from "@/lib/types";
import type { Locale } from "@/lib/i18n";
import { t } from "@/lib/i18n";

export function CopyPerfChart({
  data,
  locale,
}: {
  data: CopyTradePerformance;
  locale: Locale;
}) {
  const [period] = useState<"all">("all");
  const points = data.daily_points;

  if (points.length === 0) {
    return null;
  }

  // Build cumulative PnL
  const cumulative: number[] = [];
  let sum = 0;
  for (const p of points) {
    sum += p.pnl;
    cumulative.push(sum);
  }

  const maxVal = Math.max(...cumulative.map(Math.abs), 1);
  const W = 600;
  const H = 200;
  const padX = 0;
  const padY = 20;

  const xStep = points.length > 1 ? (W - padX * 2) / (points.length - 1) : 0;

  const pathPoints = cumulative.map((v, i) => {
    const x = padX + i * xStep;
    const y = padY + ((maxVal - v) / (maxVal * 2)) * (H - padY * 2);
    return `${x},${y}`;
  });

  const linePath = `M${pathPoints.join(" L")}`;
  const lastY = padY + ((maxVal - cumulative[cumulative.length - 1]) / (maxVal * 2)) * (H - padY * 2);
  const isPositive = cumulative[cumulative.length - 1] >= 0;
  const strokeColor = isPositive ? "var(--color-tint-green)" : "var(--color-tint-red)";
  const fillColor = isPositive ? "var(--color-tint-green)" : "var(--color-tint-red)";

  const areaPath = `${linePath} L${padX + (points.length - 1) * xStep},${H} L${padX},${H} Z`;

  const zeroY = padY + ((maxVal) / (maxVal * 2)) * (H - padY * 2);

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-headline font-semibold text-label-primary">
          {t(locale, "copyTrade.performance")}
        </h3>
        <div className="flex items-center gap-2 text-footnote text-label-tertiary tabular-nums">
          <span>{t(locale, "copyTrade.totalPnl")}: <span className={isPositive ? "text-tint-green" : "text-tint-red"}>${data.total_pnl.toFixed(2)}</span></span>
          <span>{t(locale, "copyTrade.winRate")}: {(data.win_rate * 100).toFixed(0)}%</span>
        </div>
      </div>
      <svg viewBox={`0 0 ${W} ${H}`} className="w-full h-auto" preserveAspectRatio="none">
        <line x1={padX} y1={zeroY} x2={W - padX} y2={zeroY} stroke="var(--color-separator)" strokeWidth="0.5" strokeDasharray="4 4" />
        <path d={areaPath} fill={fillColor} opacity="0.08" />
        <path d={linePath} fill="none" stroke={strokeColor} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx={padX + (points.length - 1) * xStep} cy={lastY} r="3" fill={strokeColor} />
      </svg>
    </div>
  );
}
