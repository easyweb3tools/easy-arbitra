"use client";

import { useState, useMemo } from "react";
import type { PnLHistoryPoint } from "@/lib/types";

type Period = "7d" | "30d" | "90d";

const periodLabels: Record<Period, string> = { "7d": "7D", "30d": "30D", "90d": "90D" };

function getField(point: PnLHistoryPoint, period: Period): number {
  if (period === "7d") return point.pnl_7d;
  if (period === "30d") return point.pnl_30d;
  return point.pnl_90d;
}

export function PnLChart({
  data,
  labels,
}: {
  data: PnLHistoryPoint[];
  labels: { title: string; d7: string; d30: string; d90: string };
}) {
  const [period, setPeriod] = useState<Period>("30d");

  const sorted = useMemo(() => [...data].reverse(), [data]);

  const values = useMemo(() => sorted.map((p) => getField(p, period)), [sorted, period]);

  const width = 600;
  const height = 200;
  const padX = 0;
  const padY = 16;

  const minVal = Math.min(0, ...values);
  const maxVal = Math.max(0, ...values);
  const range = maxVal - minVal || 1;

  const points = values.map((v, i) => {
    const x = padX + (i / Math.max(values.length - 1, 1)) * (width - 2 * padX);
    const y = padY + (1 - (v - minVal) / range) * (height - 2 * padY);
    return { x, y };
  });

  const lineD = points.map((p, i) => `${i === 0 ? "M" : "L"}${p.x},${p.y}`).join(" ");

  const zeroY = padY + (1 - (0 - minVal) / range) * (height - 2 * padY);

  const areaD = lineD + ` L${points[points.length - 1]?.x ?? width},${zeroY} L${points[0]?.x ?? 0},${zeroY} Z`;

  const lastValue = values[values.length - 1] ?? 0;
  const isPositive = lastValue >= 0;
  const strokeColor = isPositive ? "var(--color-tint-green)" : "var(--color-tint-red)";
  const fillColor = isPositive ? "rgba(52,199,89,0.12)" : "rgba(255,59,48,0.12)";

  const periodLabelMap: Record<Period, string> = { "7d": labels.d7, "30d": labels.d30, "90d": labels.d90 };

  if (sorted.length === 0) return null;

  return (
    <article className="rounded-2xl border border-separator/50 bg-surface-secondary p-5 shadow-elevation-1 transition-all duration-300 ease-apple sm:p-6">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-title-3 text-label-primary">{labels.title}</h3>
        <div className="flex gap-1.5">
          {(["7d", "30d", "90d"] as Period[]).map((p) => (
            <button
              key={p}
              onClick={() => setPeriod(p)}
              className={[
                "rounded-full px-3 py-1 text-caption-1 font-semibold transition-all duration-200 ease-apple",
                period === p
                  ? "bg-label-primary text-surface-primary shadow-sm"
                  : "bg-surface-tertiary/80 text-label-secondary hover:bg-surface-tertiary",
              ].join(" ")}
            >
              {periodLabelMap[p]}
            </button>
          ))}
        </div>
      </div>

      <svg viewBox={`0 0 ${width} ${height}`} className="w-full" preserveAspectRatio="none">
        {/* Zero line */}
        <line
          x1={padX}
          y1={zeroY}
          x2={width - padX}
          y2={zeroY}
          stroke="var(--color-separator)"
          strokeWidth="1"
          strokeDasharray="4 4"
          opacity="0.5"
        />
        {/* Fill area */}
        {points.length > 1 && <path d={areaD} fill={fillColor} />}
        {/* Line */}
        {points.length > 1 && (
          <path d={lineD} fill="none" stroke={strokeColor} strokeWidth="2" strokeLinejoin="round" strokeLinecap="round" />
        )}
      </svg>

      <div className="mt-2 flex items-center justify-between text-caption-1 text-label-quaternary">
        <span>{sorted[0]?.date}</span>
        <span className={`font-semibold tabular-nums ${isPositive ? "text-tint-green" : "text-tint-red"}`}>
          {lastValue >= 0 ? "+" : ""}{lastValue.toFixed(2)}
        </span>
        <span>{sorted[sorted.length - 1]?.date}</span>
      </div>
    </article>
  );
}
