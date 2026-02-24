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
    <div className="flex flex-col gap-0.5">
      <span className="text-caption-1 text-label-tertiary">{label}</span>
      <span
        className={[
          size === "large" ? "text-title-1" : "text-title-3",
          valueColor,
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
    <span className="text-footnote text-label-secondary">
      {label}{" "}
      <span className={`font-medium ${color}`}>{value}</span>
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
    <div className="flex flex-col gap-1">
      <span className="text-caption-1 text-label-tertiary">Risk</span>
      <div className="flex items-center gap-2">
        <div className="h-1.5 w-20 rounded-full bg-surface-tertiary">
          <div
            className={`h-1.5 rounded-full ${color} transition-all duration-500 ease-apple-decel`}
            style={{ width: `${pct}%` }}
          />
        </div>
        <span className="text-caption-1 font-medium text-label-secondary">{label}</span>
      </div>
    </div>
  );
}
