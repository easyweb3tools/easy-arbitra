type StatColor = "default" | "positive" | "negative";

function resolveColor(value: number | undefined, color: StatColor): string {
  if (color === "positive") return "text-tint-green";
  if (color === "negative") return "text-tint-red";
  if (color === "default" && value !== undefined) {
    if (value > 0) return "text-tint-green";
    if (value < 0) return "text-tint-red";
  }
  return "text-label-primary";
}

export function StatCell({
  label,
  value,
  numericValue,
  color = "default",
  size = "normal",
}: {
  label: string;
  value: string;
  numericValue?: number;
  color?: StatColor;
  size?: "normal" | "large";
}) {
  const valueColor = resolveColor(numericValue, color);
  return (
    <div className="flex flex-col gap-1">
      <span className="text-caption-1 font-medium text-label-tertiary tracking-wide uppercase">{label}</span>
      <span
        className={[
          size === "large" ? "text-title-1 tracking-tight" : "text-title-3",
          valueColor,
          "tabular-nums",
        ].join(" ")}
      >
        {value}
      </span>
    </div>
  );
}

/* ── Inline Stat (compact, for list rows) ── */

export function InlineStat({
  label,
  value,
  numericValue,
}: {
  label: string;
  value: string;
  numericValue?: number;
}) {
  const color =
    numericValue !== undefined
      ? numericValue > 0
        ? "text-tint-green"
        : numericValue < 0
        ? "text-tint-red"
        : "text-label-secondary"
      : "text-label-secondary";

  return (
    <span className="text-footnote text-label-tertiary">
      {label}{" "}
      <span className={`font-semibold tabular-nums ${color}`}>{value}</span>
    </span>
  );
}

/* ── Risk Level Bar ── */

export function RiskBar({
  level,
  label,
}: {
  level: "low" | "medium" | "high";
  label: string;
}) {
  const pct = level === "low" ? 30 : level === "medium" ? 60 : 85;
  const color =
    level === "low"
      ? "bg-tint-green"
      : level === "medium"
      ? "bg-tint-orange"
      : "bg-tint-red";

  return (
    <div className="flex flex-col gap-1.5">
      <span className="text-caption-1 font-medium text-label-tertiary tracking-wide uppercase">Risk</span>
      <div className="flex items-center gap-2.5">
        <div className="h-2 w-24 rounded-full bg-surface-tertiary overflow-hidden">
          <div
            className={`h-2 rounded-full ${color} transition-all duration-700 ease-apple-decel`}
            style={{ width: `${pct}%` }}
          />
        </div>
        <span className="text-caption-1 font-semibold text-label-secondary">{label}</span>
      </div>
    </div>
  );
}
