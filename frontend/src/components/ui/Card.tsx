import type { ReactNode } from "react";

/* ── Card Variants ── */

type Variant = "grouped" | "prominent" | "flat" | "glass";

const variantClasses: Record<Variant, string> = {
  grouped:
    "bg-surface-secondary rounded-xl shadow-elevation-1 overflow-hidden hover:shadow-[var(--card-hover-shadow)]",
  prominent:
    "bg-surface-secondary rounded-2xl shadow-elevation-2 overflow-hidden hover:shadow-[var(--card-hover-shadow)]",
  flat:
    "bg-surface-tertiary/60 rounded-xl overflow-hidden border border-separator/50",
  glass:
    "backdrop-blur-xl bg-surface-secondary/60 rounded-2xl shadow-elevation-2 overflow-hidden border border-separator/30",
};

export function Card({
  variant = "grouped",
  padding = true,
  className = "",
  children,
}: {
  variant?: Variant;
  padding?: boolean;
  className?: string;
  children: ReactNode;
}) {
  return (
    <article
      className={[
        variantClasses[variant],
        padding ? "p-5 sm:p-6" : "",
        "transition-all duration-300 ease-apple",
        className,
      ].join(" ")}
    >
      {children}
    </article>
  );
}

/* ── Card Row (for grouped lists) ── */

export function CardRow({
  href,
  className = "",
  children,
}: {
  href?: string;
  className?: string;
  children: ReactNode;
}) {
  const classes = [
    "flex items-center gap-3 px-5 py-3.5",
    "border-b border-separator/60 last:border-b-0",
    "min-h-[48px]",
    "transition-all duration-200 ease-apple",
    href ? "hover:bg-surface-tertiary/70 cursor-pointer active:bg-surface-tertiary" : "",
    className,
  ].join(" ");

  if (href) {
    const Link = require("next/link").default;
    return (
      <Link href={href} className={classes}>
        {children}
      </Link>
    );
  }
  return <div className={classes}>{children}</div>;
}

/* ── Section Header ── */

export function SectionHeader({
  title,
  action,
}: {
  title: string;
  action?: ReactNode;
}) {
  return (
    <div className="flex items-center justify-between pb-3">
      <h2 className="text-title-3 text-label-primary">{title}</h2>
      {action}
    </div>
  );
}
