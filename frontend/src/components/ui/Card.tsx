import type { ReactNode } from "react";

/* ── Card Variants ── */

type Variant = "grouped" | "prominent" | "flat";

const variantClasses: Record<Variant, string> = {
  grouped:
    "bg-surface-secondary rounded-lg shadow-elevation-1 overflow-hidden",
  prominent:
    "bg-surface-secondary rounded-xl shadow-elevation-2 overflow-hidden",
  flat:
    "bg-surface-tertiary rounded-lg overflow-hidden",
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
        "transition-shadow duration-200 ease-apple",
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
    "flex items-center gap-3 px-4 py-3",
    "border-b border-separator last:border-b-0",
    "min-h-[44px]",
    "transition-colors duration-150 ease-apple",
    href ? "hover:bg-surface-tertiary cursor-pointer" : "",
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
    <div className="flex items-center justify-between pb-2">
      <h2 className="text-title-3 text-label-primary">{title}</h2>
      {action}
    </div>
  );
}
