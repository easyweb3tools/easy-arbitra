import type { ReactNode } from "react";

/* ── Tier Badge (pool level) ── */

type TierLevel = "star" | "strategy" | "observation";

const tierConfig: Record<TierLevel, { label: string; labelZh: string; className: string }> = {
  star: {
    label: "Star",
    labelZh: "明星池",
    className: "bg-gradient-to-br from-tint-gold to-tint-orange text-black font-bold shadow-[0_2px_8px_rgba(255,214,10,0.3)]",
  },
  strategy: {
    label: "Strategy",
    labelZh: "策略池",
    className: "bg-surface-tertiary border border-separator text-label-secondary font-semibold",
  },
  observation: {
    label: "Observation",
    labelZh: "观察池",
    className: "text-label-tertiary font-normal",
  },
};

export function TierBadge({ tier, locale = "en" }: { tier: TierLevel; locale?: "en" | "zh" }) {
  const cfg = tierConfig[tier] || tierConfig.observation;
  return (
    <span className={`inline-flex h-6 items-center rounded-full px-2.5 text-caption-1 ${cfg.className}`}>
      {locale === "zh" ? cfg.labelZh : cfg.label}
    </span>
  );
}

/* ── Status Badge (with colored dot) ── */

type StatusColor = "green" | "red" | "orange" | "blue" | "gray";

const dotColors: Record<StatusColor, string> = {
  green: "bg-tint-green",
  red: "bg-tint-red",
  orange: "bg-tint-orange",
  blue: "bg-tint-blue",
  gray: "bg-label-quaternary",
};

export function StatusBadge({ color, children }: { color: StatusColor; children: ReactNode }) {
  return (
    <span className="inline-flex items-center gap-1.5 text-caption-1 text-label-secondary">
      <span className={`h-1.5 w-1.5 rounded-full ${dotColors[color]}`} />
      {children}
    </span>
  );
}

/* ── Category Tag ── */

export function CategoryTag({ children, color = "blue" }: { children: ReactNode; color?: StatusColor }) {
  const colorMap: Record<StatusColor, string> = {
    blue: "bg-tint-blue/10 text-tint-blue",
    green: "bg-tint-green/10 text-tint-green",
    red: "bg-tint-red/10 text-tint-red",
    orange: "bg-tint-orange/10 text-tint-orange",
    gray: "bg-surface-tertiary text-label-secondary",
  };
  return (
    <span className={`inline-flex h-6 items-center rounded-full px-2.5 text-caption-1 font-medium ${colorMap[color]}`}>
      {children}
    </span>
  );
}
